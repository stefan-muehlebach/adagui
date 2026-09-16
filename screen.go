package adagui

import (
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/stefan-muehlebach/adagui/point"
	"github.com/stefan-muehlebach/adatft"
	"github.com/stefan-muehlebach/gg/geom"
)

var (
	screen      *Screen = nil
	refreshRate         = 30 * time.Millisecond
)

// Dies ist die Datenstruktur, welche das TFT-Display aus einer hoeheren
// Abstraktion beschreibt. Diese Struktur darf es nur einmal (1) in einer
// Applikation geben.
type Screen struct {
	disp                     *adatft.Display
	point                    *Pointer
	window                   *Window
	paintTicker              *time.Ticker
	paintCloseQ, eventCloseQ chan bool
	wg                       sync.WaitGroup
	mutex                    *sync.Mutex
}

// Mit NewScreen wird ein neues Screen-Objekt erzeugt und alle technischen
// Objekte in Zusammenhang mit der Ansteuerung des Bildschirm und Point-
// Screens erzeugt. Aktuell darf es nur ein (1) solches Objekt geben - ein
// mehrfaches Aufrufen von NewScreen führt zu einem Abbruch der Applikation.
func NewScreen(rotation adatft.RotationType) *Screen {
	if screen != nil {
		log.Fatal("there is already a 'Screen' object in this application")
	}
	s := &Screen{}
	s.disp = adatft.OpenDisplay(rotation)
	s.point = OpenPointer()
	s.point.SetPosRange(geom.Rectangle{Max: s.disp.DrawBounds().Size()})
	s.point.StartEvents()
	s.paintTicker = time.NewTicker(refreshRate)
	s.window = nil
	s.paintCloseQ = make(chan bool)
	s.eventCloseQ = make(chan bool)
	s.wg.Add(2)
	s.mutex = &sync.Mutex{}

	screen = s

	go s.paintThread()

	return s
}

// Mit CurrentScreen wird die Referenz auf den aktuellen (einzigen) Bildschirm
// retourniert. Man könnte dies auch über eine globale Variable lösen.
func CurrentScreen() *Screen {
	return screen
}

func (s *Screen) SaveScreenshot(name string) {
	s.Window().SaveScreenshot(name)
}

func (s *Screen) SaveMovie(name string) {
	fileName := fmt.Sprintf(name, 0)
	s.Window().SaveScreenshot(fileName)
}

// Mit NewWindow wird ein neues Fenster erzeugt. Im Gegensatz zum Screen
// darf es in einer Applikation beliebig viele Fenster geben, von denen jedoch
// nur eines sichtbar, resp. aktiv ist.
func (s *Screen) NewWindow() *Window {
	w := newWindow(s)
	return w
}

// Liefert das aktuell angezeigte Window zurueck.
func (s *Screen) Window() *Window {
	return s.window
}

// Mit SetWindow wird das übergebene Fenster zum sichtbaren und aktiven
// Fenster. Nur aktive Fenster erhalten die Point-Events vom Pointscreen und
// nur aktive Fenster werden dargestellt.
func (s *Screen) SetWindow(w *Window) {
	if s.window == w {
		return
	}
	s.mutex.Lock()
	if s.window != nil {
		s.window.stage = StageAlive
	}
	s.window = w
	s.window.stage = StageVisible
	s.mutex.Unlock()
	s.Repaint()
}

// Mit Run schliesslich wird der MainEvent-Loop der Applikation gestartet,
// das aktive Fenster wird dargestellt und mit Point-Events beliefert.
// Wichtig: diese Methode kehrt nicht zurück, solange die Applikation läuft.
// Ein Aufruf dieser Methode via Go-Routine ist nicht sinnvoll, da sonst
// die Applikation gar nie richtig läuft (siehe auch Methode Quit).
func (s *Screen) Run() {
	s.eventThread()
	s.disp.Close()
}

// Mit Quit wird die Applikation (d.h. der MainEvent-Loop) terminiert.
// Da Run im Main-Thread gestartet wird und während der Laufzeit der
// Applikation nicht zurückkehrt, muss diese Methode aus einer weiteren
// Go-Routine (bspw. dem Callback-Handler eines Buttons) aufgerufen werden.
func (s *Screen) Quit() {
	go s.quit()
}

func (s *Screen) quit() {
	//fmt.Printf("Screen.Quit() has been called\n")
	if s.window != nil {
		//fmt.Printf("Screen.Quit()   closing an open window\n")
		s.window.Close()
		//fmt.Printf("Screen.Quit()   done!\n")
		s.window = nil
	}
	//fmt.Printf("Screen.Quit()   send close to the event thread\n")
	s.eventCloseQ <- true
	//fmt.Printf("Screen.Quit()   send close to the paint thread\n")
	s.paintCloseQ <- true
	//fmt.Printf("Screen.Quit()   wait for the threads to complete\n")
	s.wg.Wait()
	//fmt.Printf("Screen.Quit()   done!\n")
}

func (s *Screen) StopPaint() {
	s.paintTicker.Stop()
}

func (s *Screen) ContPaint() {
	s.paintTicker.Reset(refreshRate)
}

func (s *Screen) Repaint() {
	if s.window == nil || s.window.stage != StageVisible {
		return
	}
	if s.window.Repaint() {
		s.point.Draw(s.window.gc)
		s.disp.Draw(s.window.gc.Image())
	}
	
}

func (s *Screen) paintThread() {
PAINT_LOOP:
	for {
		select {
		case <-s.paintCloseQ:
			break PAINT_LOOP
		case <-s.paintTicker.C:
			s.Repaint()
		}
	}
	s.wg.Done()
}

// In dieser Methode schliesslich spielt die Musik: vom Point-Screen werden
// laufend Events empfangen, ggf. 'veredelt' (bspw. werden hier LongPress,
// Tap oder DoubleTap Events generiert) und dem aktiven Fenster zur
// Verarbeitung weitergeleitet. Die Positionsdaten aus den Point-Events
// beziehen sich auf den gesamten Bildschirm. Die Transformation von
// Koordianten in Objekt-relative Daten erfolgt im Objekt Window!
// Gestoppt wird dieser Thread durch das Schliessen der Event-Queue.
func (s *Screen) eventThread() {
	var evt, clkEvt point.Event
	var seqNumber int = 0

EVENT_LOOP:
	for {
		select {
		case <-s.eventCloseQ:
			break EVENT_LOOP
		case ptrEvt := <-s.point.EventQ:
			log.Printf("[%d]: %10s: %v", ptrEvt.Time.UnixMilli(),
				ptrEvt.Type, ptrEvt)
			switch ptrEvt.Type {
			case adatft.PointerPress:
				seqNumber++
				evt.Type = point.TypePress
				evt.SeqNumber = seqNumber
				evt.LongPressed = false
				evt.InitTime = time.Now()
				evt.InitPos = ptrEvt.Pos
				evt.Time = evt.InitTime
				evt.Pos = evt.InitPos

				// Setze eine verzoegerte Go-Routine zur Erkennung des Events
				// 'LongPress'.
				//
				go func(seqNr int) {
					time.Sleep(point.LongPressThreshold)
					if seqNr == seqNumber &&
						evt.Type != point.TypeRelease &&
						evt.InitPos.Distance(evt.Pos) <= point.NearThreshold {
						evt.LongPressed = true
						newEvent := evt
						newEvent.Type = point.TypeLongPress
						newEvent.Time = time.Now()
						s.window.eventQ <- newEvent
					}
				}(seqNumber)
				s.window.eventQ <- evt

			case adatft.PointerDrag:
				evt.Type = point.TypeDrag
				evt.Time = time.Now()
				evt.Pos = ptrEvt.Pos
				s.window.eventQ <- evt

			case adatft.PointerRelease:
				evt.Type = point.TypeRelease
				evt.Time = time.Now()
				evt.Pos = ptrEvt.Pos
				s.window.eventQ <- evt

				if evt.InitPos.Distance(evt.Pos) <= point.NearThreshold {

					// An dieser Stelle steht fest: es wurde ein korrekter Tap
					// erkannt. Die Frage ist noch: war es ein DoubleTap?
					if clkEvt.Type == point.TypeClick &&
						evt.Time.Sub(clkEvt.Time) < point.DoubleClickDuration &&
						evt.Pos.Distance(clkEvt.Pos) <= point.NearThreshold {
						clkEvt = evt
						clkEvt.Type = point.TypeDoubleClick
					} else {
						clkEvt = evt
						clkEvt.Type = point.TypeClick
					}
					s.window.eventQ <- clkEvt
				}
			}
		}
	}
	s.wg.Done()
	//fmt.Printf("Screen.eventThread()   exits\n")
}

func (s *Screen) StartAnimation(a *Animation) {

}

func (s *Screen) StopAnimation(a *Animation) {

}
