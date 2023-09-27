package adagui

import (
    "golang.org/x/image/font"
    "mju.net/utils"
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
    TextColor
    numColorProperties
)

type FontPropertyName int

const (
    RegularFont FontPropertyName = iota
    BoldFont
    numFontProperties
)

type SizePropertyName int

const (
    LineWidth SizePropertyName = iota
    PanelLineWidth
    ButtonHeight
    ButtonInsetHorizontal
    ButtonInset
    ButtonLineWidth
    ButtonCornerRad
    RadioBtnSize
    RadioBtnLineWidth
    RadioBtnDotSize
    CheckboxSize
    CheckboxLineWidth
    CheckboxCornerRad
    SliderWidth
    SliderBarWidth
    SliderCtrlWidth
    ScrollWidth
    ScrollBarWidth
    ScrollCtrlWidth
    numSizeProperties
)

// ----------------------------------------------------------------------------

type Properties interface {
    Color(n ColorPropertyName) utils.Color
    Font(n FontPropertyName) font.Face
    Size(n SizePropertyName) float64
}

// ----------------------------------------------------------------------------

type DefaultProps struct {
     colorMap map[ColorPropertyName]utils.Color
     fontMap map[FontPropertyName]font.Face
     sizeMap map[SizePropertyName]float64
}

func NewDefaultProps() (*DefaultProps) {
    p := &DefaultProps{}

    p.colorMap = map[ColorPropertyName]utils.Color{
        RedColor:      utils.Red,
        OrangeColor:   utils.Orange,
        YellowColor:   utils.Yellow,
        GreenColor:    utils.Green,
        BlueColor:     utils.Blue,
        PurpleColor:   utils.Purple,
        BrownColor:    utils.Brown,
        GrayColor:     utils.Gray,
        BlackColor:    utils.Black,
        WhiteColor:    utils.Whitesmoke,
        FillColor:     utils.Teal,
        StrokeColor:   utils.Teal,
        TextColor:     utils.Whitesmoke,
    }

    p.fontMap = map[FontPropertyName]font.Face{
        RegularFont:   utils.NewFontFace(utils.GoRegular, 15.0),
        BoldFont:      utils.NewFontFace(utils.GoBold, 15.0),
    }

    p.sizeMap = map[SizePropertyName]float64{
        LineWidth:             2.0,
        PanelLineWidth:        0.0,
        ButtonHeight:          32.0,
        ButtonInsetHorizontal: 15.0,
        ButtonInset:   	       5.0,
        ButtonCornerRad:       6.0,
        RadioBtnSize:          18.0,
        RadioBtnLineWidth:     4.0,
        RadioBtnDotSize:       8.0,
        CheckboxSize:          18.0,
        CheckboxLineWidth:     4.0,
        CheckboxCornerRad:     5.0,
        SliderWidth:           18.0,
        SliderBarWidth:        14.0,
        SliderCtrlWidth:       18.0,
        ScrollWidth:           18.0,
        ScrollBarWidth:        14.0,
        ScrollCtrlWidth:       18.0,
    }

    return p
}

func (p *DefaultProps) Color(n ColorPropertyName) utils.Color {
    return p.colorMap[n]
}

func (p *DefaultProps) Font(n FontPropertyName) font.Face {
    return p.fontMap[n]
}

func (p *DefaultProps) Size(n SizePropertyName) float64 {
    return p.sizeMap[n]
}

