package adagui

import (
    "testing"
    "github.com/stefan-muehlebach/gg/color"
    "github.com/stefan-muehlebach/gg/colornames"
    "golang.org/x/image/font/opentype"
)

var (
    defProps *Properties
    typProps *Properties
    objProps *Properties
    c color.Color
    f *opentype.Font
)

func init() {
    defProps = NewDefaultProps()
    typProps = NewButtonProps(defProps)
    objProps = NewProperties(typProps)
}

func TestGetColor(t *testing.T) {
    c = defProps.Color(Color)
    t.Logf("Def.Color: %v", c)
    c = typProps.Color(Color)
    t.Logf("Typ.Color: %v", c)
    c = objProps.Color(Color)
    t.Logf("Obj.Color: %v", c)

    typProps.SetColor(Color, colornames.FireBrick)
    objProps.SetColor(Color, colornames.Yellow)

    c = defProps.Color(Color)
    t.Logf("Def.Color: %v", c)
    c = typProps.Color(Color)
    t.Logf("Typ.Color: %v", c)
    c = objProps.Color(Color)
    t.Logf("Obj.Color: %v", c)
}

func TestGetFont(t *testing.T) {
    f = defProps.Font(BoldFont)
    t.Logf("Def.BoldFont: %T", f)
}

func BenchmarkGetColorDef(b *testing.B) {
    for i:=0; i<b.N; i++ {
        c = defProps.Color(BorderColor)
    }
}

func BenchmarkGetColorTyp(b *testing.B) {
    for i:=0; i<b.N; i++ {
        c = typProps.Color(BorderColor)
    }
}

func BenchmarkGetColorObj(b *testing.B) {
    for i:=0; i<b.N; i++ {
        c = objProps.Color(BorderColor)
    }
}

