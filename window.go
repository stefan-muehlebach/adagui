package adagui

import (
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
    s *Screen
    gc *gg.Context
    paintCloseQ chan bool
    paintTicker *time.Ticker
    eventQ chan touch.Event
    quitQ  chan bool
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
    w.gc = gg.NewContext(width, height)
    //w.paintCloseQ = make(chan bool)
    //w.paintTicker = time.NewTicker(RefreshCycle)
    w.eventQ = make(chan touch.Event)
    w.quitQ  = make(chan bool)
    w.stage  = StageAlive
    w.mutex  = &sync.Mutex{}

    go w.eventThread()
    //go w.paintThread()

    return w
}

// Schliesst das Fenster.
func (w *Window) Close() {
    close(w.eventQ)
    <- w.quitQ
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
    w.gc.SetFillColor(color.Black)
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

    for evt := range w.eventQ {

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

        w.mutex.Lock()
        target.OnInputEvent(evt)
        w.mutex.Unlock()
    }
    w.quitQ <- true
}
