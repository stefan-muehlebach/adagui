//go:generate go run gen.go

package props

import (
	"embed"
	"encoding/json"
	"errors"
    "fmt"
    "path/filepath"
	"github.com/stefan-muehlebach/gg/color"
	"github.com/stefan-muehlebach/gg/colornames"
	"github.com/stefan-muehlebach/gg/fonts"
	"log"
    "os"
)

//go:embed config/*.json
var propFiles embed.FS

var (
	DefProps = NewPropsFromFile(nil, "DefProps.json")
)

type ColorPropertyName int

const (
	Color ColorPropertyName = iota
	PushedColor
	SelectedColor
	BorderColor
	PushedBorderColor
	SelectedBorderColor
	TextColor
	PushedTextColor
	SelectedTextColor
	LineColor
	PushedLineColor
	SelectedLineColor
	BarColor
	PushedBarColor
	BackgroundColor
	MenuBackgroundColor
//	RedColor
//	OrangeColor
//	YellowColor
//	GreenColor
//	BlueColor
//	PurpleColor
//	BrownColor
//	GrayColor
//	BlackColor
//	WhiteColor
	NumColorProperties
)

var (
	ColorPropertyList = []string{
		"Color",
		"PushedColor",
		"SelectedColor",
		"BorderColor",
		"PushedBorderColor",
		"SelectedBorderColor",
		"TextColor",
		"PushedTextColor",
		"SelectedTextColor",
		"LineColor",
		"PushedLineColor",
		"SelectedLineColor",
		"BarColor",
		"PushedBarColor",
		"BackgroundColor",
		"MenuBackgroundColor",
//		"RedColor",
//		"OrangeColor",
//		"YellowColor",
//		"GreenColor",
//		"BlueColor",
//		"PurpleColor",
//		"BrownColor",
//		"GrayColor",
//		"BlackColor",
//		"WhiteColor",
	}
)

func (p ColorPropertyName) String() string {
	return ColorPropertyList[p]
}

func (p ColorPropertyName) MarshalText() ([]byte, error) {
	return []byte(ColorPropertyList[p]), nil
}

func (p *ColorPropertyName) UnmarshalText(text []byte) error {
	txt := string(text)
	for i, t := range ColorPropertyList {
		if t == txt {
			*p = ColorPropertyName(i)
			return nil
		}
	}
	return errors.New(fmt.Sprintf("Color  property '%s' in file but not in the prop list", txt))
}

var (
	EmbedColorProps = []ColorPropertyName{
		Color,
		PushedColor,
		SelectedColor,
		BorderColor,
		PushedBorderColor,
		SelectedBorderColor,
		TextColor,
		PushedTextColor,
		SelectedTextColor,
		LineColor,
		PushedLineColor,
		SelectedLineColor,
		BarColor,
		PushedBarColor,
		BackgroundColor,
		MenuBackgroundColor,
	}
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
	NumFontProperties
)

var (
	FontPropertyList = []string{
		"Font",
		"RegularFont",
		"BoldFont",
		"ItalicFont",
		"BoldItalicFont",
		"MonoFont",
		"MonoBoldFont",
	}
)

func (p FontPropertyName) String() string {
	return FontPropertyList[p]
}

func (p FontPropertyName) MarshalText() ([]byte, error) {
	return []byte(FontPropertyList[p]), nil
}

func (p *FontPropertyName) UnmarshalText(text []byte) error {
	txt := string(text)
	for i, t := range FontPropertyList {
		if t == txt {
			*p = FontPropertyName(i)
			return nil
		}
	}
	return errors.New(fmt.Sprintf("Font property '%s' in file but not in the prop list", txt))
}

var (
	EmbedFontProps = []FontPropertyName{
		Font,
	}
)

type SizePropertyName int

const (
	Width SizePropertyName = iota
	Height
	BorderWidth
	PushedBorderWidth
	SelectedBorderWidth
	LineWidth
	InnerPadding
	Padding
	CornerRadius
	FontSize
	BarWidth
	CtrlWidth
	NumSizeProperties
)

var (
	SizePropertyList = []string{
		"Width",
		"Height",
		"BorderWidth",
		"PushedBorderWidth",
		"SelectedBorderWidth",
		"LineWidth",
		"InnerPadding",
		"Padding",
		"CornerRadius",
		"FontSize",
		"BarWidth",
		"CtrlWidth",
	}
)

func (p SizePropertyName) String() string {
	return SizePropertyList[p]
}

func (p SizePropertyName) MarshalText() ([]byte, error) {
	return []byte(SizePropertyList[p]), nil
}

func (p *SizePropertyName) UnmarshalText(text []byte) error {
	txt := string(text)
	for i, t := range SizePropertyList {
		if t == txt {
			*p = SizePropertyName(i)
			return nil
		}
	}
	return errors.New(fmt.Sprintf("Size property '%s' in file but not in the prop list", txt))
}

var (
	EmbedSizeProps = []SizePropertyName{
		Width,
		Height,
		BorderWidth,
		PushedBorderWidth,
		SelectedBorderWidth,
		LineWidth,
		InnerPadding,
		Padding,
		CornerRadius,
		FontSize,
		BarWidth,
		CtrlWidth,
	}
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
	ColorMap map[ColorPropertyName]color.Color
	FontMap  map[FontPropertyName]*fonts.Font
	SizeMap  map[SizePropertyName]float64
}

// Erzeugt ein neues Property-Objekt und hinterlegt parent als Vater-Property.
func NewProperties(parent *Properties) *Properties {
	p := &Properties{}

	p.parent = parent
	p.ColorMap = make(map[ColorPropertyName]color.Color)
	p.FontMap = make(map[FontPropertyName]*fonts.Font)
	p.SizeMap = make(map[SizePropertyName]float64)

	return p
}

func NewPropsFromFile(parent *Properties, fileName string) *Properties {
    data, err := propFiles.ReadFile(filepath.Join("config", fileName))
	if err != nil {
		log.Fatal(err)
	}
    return NewPropsFromData(parent, data)
}

func NewPropsFromUser(parent *Properties, fileName string) *Properties {
    data, err := os.ReadFile(fileName)
	if err != nil {
		log.Fatal(err)
	}
    return NewPropsFromData(parent, data)
}

func NewPropsFromData(parent *Properties, data []byte) *Properties {
	prop := struct {
		ColorMap map[ColorPropertyName]color.RGBAF
		FontMap  map[FontPropertyName]*fonts.Font
		SizeMap  map[SizePropertyName]float64
	}{}

	err := json.Unmarshal(data, &prop)
	if err != nil {
		log.Fatal(err)
	}
	p := NewProperties(parent)
	for key, val := range prop.ColorMap {
		p.ColorMap[key] = val
	}
	p.FontMap = prop.FontMap
	p.SizeMap = prop.SizeMap
	return p
}

/*
func (p *Properties) Write(fileName string) {
	b, err := json.MarshalIndent(p, "", "\t")
	if err != nil {
		log.Fatal(err)
	}
	err = os.WriteFile(fileName, b, 0644)
	if err != nil {
		log.Fatal(err)
	}
}
*/

// Interne Funktion. Damit werden die Properties für Widget-Kategorien
// (Buttons, Labels, etc) erzeugt.
func newProps(parent *Properties, ColorMap map[ColorPropertyName]color.Color,
	FontMap map[FontPropertyName]*fonts.Font,
	SizeMap map[SizePropertyName]float64) *Properties {
	p := &Properties{}

	p.parent = parent

	if ColorMap == nil {
		ColorMap = make(map[ColorPropertyName]color.Color)
	}
	p.ColorMap = ColorMap
	if FontMap == nil {
		FontMap = make(map[FontPropertyName]*fonts.Font)
	}
	p.FontMap = FontMap
	if SizeMap == nil {
		SizeMap = make(map[SizePropertyName]float64)
	}
	p.SizeMap = SizeMap

	return p
}

var (
	NewProps func(*Properties, map[ColorPropertyName]color.Color,
		map[FontPropertyName]*fonts.Font,
		map[SizePropertyName]float64) *Properties = newProps
)

// Das sind die Hauptmethoden, um Farben, Font oder Groessen aus den
// Properties zu lesen. Kann ein Property nicht gefunden werden, dann
// wird (falls vorhanden) das Parent-Property angefragt.
func (p *Properties) Color(name ColorPropertyName) color.Color {
	if col, ok := p.ColorMap[name]; !ok && p.parent != nil {
		return p.parent.Color(name)
	} else {
		return col
	}
}

func (p *Properties) Font(name FontPropertyName) *fonts.Font {
	if fnt, ok := p.FontMap[name]; !ok && p.parent != nil {
		return p.parent.Font(name)
	} else {
		return fnt
	}
}

func (p *Properties) Size(name SizePropertyName) float64 {
	if siz, ok := p.SizeMap[name]; !ok && p.parent != nil {
		return p.parent.Size(name)
	} else {
		return siz
	}
}

// Über diese Methoden können einzelne Eigenschaften auf Typen- oder Objekt-
// ebene definiert werden.
func (p *Properties) SetColor(name ColorPropertyName, col color.Color) {
	p.ColorMap[name] = col
}

func (p *Properties) SetFont(name FontPropertyName, fnt *fonts.Font) {
	p.FontMap[name] = fnt
}

func (p *Properties) SetSize(name SizePropertyName, size float64) {
	p.SizeMap[name] = size
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
	delete(p.ColorMap, name)
}

func (p *Properties) DelFont(name FontPropertyName) {
	if p.parent == nil {
		return
	}
	delete(p.FontMap, name)
}

func (p *Properties) DelSize(name SizePropertyName) {
	if p.parent == nil {
		return
	}
	delete(p.SizeMap, name)
}

// ----------------------------------------------------------------------------

// Erstellt ein neues Default-Property Objekt. Die Default Properties muessen
// fur jedes Property einen Wert bereitstellen. Mit den Tests in props_test.go
// kann geprüft werden, ob dies erfüllt ist.
func NewDefaultProps() *Properties {
	var c1, c2, c3 color.Color

	p := &Properties{}

	//c1 := colornames.Navy.Dark(0.4)
	//c2 := colornames.DeepSkyBlue.Dark(0.3)
	//c3 := colornames.DeepSkyBlue.Dark(0.3)
	//c1 := colornames.DarkRed
	//c2 := colornames.Gold
	//c3 := colornames.YellowGreen
	c1 = colornames.DarkGreen
	c2 = c1.Interpolate(colornames.YellowGreen, 0.9)
	c3 = c1.Interpolate(colornames.YellowGreen, 0.7)

	p.ColorMap = map[ColorPropertyName]color.Color{
		Color:         c1,
		PushedColor:   c2,
		SelectedColor: c3,

		BorderColor:         c1,
		PushedBorderColor:   c2,
		SelectedBorderColor: c3,

		TextColor:         colornames.WhiteSmoke,
		PushedTextColor:   colornames.Black,
		SelectedTextColor: colornames.White,

		LineColor:         colornames.WhiteSmoke,
		PushedLineColor:   colornames.Black,
		SelectedLineColor: colornames.WhiteSmoke,

		BarColor:       colornames.DarkSlateGray.Dark(0.5),
		PushedBarColor: colornames.DarkSlateGray.Dark(0.5),

		BackgroundColor:     colornames.Navy.Dark(0.8),
		MenuBackgroundColor: colornames.DarkGreen.Dark(0.8),

		// Out
		//RedColor:    colornames.Red,
		//OrangeColor: colornames.Orange,
		//YellowColor: colornames.Yellow,
		//GreenColor:  colornames.Green,
		//BlueColor:   colornames.Blue,
		//PurpleColor: colornames.Purple,
		//BrownColor:  colornames.Brown,
		//GrayColor:   colornames.Gray,
		//BlackColor:  colornames.Black,
		//WhiteColor:  colornames.WhiteSmoke,
	}

	p.FontMap = map[FontPropertyName]*fonts.Font{
		Font:           fonts.GoRegular,
		RegularFont:    fonts.GoRegular,
		BoldFont:       fonts.GoBold,
		ItalicFont:     fonts.GoItalic,
		BoldItalicFont: fonts.GoBoldItalic,
		MonoFont:       fonts.GoMono,
		MonoBoldFont:   fonts.GoMonoBold,
	}

	p.SizeMap = map[SizePropertyName]float64{
		Width:               0.0,
		Height:              0.0,
		BorderWidth:         0.0,
		PushedBorderWidth:   0.0,
		SelectedBorderWidth: 0.0,
		LineWidth:           2.5,
		InnerPadding:        5.0,
		Padding:             15.0,
		CornerRadius:        6.0,
		FontSize:            15.0,
		BarWidth:            18.0,
		CtrlWidth:           18.0,
	}
	return p
}
