
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
    . "github.com/stefan-muehlebach/adagui/props"
    "github.com/stefan-muehlebach/gg"
    //"github.com/stefan-muehlebach/gg/color"
    //"github.com/stefan-muehlebach/gg/colornames"
    "github.com/stefan-muehlebach/gg/fonts"
    "github.com/stefan-muehlebach/gg/geom"
    "golang.org/x/image/font"
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
    LabelProps = NewPropsFromFile(DefProps, "LabelProps.json")
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
func (l *Label) SetFont(fontFont *fonts.Font) {
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
    Debugf(Painting, "type %T", l.Wrapper)
    gc.SetFontFace(l.fontFace)
    gc.SetStrokeColor(l.TextColor())
    gc.DrawString(l.text.Get(), l.rPt.X, l.rPt.Y)
    // Groesse des Labels als Rechteck
    gc.DrawRectangle(l.Bounds().AsCoord())
    gc.SetStrokeColor(l.BorderColor())
    gc.SetStrokeWidth(l.BorderWidth())
    gc.Stroke()
    // Markierungen um den Bereich fuer den Text
//    gc.SetStrokeColor(colornames.Crimson)
//    gc.SetStrokeWidth(2.0)
//    pt0 := l.Bounds().Min
//    pt1 := l.Bounds().Max
    // Links oben
//    gc.MoveTo(pt0.X, pt0.Y+10.0)
//    gc.LineTo(pt0.X, pt0.Y)
//    gc.LineTo(pt0.X+10.0, pt0.Y)
    // Links unten
//    gc.MoveTo(pt0.X, pt1.Y-10.0)
//    gc.LineTo(pt0.X, pt1.Y)
//    gc.LineTo(pt0.X+10.0, pt1.Y)
    // Rechts oben
//    gc.MoveTo(pt1.X-10.0, pt0.Y)
//    gc.LineTo(pt1.X, pt0.Y)
//    gc.LineTo(pt1.X, pt0.Y+10.0)
    // Rechts unten
//    gc.MoveTo(pt1.X, pt1.Y-10.0)
//    gc.LineTo(pt1.X, pt1.Y)
//    gc.LineTo(pt1.X-10.0, pt1.Y)
//    gc.Stroke()
    // Referenzpunkt fuer den Text
//    gc.SetFillColor(colornames.Crimson)
//    gc.DrawPoint(l.rPt.X, l.rPt.Y, 5.0)
//    gc.Fill()
}

// Buttons sind neutrale Knoepfe, ohne spezifischen Inhalt, d.h. ohne Text
// oder Icons. Sie werden selten direkt verwendet, sondern dienen als
// generische Grundlage fuer die weiter unten definierten Text- oder Icon-
// Buttons.
var (
    ButtonProps = NewPropsFromFile(DefProps, "ButtonProps.json")
)

type Button struct {
    LeafEmbed
    PushEmbed
    checked bool
}

func NewButton(w, h float64) (*Button) {
    b := &Button{}
    b.Wrapper = b
    b.LeafEmbed.Init()
    b.PushEmbed.Init(b, nil)
    b.PropertyEmbed.Init(ButtonProps)
    b.SetMinSize(geom.Point{w, h})
    b.checked   = false
    return b
}

func (b *Button) Paint(gc *gg.Context) {
    gc.DrawRoundedRectangle(0.0, 0.0, b.Size().X, b.Size().Y,
            b.CornerRadius())
    if b.Pushed() {
        gc.SetFillColor(b.PushedColor())
        gc.SetStrokeColor(b.PushedBorderColor())
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
    b.PushEmbed.OnInputEvent(evt)
    b.CallTouchFunc(evt)
}

// Ein TextButton verhaelt sich analog zum neutralen Button, stellt jedoch
// zusaetzlich Text dar und passt seine Groesse diesem Text an.
type TextButton struct {
    Button
    label string
    fontFace font.Face
    align AlignType
    refPt geom.Point
    ax, ay float64
    desc float64
}

func NewTextButton(label string) (*TextButton) {
    b := &TextButton{}
    b.Wrapper = b
    b.LeafEmbed.Init()
    b.PushEmbed.Init(b, nil)
    b.PropertyEmbed.Init(ButtonProps)
    b.label = label
    b.align = AlignCenter | AlignMiddle
    b.updateSize()
    return b
}

func (b *TextButton) SetSize(size geom.Point) {
    b.Button.SetSize(size)
    b.updateRefPoint()
}

func (b *TextButton) Text() (string) {
    return b.label
}
func (b *TextButton) SetText(str string) {
    b.label = str
    b.updateSize()
}

func (b *TextButton) Align() (AlignType) {
    return b.align
}
func (b *TextButton) SetAlign(a AlignType) {
    b.align = a
    b.updateRefPoint()
}

func (b *TextButton) SetFont(fontFont *fonts.Font) {
    b.PropertyEmbed.SetFont(fontFont)
    b.updateSize()
}

func (b *TextButton) SetFontSize(fontSize float64) {
    b.PropertyEmbed.SetFontSize(fontSize)
    b.updateSize()
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
    if b.align & AlignLeft != 0 {
        b.refPt.X = b.Bounds().Min.X + b.InnerPadding()
        b.ax = 0.0
    } else if b.align & AlignCenter != 0 {
        b.refPt.X = b.Bounds().Center().X
        b.ax = 0.5
    } else {
        b.refPt.X = b.Bounds().Max.X - b.InnerPadding()
        b.ax = 1.0
    }
    if b.align & AlignBottom != 0 {
        b.refPt.Y = b.Bounds().Max.Y - b.Padding()
        b.ay = 0.0
    } else if b.align & AlignMiddle != 0 {
        b.refPt.Y = b.Bounds().Center().Y
        b.ay = 0.5
    } else {
        b.refPt.Y = b.Bounds().Min.Y + b.Padding()
        b.ay = 1.0
    }
}

func (b *TextButton) Paint(gc *gg.Context) {
    b.Button.Paint(gc)
    gc.SetFontFace(b.fontFace)
    if b.Pushed() {
        gc.SetStrokeColor(b.PushedTextColor())
    } else {
        if b.checked {
            gc.SetStrokeColor(b.SelectedTextColor())
        } else {
            gc.SetStrokeColor(b.TextColor())
        }
    }
    gc.DrawStringAnchored(b.label, b.refPt.X, b.refPt.Y, b.ax, b.ay)
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
    b.LeafEmbed.Init()
    b.PushEmbed.Init(b, nil)
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
    if b.Pushed() {
        gc.SetStrokeColor(b.PushedTextColor())
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
    if b.Pushed() {
        gc.SetStrokeColor(b.PushedBorderColor())
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
    if b.Pushed() {
        gc.SetFillColor(b.PushedLineColor())
        gc.SetStrokeColor(b.PushedLineColor())
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
    b.PushEmbed.OnInputEvent(evt)
    switch evt.Type {
    case touch.TypePress, touch.TypeEnter:
        b.next()
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
var (
    IconButtonProps = NewPropsFromFile(ButtonProps, "IconBtnProps.json")
)

type IconButton struct {
    Button
    img image.Image
    btnData int
    data binding.Int
}

func NewIconButton(imgFile string) (*IconButton) {
    b := &IconButton{}
    b.Wrapper = b
    b.LeafEmbed.Init()
    b.PushEmbed.Init(b, nil)
    b.PropertyEmbed.Init(IconButtonProps)
    b.img, _ = gg.LoadPNG(imgFile)
    i := b.InnerPadding()
    //rect := geom.NewRectangleIMG(b.img.Bounds())
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
    TabButtonProps = NewPropsFromFile(ButtonProps, "TabBtnProps.json")
/*
        map[ColorPropertyName]color.Color{
            Color:             ButtonProps.Color(Color).Alpha(0.4),
            BorderColor:       ButtonProps.Color(BorderColor).Alpha(0.4),
            TextColor:         ButtonProps.Color(TextColor).Alpha(0.4),
            SelectedTextColor: colornames.Black,
        },
        nil,
        map[SizePropertyName]float64{
            Width:        30.0,
            Height:       18.0,
            BorderWidth:   0.0,
            CornerRadius: 10.0,
            Padding:      15.0,
            FontSize:     12.0,
        })
*/
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
    b.LeafEmbed.Init()
    b.PushEmbed.Init(b, nil)
    b.PropertyEmbed.Init(TabButtonProps)
    b.label     = label
    b.fontFace  = fonts.NewFace(b.Font(), b.FontSize())
    w := (float64(font.MeasureString(b.fontFace, b.label))/64.0) +
            (2.0*b.Padding())
    h := b.Height()
    b.SetMinSize(geom.Point{w, h})
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
    gc.DrawRoundedRectangle(0.0, 0.0,
            b.Size().X, b.Size().Y+b.CornerRadius(),
            b.CornerRadius())
    if b.Pushed() {
        gc.SetFillColor(b.PushedColor())
        gc.SetStrokeColor(b.PushedBorderColor())
    } else {
        if b.checked {
            gc.SetFillColor(b.SelectedColor())
            gc.SetStrokeColor(b.SelectedBorderColor())
        } else {
            gc.SetFillColor(b.Color())
            gc.SetStrokeColor(b.BorderColor())
        }
    }
    gc.FillStroke()

    mp := b.Bounds().Center()
    if b.Pushed() {
        gc.SetStrokeColor(b.PushedTextColor())
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
    b.PushEmbed.OnInputEvent(evt)
    switch evt.Type {
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
    CheckboxProps = NewPropsFromFile(ButtonProps, "CheckboxProps.json")
/*
        map[FontPropertyName]*fonts.Font{
            Font:         fonts.GoRegular,
        },
        map[SizePropertyName]float64{
            Width:        18.0,
            Height:       18.0,
            LineWidth:     4.0,
            CornerRadius:  5.0,
        })
*/
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
    c.LeafEmbed.Init()
    c.PushEmbed.Init(c, nil)
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
    if c.Pushed() {
        gc.SetFillColor(c.PushedColor())
        gc.SetStrokeColor(c.PushedBorderColor())
    } else {
        gc.SetFillColor(c.Color())
        gc.SetStrokeColor(c.BorderColor())
    }
    gc.SetStrokeWidth(c.BorderWidth())
    gc.FillStroke()
    if c.Checked() {
        gc.SetStrokeWidth(c.LineWidth())
        if c.Pushed() {
            gc.SetStrokeColor(c.PushedLineColor())
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
    RadioButtonProps = NewPropsFromFile(ButtonProps, "RadioBtnProps.json")
/*
        map[FontPropertyName]*fonts.Font{
            Font:         fonts.GoRegular,
        },
        map[SizePropertyName]float64{
            Width:        18.0,
            Height:       18.0,
            LineWidth:     8.0,
        })
*/
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
    b.LeafEmbed.Init()
    b.PushEmbed.Init(b, nil)
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
    if b.Pushed() {
        gc.SetFillColor(b.PushedColor())
        gc.SetStrokeColor(b.PushedBorderColor())
    } else {
        gc.SetFillColor(b.Color())
        gc.SetStrokeColor(b.BorderColor())
    }
    gc.SetStrokeWidth(b.BorderWidth())
    gc.FillStroke()
    if b.checked {
        if b.Pushed() {
	        gc.SetFillColor(b.PushedLineColor())
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
    ScrollbarProps = NewPropsFromFile(DefProps, "ScrollbarProps.json")
/*
        map[SizePropertyName]float64{
            Width:    18.0,
            Height:   18.0,
            BarWidth:  3.0,
        })
*/
)

type Scrollbar struct {
    LeafEmbed
    PushEmbed
    orient Orientation
    initValue, visiRange float64
    value binding.Float
    barStart, barEnd geom.Point
    ctrlStart, ctrlEnd geom.Point
    dp1, dp2, startPt, endPt1, endPt2 geom.Point
    barLen, ctrlLen float64
    isDragging bool
}

func NewScrollbar(len float64, orient Orientation) (*Scrollbar) {
    s := &Scrollbar{}
    s.Wrapper = s
    s.LeafEmbed.Init()
    s.PropertyEmbed.Init(ScrollbarProps)
    s.PushEmbed.Init(s, nil)
    s.orient = orient
    d := max(0.5*s.BarWidth(), 0.5*s.CtrlWidth())
    if s.orient == Horizontal {
        s.SetMinSize(geom.Point{len, s.Height()})
        s.barStart = geom.Point{0, d}
    } else {
        s.SetMinSize(geom.Point{s.Width(), len})
        s.barStart = geom.Point{d, 0}
    }
    s.initValue = 0.0
    s.visiRange = 0.1
    s.value     = binding.NewFloat()
    s.SetValue(s.initValue)
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
    s.updateCtrl()
}

func (s *Scrollbar) VisiRange() (float64) {
    return s.visiRange
}
func (s *Scrollbar) SetVisiRange(vr float64) {
    if vr < 0.0 || vr > 1.0 {
        return
    }
    s.visiRange = vr
    s.updateCtrl()
}

func (s *Scrollbar) Value() (float64) {
    return s.value.Get()
}
func (s *Scrollbar) SetValue(v float64) {
    if v > 1.0 { v = 1.0 }
    if v < 0.0 { v = 0.0 }
    s.value.Set(v)
    s.updateCtrl()
}

func (s *Scrollbar) updateCtrl() {
    s.barEnd = s.Size().Sub(s.barStart)
    if s.orient == Horizontal {
        s.barLen = s.barEnd.X - s.barStart.X
    } else {
        s.barLen = s.barEnd.Y - s.barStart.Y
    }
    s.ctrlLen = s.visiRange * s.barLen
    l := s.barLen - s.ctrlLen
    if s.orient == Horizontal {
        s.ctrlStart = s.barStart.AddXY(s.Value() * l, 0.0)
        s.ctrlEnd = s.ctrlStart.AddXY(s.ctrlLen, 0.0)
    } else {
        s.ctrlStart = s.barStart.AddXY(0.0, s.Value() * l)
        s.ctrlEnd = s.ctrlStart.AddXY(0.0, s.ctrlLen)
    }
}

func (s *Scrollbar) Paint(gc *gg.Context) {
//    var pt1, pt2 geom.Point
    if s.Pushed() {
        gc.SetStrokeColor(s.PushedBarColor())
    } else {
        gc.SetStrokeColor(s.BarColor())
    }
    gc.SetStrokeWidth(s.BarWidth())
    gc.DrawLine(s.barStart.X, s.barStart.Y, s.barEnd.X, s.barEnd.Y)
    gc.Stroke()

    if s.Pushed() {
        gc.SetStrokeColor(s.PushedColor())
    } else {
        gc.SetStrokeColor(s.Color())
    }
    gc.SetStrokeWidth(s.CtrlWidth())
    gc.DrawLine(s.ctrlStart.X, s.ctrlStart.Y, s.ctrlEnd.X, s.ctrlEnd.Y)
    gc.Stroke()
}

func (s *Scrollbar) OnInputEvent(evt touch.Event) {
    //log.Printf("%T: %v", s, evt)
    s.PushEmbed.OnInputEvent(evt)
    switch evt.Type {
    case touch.TypePress:
        if s.orient == Horizontal {
            if evt.Pos.X >= s.ctrlStart.X && evt.Pos.X <= s.ctrlEnd.X {
                s.isDragging = true
            } else {
                s.isDragging = false
            }
        } else {
            if evt.Pos.Y >= s.ctrlStart.Y && evt.Pos.Y <= s.ctrlEnd.Y {
                s.isDragging = true
            } else {
                s.isDragging = false
            }
        }
    case touch.TypeDrag:
        if !s.isDragging {
            break
        }
        v := 0.0
        if s.orient == Horizontal {
            r := s.Rect().Inset(s.barStart.X+s.ctrlLen/2, s.barStart.Y)
            v, _ = r.PosRel(evt.Pos)
        } else {
            r := s.Rect().Inset(s.barStart.X, s.barStart.Y+s.ctrlLen/2)
            _, v = r.PosRel(evt.Pos)
        }
        s.SetValue(v)
        s.Mark(MarkNeedsPaint)
    case touch.TypeTap:
        v := s.Value()
        if s.orient == Horizontal {
            if evt.Pos.X < s.ctrlStart.X {
                v -= s.visiRange
            }
            if evt.Pos.X > s.ctrlEnd.X {
                v += s.visiRange
            }
        } else {
            if evt.Pos.Y < s.ctrlStart.Y {
                v -= s.visiRange
            }
            if evt.Pos.Y > s.ctrlEnd.Y {
                v += s.visiRange
            }
        }
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
    SliderProps = NewPropsFromFile(DefProps, "SliderProps.json")
/*(
        DefProps,
        nil,
        nil,
        map[SizePropertyName]float64{
            Width:    18.0,
            Height:   18.0,
            BarWidth:  3.0,
        })
*/
)

type Slider struct {
    LeafEmbed
    PushEmbed
    orient Orientation
    initValue, minValue, maxValue, stepSize float64
    value binding.Float
    barStart, barEnd geom.Point
    barLen float64
    ctrlPos geom.Point
}

func NewSlider(len float64, orient Orientation) (*Slider) {
    s := &Slider{}
    s.Wrapper = s
    s.LeafEmbed.Init()
    s.PropertyEmbed.Init(SliderProps)
    s.PushEmbed.Init(s, nil)

    s.orient = orient
    d := max(0.5*s.BarWidth(), 0.5*s.CtrlWidth())
    if s.orient == Horizontal {
        s.SetMinSize(geom.Point{len, s.Height()})
        s.barStart = geom.Point{0, d}
    } else {
        s.SetMinSize(geom.Point{s.Width(), len})
        s.barStart = geom.Point{d, 0}
    }
    s.initValue = 0.0
    s.minValue  = 0.0
    s.maxValue  = 1.0
    s.stepSize  = 0.1
    s.value     = binding.NewFloat()
    s.SetValue(s.initValue)
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
    s.updateCtrl()
}

func (s *Slider) updateCtrl() {
    s.barEnd = s.Size().Sub(s.barStart)
    if s.orient == Horizontal {
        s.barLen = s.barEnd.X - s.barStart.X
        p0 := s.barStart.AddXY(0.5*s.CtrlWidth(), 0)
        p1 := s.barEnd.AddXY(-0.5*s.CtrlWidth(), 0)
        s.ctrlPos = p0.Interpolate(p1, s.Factor())
    } else {
        s.barLen = s.barEnd.Y - s.barStart.Y
        p0 := s.barStart.AddXY(0, 0.5*s.CtrlWidth())
        p1 := s.barEnd.AddXY(0, -0.5*s.CtrlWidth())
        s.ctrlPos = p0.Interpolate(p1, 1.0 - s.Factor())
    }
}

func (s *Slider) Paint(gc *gg.Context) {
    //log.Printf("Slider.Paint()")
    if s.Pushed() {
        gc.SetStrokeColor(s.PushedBarColor())
    } else {
        gc.SetStrokeColor(s.BarColor())
    }
    gc.SetStrokeWidth(s.BarWidth())
    gc.DrawLine(s.barStart.X, s.barStart.Y, s.barEnd.X, s.barEnd.Y)
    gc.Stroke()

    if s.Pushed() {
        gc.SetStrokeColor(s.PushedColor())
    } else {
        gc.SetStrokeColor(s.Color())
    }
    gc.SetStrokeWidth(s.CtrlWidth())
    if s.orient == Horizontal {
        gc.DrawLine(s.ctrlPos.X-0.5, s.ctrlPos.Y, s.ctrlPos.X+0.5, s.ctrlPos.Y)
    } else {
        gc.DrawLine(s.ctrlPos.X, s.ctrlPos.Y-0.5, s.ctrlPos.X, s.ctrlPos.Y+0.5)
    }
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
    s.updateCtrl()
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
    s.SetValue((1.0-f)*s.minValue + f*s.maxValue)
    
}

func (s *Slider) Factor() (float64) {
    return (s.Value()-s.minValue)/(s.maxValue-s.minValue)
}

func (s *Slider) OnInputEvent(evt touch.Event) {
    s.PushEmbed.OnInputEvent(evt)
    switch evt.Type {
    case touch.TypeDrag:
        v := 0.0
        r := s.Rect().Inset(0.5*s.CtrlWidth(), 0.5*s.CtrlWidth())
        if s.orient == Horizontal {
            v, _ = r.PosRel(evt.Pos)
        } else {
            _, v = r.PosRel(evt.Pos)
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

