//
// THIS FILE IS AUTO-GENERATED, PLEASE DO NOT EDTI
//

package props

import (
    "golang.org/x/image/font/opentype"
    "github.com/stefan-muehlebach/gg/color"
)

type PropertyEmbed struct {
    prop *Properties
}

func (pe *PropertyEmbed) Init(parent *Properties) {
    pe.prop = NewProperties(parent)
}


func (pe *PropertyEmbed) Color() (color.Color) {
    return pe.prop.Color(Color)
}
func (pe *PropertyEmbed) SetColor(c color.Color) {
    pe.prop.SetColor(Color, c)
}

func (pe *PropertyEmbed) PushedColor() (color.Color) {
    return pe.prop.Color(PushedColor)
}
func (pe *PropertyEmbed) SetPushedColor(c color.Color) {
    pe.prop.SetColor(PushedColor, c)
}

func (pe *PropertyEmbed) SelectedColor() (color.Color) {
    return pe.prop.Color(SelectedColor)
}
func (pe *PropertyEmbed) SetSelectedColor(c color.Color) {
    pe.prop.SetColor(SelectedColor, c)
}

func (pe *PropertyEmbed) BorderColor() (color.Color) {
    return pe.prop.Color(BorderColor)
}
func (pe *PropertyEmbed) SetBorderColor(c color.Color) {
    pe.prop.SetColor(BorderColor, c)
}

func (pe *PropertyEmbed) PushedBorderColor() (color.Color) {
    return pe.prop.Color(PushedBorderColor)
}
func (pe *PropertyEmbed) SetPushedBorderColor(c color.Color) {
    pe.prop.SetColor(PushedBorderColor, c)
}

func (pe *PropertyEmbed) SelectedBorderColor() (color.Color) {
    return pe.prop.Color(SelectedBorderColor)
}
func (pe *PropertyEmbed) SetSelectedBorderColor(c color.Color) {
    pe.prop.SetColor(SelectedBorderColor, c)
}

func (pe *PropertyEmbed) TextColor() (color.Color) {
    return pe.prop.Color(TextColor)
}
func (pe *PropertyEmbed) SetTextColor(c color.Color) {
    pe.prop.SetColor(TextColor, c)
}

func (pe *PropertyEmbed) PushedTextColor() (color.Color) {
    return pe.prop.Color(PushedTextColor)
}
func (pe *PropertyEmbed) SetPushedTextColor(c color.Color) {
    pe.prop.SetColor(PushedTextColor, c)
}

func (pe *PropertyEmbed) SelectedTextColor() (color.Color) {
    return pe.prop.Color(SelectedTextColor)
}
func (pe *PropertyEmbed) SetSelectedTextColor(c color.Color) {
    pe.prop.SetColor(SelectedTextColor, c)
}

func (pe *PropertyEmbed) LineColor() (color.Color) {
    return pe.prop.Color(LineColor)
}
func (pe *PropertyEmbed) SetLineColor(c color.Color) {
    pe.prop.SetColor(LineColor, c)
}

func (pe *PropertyEmbed) PushedLineColor() (color.Color) {
    return pe.prop.Color(PushedLineColor)
}
func (pe *PropertyEmbed) SetPushedLineColor(c color.Color) {
    pe.prop.SetColor(PushedLineColor, c)
}

func (pe *PropertyEmbed) SelectedLineColor() (color.Color) {
    return pe.prop.Color(SelectedLineColor)
}
func (pe *PropertyEmbed) SetSelectedLineColor(c color.Color) {
    pe.prop.SetColor(SelectedLineColor, c)
}

func (pe *PropertyEmbed) BarColor() (color.Color) {
    return pe.prop.Color(BarColor)
}
func (pe *PropertyEmbed) SetBarColor(c color.Color) {
    pe.prop.SetColor(BarColor, c)
}

func (pe *PropertyEmbed) PushedBarColor() (color.Color) {
    return pe.prop.Color(PushedBarColor)
}
func (pe *PropertyEmbed) SetPushedBarColor(c color.Color) {
    pe.prop.SetColor(PushedBarColor, c)
}

func (pe *PropertyEmbed) BackgroundColor() (color.Color) {
    return pe.prop.Color(BackgroundColor)
}
func (pe *PropertyEmbed) SetBackgroundColor(c color.Color) {
    pe.prop.SetColor(BackgroundColor, c)
}

func (pe *PropertyEmbed) MenuBackgroundColor() (color.Color) {
    return pe.prop.Color(MenuBackgroundColor)
}
func (pe *PropertyEmbed) SetMenuBackgroundColor(c color.Color) {
    pe.prop.SetColor(MenuBackgroundColor, c)
}

func (pe *PropertyEmbed) Font() (*opentype.Font) {
    return pe.prop.Font(Font)
}
func (pe *PropertyEmbed) SetFont(f *opentype.Font) {
    pe.prop.SetFont(Font, f)
}

func (pe *PropertyEmbed) Width() (float64) {
    return pe.prop.Size(Width)
}
func (pe *PropertyEmbed) SetWidth(s float64) {
    pe.prop.SetSize(Width, s)
}

func (pe *PropertyEmbed) Height() (float64) {
    return pe.prop.Size(Height)
}
func (pe *PropertyEmbed) SetHeight(s float64) {
    pe.prop.SetSize(Height, s)
}

func (pe *PropertyEmbed) BorderWidth() (float64) {
    return pe.prop.Size(BorderWidth)
}
func (pe *PropertyEmbed) SetBorderWidth(s float64) {
    pe.prop.SetSize(BorderWidth, s)
}

func (pe *PropertyEmbed) PushedBorderWidth() (float64) {
    return pe.prop.Size(PushedBorderWidth)
}
func (pe *PropertyEmbed) SetPushedBorderWidth(s float64) {
    pe.prop.SetSize(PushedBorderWidth, s)
}

func (pe *PropertyEmbed) SelectedBorderWidth() (float64) {
    return pe.prop.Size(SelectedBorderWidth)
}
func (pe *PropertyEmbed) SetSelectedBorderWidth(s float64) {
    pe.prop.SetSize(SelectedBorderWidth, s)
}

func (pe *PropertyEmbed) LineWidth() (float64) {
    return pe.prop.Size(LineWidth)
}
func (pe *PropertyEmbed) SetLineWidth(s float64) {
    pe.prop.SetSize(LineWidth, s)
}

func (pe *PropertyEmbed) InnerPadding() (float64) {
    return pe.prop.Size(InnerPadding)
}
func (pe *PropertyEmbed) SetInnerPadding(s float64) {
    pe.prop.SetSize(InnerPadding, s)
}

func (pe *PropertyEmbed) Padding() (float64) {
    return pe.prop.Size(Padding)
}
func (pe *PropertyEmbed) SetPadding(s float64) {
    pe.prop.SetSize(Padding, s)
}

func (pe *PropertyEmbed) CornerRadius() (float64) {
    return pe.prop.Size(CornerRadius)
}
func (pe *PropertyEmbed) SetCornerRadius(s float64) {
    pe.prop.SetSize(CornerRadius, s)
}

func (pe *PropertyEmbed) FontSize() (float64) {
    return pe.prop.Size(FontSize)
}
func (pe *PropertyEmbed) SetFontSize(s float64) {
    pe.prop.SetSize(FontSize, s)
}

func (pe *PropertyEmbed) BarWidth() (float64) {
    return pe.prop.Size(BarWidth)
}
func (pe *PropertyEmbed) SetBarWidth(s float64) {
    pe.prop.SetSize(BarWidth, s)
}

func (pe *PropertyEmbed) CtrlWidth() (float64) {
    return pe.prop.Size(CtrlWidth)
}
func (pe *PropertyEmbed) SetCtrlWidth(s float64) {
    pe.prop.SetSize(CtrlWidth, s)
}
