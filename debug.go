package adagui

import (
    "fmt"
    "log"
)

// Das ist ein Ueberbleibsel einer Debug-Session, als ich ohne klare
// Darstellung der Stack-Level einfach nicht begriffen haben, wo das Problem
// liegt. Vielleicht wird daraus mal noch etwas...?
type DebugLevel int

func (l *DebugLevel) Inc() {
    *l += 3
    log.SetPrefix(":" + l.String())
}

func (l *DebugLevel) Dec() {
    *l -= 3
    log.SetPrefix(":" + l.String())
}

func (l DebugLevel) String() (string) {
    return fmt.Sprintf("%*s", l, " ")
}

var (
    stackLevel DebugLevel = 0
)

