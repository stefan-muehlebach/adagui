package adagui

import (
    "flag"
    "fmt"
    "log"
    "path"
    "runtime"
)

//----------------------------------------------------------------------------

type DebugDomain uint64

const (
    Painting DebugDomain = 1 << iota
    Coordinates
    Events
)

var (
    debugFlag bool
    debugDomains uint64
)

func init() {
    flag.BoolVar(&debugFlag, "debug", false, "show debugging messages")
    flag.Uint64Var(&debugDomains, "domain", uint64(Painting | Coordinates),
            "whith this flag, you can filter the debug messages.")
}

func Debugf(domain DebugDomain, format string, a ...any) {
    if !debugFlag {
        return
    }
    if domain & DebugDomain(debugDomains) == 0 {
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

