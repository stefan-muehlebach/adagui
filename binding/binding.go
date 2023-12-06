package binding

import (
    "log"
//    "reflect"
    "sync"
)

// Die init-Funktion ist vorallem hier, damit der Import des Log-Packages
// nicht bei jedem Aktivieren von Debug-Meldungen (aus)kommentiert werden muss.
func init() {
    log.SetFlags(log.Lmicroseconds | log.Lmsgprefix)
    log.SetPrefix(": ")
}

// Als erstes folgen hier allgemeine, resp. generische Typen und Interfaces
//fuer das Binding.

// Mit DataItem ist das minimale Interface definiert welches ein Bind-Objekt
// implementieren muss. Einem Bind-Objekt koennen einerseits Objekte
// hinzugefuegt werden, die das DataListener Interface implementieren,
// andererseits aber auch direkt Funktionen vom Typ CallbackFunc. Diese beiden
// Moeglichkeiten koennen auch zusammen, d.h. synchron verwendet werden.
type DataItem interface {
    AddListener(DataListener)
    AddCallback(CallbackFunc)
    RemoveListener(DataListener)
    RemoveCallback(CallbackFunc)
}

// Alle Typen, welche ueber die Aenderungen von Bind-Objekten informiert werden
// wollen, muessen das DataListener-Interface implementieren. Im Wesentlichen
// eine einzige Methode (DataChanged), welche als Argument das veraenderte
// Bind-Objekt hat.
type DataListener interface {
    DataChanged(data DataItem)
}

// Mit NewDataListener kann einfach und bequem ein neues Listener-Objekt
// erzeugt werden.
func NewDataListener(fn func(data DataItem)) (DataListener) {
    return &listener{fn}
}

// Der (private) Typ listener ist nun eine der moeglichen Implementationen des
// DataListener-Interfaces...
type listener struct {
    callback func(data DataItem)
}

// ... und implementiert logischerweise das DataListener Interface.
func (l *listener) DataChanged(data DataItem) {
    l.callback(data)
}

// Neben den DataListeners gibt es noch eine schlanke Variante seinen Code bei
// Veraenderungen eines Bind-Objektes aufrufen zu lassen: man kann einfach eine
// Funktion/Methode, welche dem Typ CallbackFunc entspricht, mit AddCallback
// hinzufuegen.
type CallbackFunc func(data DataItem)

// base ist der Basistyp, welche die Methoden des DataItem-Interfaces
// implementiert. Er ist nicht oeffentlich, sondern wird von den weiter
// unten gezeigten konkreten Bind-Type verwendet.
type base struct {
    super DataItem
    // listeners und callbacks sind synchronisierte Maps, in welchen die
    // DataListener-Objekte oder Callback-Funktionen hinterlegt werden.
    listeners sync.Map
    callbacks sync.Map
    // Der Zugriff auf weitere Strukturen dieses Typs wird ueber das Mutex
    // lock gesteuert.
    lock sync.RWMutex
}

// Mit Init werden wichtige Initialisierungen vorgenommen. Init ist bei jeder
// Erstellung eines Bind-Objektes aufzurufen!
func (b *base) Init(super DataItem) {
    b.super = super
}

// AddListener fuegt ein DataListener hinzu und ruft die DataChanged-Methode
// auch gleich auf!
func (b *base) AddListener(l DataListener) {
    b.listeners.Store(l, true)
    l.DataChanged(b.super)
}

// Mit RemoveListener koennen DataListener wieder entfernt werden.
func (b *base) RemoveListener(l DataListener) {
    b.listeners.Delete(l)
}

// AddCallback fuegt eine Callback-Funktion hinzu und ruft diese auch gleich
// das erste mal auf.
func (b *base) AddCallback(f CallbackFunc) {
    b.callbacks.Store(&f, true)
    f(b.super)
}

// Mit RemoveCallback schliesslich koennen Callback-Funktionen wieder
// entfernt werden.
func (b *base) RemoveCallback(f CallbackFunc) {
    b.callbacks.Delete(&f)
}

// Die (private) Methode trigger wird immer dann aufgerufen, wenn sich der
// Wert des Bind-Objektes aendert. Diese Methode ruft die registrierten
// Listener-, resp. Callback-Methoden auf.
func (b *base) trigger() {
    b.listeners.Range(func(key, _ any) bool {
        go key.(DataListener).DataChanged(b.super)
        return true
    })
    b.callbacks.Range(func(f, _ any) bool {
        go (*f.(*CallbackFunc))(b.super)
        return true
    })
}

// Untyped---------------------------------------------------------------------

/*
type Untyped interface {
    DataItem
    Get() (interface{})
    Set(interface{})
}

type boundUntyped struct {
    base
    val reflect.Value
}

// Mit NewUntyped wird ein neues Bind-Objekt fuer einen interface{}-Wert erzeugt.
// Der Wert selber ist dabei von aussen nur ueber die Methoden Get und Set
// zugreifbar.
func NewUntyped() Untyped {
    var blank interface{} = nil
    b := &boundUntyped{val: reflect.ValueOf(&blank).Elem()}
    b.Init(b)
    return b
}

func (b *boundUntyped) Get() (interface{}) {
    b.lock.RLock()
    defer b.lock.RUnlock()
    return b.val.Interface()
}

func (b *boundUntyped) Set(val interface{}) {
    b.lock.Lock()
    defer b.lock.Unlock()
    if b.val.Interface() == val {
        return
    }
    b.val.Set(reflect.ValueOf(val))
    b.trigger()
}

// ExternalUntyped

type ExternalUntyped interface {
    Untyped
    Reload()
}

type boundExternalUntyped struct {
    boundUntyped
    old interface{}
}

// Mit BindUntyped kann ein Bind-Objekt ueber eine bereits bestehende
// interface{}-Variable erstellt werden. Dabei muss der Programmierer
// dafuer sorgen, dass Veraenderungen an der Variable mit der Methode
// Reload nach aussen bekannt gemacht werden.
func BindUntyped(v interface{}) ExternalUntyped {
    t := reflect.TypeOf(v)
    if t.Kind() != reflect.Ptr {
        log.Fatalf("Invalid type passed to BindUntyped, must be a pointer")
    }
    if v == nil {
        var blank interface{}
        v = &blank
    }
    b := &boundExternalUntyped{}
    b.val = reflect.ValueOf(v).Elem()
    b.old = b.val.Interface()
    b.Init(b)
    return b
}

func (b *boundExternalUntyped) Set(val interface{}) {
    b.lock.Lock()
    defer b.lock.Unlock()
    if b.old == val {
        return
    }
    b.val.Set(reflect.ValueOf(val))
    b.old = val
    b.trigger()
}

func (b *boundExternalUntyped) Reload() {
    b.Set(b.val.Interface())
}
*/

