package adagui

import (
    "flag"
    "fmt"
    "log"
    "path"
    "runtime"
)

//----------------------------------------------------------------------------

var (
    debugFlag bool
)

func init() {
    flag.BoolVar(&debugFlag, "debug", false, "show debugging messages")
    flag.Parse()
}

func Debugf(format string, a ...any) {
    if !debugFlag {
        return
    }
    log.Printf("%*s %s: %s", CallDepth(2), ">", CallerInfo(2),
        fmt.Sprintf(format, a...))
}

func CallerInfo(skip int) (string) {
    pc, _, _, _ := runtime.Caller(skip)
    fnc := runtime.FuncForPC(pc)
    return path.Base(fnc.Name())
}

func CallDepth(skip int) (int) {
    for n := 1 << 8; ; n *= 2 {
        buf := make([]uintptr, n)
        n := runtime.Callers(skip+3, buf)
        if n < len(buf) {
            return n
        }
    }
}

//----------------------------------------------------------------------------

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


