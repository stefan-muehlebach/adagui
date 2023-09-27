package adagui

import (
    "golang.org/x/image/font"
    "mju.net/utils"
)

// PropertyType
type PropertyType int

const (
    FillColor PropertyType = iota
    StrokeColor
    StrokeWidth
    Clipping
    RegularFont
    BoldFont
    TextColor
    ButtonSize
    ButtonInsetHorizontal
    ButtonInset
    ButtonCornerRad
    RadioBtnSize
    RadioBtnLineWidth
    RadioBtnDotSize
    CheckboxSize
    CheckboxLineWidth
    SliderSize
    SliderCtrlSize
    numPropertyTypes
)

func (pt PropertyType) String() (str string) {
    switch pt {
    case FillColor:
        str = "FillColor"
    case StrokeColor:
        str = "StrokeColor"
    case StrokeWidth:
        str = "StrokeWidth"
    case Clipping:
        str = "Clipping"
    case RegularFont:
        str = "RegularFont"
    case BoldFont:
        str = "BoldFont"
    case TextColor:
        str = "TextColor"
    case ButtonSize:
        str = "ButtonSize"
    case ButtonInsetHorizontal:
        str = "ButtonInsetHorizontal"
    case ButtonInset:
        str = "ButtonInset"
    case ButtonCornerRad:
        str = "ButtonCornerRad"
    case RadioBtnSize:
        str = "RadioBtnSize"
    case RadioBtnLineWidth:
        str = "RadioBtnLineWidth"
    case RadioBtnDotSize:
        str = "RadioBtnDotSize"
    case CheckboxSize:
        str = "CheckboxSize"
    case CheckboxLineWidth:
        str = "CheckboxLineWidth"
    case SliderSize:
        str = "SliderSize"
    case SliderCtrlSize:
        str = "SliderCtrlSize"
    default:
        str = "(unnamed property)"
    }
    return str
}

// PropertyValue
type PropertyValue interface {}

// Properties
type Properties struct {
    parent *Properties
    props map[PropertyType]PropertyValue
}

// NewProperties erstellt ein neues Objekt, in welchem beliebige Einstellungen
// abgelegt werden können. Mit parent kann ein Property-Objekt spezifiziert
// werden, welches als Default-Gefäss verwendet wird. Es können so beliebig
// viele Property-Objekte hintereinander geschaltet werden. Das Objekt an
// der Spitze der Hierarchie muss als Parent nil haben.
func NewProperties(parent *Properties) (*Properties) {
    p := &Properties{}
    p.parent = parent
    p.props = make(map[PropertyType]PropertyValue)
    return p
}

// Check prüft, ob zu jedem Property-Typ ein Wert hinterlegt ist. Dies
// ist wichtig beim Property-Objekt, welches an der Spitze der Property-
// Hierarchie steht.
func (p *Properties) Check() (bool) {
    var pt PropertyType

    for pt = 0; pt < numPropertyTypes; pt += 1 {
        _, ok := p.props[pt]
        if !ok {
            return false
        }
    }
    return true
}

// Reset löscht alle Werte in diesem Property-Objekt. Der Verweis auf das
// parent Objekt bleibt erhalten.
func (p *Properties) Reset() {
    clear(p.props)
}

// Delete kann zum spezifischen Löschen eines bestimmten Wertes verwendet
// werden.
func (p *Properties) Delete(pType PropertyType) {
    delete(p.props, pType)
}

// Mit Set kann der Wert pVal als Einstellung pType im Property-Objekt
// hinterlegt werden. Ein ggf. bestehender Wert wird dabei überschrieben.
func (p *Properties) Set(pType PropertyType, pVal PropertyValue) {
    p.props[pType] = pVal
}

// Get ermittelt in den Properties den Wert zum Typ pType. Ist in diesem
// Property-Objekt kein Wert zu diesem Typ hinterlegt, so geht die Suche
// in den Parent-Properties weiter.
func (p *Properties) Get(pType PropertyType) (pVal PropertyValue) {
    pVal, ok := p.props[pType]
    if !ok && p.parent != nil {
        pVal = p.parent.Get(pType)
    }
    return pVal
}

// GetString wird verwendet, wenn der Typ des Properties (string) bekannt
// ist. Die Konvertierung nach string erfolgt bereits in der Methode.
func (p *Properties) GetString(pType PropertyType) (string) {
    return p.Get(pType).(string)
}

// GetFloat64 funktioniert analog zu GetString für float64 Werte.
func (p *Properties) GetFloat64(pType PropertyType) (float64) {
    return p.Get(pType).(float64)
}

// GetColor funktioniert analog zu GetString für utils.Color Werte.
func (p *Properties) GetColor(pType PropertyType) (utils.Color) {
    return p.Get(pType).(utils.Color)
}

// GetBool funktioniert analog zu GetString für bool Werte.
func (p *Properties) GetBool(pType PropertyType) (bool) {
    return p.Get(pType).(bool)
}

// GetFontFace funktioniert analog zu GetString für font.Face Werte.
func (p *Properties) GetFontFace(pType PropertyType) (font.Face) {
    return p.Get(pType).(font.Face)
}

