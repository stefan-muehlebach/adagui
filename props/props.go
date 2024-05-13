//go:generate go run gen.go

package props

import (
	"embed"
	"encoding/json"
	"errors"
    "fmt"
    "path/filepath"
	"github.com/stefan-muehlebach/gg/color"
	"github.com/stefan-muehlebach/gg/fonts"
	"log"
    "os"
)

//go:embed *.json
var propFiles embed.FS

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
	BarSize
	CtrlSize
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
		"BarSize",
		"CtrlSize",
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
		BarSize,
		CtrlSize,
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
    data, err := propFiles.ReadFile(filepath.Join(fileName))
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
	var prop struct {
		Colors map[ColorPropertyName]color.RGBAF
		Fonts  map[FontPropertyName]*fonts.Font
		Sizes  map[SizePropertyName]float64
	}

	err := json.Unmarshal(data, &prop)
	if err != nil {
		log.Fatal(err)
	}
	p := NewProperties(parent)
	for key, val := range prop.Colors {
		p.ColorMap[key] = val
	}
    for key, val := range prop.Fonts {
        p.FontMap[key] = val
    }
    for key, val := range prop.Sizes {
        p.SizeMap[key] = val
    }
	return p
}

var (
    PropsMap map[string]*Properties
)

func init() {
    initPropsMapFromFile("Props.json")
}

func initPropsMapFromFile(fileName string) {
    data, err := propFiles.ReadFile(filepath.Join(fileName))
	if err != nil {
		log.Fatal(err)
	}
    initPropsMapFromData(data)
}

func initPropsMapFromData(data []byte) {
    var propList []struct {
        Name string
        ParentName string
		Colors map[ColorPropertyName]color.RGBAF
		Fonts  map[FontPropertyName]*fonts.Font
		Sizes  map[SizePropertyName]float64
    }
    var parent *Properties
    var ok bool

    PropsMap = make(map[string]*Properties)
	err := json.Unmarshal(data, &propList)
	if err != nil {
		log.Fatalf("failed unmarshalling data: %v", err)
	}
	for _, val := range propList {
        if val.ParentName == "" {
            parent = nil
        } else {
            if parent, ok = PropsMap[val.ParentName]; !ok {
		        log.Fatalf("on processing '%s', parent property '%s' not found", val.Name, val.ParentName)
            }
        }
        p := NewProperties(parent)
	    for colorName, colorValue := range val.Colors {
	    	p.ColorMap[colorName] = colorValue
	    }
        for key, val := range val.Fonts {
            p.FontMap[key] = val
        }
        for key, val := range val.Sizes {
            p.SizeMap[key] = val
        }
        PropsMap[val.Name] = p
    }
}

// Das sind die Hauptmethoden, um Farben, Font oder Groessen aus den
// Properties zu lesen. Kann ein Property nicht gefunden werden, dann
// wird (falls vorhanden) das Parent-Property angefragt.
func (p *Properties) Color(name ColorPropertyName) color.Color {
    var col color.Color
    var found bool

    if col, found = p.ColorMap[name]; !found && p.parent != nil {
        col = p.parent.Color(name)
        p.ColorMap[name] = col
    }
    return col
}

func (p *Properties) Font(name FontPropertyName) *fonts.Font {
    var fnt *fonts.Font
    var found bool

	if fnt, found = p.FontMap[name]; !found && p.parent != nil {
		fnt = p.parent.Font(name)
        p.FontMap[name] = fnt
	}
	return fnt
}

func (p *Properties) Size(name SizePropertyName) float64 {
    var siz float64
    var found bool

	if siz, found = p.SizeMap[name]; !found && p.parent != nil {
		siz = p.parent.Size(name)
        p.SizeMap[name] = siz
	}
	return siz
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

