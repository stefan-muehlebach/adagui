package adagui

import (
    "github.com/stefan-muehlebach/gg/color"
    "github.com/stefan-muehlebach/gg/colornames"
    "github.com/stefan-muehlebach/gg/fonts"
    "golang.org/x/image/font/opentype"
)

var (
    DefProp = NewDefaultProps()
    Pr = DefProp
    pr = DefProp

    panelProps  = NewPanelProps(DefProp)
    buttonProps = NewButtonProps(DefProp)
    checkProps  = NewCheckProps(buttonProps)
    radioProps  = NewRadioProps(buttonProps)
    tabProps    = NewTabProps(DefProp)
    gaugeProps  = NewGaugeProps(DefProp)
)

func init() {
}

type ColorPropertyName int

const (
    RedColor ColorPropertyName = iota
    OrangeColor
    YellowColor
    GreenColor
    BlueColor
    PurpleColor
    BrownColor
    GrayColor
    BlackColor
    WhiteColor
    ActiveColor
    TranspWhite

    Color
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

    SeparatorColor

    StrokeColor

    TextFocusColor
    TextDimColor

//    ScrollBarColor
//    ScrollBarFocusColor
//    ScrollCtrlColor
//    ScrollCtrlFocusColor

//    SliderBarColor
//    SliderBarFocusColor
//    SliderCtrlColor
//    SliderCtrlFocusColor

    ArrowColor

    numColorProperties
)

type FontPropertyName int

const (
    RegularFont FontPropertyName = iota
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

    LineWidth

    BorderWidth
    PressedBorderWidth
    SelectedBorderWidth

    InnerPadding
    Padding

    CornerRadius

    FontSize





    PanelBorderSize

    InnerPaddingSize
    PaddingSize

    ButtonSize
    ButtonBorderSize
    ButtonCornerRad

    TextButtonPaddingSize

    RadioSize
    RadioBorderSize
    RadioDotSize

    CheckSize
    CheckLineSize
    CheckCornerRad

    IconInlineSize

    TabButtonWidth
    TabButtonHeight
    TabButtonBorderSize
    TabButtonCornerRad
    TabButtonTextSize

    TextSize
    TextHeadingSize
    TextSubHeadingSize

    ScrollSize
    ScrollBarSize
    ScrollCtrlSize

    SliderSize
    SliderBarSize
    SliderCtrlSize

    numSizeProperties
)

// ----------------------------------------------------------------------------

type Properties struct {
    parent   *Properties
    colorMap map[ColorPropertyName]color.Color
    fontMap  map[FontPropertyName]*opentype.Font
    sizeMap  map[SizePropertyName]float64
}

// Erzeugt ein neues Property-Objekt und hinterlegt parent als Vater-Property
//
func NewProperties(parent *Properties) (*Properties) {
    p := &Properties{}

    p.parent   = parent
    p.colorMap = make(map[ColorPropertyName]color.Color)
    p.fontMap  = make(map[FontPropertyName]*opentype.Font)
    p.sizeMap  = make(map[SizePropertyName]float64)

    return p
}

// Das sind die Hauptmethoden, um Farben, Font oder Groessen aus den
// Properties zu lesen. Kann ein Property nicht gefunden werden, dann
// wird (falls vorhanden) das Parent-Property angefragt.
//
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

// Ueber diese Methoden koennen einzelne Properties auf Objekt-Ebene
// veraendert werden.
//
func (p *Properties) SetColor(name ColorPropertyName, col color.Color) {
    p.colorMap[name] = col
}

func (p *Properties) SetFont(name FontPropertyName, fnt *opentype.Font) {
    p.fontMap[name] = fnt
}

func (p *Properties) SetSize(name SizePropertyName, size float64) {
    p.sizeMap[name] = size
}

// ----------------------------------------------------------------------------

func NewDefaultProps() (*Properties) {
    p := &Properties{}

    p.parent = nil
    p.colorMap = map[ColorPropertyName]color.Color{
        Color:               colornames.DarkGreen,
        PressedColor:        colornames.GreenYellow,
        SelectedColor:       colornames.LimeGreen,

        BorderColor:         colornames.DarkGreen,
        PressedBorderColor:  colornames.GreenYellow,
        SelectedBorderColor: colornames.LimeGreen,

        TextColor:           colornames.WhiteSmoke,
        PressedTextColor:    colornames.Black,
        SelectedTextColor:   colornames.White,

        LineColor:           colornames.WhiteSmoke,
        PressedLineColor:    colornames.Black,

        // Out
        RedColor:    colornames.Red,
        OrangeColor: colornames.Orange,
        YellowColor: colornames.Yellow,
        GreenColor:  colornames.Green,
        BlueColor:   colornames.Blue,
        PurpleColor: colornames.Purple,
        BrownColor:  colornames.Brown,
        GrayColor:   colornames.Gray,
        BlackColor:  colornames.Black,
        WhiteColor:  colornames.WhiteSmoke,

        TranspWhite: colornames.WhiteSmoke.Alpha(0.5),


        SeparatorColor:      colornames.Gray,

        ArrowColor:          colornames.WhiteSmoke,

        TextFocusColor:      colornames.White,
        TextDimColor:        colornames.Silver,
    }

    p.fontMap = map[FontPropertyName]*opentype.Font{
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

        LineWidth:           2.0,

        BorderWidth:         0.0,
        PressedBorderWidth:  0.0,
        SelectedBorderWidth: 0.0,

        InnerPadding:        5.0,
        Padding:            15.0,

        CornerRadius:        6.0,

        FontSize:           15.0,








        PanelBorderSize: 0.0,

        InnerPaddingSize: 5.0,
        PaddingSize:      5.0,

        ButtonSize:       32.0,
        ButtonBorderSize: 0.0,
        ButtonCornerRad:  6.0,

        TextButtonPaddingSize: 15.0,

        IconInlineSize: 24.0,

        TextSize:           15.0,
        TextHeadingSize:    15.0,
        TextSubHeadingSize: 15.0,

        ScrollSize:     18.0,
        ScrollBarSize:  18.0,
        ScrollCtrlSize: 18.0,

        SliderSize:     18.0,
        SliderBarSize:  18.0,
        SliderCtrlSize: 18.0,



    }
    return p
}

func NewPanelProps(parent *Properties) (*Properties) {
    p := NewProperties(parent)
    p.colorMap = map[ColorPropertyName]color.Color{
        Color:         colornames.Black,
        BorderColor:         colornames.Black,
    }
    return p
}

func NewButtonProps(parent *Properties) (*Properties) {
    p := NewProperties(parent)
    p.sizeMap = map[SizePropertyName]float64{
        Size:         32.0,
    }
    return p
}

func NewCheckProps(parent *Properties) (*Properties) {
    p := NewProperties(parent)
    p.sizeMap = map[SizePropertyName]float64{
        Size:         18.0,
        LineWidth:     4.0,
        CornerRadius:  5.0,
    }
    return p
}

func NewRadioProps(parent *Properties) (*Properties) {
    p := NewProperties(parent)
    p.sizeMap = map[SizePropertyName]float64{
        Size:         18.0,
        LineWidth:     8.0,
    }
    return p
}

func NewTabProps(parent *Properties) (*Properties) {
    p := NewProperties(parent)
    p.sizeMap = map[SizePropertyName]float64{
        Width:        32.0,
        Height:       20.0,
        CornerRadius:  8.0,
        FontSize:     12.0,
    }
    return p
}

func NewGaugeProps(parent *Properties) (*Properties) {
    p := NewProperties(parent)
    p.colorMap = map[ColorPropertyName]color.Color{
        BorderColor:         colornames.DarkSlateGray.Dark(0.5),
        PressedBorderColor:  colornames.DarkSlateGray.Dark(0.5),
    }
    return p
}

