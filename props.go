package adagui

import (
    "github.com/stefan-muehlebach/gg/color"
    "github.com/stefan-muehlebach/gg/colornames"
    "github.com/stefan-muehlebach/gg/fonts"
    "golang.org/x/image/font/opentype"
)

var (
    DefProps = NewDefaultProps()
)

type ColorPropertyName int

const (
    Color ColorPropertyName = iota
    PressedColor
    SelectedColor
    BorderColor
    PressedBorderColor
    SelectedBorderColor
    TextColor
    PressedTextColor
    SelectedTextColor
    LineColor
    PressedLineColor
    SelectedLineColor
    BarColor
    PressedBarColor
    // Old
    RedColor
    OrangeColor
    YellowColor
    GreenColor
    BlueColor
    PurpleColor
    BrownColor
    GrayColor
    BlackColor
    WhiteColor
    numColorProperties
)

type FontPropertyName int

const (
    Font FontPropertyName = iota
    RegularFont
    BoldFont
    ItalicFont
    BoldItalicFont
    MonoFont
    MonoBoldFont
    numFontProperties
)

type SizePropertyName int

const (
    Width SizePropertyName = iota
    Height
    Size
    BorderWidth
    PressedBorderWidth
    SelectedBorderWidth
    LineWidth
    InnerPadding
    Padding
    CornerRadius
    FontSize
    BarSize
    CtrlSize
    numSizeProperties
)

// ----------------------------------------------------------------------------

// Properties dienen dazu, graphische Eigenschaften von Widgets hierarchisch
// zu verwalten. In einem Properties-Objekt können drei Arten von Eigenschaften
// verwaltet werden: Farben (Datentyp: color.Color), Schriftarten (Datentyp:
// *opentype.Font) und Zahlen (Datentyp: float64). Durch die Hierarchie ist
// es möglich für einzelne Widgets vom Standard abweichende Eigenschaften
// zu definieren.
type Properties struct {
    parent   *Properties
    colorMap map[ColorPropertyName]color.Color
    fontMap  map[FontPropertyName]*opentype.Font
    sizeMap  map[SizePropertyName]float64
}

// Erzeugt ein neues Property-Objekt und hinterlegt parent als Vater-Property.
func NewProperties(parent *Properties) (*Properties) {
    p := &Properties{}

    p.parent   = parent
    p.colorMap = make(map[ColorPropertyName]color.Color)
    p.fontMap  = make(map[FontPropertyName]*opentype.Font)
    p.sizeMap  = make(map[SizePropertyName]float64)

    return p
}

// Interne Funktion. Damit werden die Properties für Widget-Kategorien
// (Buttons, Labels, etc) erzeugt.
func newProps(parent *Properties, colorMap map[ColorPropertyName]color.Color,
        fontMap map[FontPropertyName]*opentype.Font,
        sizeMap map[SizePropertyName]float64) (*Properties) {
    p := &Properties{}

    p.parent   = parent

    if colorMap == nil {
        colorMap = make(map[ColorPropertyName]color.Color)
    }
    p.colorMap = colorMap
    if fontMap == nil {
        fontMap = make(map[FontPropertyName]*opentype.Font)
    }
    p.fontMap  = fontMap
    if sizeMap == nil {
        sizeMap = make(map[SizePropertyName]float64)
    }
    p.sizeMap  = sizeMap

    return p
}

var (
    NewProps func(*Properties, map[ColorPropertyName]color.Color,
        map[FontPropertyName]*opentype.Font,
        map[SizePropertyName]float64) (*Properties) = newProps
)

// Das sind die Hauptmethoden, um Farben, Font oder Groessen aus den
// Properties zu lesen. Kann ein Property nicht gefunden werden, dann
// wird (falls vorhanden) das Parent-Property angefragt.
func (p *Properties) Color(name ColorPropertyName) (color.Color) {
    if col, ok := p.colorMap[name]; !ok && p.parent != nil {
        return p.parent.Color(name)
    } else {
        return col
    }
}

func (p *Properties) Font(name FontPropertyName) (*opentype.Font) {
    if fnt, ok := p.fontMap[name]; !ok && p.parent != nil {
        return p.parent.Font(name)
    } else {
        return fnt
    }
}

func (p *Properties) Size(name SizePropertyName) (float64) {
    if siz, ok := p.sizeMap[name]; !ok && p.parent != nil {
        return p.parent.Size(name)
    } else {
        return siz
    }
}

// Über diese Methoden können einzelne Eigenschaften auf Typen- oder Objekt-
// ebene definiert werden.
func (p *Properties) SetColor(name ColorPropertyName, col color.Color) {
    p.colorMap[name] = col
}

func (p *Properties) SetFont(name FontPropertyName, fnt *opentype.Font) {
    p.fontMap[name] = fnt
}

func (p *Properties) SetSize(name SizePropertyName, size float64) {
    p.sizeMap[name] = size
}

// Auf Typen- oder Objekt-Stufe definierte Eigenschaften können mit den
// Del-Methoden wieder entfernt werden, so dass der Eintrag des Parents wieder
// aktiviert wird. Existiert die Eigenschaft in den Properties nicht, sind
// die Methoden no-op. Auf Properties der obersten Hierarchiestufe (d.h. mit
// parent == nil) haben die Methoden keinen Einfluss.
func (p *Properties) DelColor(name ColorPropertyName) {
    if p.parent == nil {
        return
    }
    delete(p.colorMap, name)
}

func (p *Properties) DelFont(name FontPropertyName) {
    if p.parent == nil {
        return
    }
    delete(p.fontMap, name)
}

func (p *Properties) DelSize(name SizePropertyName) {
    if p.parent == nil {
        return
    }
    delete(p.sizeMap, name)
}

// ----------------------------------------------------------------------------

// Erstellt ein neues Default-Property Objekt. Die Default Properties muessen
// fur jedes Property einen Wert bereitstellen. Mit den Tests in props_test.go
// kann geprüft werden, ob dies erfüllt ist.
func NewDefaultProps() (*Properties) {
    p := &Properties{}

    c1 := colornames.DarkGreen
    c2 := c1.Interpolate(colornames.YellowGreen, 0.9)
    c3 := c1.Interpolate(colornames.YellowGreen, 0.7)

    p.colorMap = map[ColorPropertyName]color.Color{
        Color:               c1,
        PressedColor:        c2,
        SelectedColor:       c3,

        BorderColor:         c1,
        PressedBorderColor:  c2,
        SelectedBorderColor: c3,

        TextColor:           colornames.WhiteSmoke,
        PressedTextColor:    colornames.Black,
        SelectedTextColor:   colornames.White,

        LineColor:           colornames.WhiteSmoke,
        PressedLineColor:    colornames.Black,
        SelectedLineColor:   colornames.WhiteSmoke,

        BarColor:            colornames.DarkSlateGray.Dark(0.5),
        PressedBarColor:     colornames.DarkSlateGray.Dark(0.5),

        // Out
        RedColor:            colornames.Red,
        OrangeColor:         colornames.Orange,
        YellowColor:         colornames.Yellow,
        GreenColor:          colornames.Green,
        BlueColor:           colornames.Blue,
        PurpleColor:         colornames.Purple,
        BrownColor:          colornames.Brown,
        GrayColor:           colornames.Gray,
        BlackColor:          colornames.Black,
        WhiteColor:          colornames.WhiteSmoke,
    }

    p.fontMap = map[FontPropertyName]*opentype.Font{
        Font:           fonts.GoRegular,
        RegularFont:    fonts.GoRegular,
        BoldFont:       fonts.GoBold,
        ItalicFont:     fonts.GoItalic,
        BoldItalicFont: fonts.GoBoldItalic,
        MonoFont:       fonts.GoMono,
        MonoBoldFont:   fonts.GoMonoBold,
    }

    p.sizeMap = map[SizePropertyName]float64{
        Width:               0.0,
        Height:              0.0,
        Size:                0.0,

        BorderWidth:         0.0,
        PressedBorderWidth:  0.0,
        SelectedBorderWidth: 0.0,

        LineWidth:           2.0,
        InnerPadding:        5.0,
        Padding:            15.0,
        CornerRadius:        6.0,
        FontSize:           15.0,
        BarSize:            18.0,
        CtrlSize:           18.0,
    }
    return p
}

