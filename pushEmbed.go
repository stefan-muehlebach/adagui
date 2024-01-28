package adagui

import (
    "github.com/stefan-muehlebach/adagui/binding"
    "github.com/stefan-muehlebach/adagui/touch"
)

// Mit diesem Embed erhaelt ein Widget die Moeglichkeit, "gedrueckt"
// zu werden.
type PushEmbed struct {
    pushed binding.Bool
    node Node
}

// Initialisert wird das Embed mit einem Verweis auf das eigentliche Widget
// und der Moeglichkeit, den Status mit anderen Widgets zu teilen.
func (e *PushEmbed) Init(node Node, extData binding.Bool) {
    e.node = node
    if extData == nil {
        e.pushed = binding.NewBool()
        e.pushed.Set(false)
    } else {
        e.pushed = extData
    }
    e.pushed.AddListener(e)
}

// Ermittelt den Status des Embed.
func (e *PushEmbed) Pushed() (bool) {
    return e.pushed.Get()
}

// Muss vom umschliessenden Widget aufgerufen werden.
func (e *PushEmbed) OnInputEvent(evt touch.Event) {
    Debugf(Events, "evt: %v", evt)
    switch evt.Type {
    case touch.TypePress, touch.TypeEnter:
        e.pushed.Set(true)
    case touch.TypeRelease, touch.TypeLeave:
        e.pushed.Set(false)
    }
}

// Wird autom. aufgerufen, sobald der Wert von 'pushed' veraendert wird.
func (e *PushEmbed) DataChanged(pushed binding.DataItem) {
    e.node.Mark(MarkNeedsPaint)
}
