package binding

import (
    "fmt"
)

func FloatToStringWithFormat(v Float, format string) (String) {
    if format == "%f" {
        return FloatToString(v)
    }
    s := NewString()
    v.AddCallback(func (data DataItem) {
        val := data.(Float).Get()
        s.Set(fmt.Sprintf(format, val))
    })
    return s
}

func FloatToString(v Float) (String) {
    s := NewString()
    v.AddCallback(func (data DataItem) {
        val := data.(Float).Get()
        s.Set(fmt.Sprintf("%f", val))
    })
    return s
}

