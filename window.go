package adagui

import (
//    "fmt"
    "image/png"
    "log"
    "os"
    "sync"
    "time"
    "github.com/stefan-muehlebach/adatft"
    "github.com/stefan-muehlebach/adagui/touch"
    "github.com/stefan-muehlebach/gg"
    "github.com/stefan-muehlebach/gg/color"
    "github.com/stefan-muehlebach/gg/geom"
)

//----------------------------------------------------------------------------

type WindowStage uint32

const (
    StageDead WindowStage = iota
    StageAlive
    StageVisible
    StageFocused
    RefreshCycle = 30 * time.Millisecond
)

type Window struct {
    Rect geom.Rectangle
	Color color.Color
    s *Screen
    gc *gg.Context
    eventQ chan touch.Event
    eventCloseQ chan bool
    wg sync.WaitGroup
    root Node
    stage WindowStage
    mutex *sync.Mutex
}

// Dies ist die interne Funktion, welche der Screen beim Erzeugen einer neuen
// Window-Struktur aufruft.
func newWindow(s *Screen) (*Window) {
    w := &Window{}
    w.s = s
    width  := adatft.Width
    height := adatft.Height
    w.Rect = geom.NewRectangleWH(0.0, 0.0, float64(width), float64(height))
	w.Color = color.Black
    w.gc = gg.NewContext(width, height)
    w.eventQ = make(chan touch.Event)
    w.eventCloseQ = make(chan bool)
    w.wg.Add(1)
    w.stage  = StageAlive
    w.mutex  = &sync.Mutex{}

    go w.eventThread()

    return w
}

// Schliesst das Fenster.
func (w *Window) Close() {
    //fmt.Printf("Window.Close() has been called\n")
    //fmt.Printf("Window.Close()   send close to the event thread\n")
    w.eventCloseQ <- true
    //fmt.Printf("Window.Close()   wait for the threads to complete\n")
    w.wg.Wait()
    //fmt.Printf("Window.Close()   done!\n")
}

// In jedem Fenster muss es ein GUI-Element geben, welches an der obersten
// Stelle (der Wurzel, 'root') des SceneGraphs steht. Dies ist ueblicherweise
// ein Container-Widget.
func (w *Window) Root() (Node) {
    return w.root
}

func (w *Window) SetRoot(root Node) {
    n := root.Wrappee()
    if w.root != nil {
        n.Win = nil
    }
    w.root = root
    n.Win = w
    root.SetPos(w.Rect.Min)
    root.SetSize(w.Rect.Size())
}

func (w *Window) SaveScreenshot(fileName string) {
    fh, err := os.Create(fileName)
    if err != nil {
        log.Fatal(err)
    }
    defer fh.Close()

    w.mutex.Lock()
    defer w.mutex.Unlock()
    if err = png.Encode(fh, w.gc.Image()); err != nil {
        log.Fatal(err)
    }
}

// Mit dieser Methode wird ein Neuaufbau des Bildschirms angestossen. Ueber die
// interne Queue paintQ wird dem paintThread der Auftrag fuer den Neuaufbau
// gegeben. Diese Methode blockiert nie! Ist bereits ein Auftrag fuer den
// Neuaufbau in der Queue, dann ist soweit alles i.O. und wir sind sicher,
// dass auch unser Auftrag behandelt wird.
func (w *Window) Repaint(disp *adatft.Display) {
    w.gc.SetFillColor(w.Color)
    w.gc.Clear()
    w.mutex.Lock()
    w.root.Wrappee().Paint(w.gc)
    w.mutex.Unlock()
    disp.Draw(w.gc.Image())
}

// Mit dieser Go-Routine werden die Events vom Screen-Objekt empfangen und
// weiterverarbeitet.
func (w *Window) eventThread() {
    var target Node
    var onTarget bool

LOOP:
    for {
        //fmt.Printf("Window.eventThread() next iteration\n")
        select {
        case <- w.eventCloseQ:
            //fmt.Printf("Window.eventThread() got called on eventCloseQ\n")
            break LOOP
        case evt := <- w.eventQ:
            //fmt.Printf("Window.eventThread() new event received\n")
			// Ist kein root-Element vorhanden, dann wird das Event nicht weiter
			// verarbeitet und die Go-Routine wartet auf das naechste Event.
			if w.root == nil {
				continue
			}
			Debugf(Events, "event received: %v", evt)
			if evt.Type == touch.TypePress {
				target = w.root.SelectTarget(evt.Pos)
				Debugf(Events, "new target    : %T", target)
				onTarget = true
			}
			if target == nil {
				continue
			}
			evt.InitPos = target.Screen2Local(evt.InitPos)
			evt.Pos = target.Screen2Local(evt.Pos)
			Debugf(Events, "relative pos  : %v", evt.Pos)

			if evt.Type == touch.TypeDrag {
				if !target.Contains(evt.Pos) {
					if onTarget {
						onTarget = false
						newEvent := evt
						newEvent.Type = touch.TypeLeave
						w.mutex.Lock()
						target.OnInputEvent(newEvent)
						w.mutex.Unlock()
					}
				} else {
					if !onTarget {
						onTarget = true
						newEvent := evt
						newEvent.Type = touch.TypeEnter
						w.mutex.Lock()
						target.OnInputEvent(newEvent)
						w.mutex.Unlock()
					}
				}
			}

            //fmt.Printf("Window.eventThread() about to enter crit. section\n")
			w.mutex.Lock()
			target.OnInputEvent(evt)
			w.mutex.Unlock()
            //fmt.Printf("Window.eventThread()   leave crit. section\n")
        }
    }
    w.wg.Done()
    //fmt.Printf("Window.eventThread()   exits\n")
}
