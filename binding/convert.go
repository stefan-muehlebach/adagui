package binding

import (
    "fmt"
)

func IntToStringWithFormat(v Int, format string) (String) {
    if format == "%d" {
        return IntToString(v)
    }
    s := NewString()
    v.AddCallback(func (data DataItem) {
        val := data.(Int).Get()
        s.Set(fmt.Sprintf(format, val))
    })
    return s
}

func IntToString(v Int) (String) {
    s := NewString()
    v.AddCallback(func (data DataItem) {
        val := data.(Int).Get()
        s.Set(fmt.Sprintf("%d", val))
    })
    return s
}

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

