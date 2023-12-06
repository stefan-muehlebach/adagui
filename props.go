package adagui

import (
    // "image/color"
    "golang.org/x/image/font/opentype"
    "github.com/stefan-muehlebach/gg/color"
    "github.com/stefan-muehlebach/gg/colornames"
    "github.com/stefan-muehlebach/gg/fonts"
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
    FillColor
    StrokeColor
    ActiveColor
    TranspWhite

    TextColor
    TextFocusColor

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
    //PanelLineWidth
    //ButtonHeight
    //ButtonInsetHorizontal
    //ButtonInset
    //ButtonLineWidth
    //ButtonCornerRad
    //RadioBtnSize
    //RadioBtnLineWidth
    //RadioBtnDotSize
    //CheckboxSize
    //CheckboxLineWidth
    //CheckboxCornerRad
    //SliderWidth
    //SliderBarWidth
    //SliderCtrlWidth
    //ScrollWidth
    //ScrollBarWidth
    //ScrollCtrlWidth

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
     fontMap map[FontPropertyName]*opentype.Font
     sizeMap map[SizePropertyName]float64
}

func NewDefaultProps() (*DefaultProps) {
    p := &DefaultProps{}

    p.colorMap = map[ColorPropertyName]color.Color{
        RedColor:      colornames.Red,
        OrangeColor:   colornames.Orange,
        YellowColor:   colornames.Yellow,
        GreenColor:    colornames.Green,
        BlueColor:     colornames.Blue,
        PurpleColor:   colornames.Purple,
        BrownColor:    colornames.Brown,
        GrayColor:     colornames.Gray,
        BlackColor:    colornames.Black,
        WhiteColor:    colornames.Whitesmoke,
        FillColor:     colornames.Darkcyan,
        StrokeColor:   colornames.Darkcyan,
        ActiveColor:   colornames.Darkcyan.Bright(2),
        TranspWhite:   colornames.Whitesmoke.Alpha(0.5),

        ArrowColor:             colornames.Whitesmoke,

        TextColor:              colornames.Whitesmoke,
        TextFocusColor:         colornames.White,
    
        ButtonColor:            colornames.Darkcyan,
        ButtonFocusColor:       colornames.Darkcyan.Bright(2),
        ButtonBorderColor:      colornames.Darkcyan,
        ButtonBorderFocusColor: colornames.Darkcyan.Bright(2),
    
        IconButtonColor:        colornames.Darkcyan,
        IconButtonFocusColor:   colornames.Darkcyan.Bright(2),
        IconButtonSelColor:     colornames.Darkcyan.Bright(2),
        IconButtonBorderColor:  colornames.Darkcyan,
        IconButtonBorderFocusColor: colornames.Darkcyan.Bright(2),
        IconButtonBorderSelColor: colornames.Darkcyan.Bright(2),

        //IconButtonColor:        colornames.White,
        //IconButtonFocusColor:   colornames.Darkgray,
        //IconButtonSelColor:     color.RGBAF64{0.353, 0.627, 1.0, 1.0},
        //IconButtonBorderColor:  colornames.Silver,
        //IconButtonBorderFocusColor: colornames.Darkgray,
        //IconButtonBorderSelColor: colornames.Silver,

        ScrollBarColor:         colornames.Silver.Alpha(0.3),
        ScrollBarFocusColor:    colornames.Silver,
        ScrollCtrlColor:        colornames.Darkcyan,
        ScrollCtrlFocusColor:   colornames.Darkcyan.Bright(2),
    
        SliderBarColor:         colornames.Silver.Alpha(0.3),
        SliderBarFocusColor:    colornames.Silver,
        SliderCtrlColor:        colornames.Darkcyan,
        SliderCtrlFocusColor:   colornames.Darkcyan.Bright(2),
    }

    p.fontMap = map[FontPropertyName]*opentype.Font{
        RegularFont:            font.GoRegular,
        BoldFont:               font.GoBold,
        ItalicFont:             font.GoItalic,
        BoldItalicFont:         font.GoBoldItalic,
        MonoFont:               font.GoMono,
        MonoBoldFont:           font.GoMonoBold,
    }

    p.sizeMap = map[SizePropertyName]float64{
        LineWidth:              2.0,
//        PanelLineWidth:         0.0,
//        ButtonHeight:          32.0,
//        ButtonInsetHorizontal: 15.0,
//        ButtonInset:   	        5.0,
//        ButtonCornerRad:        6.0,
//        RadioBtnSize:          18.0,
//        RadioBtnLineWidth:      4.0,
//        RadioBtnDotSize:        8.0,
//        CheckboxSize:          18.0,
//        CheckboxLineWidth:      4.0,
//        CheckboxCornerRad:      5.0,
//        SliderWidth:           18.0,
//        SliderBarWidth:        14.0,
//        SliderCtrlWidth:       18.0,
//        ScrollWidth:           18.0,
//        ScrollBarWidth:        14.0,
//        ScrollCtrlWidth:       18.0,


        PanelBorderSize:        0.0,

        InnerPaddingSize:       5.0,
        PaddingSize:            5.0,
    
        ButtonSize:            32.0,
        ButtonBorderSize:       0.0,
        ButtonCornerRad:        6.0,

        TextButtonPaddingSize: 15.0,
    
        RadioSize:             18.0,
        RadioBorderSize:        4.0,
        RadioDotSize:           8.0,

        CheckSize:             18.0,
        CheckLineSize:          4.0,
        CheckCornerRad:         5.0,

        IconInlineSize:        24.0,
    
        TextSize:              15.0,
        TextHeadingSize:       15.0,
        TextSubHeadingSize:    15.0,

        ScrollSize:            18.0,
        ScrollBarSize:         14.0,
        ScrollCtrlSize:        18.0,
    
        SliderSize:            18.0,
        SliderBarSize:         14.0,
        SliderCtrlSize:        18.0,
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

