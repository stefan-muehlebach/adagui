package adagui

type DebugDomain uint64

const (
    Painting DebugDomain = 1 << iota
    Coordinates
    Events
    Layout
)

var (
    debugDomains uint64
    lastDomain DebugDomain = Events
    Debugf = func(domain DebugDomain, format string, a ...any){}
)

func NewDebugDomain() (DebugDomain) {
    lastDomain <<= 1
    return lastDomain
}

func SetDebugDomain(domain DebugDomain) {
    debugDomains = uint64(domain)
}

func AddDebugDomain(domain DebugDomain) {
    debugDomains |= uint64(domain)
}

