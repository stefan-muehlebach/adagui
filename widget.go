
// In diesem File befinden sich alle Widgets, die im Zusammenhang mit adagui
// existieren. Aktuell sind dies:
//
// Container Widgets
// -----------------
//   Group
//   Panel
//   Drawer
//   Scroller   (Geplant)
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
// Leaf Widgets (Graphik bezogen)
// ------------------------------
//   Circle
//   Rectangle
//   Line       (Geplant)
//
//   Canvas     Ein ziemlicher Exot! die beiden oben aufgefuehrten
//              Graphikelemente stehen in KEINEM Zusammenhang zu diesem Typ.
//
package adagui

import (
    "container/list"
    "image"
    "log"
    "math"
//    "os"
    "sync"
//    "time"
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

// Eine Group ist die einfachste Form eines Containers. Es dient bloss als
// logisches Sammelbecken fuer Widgets und hat ein eigenes Koordinatensystem
// sowie eine Position im darüberliegenden (Parent) Widget.
type Group struct {
    ContainerEmbed
}

func NewGroup() (*Group) {
    g := &Group{}
    g.Wrapper = g
    g.Init(DefProps)
    return g
}

func (g *Group) Paint(gc *gg.Context) {
    //log.Printf("Group.Paint()")
    g.ContainerEmbed.Paint(gc)
}

// Ein Panel ist eine komplexere Variante eines Containers. Er kann eine
// Hintergrundfarbe haben und ordnet seine Kinder gem. ihren Koordinaten an. 
var (
    PanelProps = newProps(DefProps,
        map[ColorPropertyName]color.Color{
            Color:        colornames.Black,
            BorderColor:  colornames.Black,
        }, nil, nil)
)

type Panel struct {
    ContainerEmbed
    Clip bool
    virtSize, sizeDiff, viewPort, refPt geom.Point
}

func NewPanel(w, h float64) (*Panel) {
    g := &Panel{}
    g.Wrapper = g
    g.Init(PanelProps)

    g.SetMinSize(geom.Point{w, h})
    g.SetVirtualSize(g.Size())
    g.Clip        = true
    g.viewPort    = geom.Point{0, 0}
    g.refPt       = geom.Point{0, 0}
    return g
}

func (g *Panel) Paint(gc *gg.Context) {
    //log.Printf("Panel.Paint()")
    gc.SetFillColor(g.Prop.Color(Color))
    gc.SetStrokeColor(g.Prop.Color(BorderColor))
    gc.SetStrokeWidth(g.Prop.Size(BorderWidth))
    gc.DrawRectangle(g.Bounds().AsCoord())
    if g.Clip {
        gc.ClipPreserve()
    }
    gc.FillStroke()

    gc.Translate(g.refPt.AsCoord())
    g.ContainerEmbed.Paint(gc)
    if g.Clip {
        gc.ResetClip()
    }
}

func (p *Panel) VisibleRange() (geom.Point) {
    if p.virtSize.X == 0.0 && p.virtSize.Y == 0.0 {
        return geom.Point{1, 1}
    }
    vis := p.Wrapper.Size()
    vis.X /= p.virtSize.X
    vis.Y /= p.virtSize.Y
    return vis
}

func (p *Panel) SetXView(vx float64) {
    p.viewPort.X = vx
    p.refPt.X = p.sizeDiff.X * p.viewPort.X
}

func (p *Panel) SetYView(vy float64) {
    p.viewPort.Y = vy
    p.refPt.Y = p.sizeDiff.Y * p.viewPort.Y
}

func (p *Panel) ViewPort() (geom.Point) {
    return p.viewPort
}

func (p *Panel) SetVirtualSize(sz geom.Point) {
    p.virtSize = sz
    p.sizeDiff = p.Size().Sub(p.virtSize)
}

func (p *Panel) VirtualSize() (geom.Point) {
    return p.virtSize
}

// Fuer die visuelle Abgrenzung in Box-Layouts.
type Separator struct {
    LeafEmbed
    orient Orientation
}

func NewSeparator(orient Orientation) (*Separator) {
    s := &Separator{}
    s.Wrapper = s
    s.Init(DefProps)
    s.orient = orient
    s.SetMinSize(geom.Point{s.Prop.Size(LineWidth), s.Prop.Size(LineWidth)})
    return s
}

func (s *Separator) Paint(gc *gg.Context) {
    gc.SetStrokeColor(s.Prop.Color(BarColor))
    gc.SetStrokeWidth(s.Prop.Size(LineWidth))
    gc.MoveTo(s.Bounds().W().AsCoord())
    gc.LineTo(s.Bounds().E().AsCoord())
    gc.Stroke()
}

type AlignType int

const (
    AlignLeft AlignType = (1 << 0)
    AlignCenter         = (1 << 1)
    AlignRight          = (1 << 2)
    AlignTop            = (1 << 3)
    AlignMiddle         = (1 << 4)
    AlignBottom         = (1 << 5)
)

// Unter einem Label verstehen wir einfach eine Konserve für Text. Kurzen
// Text! Für die Darstellung von grösseren Textmengen, bitte Widget Text
// berücksichtigen.
var (
    LabelProps = newProps(DefProps, nil, nil, nil)
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
    l.Init(LabelProps)
    l.align = AlignLeft | AlignMiddle
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

func (l *Label) Font() (*opentype.Font) {
    return l.Prop.Font(Font)
}
func (l *Label) SetFont(fontFont *opentype.Font) {
    l.Prop.SetFont(Font, fontFont)
    l.updateSize()
}

func (l *Label) FontSize() (float64) {
    return l.Prop.Size(FontSize)
}
func (l *Label) SetFontSize(fontSize float64) {
    l.Prop.SetSize(FontSize, fontSize)
    l.updateSize()
}

func (l *Label) updateSize() {
    l.fontFace = fonts.NewFace(l.Prop.Font(Font), l.Prop.Size(FontSize))
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
    //log.Printf("Label.Paint()")
    gc.SetFontFace(l.fontFace)
    gc.SetStrokeColor(l.Prop.Color(TextColor))
    gc.DrawString(l.text.Get(), l.rPt.X, l.rPt.Y)
    // Groesse des Labels als graues Rechteck
    gc.DrawRectangle(l.Bounds().AsCoord())
    gc.SetStrokeColor(l.Prop.Color(BorderColor))
    gc.SetStrokeWidth(l.Prop.Size(BorderWidth))
    gc.Stroke()
    // Referenzpunkt fuer den Text
    //gc.SetFillColor(colornames.Lightgray)
    //gc.DrawPoint(l.rPt.X, l.rPt.Y, 5.0)
    //gc.Fill()
}

// Buttons sind neutrale Knoepfe, ohne spezifischen Inhalt, d.h. ohne Text
// oder Icons. Sie werden selten direkt verwendet, sondern dienen als
// generische Grundlage fuer die weiter unten definierten Text- oder Icon-
// Buttons.
var (
    ButtonProps = newProps(DefProps, nil,
        map[FontPropertyName]*opentype.Font{
            Font:  fonts.GoBold,
        },
        map[SizePropertyName]float64{
            Size:  32.0,
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
    b.Init(ButtonProps)
    b.SetMinSize(geom.Point{w, h})
    b.pushed    = false
    b.checked   = false
    return b
}

func (b *Button) Paint(gc *gg.Context) {
    gc.DrawRoundedRectangle(0.0, 0.0, b.Size().X, b.Size().Y,
            b.Prop.Size(CornerRadius))
    if b.pushed {
        gc.SetFillColor(b.Prop.Color(PressedColor))
        gc.SetStrokeColor(b.Prop.Color(PressedBorderColor))
    } else {
        if b.checked {
            gc.SetFillColor(b.Prop.Color(SelectedColor))
            gc.SetStrokeColor(b.Prop.Color(SelectedBorderColor))
        } else {
            gc.SetFillColor(b.Prop.Color(Color))
            gc.SetStrokeColor(b.Prop.Color(BorderColor))
        }
    }
    gc.SetStrokeWidth(b.Prop.Size(BorderWidth))
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
    b.Init(ButtonProps)
    b.label = label
    b.updateSize()
    return b
}

func (b *TextButton) SetSize(size geom.Point) {
    b.Button.SetSize(size)
    b.updateRefPoint()
}

func (b *TextButton) updateSize() {
    b.fontFace = fonts.NewFace(b.Prop.Font(Font), b.Prop.Size(FontSize))
    w := float64(font.MeasureString(b.fontFace, b.label)) / 64.0
    h := b.Prop.Size(Size)
    b.desc = float64(b.fontFace.Metrics().Descent) / 64.0
    b.SetMinSize(geom.Point{w+2*b.Prop.Size(Padding), h})
    b.updateRefPoint()
}

func (b *TextButton) updateRefPoint() {
    b.rPt = b.Bounds().Center()
}

func (b *TextButton) Paint(gc *gg.Context) {
    b.Button.Paint(gc)
    gc.SetFontFace(b.fontFace)
    if b.pushed {
        gc.SetStrokeColor(b.Prop.Color(PressedTextColor))
    } else {
        if b.checked {
            gc.SetStrokeColor(b.Prop.Color(SelectedTextColor))
        } else {
            gc.SetStrokeColor(b.Prop.Color(TextColor))
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
    b.Init(ButtonProps)
    b.fontFace  = fonts.NewFace(b.Prop.Font(Font), b.Prop.Size(FontSize))
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
    w := maxWidth/64.0 + 2.0*b.Prop.Size(Padding) + b.Prop.Size(Size)
    h := b.Prop.Size(Size)
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
        gc.SetStrokeColor(b.Prop.Color(PressedTextColor))
    } else {
        if b.checked {
            gc.SetStrokeColor(b.Prop.Color(SelectedTextColor))
        } else {
            gc.SetStrokeColor(b.Prop.Color(TextColor))
        }
    }
    pt := geom.Point{0.6*b.Size().Y+2*b.Prop.Size(InnerPadding), 0.5*b.Size().Y}
    gc.DrawStringAnchored(b.Selected, pt.X, pt.Y, 0.0, 0.5)

    gc.SetStrokeWidth(b.Prop.Size(BorderWidth))
    if b.pushed {
        gc.SetStrokeColor(b.Prop.Color(PressedBorderColor))
    } else {
        gc.SetStrokeColor(b.Prop.Color(BorderColor))
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
        gc.SetFillColor(b.Prop.Color(PressedLineColor))
        gc.SetStrokeColor(b.Prop.Color(PressedLineColor))
    } else {
        gc.SetFillColor(b.Prop.Color(LineColor))
        gc.SetStrokeColor(b.Prop.Color(LineColor))
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
    data binding.Untyped
    btnData interface {}
    UserData int
}

func NewIconButton(imgFile string) (*IconButton) {
    b := &IconButton{}
    b.Wrapper = b
    b.Init(ButtonProps)
    b.img, _ = gg.LoadPNG(imgFile)
    i := b.Prop.Size(InnerPadding)
    rect := geom.NewRectangleIMG(b.img.Bounds()).Inset(-i, -i)
    b.SetMinSize(rect.Size())
    b.data = binding.NewUntyped()
    return b
}

func NewIconButtonWithCallback(imgFile string, btnData interface {}, callback func(interface {})) (*IconButton) {
    b := NewIconButton(imgFile)
    b.data.AddCallback(func (data binding.DataItem) {
        callback(data.(binding.Untyped).Get())
    })
    return b
}

func NewIconButtonWithData(imgFile string, btnData interface {}, data binding.Untyped) (*IconButton) {
    b := NewIconButton(imgFile)
    b.data = data
    b.data.AddListener(b)
    b.btnData = btnData
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
            b.data.Set(nil)
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
    val := data.(binding.Untyped).Get()
    if b.btnData == val {
        b.checked = true
    } else {
        b.checked = false
    }
}

// TabPanel und TabButton sind fuer Tabbed Windows gedacht.
type TabPanel struct {
    ContainerEmbed
    data binding.Int
    contentList []Node
    menu  *Group
    panel *Group
}

func NewTabPanel(w, h float64) (*TabPanel) {
    p := &TabPanel{}
    p.Wrapper      = p
    p.Init(DefProps)
    p.SetMinSize(geom.Point{w, h})
    p.Layout       = NewVBoxLayout(0)
    p.data         = binding.NewInt()
    p.data.Set(-1)
    p.contentList  = make([]Node, 0)
    p.menu         = NewGroup()
    p.menu.Layout  = NewHBoxLayout(0)
    p.panel        = NewGroup()
    p.panel.Layout = NewPaddedLayout(0)
    p.data.AddCallback(func (d binding.DataItem) {
        idx := d.(binding.Int).Get()
        if (idx < 0) || (idx >= len(p.contentList)) ||
                (p.contentList[idx] == nil) {
            return
        }
        p.panel.DelAll()
        p.panel.Add(p.contentList[idx])
        p.layout()
    })
    p.Add(p.menu, p.panel)
    return p
}

func (p *TabPanel) AddTab(label string, content Node) {
    tabIndex := len(p.contentList)
    p.contentList = append(p.contentList, content)
    b := NewTabButtonWithData(label, tabIndex, p.data)
    p.menu.Add(b)
    p.layout()
}

func (p *TabPanel) SetTab(idx int) {
    p.data.Set(idx)
}

//func (p *TabPanel) Paint(gc *gg.Context) {
//    p.ContainerEmbed.Paint(gc)
//}

var (
    TabButtonProps = newProps(ButtonProps,
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
    b.Init(TabButtonProps)
    b.SetMinSize(geom.Point{b.Prop.Size(Width), b.Prop.Size(Height)})
    b.label     = label
    b.fontFace  = fonts.NewFace(b.Prop.Font(Font), b.Prop.Size(FontSize))
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
            b.Size().X, b.Size().Y+b.Prop.Size(CornerRadius),
            b.Prop.Size(CornerRadius))
    if b.pushed {
        gc.SetFillColor(b.Prop.Color(PressedColor))
        gc.SetStrokeColor(b.Prop.Color(PressedBorderColor))
    } else {
        if b.checked {
            gc.SetFillColor(b.Prop.Color(SelectedColor))
            gc.SetStrokeColor(b.Prop.Color(SelectedBorderColor))
        } else {
            gc.SetFillColor(b.Prop.Color(Color))
            gc.SetStrokeColor(b.Prop.Color(BorderColor))
        }
    }
    gc.SetStrokeWidth(b.Prop.Size(BorderWidth))
    gc.FillStroke()
    gc.ResetClip()

    mp := b.Bounds().Center()
    if b.pushed {
        gc.SetStrokeColor(b.Prop.Color(PressedTextColor))
    } else {
        if b.checked {
            gc.SetStrokeColor(b.Prop.Color(SelectedTextColor))
        } else {
            gc.SetStrokeColor(b.Prop.Color(TextColor))
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
    CheckboxProps = newProps(ButtonProps, nil,
        map[FontPropertyName]*opentype.Font{
            Font:         fonts.GoRegular,
        },
        map[SizePropertyName]float64{
            Size:         18.0,
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
    c.Init(CheckboxProps)
    c.label     = label
    c.fontFace  = fonts.NewFace(c.Prop.Font(Font), c.Prop.Size(FontSize))
    w := float64(font.MeasureString(c.fontFace, label))/64.0
    c.SetMinSize(geom.Point{c.Prop.Size(Size)+c.Prop.Size(InnerPadding)+w,
        c.Prop.Size(Size)})
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
    gc.DrawRoundedRectangle(0.0, 0.0, c.Prop.Size(Size), c.Prop.Size(Size),
            c.Prop.Size(CornerRadius))
    if c.pushed {
        gc.SetFillColor(c.Prop.Color(PressedColor))
        gc.SetStrokeColor(c.Prop.Color(PressedBorderColor))
    } else {
        gc.SetFillColor(c.Prop.Color(Color))
        gc.SetStrokeColor(c.Prop.Color(BorderColor))
    }
    gc.SetStrokeWidth(c.Prop.Size(BorderWidth))
    gc.FillStroke()
    if c.Checked() {
        gc.SetStrokeWidth(c.Prop.Size(LineWidth))
        if c.pushed {
            gc.SetStrokeColor(c.Prop.Color(PressedLineColor))
        } else {
            gc.SetStrokeColor(c.Prop.Color(LineColor))
        }
        gc.MoveTo(4, 9)
        gc.LineTo(8, 14)
        gc.LineTo(14, 5)
        gc.Stroke()
    }
    x := c.Prop.Size(Size) + c.Prop.Size(InnerPadding)
    y := 0.5*c.Size().Y
    gc.SetStrokeColor(c.Prop.Color(TextColor))
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
    RadioButtonProps = newProps(ButtonProps, nil,
        map[FontPropertyName]*opentype.Font{
            Font:         fonts.GoRegular,
        },
        map[SizePropertyName]float64{
            Size:         18.0,
            LineWidth:     8.0,
        })
)

type RadioButton struct {
    Button
    label string
    fontFace font.Face
    checked bool
    value int
    data binding.Int
}

func NewRadioButtonWithData(label string, value int, data binding.Int) (*RadioButton) {
    b := &RadioButton{}
    b.Wrapper  = b
    b.Init(RadioButtonProps)
    b.label    = label
    b.fontFace = fonts.NewFace(b.Prop.Font(Font), b.Prop.Size(FontSize))
    w := float64(font.MeasureString(b.fontFace, label))/64.0
    b.SetMinSize(geom.Point{b.Prop.Size(Size)+b.Prop.Size(InnerPadding)+w,
        b.Prop.Size(Size)})
    b.value = value
    b.data = data
    b.data.AddListener(b)
    return b
}

func (b *RadioButton) Paint(gc *gg.Context) {
    //log.Printf("RadioButton.Paint()")
    mp := geom.Point{0.5*b.Prop.Size(Size), 0.5*b.Prop.Size(Size)}
    gc.DrawCircle(mp.X, mp.Y, 0.5*b.Prop.Size(Size))
    if b.pushed {
        gc.SetFillColor(b.Prop.Color(PressedColor))
        gc.SetStrokeColor(b.Prop.Color(PressedBorderColor))
    } else {
        gc.SetFillColor(b.Prop.Color(Color))
        gc.SetStrokeColor(b.Prop.Color(BorderColor))
    }
    gc.SetStrokeWidth(b.Prop.Size(BorderWidth))
    gc.FillStroke()
    if b.checked {
        if b.pushed {
	    gc.SetFillColor(b.Prop.Color(PressedLineColor))
        } else {
  	    gc.SetFillColor(b.Prop.Color(LineColor))
        }
        gc.DrawCircle(mp.X, mp.Y, 0.5*b.Prop.Size(LineWidth))
	gc.Fill()
    }
    x := b.Prop.Size(Size) + b.Prop.Size(InnerPadding)
    y := 0.5*b.Size().Y
    gc.SetStrokeColor(b.Prop.Color(TextColor))
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
    ScrollbarProps =  newProps(DefProps, nil, nil,
        map[SizePropertyName]float64{
            Size: 18.0,
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
    s.Init(ScrollbarProps)
    s.orient = orient
    if s.orient == Horizontal {
        s.SetMinSize(geom.Point{len, s.Prop.Size(Size)})
        s.barStart = geom.Point{0.5*s.Prop.Size(BarSize), 0.5*s.Prop.Size(Size)}
        s.ctrlStart = geom.Point{0.5*max(s.Prop.Size(CtrlSize), s.Prop.Size(BarSize)),
            0.5*s.Prop.Size(Size)}
    } else {
        s.SetMinSize(geom.Point{s.Prop.Size(Size), len})
        s.barStart = geom.Point{0.5*s.Prop.Size(Size), 0.5*s.Prop.Size(BarSize)}
        s.ctrlStart = geom.Point{0.5*s.Prop.Size(Size), 0.5*max(s.Prop.Size(CtrlSize),
            s.Prop.Size(BarSize))}
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
    //log.Printf("Scrollbar.Paint()")
    if s.pushed {
        gc.SetStrokeColor(s.Prop.Color(PressedBarColor))
    } else {
        gc.SetStrokeColor(s.Prop.Color(BarColor))
    }
    gc.SetStrokeWidth(s.Prop.Size(BarSize))
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
        gc.SetStrokeColor(s.Prop.Color(PressedColor))
    } else {
        gc.SetStrokeColor(s.Prop.Color(Color))
    }
    gc.SetStrokeWidth(s.Prop.Size(CtrlSize))
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
    s.Init(SliderProps)
    s.orient = orient
    if s.orient == Horizontal {
        s.SetMinSize(geom.Point{len, s.Prop.Size(Size)})
        s.barStart = geom.Point{0.5*s.Prop.Size(BarSize), 0.5*s.Prop.Size(Size)}
        s.ctrlStart = geom.Point{0.5*max(s.Prop.Size(CtrlSize), s.Prop.Size(BarSize)),
            0.5*s.Prop.Size(Size)}
    } else {
        s.SetMinSize(geom.Point{s.Prop.Size(Size), len})
        s.barStart = geom.Point{0.5*s.Prop.Size(Size), 0.5*s.Prop.Size(BarSize)}
        s.ctrlStart = geom.Point{0.5*s.Prop.Size(Size), 0.5*max(s.Prop.Size(CtrlSize),
            s.Prop.Size(BarSize))}
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
        gc.SetStrokeColor(s.Prop.Color(PressedBarColor))
    } else {
        gc.SetStrokeColor(s.Prop.Color(BarColor))
    }
    gc.SetStrokeWidth(s.Prop.Size(BarSize))
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
        gc.SetStrokeColor(s.Prop.Color(PressedColor))
    } else {
        gc.SetStrokeColor(s.Prop.Color(Color))
    }
    gc.SetStrokeWidth(s.Prop.Size(CtrlSize))
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

// Schoene Kreise fuer Spiele oder was auch immer lassen sich mit diesem
// Widget-Typ auf den Schirm zaubern.
var (
    GeomShapeProps = newProps(DefProps,
        map[ColorPropertyName]color.Color{
            Color:              DefProps.Color(Color),
            PressedColor:       DefProps.Color(Color).Alpha(0.5),
            BorderColor:        DefProps.Color(WhiteColor),
            PressedBorderColor: DefProps.Color(WhiteColor).Alpha(0.5),
        },
        nil,
        map[SizePropertyName]float64{
            BorderWidth:        2.0,
        })
)

type Circle struct {
    LeafEmbed
    Pressed bool
}

func NewCircle(r float64) (*Circle) {
    c := &Circle{}
    c.Wrapper = c
    c.Init(GeomShapeProps)
    c.SetMinSize(geom.Point{2*r, 2*r})
    return c
}

func (c *Circle) Paint(gc *gg.Context) {
    //log.Printf("Circle.Paint()")
    w := c.Size().X
    gc.DrawCircle(0.5*w, 0.5*w, 0.5*w)
    gc.SetStrokeWidth(c.Prop.Size(BorderWidth))
    if c.Pressed {
        gc.SetFillColor(c.Prop.Color(PressedColor))
        gc.SetStrokeColor(c.Prop.Color(PressedBorderColor))
    } else {
        gc.SetFillColor(c.Prop.Color(Color))
        gc.SetStrokeColor(c.Prop.Color(BorderColor))
    }
    gc.FillStroke()
}

func (c *Circle) Contains(pt geom.Point) (bool) {
    if !c.Embed.Contains(pt) {
        return false
    }
    x, y, w, h := c.Rect().AsCoord()
    rx := 2*(pt.X-x)/w - 1.0
    ry := 2*(pt.Y-y)/h - 1.0
    return rx*rx + ry*ry <= 1.0
}

func (c *Circle) SetCenter(p geom.Point) {
    c.SetPos(p.Sub(c.Size().Mul(0.5)))
}

func (c *Circle) Center() (geom.Point) {
    return c.Rect().Center()
}

func (c *Circle) SetRadius(r float64) {
    mp := c.Center()
    c.SetMinSize(geom.Point{2*r, 2*r})
    c.SetPos(mp.Sub(c.Size().Mul(0.5)))
}

func (c *Circle) Radius() (float64) {
    return 0.5 * c.Size().X
}

// Ein allgemeinerer Widget Typ ist die Ellipse.
type Ellipse struct {
    LeafEmbed
    Pressed bool
}

func NewEllipse(rx, ry float64) (*Ellipse) {
    e := &Ellipse{}
    e.Wrapper = e
    e.Init(GeomShapeProps)
    e.SetMinSize(geom.Point{2*rx, 2*ry})
    return e
}

func (e *Ellipse) Paint(gc *gg.Context) {
    //log.Printf("Circle.Paint()")
    w, h := e.Size().AsCoord()
    gc.DrawEllipse(0.5*w, 0.5*h, 0.5*w, 0.5*h)
    gc.SetStrokeWidth(e.Prop.Size(BorderWidth))
    if e.Pressed {
        gc.SetFillColor(e.Prop.Color(PressedColor))
        gc.SetStrokeColor(e.Prop.Color(PressedBorderColor))
    } else {
        gc.SetFillColor(e.Prop.Color(Color))
        gc.SetStrokeColor(e.Prop.Color(BorderColor))
    }
    gc.FillStroke()
}

func (e *Ellipse) Contains(pt geom.Point) (bool) {
    if !e.Embed.Contains(pt) {
        return false
    }
    x, y, w, h := e.Rect().AsCoord()
    rx := 2*(pt.X-x)/w - 1.0
    ry := 2*(pt.Y-y)/h - 1.0
    return rx*rx + ry*ry <= 1.0
}

func (e *Ellipse) SetCenter(p geom.Point) {
    e.SetPos(p.Sub(e.Size().Mul(0.5)))
}

func (e *Ellipse) Center() (geom.Point) {
    return e.Rect().Center()
}

func (e *Ellipse) SetRadius(rx, ry float64) {
    mp := e.Center()
    e.SetMinSize(geom.Point{2*rx, 2*ry})
    e.SetPos(mp.Sub(e.Size().Mul(0.5)))
}

func (e *Ellipse) Radius() (float64, float64) {
    return e.Size().Mul(0.5).AsCoord()
}

// Und wo es Kreise gibt, da sind auch die Rechtecke nicht weit.
type Rectangle struct {
    LeafEmbed
    Pressed bool
}

func NewRectangle(w, h float64) (*Rectangle) {
    r := &Rectangle{}
    r.Wrapper = r
    r.Init(GeomShapeProps)
    r.SetMinSize(geom.Point{w, h})
    return r
}

func (r *Rectangle) Paint(gc *gg.Context) {
    //log.Printf("Rectangle.Paint()")
    gc.DrawRectangle(r.Bounds().AsCoord())
    gc.SetStrokeWidth(r.Prop.Size(BorderWidth))
    if r.Pressed {
        gc.SetFillColor(r.Prop.Color(PressedColor))
        gc.SetStrokeColor(r.Prop.Color(PressedBorderColor))
    } else {
        gc.SetFillColor(r.Prop.Color(Color))
        gc.SetStrokeColor(r.Prop.Color(BorderColor))
    }
    gc.FillStroke()
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
//    s.Init(DefProps)
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

type Canvas struct {
    LeafEmbed
    Color, StrokeColor color.Color
    LineWidth float64
    Clip bool
    ObjList *list.List
    mutex *sync.Mutex
}

func NewCanvas(w, h float64) (*Canvas) {
    c := &Canvas{}
    c.Wrapper = c
    c.Init(DefProps)
    c.SetSize(geom.Point{w, h})
    c.Color   = DefProps.Color(BlackColor)
    c.StrokeColor = DefProps.Color(WhiteColor)
    c.LineWidth   = 0.0
    c.Clip        = false
    c.ObjList     = list.New()
    c.mutex       = &sync.Mutex{}
    return c
}

/*
func (c *Canvas) Paint(gc *gg.Context) {
    //log.Printf("Canvas.Paint()")
    c.Marks.UnmarkNeedsPaint()
    transformer := draw.BiLinear
    m := c.Matrix()
    s2d := f64.Aff3{m.M11, m.M12, m.M13, m.M21, m.M22, m.M23}
    c.Lock()
    if gc.Mask() == nil {
        transformer.Transform(gc.Image().(*image.RGBA), s2d, c.Gc.Image(),
                c.Gc.Image().Bounds(), draw.Over, nil)
    } else {
        transformer.Transform(gc.Image().(*image.RGBA), s2d, c.Gc.Image(),
                c.Gc.Image().Bounds(), draw.Over, &draw.Options{
                    DstMask : gc.Mask(),
                    DstMaskP: image.Point{},
                })
    }
    c.Unlock()
}
*/

func (c *Canvas) Paint(gc *gg.Context) {
    c.Marks.UnmarkNeedsPaint()
    gc.Push()
    gc.Multiply(c.Matrix())
    gc.SetFillColor(c.Color)
    gc.SetStrokeColor(c.StrokeColor)
    gc.SetStrokeWidth(c.LineWidth)
    gc.DrawRectangle(c.LocalBounds().AsCoord())
    if c.Clip {
        gc.ClipPreserve()
    }
    gc.FillStroke()
    for e := c.ObjList.Front(); e != nil; e = e.Next() {
        o := e.Value.(CanvasObject)
        o.Paint(gc)
    }
    if c.Clip {
        gc.ResetClip()
    }
    gc.Pop()
}

func (c *Canvas) Add(obj CanvasObject) {
    c.ObjList.PushBack(obj)
}

func (c *Canvas) OnInputEvent(evt touch.Event) {
    c.CallTouchFunc(evt)
}

