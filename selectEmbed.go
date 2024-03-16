package adagui

import (
    "github.com/stefan-muehlebach/adagui/binding"
    "github.com/stefan-muehlebach/adagui/touch"
)

// Mit diesem Embed erhaelt ein Widget die Moeglichkeit, "ausgewaehlt"
// oder "selektiert" zu werden.
type SelectEmbed struct {
    BindVar binding.Untyped
    BindValue any
    node Node
}

// Initialisert wird das Embed mit einem Verweis auf das eigentliche Widget
// und der Moeglichkeit, den Status mit anderen Widgets zu teilen.
func (e *SelectEmbed) Init(node Node, extData binding.Untyped, value any) {
    e.node = node
    if extData == nil {
        e.BindVar = binding.NewUntyped()
        e.BindVar.Set(nil)
    } else {
        e.BindVar = extData
    }
    e.BindValue = value
    e.BindVar.AddListener(e)
}

// Ermittelt den Status des Embed.
func (e *SelectEmbed) Selected() (bool) {
    if e.node == nil {
        return false
    }
    return e.BindVar.Get() == e.BindValue
}

// Muss vom umschliessenden Widget aufgerufen werden.
func (e *SelectEmbed) OnInputEvent(evt touch.Event) {
    Debugf(Events, "evt: %v", evt)
    if e.node == nil {
        return
    }
    switch evt.Type {
    case touch.TypeTap:
//        if e.Selected() {
//            e.BindVar.Set(nil)
//        } else {
            e.BindVar.Set(e.BindValue)
//        }
    }
}

// Wird autom. aufgerufen, sobald der Wert von 'BindVar' veraendert wird.
func (e *SelectEmbed) DataChanged(BindVar binding.DataItem) {
    if e.node == nil {
        return
    }
    e.node.Mark(MarkNeedsPaint)
}

