package adagui

import (
    "image"
    "image/color"
    //"log"
    "sync"
    _ "time"
    "github.com/stefan-muehlebach/adatft"
    "github.com/stefan-muehlebach/adagui/touch"
    "github.com/stefan-muehlebach/gg/geom"
    "github.com/stefan-muehlebach/gg"
)

//----------------------------------------------------------------------------

type WindowStage uint32

const (
    StageDead WindowStage = iota
    StageAlive
    StageVisible
    StageFocused
)

type Window struct {
    Rect geom.Rectangle
    s *Screen
    gc *gg.Context
    paintQ chan bool
    eventQ chan touch.Event
    quitQ  chan bool
    root Node
    stage WindowStage
    mutex *sync.Mutex
}

func newWindow(s *Screen) (*Window) {
    w := &Window{}
    w.s = s
    width  := adatft.Width
    height := adatft.Height
    w.Rect = geom.NewRectangleWH(0.0, 0.0, float64(width), float64(height))
    w.gc = gg.NewContext(width, height)
    w.paintQ = make(chan bool, 1)
    w.eventQ = make(chan touch.Event)
    w.quitQ  = make(chan bool)
    w.stage  = StageAlive
    w.mutex  = &sync.Mutex{}

    go w.eventThread()
    go w.paintThread()

    return w
}

func (w *Window) Close() {
    close(w.eventQ)
    <- w.quitQ
    close(w.paintQ)
    <- w.quitQ
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

func (w *Window) Root() (Node) {
    return w.root
}

// Mit dieser Methode wird ein Neuaufbau des Bildschirms angestossen. Ueber die
// interne Queue paintQ wird dem paintThread der Auftrag fuer den Neuaufbau
// gegeben. Diese Methode blockiert nie! Ist bereits ein Auftrag fuer den
// Neuaufbau in der Queue, dann ist soweit alles i.O. und wir sind sicher,
// dass auch unser Auftrag behandelt wird.
func (w *Window) Repaint() {
    //log.Printf("%T: Repaint()", w)
    select {
        case w.paintQ <- true:
        default:
    }
}

// Paint wird aufgerufen, um das Fenster und alle Objekte, die damit verbunden
// sind, auf dem Zeichen-Kontext des gg-Packages darzustellen.
func (w *Window) paintThread() {
    for range w.paintQ {
        if w.stage != StageVisible {
            continue
        }
        w.gc.SetFillColor(color.Black)
        w.gc.Clear()
        w.gc.Identity()
        w.mutex.Lock()
        w.root.Wrappee().Paint(w.gc)
        w.mutex.Unlock()
        w.s.disp.Draw(w.s.disp.Bounds(), w.gc.Image(), image.Point{})
    }
    w.quitQ <- true
}

// Mit dieser Go-Routine werden die Events vom Screen-Objekt empfangen und
// weiterverarbeitet.
func (w *Window) eventThread() {
    var target Node
    var onTarget bool

    for evt := range w.eventQ {

        //log.Printf("window: event from screen received\n")
        if evt.Type == touch.TypePress {
            target = w.root.SelectTarget(evt.Pos)
            //log.Printf("SelectTarget: %T, %v\n", target, evt.Pos)
            if target == nil {
                continue
            }
            onTarget = true
        }
        evt.InitPos = target.Screen2Local(evt.InitPos)
        evt.Pos = target.Screen2Local(evt.Pos)
        //log.Printf("SelectTarget: local coord %v\n", evt.Pos)

        if evt.Type == touch.TypeDrag {
            if !target.Contains(evt.Pos) {
                if onTarget {
                    onTarget = false
                    newEvent := evt
                    newEvent.Type = touch.TypeLeave
                    target.OnInputEvent(newEvent)
                }
            } else {
                if !onTarget {
                    onTarget = true
                    newEvent := evt
                    newEvent.Type = touch.TypeEnter
                    target.OnInputEvent(newEvent)
                }
            }
        }

        //log.Printf("SelectTarget: sending %v to %T", evt, target)
        w.mutex.Lock()
        target.OnInputEvent(evt)
        w.mutex.Unlock()

        if w.root.Wrappee().Marks.NeedsPaint() {
            w.Repaint()
        }
    }
    w.quitQ <- true
}

