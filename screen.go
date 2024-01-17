package adagui

import (
    "flag"
    "log"
    "time"
    "github.com/stefan-muehlebach/adatft"
    "github.com/stefan-muehlebach/adagui/touch"
    "github.com/stefan-muehlebach/gg/geom"
)

var (
    numScreen int = 0
    screen *Screen
    rotation adatft.RotationType = adatft.Rotate000
)

func init() {
    flag.Var(&rotation, "rotation", "display rotation")
}

// Dies ist die Datenstruktur, welche das TFT-Display aus einer hoeheren
// Abstraktion beschreibt. Diese Struktur darf es nur einmal (1) in einer
// Applikation geben.
type Screen struct {
    disp   *adatft.Display
    touch  *adatft.Touch
    window *Window
}

// Mit NewScreen wird ein neues Screen-Objekt erzeugt und alle technischen
// Objekte in Zusammenhang mit der Ansteuerung des Bildschirm und Touch-
// Screens erzeugt. Aktuell darf es nur ein (1) solches Objekt geben - ein
// mehrfaches Aufrufen von NewScreen führt zu einem Abbruch der Applikation.
func NewScreen() (*Screen) {
    if numScreen > 0 {
        log.Fatal("there is already a 'Screen' object in this application")
    }
    numScreen += 1
    s := &Screen{}
//    adatft.Init()
    s.disp  = adatft.OpenDisplay(rotation)
    s.touch = adatft.OpenTouch()
    s.touch.ReadConfig()
    s.window = nil
    screen = s
    return s
}

// Mit CurrentScreen wird die Referenz auf den aktuellen (einzigen) Bildschirm
// retourniert. Man könnte dies auch über eine globale Variable lösen.
func CurrentScreen() (*Screen) {
    return screen
}

// Mit NewWindow wird ein neues Fenster erzeugt. Im Gegensatz zum Screen
// darf es in einer Applikation beliebig viele Fenster geben, von denen jedoch
// nur eines sichtbar, resp. aktiv ist.
func (s *Screen) NewWindow() (*Window) {
    w := newWindow(s)
    return w
}

// Mit SetWindow wird das übergebene Fenster zum sichtbaren und aktiven
// Fenster. Nur aktive Fenster erhalten die Touch-Events vom Touchscreen und
// nur aktive Fenster werden dargestellt.
func (s *Screen) SetWindow(w *Window) {
    if s.window == w {
        return
    }
    if s.window != nil {
        s.window.stage = StageAlive
    }
    s.window = w
    s.window.stage = StageVisible
    s.window.Repaint()
}

func (s *Screen) Window() (*Window) {
    return s.window
}

// Mit Run schliesslich wird der MainEvent-Loop der Applikation gestartet,
// das aktive Fenster wird dargestellt und mit Touch-Events beliefert.
// Wichtig: diese Methode kehrt nicht zurück, solange die Applikation läuft.
// Ein Aufruf dieser Methode via Go-Routine ist nicht sinnvoll, da sonst
// die Applikation gar nie richtig läuft (siehe auch Methode Quit).
func (s *Screen) Run() {
    s.eventThread()
    if s.window != nil {
        s.window.Close()
    }
    s.disp.Close()
}

// Mit Quit wird die Applikation (d.h. der MainEvent-Loop) terminiert.
// Da Run im Main-Thread gestartet wird und während der Laufzeit der
// Applikation nicht zurückkehrt, muss diese Methode aus einer weiteren
// Go-Routine (bspw. dem Callback-Handler eines Buttons) aufgerufen werden.
func (s *Screen) Quit() {
    s.touch.Close()
}

// In dieser Methode schliesslich spielt die Musik: vom Touch-Screen werden
// laufend Events empfangen, ggf. 'veredelt' (bspw. werden hier LongPress,
// Tap oder DoubleTap Events generiert) und dem aktiven Fenster zur
// Verarbeitung weitergeleitet. Die Positionsdaten aus den Touch-Events
// beziehen sich auf den gesamten Bildschirm. Die Transformation von
// Koordianten in Objekt-relative Daten erfolgt im Objekt Window!
func (s *Screen) eventThread() {
    var evt, tapEvt touch.Event
    var seqNumber int = 0

    for tchEvt := range s.touch.EventQ {
        //log.Printf("screen: receive new event from queue\n")
        switch tchEvt.Type {
        case adatft.PenPress:
            seqNumber++
            evt.Type        = touch.TypePress
            evt.SeqNumber   = seqNumber
            evt.LongPressed = false
            evt.InitTime    = time.Now()
            evt.InitPos     = geom.NewPoint(tchEvt.X, tchEvt.Y)
            evt.Time        = evt.InitTime
            evt.Pos         = evt.InitPos

            // Setze eine verzoegerte Go-Routine zur Erkennung des Events
            // 'LongPress'.
            //
            go func(seqNr int) {
                time.Sleep(touch.LongPressThreshold)
                if seqNr == seqNumber &&
                        evt.Type != touch.TypeRelease &&
                        evt.InitPos.Distance(evt.Pos) <= touch.NearThreshold {
                    evt.LongPressed = true
                    newEvent := evt
                    newEvent.Type = touch.TypeLongPress
                    newEvent.Time = time.Now()
                    s.window.eventQ <- newEvent
                }
            }(seqNumber)
            s.window.eventQ <- evt

        case adatft.PenDrag:
            evt.Type = touch.TypeDrag
            evt.Time = time.Now()
            evt.Pos  = geom.NewPoint(tchEvt.X, tchEvt.Y)
            s.window.eventQ <- evt

        case adatft.PenRelease:
            evt.Type = touch.TypeRelease
            evt.Time = time.Now()
            evt.Pos  = geom.NewPoint(tchEvt.X, tchEvt.Y)
            s.window.eventQ <- evt

            if evt.InitPos.Distance(evt.Pos) <= touch.NearThreshold {

                // An dieser Stelle steht fest: es wurde ein korrekter Tap
                // erkannt. Die Frage ist noch: war es ein DoubleTap?
                if tapEvt.Type == touch.TypeTap &&
                        evt.Time.Sub(tapEvt.Time) < touch.DoubleTapDuration &&
                        evt.Pos.Distance(tapEvt.Pos) <= touch.NearThreshold {
                    tapEvt = evt
                    tapEvt.Type = touch.TypeDoubleTap
                } else {
                    tapEvt = evt
                    tapEvt.Type = touch.TypeTap
                }
                s.window.eventQ <- tapEvt
            }
        }
    }
}

func (s *Screen) StartAnimation(a *Animation) {

}

func (s *Screen) StopAnimation(a *Animation) {

}


