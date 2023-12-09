package adagui

import (
    "github.com/stefan-muehlebach/gg/color"
    "github.com/stefan-muehlebach/gg/colornames"
    "github.com/stefan-muehlebach/gg/fonts"
    "golang.org/x/image/font/opentype"
)

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

    MainColor
    FillColor
    StrokeColor

    TextColor
    TextFocusColor
    TextDimColor

    ButtonColor
    ButtonFocusColor
    ButtonBorderColor
    ButtonBorderFocusColor

    IconButtonColor
    IconButtonFocusColor
    IconButtonSelColor
    IconButtonBorderColor
    IconButtonBorderFocusColor
    IconButtonBorderSelColor

    TabButtonColor
    TabButtonFocusColor
    TabButtonSelColor
    TabButtonBorderColor
    TabButtonBorderFocusColor
    TabButtonBorderSelColor

    ScrollBarColor
    ScrollBarFocusColor
    ScrollCtrlColor
    ScrollCtrlFocusColor

    SliderBarColor
    SliderBarFocusColor
    SliderCtrlColor
    SliderCtrlFocusColor

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
    LineWidth SizePropertyName = iota

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

type Properties interface {
    Color(n ColorPropertyName) color.Color
    Font(n FontPropertyName) *opentype.Font
    Size(n SizePropertyName) float64
}

// ----------------------------------------------------------------------------

type DefaultProps struct {
    colorMap map[ColorPropertyName]color.Color
    fontMap  map[FontPropertyName]*opentype.Font
    sizeMap  map[SizePropertyName]float64
}

func NewDefaultProps() *DefaultProps {
    p := &DefaultProps{}

    mainColor := colornames.Darkcyan

    p.colorMap = map[ColorPropertyName]color.Color{
        RedColor:    colornames.Red,
        OrangeColor: colornames.Orange,
        YellowColor: colornames.Yellow,
        GreenColor:  colornames.Green,
        BlueColor:   colornames.Blue,
        PurpleColor: colornames.Purple,
        BrownColor:  colornames.Brown,
        GrayColor:   colornames.Gray,
        BlackColor:  colornames.Black,
        WhiteColor:  colornames.Whitesmoke,
        ActiveColor: mainColor.Bright(0.2),
        TranspWhite: colornames.Whitesmoke.Alpha(0.5),

        MainColor:   mainColor,
        FillColor:   mainColor,
        StrokeColor: mainColor,

        ArrowColor: colornames.Whitesmoke,

        TextColor:      colornames.Whitesmoke,
        TextFocusColor: colornames.White,
        TextDimColor:   colornames.Silver,

        ButtonColor:            mainColor,
        ButtonFocusColor:       mainColor.Bright(0.2),
        ButtonBorderColor:      mainColor,
        ButtonBorderFocusColor: mainColor.Bright(0.2),

        IconButtonColor:            mainColor,
        IconButtonFocusColor:       mainColor.Bright(0.2),
        IconButtonSelColor:         mainColor.Bright(0.2),
        IconButtonBorderColor:      mainColor,
        IconButtonBorderFocusColor: mainColor.Bright(0.2),
        IconButtonBorderSelColor:   mainColor.Bright(0.2),

        TabButtonColor:            mainColor.Alpha(0.3),
        TabButtonFocusColor:       mainColor.Bright(0.2),
        TabButtonSelColor:         mainColor,
        TabButtonBorderColor:      colornames.Whitesmoke,
        TabButtonBorderFocusColor: colornames.White,
        TabButtonBorderSelColor:   colornames.White,

        ScrollBarColor:       colornames.Silver.Alpha(0.3),
        ScrollBarFocusColor:  colornames.Silver,
        ScrollCtrlColor:      mainColor,
        ScrollCtrlFocusColor: mainColor.Bright(0.2),

        SliderBarColor:       colornames.Silver.Alpha(0.3),
        SliderBarFocusColor:  colornames.Silver,
        SliderCtrlColor:      mainColor,
        SliderCtrlFocusColor: mainColor.Bright(0.2),
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
        LineWidth: 2.0,

        PanelBorderSize: 0.0,

        InnerPaddingSize: 5.0,
        PaddingSize:      5.0,

        ButtonSize:       32.0,
        ButtonBorderSize: 0.0,
        ButtonCornerRad:  6.0,

        TextButtonPaddingSize: 15.0,

        RadioSize:       18.0,
        RadioBorderSize: 4.0,
        RadioDotSize:    8.0,

        CheckSize:      18.0,
        CheckLineSize:  4.0,
        CheckCornerRad: 5.0,

        IconInlineSize: 24.0,

        TabButtonWidth:      32.0,
        TabButtonHeight:     20.0,
        TabButtonBorderSize:  0.0,
        TabButtonCornerRad:   8.0,
        TabButtonTextSize:   12.0,

        TextSize:           15.0,
        TextHeadingSize:    15.0,
        TextSubHeadingSize: 15.0,

        ScrollSize:     18.0,
        ScrollBarSize:  14.0,
        ScrollCtrlSize: 18.0,

        SliderSize:     18.0,
        SliderBarSize:  14.0,
        SliderCtrlSize: 18.0,
    }

    return p
}

func (p *DefaultProps) Color(n ColorPropertyName) color.Color {
    return p.colorMap[n]
}

func (p *DefaultProps) Font(n FontPropertyName) *opentype.Font {
    return p.fontMap[n]
}

func (p *DefaultProps) Size(n SizePropertyName) float64 {
    return p.sizeMap[n]
}
