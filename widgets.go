
// In diesem File befinden sich alle Widgets, die im Zusammenhang mit adagui
// existieren. Aktuell sind dies:
//
// Leaf Widgets (GUI bezogen)
// --------------------------
//   Button
//   TextButton
//   IconButton
//   RadioButton
//   Checkbox
//   Slider     A.k.a. Scrollbar
//   PageButton
//   Label      Nur fuer kurze, einzeilige Texte
//   Text       (Geplant, wie Label aber fuer groessere Textmengen mit spez.
//              Ausrichtung, waehlbarer Schrift und ev. Scrollbalken)
//
package adagui

import (
    "image"
    "log"
    "math"
    "github.com/stefan-muehlebach/adagui/binding"
    "github.com/stefan-muehlebach/adagui/touch"
    "github.com/stefan-muehlebach/gg"
    "github.com/stefan-muehlebach/gg/color"
    "github.com/stefan-muehlebach/gg/colornames"
    "github.com/stefan-muehlebach/gg/fonts"
    "github.com/stefan-muehlebach/gg/geom"
    "golang.org/x/image/font"
    "golang.org/x/image/font/opentype"
)

// Die init-Funktion ist vorallem hier, damit der Import des Log-Packages
// nicht bei jedem Aktivieren von Debug-Meldungen (aus)kommentiert werden muss.
func init() {
    log.SetFlags(log.Lmicroseconds | log.Lmsgprefix)
    log.SetPrefix(": ")
}

// Der Typ Border wird fuer die Bezeichnung der vier Bildschirmseiten oder
// Richtungen verwendet.
type Border int

const (
    Left Border = iota
    Top
    Right
    Bottom
)

// Mit dem Typ Orientation koennen horizontale Ausrichtungen gegenueber
// vertikalen abgegrenzt werden.
type Orientation int

const (
    Horizontal Orientation = iota
    Vertical
)

// Der Typ AlignType dient der Ausrichtung von Text.
type AlignType int

const (
    AlignLeft AlignType = 1 << iota
    AlignCenter
    AlignRight
    AlignTop
    AlignMiddle
    AlignBottom
)

// Fuer die visuelle Abgrenzung in Box-Layouts.
type Separator struct {
    LeafEmbed
    orient Orientation
}

func NewSeparator(orient Orientation) (*Separator) {
    s := &Separator{}
    s.Wrapper = s
    s.Init()
    s.PropertyEmbed.Init(DefProps)
    s.orient = orient
    s.SetMinSize(geom.Point{s.LineWidth(), s.LineWidth()})
    return s
}

func (s *Separator) Paint(gc *gg.Context) {
    gc.SetStrokeColor(s.BarColor())
    gc.SetStrokeWidth(s.LineWidth())
    gc.MoveTo(s.Bounds().W().AsCoord())
    gc.LineTo(s.Bounds().E().AsCoord())
    gc.Stroke()
}

// Unter einem Label verstehen wir einfach eine Konserve für Text,
// kurzen Text!
var (
    LabelProps = NewProps(DefProps, nil, nil, nil)
)

type Label struct {
    LeafEmbed
    text binding.String
    fontFace font.Face
    align AlignType
    rPt geom.Point
    desc float64
}

func newLabel() (*Label) {
    l := &Label{}
    l.Wrapper = l
    l.Init()
    l.PropertyEmbed.Init(LabelProps)
    l.align = AlignLeft | AlignBottom
    return l
}

func NewLabel(txt string) (*Label) {
    l := newLabel()
    l.text = binding.NewString()
    l.text.Set(txt)
    l.updateSize()
    return l
}

func NewLabelWithData(data binding.String) (*Label) {
    l := newLabel()
    l.text = data
    l.updateSize()
    return l
}

func (l *Label) SetSize(size geom.Point) {
    l.LeafEmbed.SetSize(size)
    l.updateRefPoint()
}

func (l *Label) Align() (AlignType) {
    return l.align
}
func (l *Label) SetAlign(a AlignType) {
    l.align = a
    l.updateRefPoint()
}

func (l *Label) Text() (string) {
    return l.text.Get()
}
func (l *Label) SetText(str string) {
    l.text.Set(str)
    l.updateSize()
}

// Die Property-Funktionen SetFont und SetFontSize muessen ueberschrieben
// werden, da sie ggf. die Groesse des Widgets beeinflussen.
func (l *Label) SetFont(fontFont *opentype.Font) {
    l.PropertyEmbed.SetFont(fontFont)
    l.updateSize()
}
func (l *Label) SetFontSize(fontSize float64) {
    l.PropertyEmbed.SetFontSize(fontSize)
    l.updateSize()
}

func (l *Label) updateSize() {
    l.fontFace = fonts.NewFace(l.Font(), l.FontSize())
    w := float64(font.MeasureString(l.fontFace, l.Text())) / 64.0
    h := float64(l.fontFace.Metrics().Ascent +
            l.fontFace.Metrics().Descent) / 64.0
    l.desc = float64(l.fontFace.Metrics().Descent) / 64.0
    l.SetMinSize(geom.Point{w, h})
    l.updateRefPoint()
}

func (l *Label) updateRefPoint() {
    switch {
    case l.align & AlignLeft != 0:
        l.rPt.X = 0.0
    case l.align & AlignCenter != 0:
        l.rPt.X = 0.5*(l.Size().X - l.MinSize().X)
    case l.align & AlignRight != 0:
        l.rPt.X = l.Size().X - l.MinSize().X
    }
    switch {
    case l.align & AlignTop != 0:
        l.rPt.Y = l.MinSize().Y - l.desc
    case l.align & AlignMiddle != 0:
        l.rPt.Y = 0.5*l.MinSize().Y + 0.5*l.Size().Y - l.desc
    case l.align & AlignBottom != 0:
        l.rPt.Y = l.Size().Y - l.desc
    }
}

func (l *Label) Paint(gc *gg.Context) {
    Debugf("type %T", l.Wrapper)
    gc.SetFontFace(l.fontFace)
    gc.SetStrokeColor(l.TextColor())
    gc.DrawString(l.text.Get(), l.rPt.X, l.rPt.Y)
    // Groesse des Labels als graues Rechteck
    gc.DrawRectangle(l.Bounds().AsCoord())
    gc.SetStrokeColor(l.BorderColor())
    gc.SetStrokeWidth(l.BorderWidth())
    gc.Stroke()
    // Markierungen um den Bereich fuer den Text
    gc.SetStrokeColor(colornames.Crimson)
    gc.SetStrokeWidth(2.0)
    pt0 := l.Bounds().Min
    pt1 := l.Bounds().Max
    // Links oben
    gc.MoveTo(pt0.X, pt0.Y+10.0)
    gc.LineTo(pt0.X, pt0.Y)
    gc.LineTo(pt0.X+10.0, pt0.Y)
    // Links unten
    gc.MoveTo(pt0.X, pt1.Y-10.0)
    gc.LineTo(pt0.X, pt1.Y)
    gc.LineTo(pt0.X+10.0, pt1.Y)
    // Rechts oben
    gc.MoveTo(pt1.X-10.0, pt0.Y)
    gc.LineTo(pt1.X, pt0.Y)
    gc.LineTo(pt1.X, pt0.Y+10.0)
    // Rechts unten
    gc.MoveTo(pt1.X, pt1.Y-10.0)
    gc.LineTo(pt1.X, pt1.Y)
    gc.LineTo(pt1.X-10.0, pt1.Y)
    gc.Stroke()
    // Referenzpunkt fuer den Text
    gc.SetFillColor(colornames.Crimson)
    gc.DrawPoint(l.rPt.X, l.rPt.Y, 5.0)
    gc.Fill()
}

// Buttons sind neutrale Knoepfe, ohne spezifischen Inhalt, d.h. ohne Text
// oder Icons. Sie werden selten direkt verwendet, sondern dienen als
// generische Grundlage fuer die weiter unten definierten Text- oder Icon-
// Buttons.
var (
    ButtonProps = NewProps(DefProps, nil,
        map[FontPropertyName]*opentype.Font{
            Font:  fonts.GoBold,
        },
        map[SizePropertyName]float64{
            Width :  32.0,
            Height:  32.0,
        })
)

type Button struct {
    LeafEmbed
    pushed bool
    checked bool
}

func NewButton(w, h float64) (*Button) {
    b := &Button{}
    b.Wrapper = b
    b.Init()
    b.PropertyEmbed.Init(ButtonProps)
    b.SetMinSize(geom.Point{w, h})
    b.pushed    = false
    b.checked   = false
    return b
}

func (b *Button) Paint(gc *gg.Context) {
    gc.DrawRoundedRectangle(0.0, 0.0, b.Size().X, b.Size().Y,
            b.CornerRadius())
    if b.pushed {
        gc.SetFillColor(b.PressedColor())
        gc.SetStrokeColor(b.PressedBorderColor())
    } else {
        if b.checked {
            gc.SetFillColor(b.SelectedColor())
            gc.SetStrokeColor(b.SelectedBorderColor())
        } else {
            gc.SetFillColor(b.Color())
            gc.SetStrokeColor(b.BorderColor())
        }
    }
    gc.SetStrokeWidth(b.BorderWidth())
    gc.FillStroke()
}

func (b *Button) OnInputEvent(evt touch.Event) {
    //log.Printf("%T: %v", b, evt)
    switch evt.Type {
    case touch.TypePress, touch.TypeEnter:
        b.pushed = true
        b.Mark(MarkNeedsPaint)
    case touch.TypeRelease, touch.TypeLeave:
        b.pushed = false
        b.Mark(MarkNeedsPaint)
    }
    b.CallTouchFunc(evt)
}

// Ein TextButton verhaelt sich analog zum neutralen Button, stellt jedoch
// zusaetzlich Text dar und passt seine Groesse diesem Text an.
type TextButton struct {
    Button
    label string
    fontFace font.Face
    rPt geom.Point
    desc float64
}

func NewTextButton(label string) (*TextButton) {
    b := &TextButton{}
    b.Wrapper = b
    b.Init()
    b.PropertyEmbed.Init(ButtonProps)
    b.label = label
    b.updateSize()
    return b
}

func (b *TextButton) SetSize(size geom.Point) {
    b.Button.SetSize(size)
    b.updateRefPoint()
}

func (b *TextButton) updateSize() {
    b.fontFace = fonts.NewFace(b.Font(), b.FontSize())
    w := float64(font.MeasureString(b.fontFace, b.label)) / 64.0
    h := b.Height()
    b.desc = float64(b.fontFace.Metrics().Descent) / 64.0
    b.SetMinSize(geom.Point{w+2*b.Padding(), h})
    b.updateRefPoint()
}

func (b *TextButton) updateRefPoint() {
    b.rPt = b.Bounds().Center()
}

func (b *TextButton) Paint(gc *gg.Context) {
    b.Button.Paint(gc)
    gc.SetFontFace(b.fontFace)
    if b.pushed {
        gc.SetStrokeColor(b.PressedTextColor())
    } else {
        if b.checked {
            gc.SetStrokeColor(b.SelectedTextColor())
        } else {
            gc.SetStrokeColor(b.TextColor())
        }
    }
    gc.DrawStringAnchored(b.label, b.rPt.X, b.rPt.Y, 0.5, 0.5)
}

func (b *TextButton) SetText(str string) {
    b.label = str
    b.updateSize()
}

func (b *TextButton) Text() (string) {
    return b.label
}

// Der Versuch, ein ListButton zu implementieren...
type ListButton struct {
    Button
    fontFace font.Face
    Options []string
    selIdx int
    Selected string
}

func NewListButton(options []string) (*ListButton) {
    b := &ListButton{}
    b.Wrapper     = b
    b.Init()
    b.PropertyEmbed.Init(ButtonProps)
    b.fontFace  = fonts.NewFace(b.Font(), b.FontSize())
    b.Options   = options
    b.selIdx    = 0
    b.Selected  = b.Options[b.selIdx]
    b.updateSize()
    return b
}

func (b *ListButton) updateSize() {
    maxWidth := 0.0
    for _, option := range b.Options {
        width := font.MeasureString(b.fontFace, option)
        if float64(width) > maxWidth {
            maxWidth = float64(width)
        }
    }
    w := maxWidth/64.0 + 2.0*b.Padding() + b.Width()
    h := b.Height()
    b.SetMinSize(geom.Point{w, h})
}

func (b *ListButton) updateSelected() {
    if b.Selected != b.Options[b.selIdx] {
        b.Selected = b.Options[b.selIdx]
        b.Mark(MarkNeedsPaint)
    }
}

func (b *ListButton) Paint(gc *gg.Context) {
    b.Button.Paint(gc)
    gc.SetFontFace(b.fontFace)
    if b.pushed {
        gc.SetStrokeColor(b.PressedTextColor())
    } else {
        if b.checked {
            gc.SetStrokeColor(b.SelectedTextColor())
        } else {
            gc.SetStrokeColor(b.TextColor())
        }
    }
    pt := geom.Point{0.6*b.Size().Y+2*b.InnerPadding(), 0.5*b.Size().Y}
    gc.DrawStringAnchored(b.Selected, pt.X, pt.Y, 0.0, 0.5)

    gc.SetStrokeWidth(b.BorderWidth())
    if b.pushed {
        gc.SetStrokeColor(b.PressedBorderColor())
    } else {
        gc.SetStrokeColor(b.BorderColor())
    }
    gc.SetLineCapButt()
    // Trennlinie zwischen Text und Pfeil (links)
    p1l := geom.Point{0.6*b.Size().Y, 0.0}
    gc.DrawLine(p1l.X, p1l.Y, p1l.X, p1l.Y+b.Size().Y)
    gc.Stroke()
    // Trennlinie zwischen Text und Pfeil (rechts)
    p1r := geom.Point{b.Size().X-0.6*b.Size().Y, 0.0}
    gc.DrawLine(p1r.X, p1r.Y, p1r.X, p1r.Y+b.Size().Y)
    gc.Stroke()

    gc.SetLineCapRound()
    if b.pushed {
        gc.SetFillColor(b.PressedLineColor())
        gc.SetStrokeColor(b.PressedLineColor())
    } else {
        gc.SetFillColor(b.LineColor())
        gc.SetStrokeColor(b.LineColor())
    }
    // Pfeil nach links
    pa := p1l.Add(geom.Point{-0.2*b.Size().Y, 0.25*b.Size().Y})
    pb := p1l.Add(geom.Point{-0.4*b.Size().Y, 0.5*b.Size().Y})
    pc := pa.Add(geom.Point{0.0, 0.5*b.Size().Y})
    gc.MoveTo(pa.AsCoord())
    gc.LineTo(pb.AsCoord())
    gc.LineTo(pc.AsCoord())
    gc.ClosePath()
    gc.FillStroke()

    // Pfeil nach rechts
    pa = p1r.Add(geom.Point{0.2*b.Size().Y, 0.25*b.Size().Y})
    pb = p1r.Add(geom.Point{0.4*b.Size().Y, 0.5*b.Size().Y})
    pc = pa.Add(geom.Point{0.0, 0.5*b.Size().Y})
    gc.MoveTo(pa.AsCoord())
    gc.LineTo(pb.AsCoord())
    gc.LineTo(pc.AsCoord())
    gc.ClosePath()
    gc.FillStroke()
}

func (b *ListButton) OnInputEvent(evt touch.Event) {
    //log.Printf("%T: %v", b, evt)
    switch evt.Type {
    case touch.TypePress, touch.TypeEnter:
        b.pushed = true
        b.next()
        b.Mark(MarkNeedsPaint)
    case touch.TypeRelease, touch.TypeLeave:
        b.pushed = false
        b.Mark(MarkNeedsPaint)
    }
    b.CallTouchFunc(evt)
}

func (b *ListButton) SetSelectedIndex(i int) {
    if i < 0 || i >= len(b.Options) {
        return
    }
    b.selIdx = i
    b.updateSelected()
}

func (b *ListButton) SelectedIndex() (int) {
    return b.selIdx
}

func (b *ListButton) SetOptions(options []string) {
    b.Options = options
    b.updateSize()
    b.selIdx = 0
    b.updateSelected()
}

func (b *ListButton) next() {
    b.selIdx++
    if b.selIdx == len(b.Options) {
        b.selIdx -= len(b.Options)
    }
    b.updateSelected()
}

func (b *ListButton) prev() {
    if b.selIdx == 0 {
        b.selIdx += len(b.Options)
    }
    b.selIdx--
    b.updateSelected()
}

// Der IconButton stellt ein kleines Bild dar, welches als PNG-Datei beim
// Erstellen des Buttons angegeben wird. Die Groesse des Buttons passt sich
// der Groess der Bilddatei an.
type IconButton struct {
    Button
    img image.Image
    btnData int
    data binding.Int
}

func NewIconButton(imgFile string) (*IconButton) {
    b := &IconButton{}
    b.Wrapper = b
    b.Init()
    b.PropertyEmbed.Init(ButtonProps)
    b.img, _ = gg.LoadPNG(imgFile)
    i := b.InnerPadding()
    rect := geom.NewRectangleIMG(b.img.Bounds()).Inset(-i, -i)
    b.SetMinSize(rect.Size())
    b.data = binding.NewInt()
    return b
}

func NewIconButtonWithCallback(imgFile string, btnData int, callback func(int)) (*IconButton) {
    b := NewIconButton(imgFile)
    b.data.AddCallback(func (data binding.DataItem) {
        callback(data.(binding.Int).Get())
    })
    return b
}

func NewIconButtonWithData(imgFile string, btnData int, data binding.Int) (*IconButton) {
    b := NewIconButton(imgFile)
    b.btnData = btnData
    b.data = data
    b.data.AddListener(b)
    return b
}

func (b *IconButton) OnInputEvent(evt touch.Event) {
    //log.Printf("%T: %v", b, evt)
    b.Button.OnInputEvent(evt)
    switch evt.Type {
    case touch.TypeTap:
        if !b.checked {
            b.data.Set(b.btnData)
        } else {
            b.data.Set(-1)
        }
        b.Mark(MarkNeedsPaint)
    }
}

func (b *IconButton) Paint(gc *gg.Context) {
    //log.Printf("IconButton.Paint()")
    b.Button.Paint(gc)
    cp := b.Bounds().Center()
    gc.DrawImageAnchored(b.img, cp.X, cp.Y, 0.5, 0.5)
}

func (b *IconButton) DataChanged(data binding.DataItem) {
    val := data.(binding.Int).Get()
    if b.btnData == val {
        b.checked = true
    } else {
        b.checked = false
    }
}

var (
    TabButtonProps = NewProps(ButtonProps,
        map[ColorPropertyName]color.Color{
            Color:             DefProps.Color(Color).Alpha(0.4),
            BorderColor:       DefProps.Color(BorderColor).Alpha(0.4),
            TextColor:         DefProps.Color(TextColor).Alpha(0.4),
            SelectedTextColor: colornames.Black,
        },
        nil,
        map[SizePropertyName]float64{
            Width:        30.0,
            Height:       18.0,
            CornerRadius:  8.0,
            FontSize:     12.0,
        })
)

type TabButton struct {
    Button
    label string
    fontFace font.Face
    idx int
    data binding.Int
}

func NewTabButton(label string, idx int) (*TabButton) {
    b := &TabButton{}
    b.Wrapper = b
    b.Init()
    b.PropertyEmbed.Init(TabButtonProps)
    b.SetMinSize(geom.Point{b.Width(), b.Height()})
    b.label     = label
    b.fontFace  = fonts.NewFace(b.Font(), b.FontSize())
    b.data      = binding.NewInt()
    b.idx       = idx
    return b
}

func NewTabButtonWithData(label string, idx int, data binding.Int) (*TabButton) {
    b := NewTabButton(label, idx)
    b.data = data
    b.data.AddListener(b)
    return b
}

func (b *TabButton) Paint(gc *gg.Context) {
    //log.Printf("Button.Paint()")
    gc.DrawRectangle(b.Bounds().AsCoord())
    gc.Clip()
    gc.DrawRoundedRectangle(0.0, 0.0,
            b.Size().X, b.Size().Y+b.CornerRadius(),
            b.CornerRadius())
    if b.pushed {
        gc.SetFillColor(b.PressedColor())
        gc.SetStrokeColor(b.PressedBorderColor())
    } else {
        if b.checked {
            gc.SetFillColor(b.SelectedColor())
            gc.SetStrokeColor(b.SelectedBorderColor())
        } else {
            gc.SetFillColor(b.Color())
            gc.SetStrokeColor(b.BorderColor())
        }
    }
    gc.SetStrokeWidth(b.BorderWidth())
    gc.FillStroke()
    gc.ResetClip()

    mp := b.Bounds().Center()
    if b.pushed {
        gc.SetStrokeColor(b.PressedTextColor())
    } else {
        if b.checked {
            gc.SetStrokeColor(b.SelectedTextColor())
        } else {
            gc.SetStrokeColor(b.TextColor())
        }
    }
    gc.SetFontFace(b.fontFace)
    gc.DrawStringAnchored(b.label, mp.X, mp.Y, 0.5, 0.5)
}

func (b *TabButton) OnInputEvent(evt touch.Event) {
    //log.Printf("%T: %v", b, evt)
    switch evt.Type {
    case touch.TypePress, touch.TypeEnter:
        b.pushed = true
        b.Mark(MarkNeedsPaint)
    case touch.TypeRelease, touch.TypeLeave:
        b.pushed = false
        b.Mark(MarkNeedsPaint)
    case touch.TypeTap:
        if !b.checked {
            b.data.Set(b.idx)
        }
    }
}

func (b *TabButton) TabIndex() (int) {
    return b.idx
}

func (b *TabButton) SetTabIndex(idx int) {
    b.idx = idx
}

func (b *TabButton) DataChanged(data binding.DataItem) {
    newIndex := data.(binding.Int).Get()
    if b.idx == newIndex {
        if !b.checked {
            b.checked = true
            b.Mark(MarkNeedsPaint)
        }
    } else {
        if b.checked {
            b.checked = false
            b.Mark(MarkNeedsPaint)
        }
    }
}

// Checkboxen verhalten sich sehr aehnlich zu RadioButtons, sind jedoch eigen-
// staendig und nicht Teil einer Gruppe.
var (
    CheckboxProps = NewProps(ButtonProps, nil,
        map[FontPropertyName]*opentype.Font{
            Font:         fonts.GoRegular,
        },
        map[SizePropertyName]float64{
            Width:        18.0,
            Height:       18.0,
            LineWidth:     4.0,
            CornerRadius:  5.0,
        })
)

type Checkbox struct {
    Button
    label string
    fontFace font.Face
    value binding.Bool
}

func NewCheckbox(label string) (*Checkbox) {
    c := &Checkbox{}
    c.Wrapper = c
    c.Init()
    c.PropertyEmbed.Init(CheckboxProps)
    c.label     = label
    c.fontFace  = fonts.NewFace(c.Font(), c.FontSize())
    w := float64(font.MeasureString(c.fontFace, label))/64.0
    c.SetMinSize(geom.Point{c.Width()+c.InnerPadding()+w,
        c.Height()})
    c.value = binding.NewBool()
    return c
}

func NewCheckboxWithCallback(label string, callback func(bool)) (*Checkbox) {
    c := NewCheckbox(label)
    c.value.AddCallback(func (data binding.DataItem) {
        callback(data.(binding.Bool).Get())
    })
    return c
}

func NewCheckboxWithData(label string, data binding.Bool) (*Checkbox) {
    c := NewCheckbox(label)
    c.value = data
    return c
}

func (c *Checkbox) Paint(gc *gg.Context) {
    gc.DrawRoundedRectangle(0.0, 0.0, c.Width(), c.Height(),
            c.CornerRadius())
    if c.pushed {
        gc.SetFillColor(c.PressedColor())
        gc.SetStrokeColor(c.PressedBorderColor())
    } else {
        gc.SetFillColor(c.Color())
        gc.SetStrokeColor(c.BorderColor())
    }
    gc.SetStrokeWidth(c.BorderWidth())
    gc.FillStroke()
    if c.Checked() {
        gc.SetStrokeWidth(c.LineWidth())
        if c.pushed {
            gc.SetStrokeColor(c.PressedLineColor())
        } else {
            gc.SetStrokeColor(c.LineColor())
        }
        gc.MoveTo(4, 9)
        gc.LineTo(8, 14)
        gc.LineTo(14, 5)
        gc.Stroke()
    }
    x := c.Width() + c.InnerPadding()
    y := 0.5*c.Height()
    gc.SetStrokeColor(c.TextColor())
    gc.SetFontFace(c.fontFace)
    gc.DrawStringAnchored(c.label, x, y, 0.0, 0.5)
}

func (c *Checkbox) OnInputEvent(evt touch.Event) {
    c.Button.OnInputEvent(evt)
    if evt.Type == touch.TypeTap {
        c.SetChecked(!c.Checked())
        c.Mark(MarkNeedsPaint)
    }
}

func (c *Checkbox) Checked() bool {
    return c.value.Get()
}

func (c *Checkbox) SetChecked(val bool) {
    c.value.Set(val)
}

// Der RadioButton ist insofern ein Spezialfall, als er erstens zwei Zustaende
// haben kann (aktiv und nicht aktiv) und moeglicherweise einer Gruppe von
// RadioButtons angehoert, von denen nur einer aktiviert sein kann.
var (
    RadioButtonProps = NewProps(ButtonProps, nil,
        map[FontPropertyName]*opentype.Font{
            Font:         fonts.GoRegular,
        },
        map[SizePropertyName]float64{
            Width:        18.0,
            Height:       18.0,
            LineWidth:     8.0,
        })
)

type RadioButton struct {
    Button
    label string
    fontFace font.Face
    value int
    data binding.Int
}

func NewRadioButtonWithData(label string, value int, data binding.Int) (*RadioButton) {
    b := &RadioButton{}
    b.Wrapper  = b
    b.Init()
    b.PropertyEmbed.Init(RadioButtonProps)
    b.label    = label
    b.fontFace = fonts.NewFace(b.Font(), b.FontSize())
    w := float64(font.MeasureString(b.fontFace, label))/64.0
    b.SetMinSize(geom.Point{b.Width()+b.InnerPadding()+w,
        b.Height()})
    b.value = value
    b.data = data
    b.data.AddListener(b)
    return b
}

func (b *RadioButton) Paint(gc *gg.Context) {
    //log.Printf("RadioButton.Paint()")
    mp := geom.Point{0.5*b.Width(), 0.5*b.Height()}
    gc.DrawCircle(mp.X, mp.Y, 0.5*b.Width())
    if b.pushed {
        gc.SetFillColor(b.PressedColor())
        gc.SetStrokeColor(b.PressedBorderColor())
    } else {
        gc.SetFillColor(b.Color())
        gc.SetStrokeColor(b.BorderColor())
    }
    gc.SetStrokeWidth(b.BorderWidth())
    gc.FillStroke()
    if b.checked {
        if b.pushed {
	    gc.SetFillColor(b.PressedLineColor())
        } else {
  	    gc.SetFillColor(b.LineColor())
        }
        gc.DrawCircle(mp.X, mp.Y, 0.5*b.LineWidth())
	gc.Fill()
    }
    x := b.Width() + b.InnerPadding()
    y := 0.5*b.Height()
    gc.SetStrokeColor(b.TextColor())
    gc.SetFontFace(b.fontFace)
    gc.DrawStringAnchored(b.label, x, y, 0.0, 0.5)
}

func (b *RadioButton) OnInputEvent(evt touch.Event) {
    b.Button.OnInputEvent(evt)
    if evt.Type == touch.TypeTap {
        if b.checked == true {
            return
        }
        b.data.Set(b.value)
        b.Mark(MarkNeedsPaint)
    }
}

func (b *RadioButton) DataChanged(data binding.DataItem) {
    value := data.(binding.Int).Get()
    if b.value == value {
        b.checked = true
    } else {
        b.checked = false
    }
}

// Mit Slider kann man einen Schieberegler beliebiger Laenge horizontal oder
// vertikal im GUI positionieren. Als Werte sind aktuell nur Fliesskommazahlen
// vorgesehen.
var (
    ScrollbarProps =  NewProps(DefProps, nil, nil,
        map[SizePropertyName]float64{
            Width:  18.0,
            Height: 18.0,
        })
)

type Scrollbar struct {
    LeafEmbed
    orient Orientation
    initValue, visiRange float64
    pushed bool
    value binding.Float
    barStart, barEnd geom.Point
    ctrlStart, ctrlEnd geom.Point
    dp1, dp2, startPt, endPt1, endPt2 geom.Point
}

func NewScrollbar(len float64, orient Orientation) (*Scrollbar) {
    s := &Scrollbar{}
    s.Wrapper = s
    s.Init()
    s.PropertyEmbed.Init(ScrollbarProps)
    s.orient = orient
    if s.orient == Horizontal {
        s.SetMinSize(geom.Point{len, s.Height()})
        s.barStart = geom.Point{0.5*s.BarSize(), 0.5*s.Height()}
        s.ctrlStart = geom.Point{0.5*max(s.CtrlSize(), s.BarSize()),
            0.5*s.Height()}
    } else {
        s.SetMinSize(geom.Point{s.Width(), len})
        s.barStart = geom.Point{0.5*s.Width(), 0.5*s.BarSize()}
        s.ctrlStart = geom.Point{0.5*s.Width(), 0.5*max(s.CtrlSize(),
            s.BarSize())}
    }
    s.initValue = 0.0
    s.visiRange = 0.1
    s.value     = binding.NewFloat()
    s.updateValues()
    return s
}

func NewScrollbarWithData(len float64, orient Orientation, dat binding.Float) (*Scrollbar) {
    s := NewScrollbar(len, orient)
    s.value = dat
    return s
}

func NewScrollbarWithCallback(len float64, orient Orientation,
        callback func(float64)) (*Scrollbar) {
    s := NewScrollbar(len, orient)
    s.value.AddCallback(func (data binding.DataItem) {
        callback(data.(binding.Float).Get())
    })
    return s
}

func (s *Scrollbar) SetSize(size geom.Point) {
    s.LeafEmbed.SetSize(size)
    s.updateValues()
}

func (s *Scrollbar) updateValues() {
    s.barEnd = s.Size().Sub(s.barStart)
    s.ctrlEnd = s.Size().Sub(s.ctrlStart)
}

func (s *Scrollbar) SetVisiRange(vr float64) {
    if vr < 0.0 || vr > 1.0 {
        return
    }
    s.visiRange = vr
}

func (s *Scrollbar) VisiRange() (float64) {
    return s.visiRange
}

func (s *Scrollbar) SetValue(v float64) {
    if v > 1.0 { v = 1.0 }
    if v < 0.0 { v = 0.0 }
    s.value.Set(v)
}

func (s *Scrollbar) Value() (float64) {
    return s.value.Get()
}

func (s *Scrollbar) Paint(gc *gg.Context) {
    var pt1, pt2 geom.Point
    if s.pushed {
        gc.SetStrokeColor(s.PressedBarColor())
    } else {
        gc.SetStrokeColor(s.BarColor())
    }
    gc.SetStrokeWidth(s.BarSize())
    gc.DrawLine(s.barStart.X, s.barStart.Y, s.barEnd.X, s.barEnd.Y)
    gc.Stroke()

    newVal     := 0.5*s.visiRange + s.Value()*(1.0-s.visiRange)
    startValue := newVal - 0.5*s.visiRange
    endValue   := newVal + 0.5*s.visiRange

    r := s.Bounds().Inset(s.ctrlStart.X, s.ctrlStart.Y)
    if s.orient == Horizontal {
        pt1 = r.RelPos(startValue, 0.0)
        pt2 = r.RelPos(endValue, 0.0)
    } else {
        pt1 = r.RelPos(0.0, startValue)
        pt2 = r.RelPos(0.0, endValue)
    }

    if s.pushed {
        gc.SetStrokeColor(s.PressedColor())
    } else {
        gc.SetStrokeColor(s.Color())
    }
    gc.SetStrokeWidth(s.CtrlSize())
    gc.DrawLine(pt1.X, pt1.Y, pt2.X, pt2.Y)
    gc.Stroke()
}

func (s *Scrollbar) OnInputEvent(evt touch.Event) {
    //log.Printf("%T: %v", s, evt)
    switch evt.Type {
    case touch.TypePress:
        s.pushed = true
        s.Mark(MarkNeedsPaint)
    case touch.TypeRelease:
        s.pushed = false
        s.Mark(MarkNeedsPaint)
    case touch.TypeDrag:
        r := s.Rect().Inset(s.ctrlStart.X, s.ctrlStart.Y)
        fx, fy := r.PosRel(evt.Pos)
        v := 0.0
        if s.orient == Horizontal {
            v = fx
        } else {
            v = fy
        }
        v = (v-0.5*s.visiRange)/(1.0-s.visiRange)
        s.SetValue(v)
        s.Mark(MarkNeedsPaint)
    case touch.TypeDoubleTap:
        s.SetValue(s.initValue)
        s.Mark(MarkNeedsPaint)
    }
}

// Mit Slider kann man einen Schieberegler beliebiger Laenge horizontal oder
// vertikal im GUI positionieren. Als Werte sind aktuell nur Fliesskommazahlen
// vorgesehen.
var (
    SliderProps = ScrollbarProps
)

type Slider struct {
    LeafEmbed
    orient Orientation
    initValue, minValue, maxValue, stepSize float64
    pushed bool
    value binding.Float
    barStart, barEnd geom.Point
    ctrlStart, ctrlEnd geom.Point
}

func NewSlider(len float64, orient Orientation) (*Slider) {
    s := &Slider{}
    s.Wrapper = s
    s.Init()
    s.PropertyEmbed.Init(SliderProps)
    s.orient = orient
    if s.orient == Horizontal {
        s.SetMinSize(geom.Point{len, s.Height()})
        s.barStart = geom.Point{0.5*s.BarSize(), 0.5*s.Height()}
        s.ctrlStart = geom.Point{0.5*max(s.CtrlSize(), s.BarSize()),
            0.5*s.Height()}
    } else {
        s.SetMinSize(geom.Point{s.Width(), len})
        s.barStart = geom.Point{0.5*s.Width(), 0.5*s.BarSize()}
        s.ctrlStart = geom.Point{0.5*s.Width(), 0.5*max(s.CtrlSize(),
            s.BarSize())}
    }
    s.initValue = 0.0
    s.minValue  = 0.0
    s.maxValue  = 1.0
    s.stepSize  = 0.1
    s.value     = binding.NewFloat()
    s.updateValues()
    return s
}

func NewSliderWithData(len float64, orient Orientation, dat binding.Float) (*Slider) {
    s := NewSlider(len, orient)
    s.value = dat
    return s
}

func NewSliderWithCallback(len float64, orient Orientation,
        callback func(float64)) (*Slider) {
    s := NewSlider(len, orient)
    s.value.AddCallback(func (data binding.DataItem) {
        callback(data.(binding.Float).Get())
    })
    return s
}

func (s *Slider) SetSize(size geom.Point) {
    s.LeafEmbed.SetSize(size)
    s.updateValues()
}

func (s *Slider) updateValues() {
    s.barEnd = s.Size().Sub(s.barStart)
    s.ctrlEnd = s.Size().Sub(s.ctrlStart)
}

func (s *Slider) Paint(gc *gg.Context) {
    var pt0, pt1 geom.Point
    //log.Printf("Slider.Paint()")
    if s.pushed {
        gc.SetStrokeColor(s.PressedBarColor())
    } else {
        gc.SetStrokeColor(s.BarColor())
    }
    gc.SetStrokeWidth(s.BarSize())
    gc.DrawLine(s.barStart.X, s.barStart.Y, s.barEnd.X, s.barEnd.Y)
    gc.Stroke()

    if s.orient == Horizontal {
        pt0 = s.ctrlStart.Interpolate(s.ctrlEnd, s.Factor())
        pt1 = pt0.AddXY(0.5, 0)
    } else {
        pt0 = s.ctrlStart.Interpolate(s.ctrlEnd, 1.0-s.Factor())
        pt1 = pt0.AddXY(0, 0.5)
    }

    if s.pushed {
        gc.SetStrokeColor(s.PressedColor())
    } else {
        gc.SetStrokeColor(s.Color())
    }
    gc.SetStrokeWidth(s.CtrlSize())
    gc.DrawLine(pt0.X, pt0.Y, pt1.X, pt1.Y)
    gc.Stroke()
}

func (s *Slider) SetRange(min, max, step float64) {
    s.minValue = min
    s.maxValue = max
    s.stepSize = step
    if s.Value() < s.minValue {
        s.SetValue(min)
    }
    if s.Value() > s.maxValue {
        s.SetValue(max)
    }
}

func (s *Slider) Range() (float64, float64, float64) {
    return s.minValue, s.maxValue, s.stepSize
}

func (s *Slider) SetValue(v float64) {
    v = math.Round(v/s.stepSize)*s.stepSize
    if v > s.maxValue { v = s.maxValue }
    if v < s.minValue { v = s.minValue }
    s.value.Set(v)
}

func (s *Slider) Value() (float64) {
    return s.value.Get()
}

func (s *Slider) SetInitValue(v float64) {
    s.initValue = v
    s.SetValue(v)
}

func (s *Slider) InitValue() (float64) {
    return s.initValue
}

func (s *Slider) SetFactor(f float64) {
    if f > 1.0 { f = 1.0 }
    if f < 0.0 { f = 0.0 }
    v := (1.0-f)*s.minValue + f*s.maxValue
    s.SetValue(v)
}

func (s *Slider) Factor() (float64) {
    return (s.Value()-s.minValue)/(s.maxValue-s.minValue)
}

func (s *Slider) OnInputEvent(evt touch.Event) {
    //log.Printf("%T: %v", s, evt)
    switch evt.Type {
    case touch.TypePress:
        s.pushed = true
        s.Mark(MarkNeedsPaint)
    case touch.TypeRelease:
        s.pushed = false
        s.Mark(MarkNeedsPaint)
    case touch.TypeDrag:
        r := s.Rect().Inset(s.ctrlStart.X, s.ctrlStart.Y)
        fx, fy := r.PosRel(evt.Pos)
        v := 0.0
        if s.orient == Horizontal {
            v = fx
        } else {
            v = 1.0 - fy
        }
        s.SetFactor(v)
        s.Mark(MarkNeedsPaint)
    case touch.TypeDoubleTap:
        s.SetValue(s.initValue)
        s.Mark(MarkNeedsPaint)
    }
}

// Wir wollen es wissen und machen einen auf Game-Entwickler
//type Sprite struct {
//    LeafEmbed
//    imgList []image.Image
//    curImg int
//    ticker *time.Ticker
//}
//
//func NewSprite(imgFiles ...string) (*Sprite) {
//    s := &Sprite{}
//    s.Wrapper = s
//    s.PropertyEmbed.Init(DefProps)
//    s.imgList = make([]image.Image, 0)
//    s.curImg = 0
//    s.AddImages(imgFiles...)
//    pt := geom.NewPointIMG(s.imgList[0].Bounds().Size())
//    s.SetSize(pt)
//    return s
//}
//
//func (s *Sprite) AddImages(imgFiles ...string) {
//    for _, fileName := range imgFiles {
//        fh, err := os.Open(fileName)
//        check(err)
//        img, _, err := image.Decode(fh)
//        check(err)
//        s.imgList = append(s.imgList, img)
//        fh.Close()
//    }
//}
//
//func (s *Sprite) Paint(gc *gg.Context) {
//    s.Marks.UnmarkNeedsPaint()
//    gc.Push()
//    gc.Multiply(s.Matrix())
//    gc.DrawImage(s.imgList[s.curImg], s.Rect().Min.X, s.Rect().Min.Y)
//    gc.Pop()
//}
//
//func (s *Sprite) StartAnim(dt time.Duration) {
//    s.ticker = time.NewTicker(dt)
//    go func() {
//        for {
//            <- s.ticker.C
//            s.curImg = (s.curImg + 1) % len(s.imgList)
//            //s.Mark(MarkNeedsPaint)
//            s.Win.Repaint()
//        }
//    }()
//}
//
//func (s *Sprite) StopAnim() {
//    s.ticker.Stop()
//}
//
//-----------------------------------------------------------------------------
