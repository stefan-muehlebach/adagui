//go:build debug

package adagui

import (
    "flag"
    "fmt"
    "log"
    "path"
    "runtime"
    "strings"
)

var (
    callerBuffer []uintptr
)

func init() {
    flag.Uint64Var(&debugDomains, "debugDomain",
            uint64(Painting | Coordinates | Events),
            "whith this flag, you can filter the debug messages.")

    Debugf = debugf
    callerBuffer = make([]uintptr, 5)
}

func debugf(domain DebugDomain, format string, a ...any) {
    if domain & DebugDomain(debugDomains) == 0 {
        return
    }
    log.Printf("%*s %s: %s", 2*callDepth(2), ">", callerInfo(2),
        fmt.Sprintf(format, a...))
}

func callerInfo(skip int) (string) {
    n := runtime.Callers(skip+1, callerBuffer)
    buf := callerBuffer[:n]
    frames := runtime.CallersFrames(buf)
    frame, _ := frames.Next()
    return strings.ReplaceAll(path.Base(frame.Function), "adagui.", "")
    
//    pc, _, _, _ := runtime.Caller(skip)
//    fnc := runtime.FuncForPC(pc)
//    return strings.ReplaceAll(path.Base(fnc.Name()), "adagui.", "")
}

func callDepth(skip int) (int) {
    for n := 1 << 8; ; n *= 2 {
        buf := make([]uintptr, n)
        n := runtime.Callers(skip+3, buf)
        if n < len(buf) {
            return n
        }
    }
}

