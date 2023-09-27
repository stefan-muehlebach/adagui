package adagui

import (
    "fmt"
    //"testing"
    "mju.net/utils"
)

var (
    defProp *Properties
)

func init() {
    defProp = NewProperties(nil)
    defProp.Set(FillColor, utils.Blue)
    defProp.Set(StrokeColor, utils.Blue)
    defProp.Set(StrokeWidth, 0.0)
    defProp.Set(TextFont, "Seaford")
    defProp.Set(TextColor, utils.WhiteLight)
    defProp.Set(ButtonSize, 25.0)
    if !defProp.Check() {
        fmt.Printf("Properties not complete!\n")
    }
}

func ExampleGetDefault() {
    p := NewProperties(defProp)
    fmt.Printf("%.3f\n", p.GetFloat64(StrokeWidth))
    fmt.Printf("%.3f\n", p.GetFloat64(ButtonSize))
    fmt.Printf("%v\n", p.GetColor(FillColor))
    // Output:
    // 0.000
    // 25.000
    // {41 128 185 255}
}

func ExampleSetPrivate() {
    p := NewProperties(defProp)
    p.Set(ButtonSize, 32.0)
    fmt.Printf("%.3f\n", p.GetFloat64(ButtonSize))
    p.Reset()
    fmt.Printf("%.3f\n", p.GetFloat64(ButtonSize))
    // Output:
    // 32.000
    // 25.000
}

