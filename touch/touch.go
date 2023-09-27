package touch

import (
    "fmt"
    //"log"
    "time"
    "mju.net/geom"
)

// Mit diesem Datentyp werden die unterschiedlichen Event-Arten abgebildet,
// welche durch das Druecken auf den Bildschirm an die GUI-Elemente gesendet
// werden koennen. Es ist Aufgabe der Objekte 'Screen' und 'Window', aus den
// rohen Ereginisse vom Touchscreen (Press, Drag, Release) diese Events zu
// erzeugen.
// Eine Interaktion auf dem Bildschirm laeuft grundsaetzlich wie folgt ab.
// In Klammern sind Events gesetzt, die optional erzeugt werden:
//
//   TypePress
//    |
//   TypeDrag
//   ...
//    |
//   (TypeLongPress)
//    |
//   TypeDrag
//   ...
//    |
//   TypeRelease
//
type Type uint32

const (
    // TypePress wird erzeugt, wenn der Finger oder Stift auf ein Objekt
    // gedrueckt wird. Dieses Event ist vergleichbar mit einem Press-Event
    // durch die Maus.
    TypePress Type = iota
    // TypeRelease wird erzeugt, wenn man den Finger oder Stift wieder vom
    // Bildschirm hebt - vergleichbar mit dem Loslassen der Maustaste.
    TypeRelease
    // Im Unterschied zur Maus, kann ein 'Wandern' des Druckpunktes (Finger
    // oder Maus) nur erkannt werden, wenn man den Bildschirm beruehrt.
    // D.h. es gibt nur Drag-Events - ein Move-Event ist unbekannt.
    TypeDrag
    // Wird laengere Zeit auf ein bestimmtes Objekt gedrueckt, erzeugt dies
    // das TypeLongPress-Event (siehe auch die Konstanten LongPressThreshold
    // und NearThreshold).
    TypeLongPress
    // Verlaesst oder Betritt der Druckpunkt beim Wandern ein Objekt, dann
    // werden die Events TypeEnter, resp. TypeLeave gesendet.
    TypeEnter
    TypeLeave
    // Ein Tap, bzw. DoubleTap entspricht dem Klicken resp. Doppelklicken mit
    // der Maus. Es gibt sowohl zeitliche, als auch raeumliche Grenzen, wann
    // ein Tap, resp. DoubleTap erzeugt wird (siehe Konstanten weiter unten).
    TypeTap
    TypeDoubleTap
    numEvents

    TapDuration          = 200 * time.Millisecond
    DoubleTapDuration    = 200 * time.Millisecond
    LongPressThreshold   = 400 * time.Millisecond
    NearThreshold        = 8.0
    DragRadThreshold     = 3.0
)

func (t Type) String() (string) {
    switch t {
    case TypePress:
        return "Press"
    case TypeRelease:
        return "Release"
    case TypeDrag:
        return "Drag"
    case TypeLongPress:
        return "LongPress"
    case TypeEnter:
        return "Enter"
    case TypeLeave:
        return "Leave"
    case TypeTap:
        return "Tap"
    case TypeDoubleTap:
        return "DoubleTap"
    default:
        return "(Unknown Type)"
    }
}

// In diesem Datentyp ist alles zusammengefasst, was an Touch-Events an die
// applikatorischen Elemente gesendet werden kann.
// (to think about): Ev. sollte das Ziel des Events auch in dieser Struktur
// abgelegt werden..
type Event struct {
    // Der Typ des Events (moegliche Typen: die die Konstanten TypeXXX).
    Type Type
    // Alle Events werden sequenziell durchnumeriert, damit man bspw. das
    // Release-Event mit einem vorgaengigen Press-Event in Verbindung bringen
    // kann.
    SeqNumber int
    // Dieses Feld wird auf true gesetzt, sobald ein LongPressed-Ereignis
    // erkannt wird.
    LongPressed bool
    // In InitTime und InitPos werden Zeitpunkt und Position des Press-Events
    // festgehalten.
    InitTime time.Time
    InitPos  geom.Point
    // Wohingegen Time und Pos die Zeit und die Position des aktuellen Events
    // enthalten.
    Time time.Time
    Pos  geom.Point
}

// Fuer Debugging implementiert dieser Datentyp das Stringer-Interface.
func (evt Event) String() (string) {
    return fmt.Sprintf("%v %d %v %s %v %s %v", evt.Type, evt.SeqNumber,
            evt.LongPressed, evt.InitTime.Format("15:04:05.000000"),
            evt.InitPos, evt.Time.Format("15:04:05.000000"), evt.Pos)
}

type TouchFunction func(evt Event)

type TouchEmbed struct {
    touchFuncList [numEvents]TouchFunction
}

func (m *TouchEmbed) OnInputEvent(evt Event) {
    // log.Printf("TouchEmbed.OnInputEvent(): %v", evt)
    m.CallTouchFunc(evt)
}

func (m *TouchEmbed) SetTouchFunc(typ Type, fnc TouchFunction) {
    m.touchFuncList[typ] = fnc
}

func (m *TouchEmbed) CallTouchFunc(evt Event) {
    if fnc := m.touchFuncList[evt.Type]; fnc != nil {
        fnc(evt)
    }
}

func (m *TouchEmbed) SetOnPress(fnc TouchFunction) {
    m.SetTouchFunc(TypePress, fnc)
}

func (m *TouchEmbed) SetOnRelease(fnc TouchFunction) {
    m.SetTouchFunc(TypeRelease, fnc)
}

func (m *TouchEmbed) SetOnDrag(fnc TouchFunction) {
    m.SetTouchFunc(TypeDrag, fnc)
}

func (m *TouchEmbed) SetOnLongPress(fnc TouchFunction) {
    m.SetTouchFunc(TypeLongPress, fnc)
}

func (m *TouchEmbed) SetOnEnter(fnc TouchFunction) {
    m.SetTouchFunc(TypeEnter, fnc)
}

func (m *TouchEmbed) SetOnLeave(fnc TouchFunction) {
    m.SetTouchFunc(TypeLeave, fnc)
}

func (m *TouchEmbed) SetOnTap(fnc TouchFunction) {
    m.SetTouchFunc(TypeTap, fnc)
}

func (m *TouchEmbed) SetOnDoubleTap(fnc TouchFunction) {
    m.SetTouchFunc(TypeDoubleTap, fnc)
}

