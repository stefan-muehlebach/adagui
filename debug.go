package adagui

import (
    "fmt"
    "log"
)

// Das ist ein Ueberbleibsel einer Debug-Session, als ich ohne klare
// Darstellung der Stack-Level einfach nicht begriffen haben, wo das Problem
// liegt. Vielleicht wird daraus mal noch etwas...?
type LevelType int

func (l *LevelType) Inc() {
    *l += 3
    log.SetPrefix(":" + l.String())
}

func (l *LevelType) Dec() {
    *l -= 3
    log.SetPrefix(":" + l.String())
}

func (l LevelType) String() (string) {
    return fmt.Sprintf("%*s", l, " ")
}

var (
    stackLevel LevelType = 0
)

