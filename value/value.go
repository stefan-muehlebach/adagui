package value

type Type uint32

const (
    TypeChanging Type = iota
    TypeChanged
    NumEvents
)

func (t Type) String() (string) {
    switch t {
    case TypeChanging:
        return "Changing"
    case TypeChanged:
        return "Changed"
    default:
        return "(Unknown Type)"
    }
}

type Event struct {
    Type Type

    OldValue, NewValue float64
}

