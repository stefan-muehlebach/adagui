
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
    //"fmt"
    "image"
    //"image/color"
    //"image/draw"
    //"image/png"
    "log"
    "math"
    "os"
    "sync"
    "time"
    //"mju.net/adatft"
    "mju.net/adagui/binding"
    "mju.net/adagui/touch"
    "mju.net/geom"
    "mju.net/gg"
    "mju.net/utils"
    //"github.com/golang/freetype/truetype"
    //"golang.org/x/image/draw"
    "golang.org/x/image/font"
    //"golang.org/x/image/font/opentype"
    //"golang.org/x/image/font/gofont/goregular"
    //"golang.org/x/image/font/gofont/gobold"
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
    g.Init()
    return g
}

func (g *Group) Paint(gc *gg.Context) {
    //log.Printf("Group.Paint()")
    g.ContainerEmbed.Paint(gc)
}

// Ein Panel ist eine komplexere Variante eines Containers. Er kann eine
// Hintergrundfarbe haben und ordnet seine Kinder gem. ihren Koordinaten an. 
type Panel struct {
    ContainerEmbed
    FillColor, StrokeColor utils.Color
    LineWidth float64
    Clip bool
    realSize, viewPort geom.Point
}

func NewPanel(w, h float64) (*Panel) {
    g := &Panel{}
    g.Wrapper = g
    g.Init()
    g.SetMinSize(geom.Point{w, h})
    g.FillColor   = pr.Color(BlackColor)
    g.StrokeColor = pr.Color(BlackColor)
    g.LineWidth   = pr.Size(PanelLineWidth)
    g.Clip        = true
    g.realSize    = geom.Point{0, 0}
    g.viewPort    = geom.Point{0, 0}
    return g
}

func (g *Panel) Paint(gc *gg.Context) {
    //log.Printf("Panel.Paint()")
    gc.SetFillColor(g.FillColor)
    gc.SetStrokeColor(g.StrokeColor)
    gc.SetLineWidth(g.LineWidth)
    gc.DrawRectangle(g.Bounds().AsCoord())
    if g.Clip {
        gc.ClipPreserve()
    }
    gc.FillStroke()
    if g.realSize.X != 0.0 && g.realSize.Y != 0.0 {
        refPt := g.Size().Sub(g.realSize)
        refPt.X *= g.viewPort.X
        refPt.Y *= g.viewPort.Y
        gc.Translate(refPt.AsCoord())
    }
    g.ContainerEmbed.Paint(gc)
    if g.Clip {
        gc.ResetClip()
    }
}

func (p *Panel) VisibleRange() (geom.Point) {
    if p.realSize.X == 0.0 && p.realSize.Y == 0.0 {
        return geom.Point{1, 1}
    }
    vis := p.Wrapper.Size()
    vis.X /= p.realSize.X
    vis.Y /= p.realSize.Y
    return vis
}

func (p *Panel) SetXView(vx float64) {
    p.viewPort.X = vx
}

func (p *Panel) SetYView(vy float64) {
    p.viewPort.Y = vy
}

func (p *Panel) ViewPort() (geom.Point) {
    return p.viewPort
}

func (p *Panel) SetRealSize(sz geom.Point) {
    p.realSize = sz
}

func (p *Panel) RealSize() (geom.Point) {
    return p.realSize
}

// Fuer die visuelle Abgrenzung in Box-Layouts.
type Separator struct {
    LeafEmbed
    StrokeColor utils.Color
    LineWidth float64
    orient Orientation
}

func NewSeparator(orient Orientation) (*Separator) {
    s := &Separator{}
    s.Wrapper = s
    s.Init()
    s.StrokeColor = pr.Color(GrayColor)
    s.LineWidth   = 4.0
    s.orient = orient
    s.SetMinSize(geom.Point{s.LineWidth, s.LineWidth})
    return s
}

func (s *Separator) Paint(gc *gg.Context) {
    gc.SetStrokeColor(s.StrokeColor)
    gc.SetLineWidth(s.LineWidth)
    gc.MoveTo(s.Bounds().W().AsCoord())
    gc.LineTo(s.Bounds().E().AsCoord())
    gc.Stroke()
}

// Unter einem Label verstehen wird einfach eine Konserve fuer Text. Kurzen
// Text! Fuer die Darstellung von groesseren Textmengen, bitte Widget Text
// beruecksichtigen.
type AlignType int

const (
    AlignLeft AlignType = (1 << 0)
    AlignCenter         = (1 << 1)
    AlignRight          = (1 << 2)
    AlignTop            = (1 << 3)
    AlignMiddle         = (1 << 4)
    AlignBottom         = (1 << 5)
)

type Label struct {
    LeafEmbed
    text binding.String
    fontFace font.Face
    TextColor utils.Color
    align AlignType
    rPt geom.Point
    desc float64
}

func newLabel() (*Label) {
    l := &Label{}
    l.Wrapper = l
    l.Init()
    l.fontFace  = pr.Font(RegularFont)
    l.TextColor = pr.Color(WhiteColor)
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

func (l *Label) SetAlign(a AlignType) {
    l.align = a
    l.updateRefPoint()
}

func (l *Label) Align() (AlignType) {
    return l.align
}

func (l *Label) SetText(str string) {
    l.text.Set(str)
    l.updateSize()
}

func (l *Label) Text() (string) {
    return l.text.Get()
}

func (l *Label) SetFontFace(fontFace font.Face) {
    l.fontFace = fontFace
    l.updateSize()
}

func (l *Label) FontFace() (font.Face) {
    return l.fontFace
}

func (l *Label) updateSize() {
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
    gc.SetStrokeColor(l.TextColor)
    gc.DrawString(l.text.Get(), l.rPt.X, l.rPt.Y)
    // Groesse des Labels als graues Rechteck
    //gc.SetStrokeColor(utils.Lightgray)
    //gc.SetLineWidth(1.0)
    //gc.DrawRectangle(l.Bounds().AsCoord())
    //gc.Stroke()
    // Referenzpunkt fuer den Text
    //gc.SetFillColor(utils.Lightgray)
    //gc.DrawPoint(l.rPt.X, l.rPt.Y, 5.0)
    //gc.Fill()
}

// Button, TextButton und IconButton funktionieren alle ungefaehr nach dem
// gleichen Prinzip - mit dem Unterschied, dass bei Button als Inhalt bloss
// eine Farbe, bei TextButton eben Text und bei IconButton ein ganzes
// image.Image als Inhalt verwendet werden kann.
/*
var (
    btnHeight          = 32.0
    btnInsetHorizontal = 15.0
    btnInset           =  5.0
    btnLineWidth       =  0.0
    btnRectRoundRad    =  6.0

    radBtnSize         = 20.0
    radBtnLineWidth    =  4.0
    radBtnDotSize      =  8.0

    chkBoxSize         = 20.0
    chkBoxLineWidth    =  3.0
    chkBoxRoundRectRad =  6.0
)
*/

// ----------------------------------------------------------------------------

// Buttons sind neutrale Knoepfe, ohne spezifischen Inhalt, d.h. ohne Text
// oder Icons. Sie werden selten direkt verwendet, sondern dienen als
// generische Grundlage fuer die weiter unten definierten Text- oder Icon-
// Buttons.
type Button struct {
    LeafEmbed
    StrokeColor, FillColor utils.Color
    LineWidth float64
    pushed bool
}

func NewButton(w, h float64) (*Button) {
    b := &Button{}
    b.Wrapper = b
    b.Init()
    b.SetMinSize(geom.Point{w, h})
    b.FillColor   = pr.Color(FillColor)
    b.StrokeColor = pr.Color(StrokeColor)
    b.LineWidth   = pr.Size(LineWidth)
    b.pushed      = false
    return b
}

func (b *Button) Paint(gc *gg.Context) {
    //log.Printf("Button.Paint()")
    gc.DrawRoundedRectangle(0.0, 0.0, b.Size().X, b.Size().Y,
            pr.Size(ButtonCornerRad))
    if b.pushed {
        gc.SetStrokeColor(b.StrokeColor.Bright())
        gc.SetFillColor(b.FillColor.Bright())
    } else {
        gc.SetStrokeColor(b.StrokeColor)
        gc.SetFillColor(b.FillColor)
    }
    gc.SetLineWidth(b.LineWidth)
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
    TextColor utils.Color
    align AlignType
    rPt geom.Point
    desc float64
}

func NewTextButton(label string) (*TextButton) {
    b := &TextButton{}
    b.Wrapper     = b
    b.Init()
    b.label       = label
    b.fontFace    = pr.Font(BoldFont)
    b.FillColor   = pr.Color(FillColor)
    b.StrokeColor = pr.Color(StrokeColor)
    b.LineWidth   = pr.Size(LineWidth)
    b.TextColor   = pr.Color(TextColor)
    b.align       = AlignCenter | AlignMiddle
    b.updateSize()
    return b
}

func (b *TextButton) SetSize(size geom.Point) {
    b.Button.SetSize(size)
    b.updateSize()
}

func (b *TextButton) updateSize() {
    w := float64(font.MeasureString(b.fontFace, b.label)) / 64.0
    h := pr.Size(ButtonHeight)
    //h := float64(b.fontFace.Metrics().Ascent +
    //        b.fontFace.Metrics().Descent) / 64.0
    b.desc = float64(b.fontFace.Metrics().Descent) / 64.0
    b.SetMinSize(geom.Point{w+2*pr.Size(ButtonInsetHorizontal), h})
    b.updateRefPoint()
}

func (b *TextButton) updateRefPoint() {
    b.rPt = b.Bounds().Center()
}

func (b *TextButton) Paint(gc *gg.Context) {
    //log.Printf("TextButton.Paint()")
    b.Button.Paint(gc)
    gc.SetFontFace(b.fontFace)
    gc.SetStrokeColor(b.TextColor)
    gc.DrawStringAnchored(b.label, b.rPt.X, b.rPt.Y, 0.5, 0.5)
}

// Der Versuch, ein ListButton zu implementieren...
type ListButton struct {
    Button
    fontFace font.Face
    TextColor utils.Color
    Options []string
    selIdx int
    Selected string
}

func NewListButton(options []string) (*ListButton) {
    b := &ListButton{}
    //b.TextButton  = *NewTextButton(label)
    b.Wrapper     = b
    b.Init()
    b.FillColor   = pr.Color(BlackColor)
    b.StrokeColor = pr.Color(StrokeColor)
    b.LineWidth   = pr.Size(LineWidth)
    b.fontFace    = pr.Font(RegularFont)
    b.TextColor   = pr.Color(TextColor)
    b.Options     = options
    b.selIdx      = 0
    b.Selected    = b.Options[b.selIdx]
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
    w := maxWidth/64.0 + 2.0*pr.Size(ButtonInsetHorizontal) + pr.Size(ButtonHeight)
    h := pr.Size(ButtonHeight)
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
    gc.SetStrokeColor(b.TextColor)
    pt := geom.Point{0.6*b.Size().Y+2*pr.Size(ButtonInset), 0.5*b.Size().Y}
    gc.DrawStringAnchored(b.Selected, pt.X, pt.Y, 0.0, 0.5)

    if b.pushed {
        gc.SetStrokeColor(b.StrokeColor.Bright())
    } else {
        gc.SetStrokeColor(b.StrokeColor)
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

    gc.SetFillColor(pr.Color(WhiteColor))
    gc.SetStrokeColor(pr.Color(WhiteColor))
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
    checked bool
}

func NewIconButton(imgFile string) (*IconButton) {
    b := &IconButton{}
    b.Wrapper = b
    b.Init()
    b.pushed = false
    b.img, _ = gg.LoadPNG(imgFile)
    i := pr.Size(ButtonInset)
    rect := geom.NewRectangleIMG(b.img.Bounds()).Inset(-i, -i)
    b.SetMinSize(rect.Size())
    b.FillColor   = pr.Color(FillColor)
    b.StrokeColor = pr.Color(StrokeColor)
    b.LineWidth   = pr.Size(LineWidth)
    return b
}

func (b *IconButton) OnInputEvent(evt touch.Event) {
    b.Button.OnInputEvent(evt)
    if evt.Type == touch.TypeTap {
        b.checked = !b.checked
        if b.checked {
            b.FillColor = pr.Color(FillColor).Dark()
            b.StrokeColor = pr.Color(StrokeColor).Bright()
        } else {
            b.FillColor = pr.Color(FillColor)
            b.StrokeColor = pr.Color(StrokeColor)
        }
        b.Mark(MarkNeedsPaint)
    }
}

func (b *IconButton) Paint(gc *gg.Context) {
    //log.Printf("IconButton.Paint()")
    b.Button.Paint(gc)
    cp := b.Bounds().Center()
    gc.DrawImageAnchored(b.img, int(cp.X), int(cp.Y), 0.5, 0.5)
}

// Checkboxen verhalten sich sehr aehnlich zu RadioButtons, sind jeoch eigen-
// staendig und nicht Teil einer Gruppe.
type Checkbox struct {
    Button
    label string
    fontFace font.Face
    TextColor utils.Color
    value binding.Bool
}

func NewCheckbox(label string) (*Checkbox) {
    c := &Checkbox{}
    c.Wrapper = c
    c.Init()
    c.FillColor   = pr.Color(FillColor)
    c.StrokeColor = pr.Color(StrokeColor)
    c.LineWidth   = pr.Size(LineWidth)
    c.label       = label
    c.fontFace    = pr.Font(RegularFont)
    c.TextColor   = pr.Color(TextColor)
    w := float64(font.MeasureString(c.fontFace, label))/64.0
    c.SetMinSize(geom.Point{pr.Size(CheckboxSize)+pr.Size(ButtonInset)+w, pr.Size(CheckboxSize)})
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
    gc.DrawRoundedRectangle(0.0, 0.0, pr.Size(CheckboxSize),
            pr.Size(CheckboxSize), pr.Size(CheckboxCornerRad))
    if c.pushed {
        gc.SetStrokeColor(c.StrokeColor.Bright())
        gc.SetFillColor(c.FillColor.Bright())
    } else {
        gc.SetStrokeColor(c.StrokeColor)
        gc.SetFillColor(c.FillColor)
    }
    gc.SetLineWidth(c.LineWidth)
    gc.FillStroke()
    if c.Checked() {
        gc.SetLineWidth(pr.Size(CheckboxLineWidth))
        gc.SetStrokeColor(pr.Color(WhiteColor))
        gc.MoveTo(4, 9)
        gc.LineTo(8, 14)
        gc.LineTo(14, 5)
        gc.Stroke()
    }
    x := pr.Size(CheckboxSize) + pr.Size(ButtonInset)
    y := 0.5*c.Size().Y
    gc.SetStrokeColor(c.TextColor)
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
type RadioButton struct {
    Button
    label string
    fontFace font.Face
    TextColor utils.Color
    checked bool
    value int
    data binding.Int
}

func NewRadioButtonWithData(label string, value int, data binding.Int) (*RadioButton) {
    b := &RadioButton{}
    b.Wrapper = b
    b.Init()
    b.FillColor   = pr.Color(FillColor)
    b.StrokeColor = pr.Color(StrokeColor)
    b.LineWidth   = pr.Size(LineWidth)
    b.label       = label
    b.fontFace    = pr.Font(RegularFont)
    b.TextColor   = pr.Color(TextColor)
    w := float64(font.MeasureString(b.fontFace, label))/64.0
    b.SetMinSize(geom.Point{pr.Size(RadioBtnSize)+pr.Size(ButtonInset)+w, pr.Size(RadioBtnSize)})
    b.value = value
    b.data = data
    b.data.AddListener(b)
    return b
}

func (b *RadioButton) Paint(gc *gg.Context) {
    //log.Printf("RadioButton.Paint()")
    mp := geom.Point{0.5*pr.Size(RadioBtnSize), 0.5*pr.Size(RadioBtnSize)}
    gc.DrawCircle(mp.X, mp.Y, 0.5*pr.Size(RadioBtnSize))
    if b.pushed {
        gc.SetStrokeColor(b.StrokeColor.Bright())
        gc.SetFillColor(b.FillColor.Bright())
    } else {
        gc.SetStrokeColor(b.StrokeColor)
        gc.SetFillColor(b.FillColor)
    }
    gc.SetLineWidth(b.LineWidth)
    gc.FillStroke()
    if b.checked {
        gc.DrawCircle(mp.X, mp.Y, 0.5*pr.Size(RadioBtnDotSize))
	gc.SetFillColor(pr.Color(WhiteColor))
	gc.Fill()
    }
    x := pr.Size(RadioBtnSize)+pr.Size(ButtonInset)
    y := 0.5*b.Size().Y
    gc.SetStrokeColor(b.TextColor)
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
/*
var (
    scrWidth     = chkBoxSize
    scrBarWidth  = scrWidth - 4.0
    scrCtrlWidth = scrWidth
)
*/

type Scrollbar struct {
    LeafEmbed
    BarColor, CtrlColor utils.Color
    orient Orientation
    initValue, visiRange float64
    //len float64
    pushed bool
    value binding.Float
    dp1, dp2, startPt, endPt1, endPt2 geom.Point
}

func NewScrollbar(len float64, orient Orientation) (*Scrollbar) {
    s := &Scrollbar{}
    s.Wrapper = s
    s.Init()
    s.orient = orient
    if s.orient == Horizontal {
        s.SetMinSize(geom.Point{len, pr.Size(ScrollWidth)})
        s.dp1 = geom.Point{0.5*pr.Size(ScrollBarWidth), 0.5*pr.Size(ScrollWidth)}
        s.dp2 = geom.Point{0.5*pr.Size(ScrollCtrlWidth), 0.5*pr.Size(ScrollWidth)}
    } else {
        s.SetMinSize(geom.Point{pr.Size(ScrollWidth), len})
        s.dp1 = geom.Point{0.5*pr.Size(ScrollWidth), 0.5*pr.Size(ScrollBarWidth)}
        s.dp2 = geom.Point{0.5*pr.Size(ScrollWidth), 0.5*pr.Size(ScrollCtrlWidth)}
    }
    s.initValue = 0.0
    s.visiRange = 0.1
    s.value     = binding.NewFloat()
    s.BarColor  = pr.Color(WhiteColor)
    s.BarColor.A = 150
    s.CtrlColor = pr.Color(FillColor)
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
    //s.len = max(s.Size().X, s.Size().Y)
    s.endPt1 = s.Size().Sub(s.dp1)
    s.endPt2 = s.Size().Sub(s.dp2)
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
    gc.DrawLine(s.dp1.X, s.dp1.Y, s.endPt1.X, s.endPt1.Y)
    gc.SetLineWidth(pr.Size(ScrollBarWidth))
    gc.SetStrokeColor(s.BarColor)
    gc.Stroke()

    newVal     := 0.5*s.visiRange + s.Value()*(1.0-s.visiRange)
    startValue := newVal - 0.5*s.visiRange
    endValue   := newVal + 0.5*s.visiRange

    r := s.Bounds().Inset(s.dp2.X, s.dp2.Y)
    if s.orient == Horizontal {
        pt1 = r.RelPos(startValue, 0.0)
        pt2 = r.RelPos(endValue, 0.0)
    } else {
        pt1 = r.RelPos(0.0, startValue)
        pt2 = r.RelPos(0.0, endValue)
    }
    if s.pushed {
        gc.SetStrokeColor(s.CtrlColor.Bright())
    } else {
        gc.SetStrokeColor(s.CtrlColor)
    }
    gc.SetLineWidth(pr.Size(ScrollCtrlWidth))
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
        r := s.Rect().Inset(s.dp2.X, s.dp2.Y)
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
/*
var (
    pr.Size(SliderWidth)       = chkBoxSize
    pr.Size(SliderBarWidth)    = pr.Size(SliderWidth) - 4.0
    pr.Size(SliderCtrlWidth)   = pr.Size(SliderWidth)
)
*/

type Slider struct {
    LeafEmbed
    BarColor, CtrlColor utils.Color
    orient Orientation
    initValue, minValue, maxValue, stepSize float64
    pushed bool
    value binding.Float
    dp1, dp2, endPt1, endPt2 geom.Point
}

func NewSlider(len float64, orient Orientation) (*Slider) {
    s := &Slider{}
    s.Wrapper = s
    s.Init()
    s.orient = orient
    if s.orient == Horizontal {
        s.SetMinSize(geom.Point{len, pr.Size(SliderWidth)})
        s.dp1 = geom.Point{0.5*pr.Size(SliderBarWidth), 0.5*pr.Size(SliderWidth)}
        s.dp2 = geom.Point{0.5*pr.Size(SliderCtrlWidth), 0.5*pr.Size(SliderWidth)}
    } else {
        s.SetMinSize(geom.Point{pr.Size(SliderWidth), len})
        s.dp1 = geom.Point{0.5*pr.Size(SliderWidth), 0.5*pr.Size(SliderBarWidth)}
        s.dp2 = geom.Point{0.5*pr.Size(SliderWidth), 0.5*pr.Size(SliderCtrlWidth)}
    }
    s.initValue = 0.0
    s.minValue  = 0.0
    s.maxValue  = 1.0
    s.stepSize  = 0.1
    s.value     = binding.NewFloat()
    s.BarColor  = pr.Color(WhiteColor)
    s.BarColor.A = 150
    s.CtrlColor = pr.Color(FillColor)
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
    s.endPt1 = s.Size().Sub(s.dp1)
    s.endPt2 = s.Size().Sub(s.dp2)
}

func (s *Slider) Paint(gc *gg.Context) {
    var pt1 geom.Point
    //log.Printf("Slider.Paint()")
    gc.DrawLine(s.dp1.X, s.dp1.Y, s.endPt1.X, s.endPt1.Y)
    gc.SetLineWidth(pr.Size(SliderBarWidth))
    gc.SetStrokeColor(s.BarColor)
    gc.Stroke()

    r := s.Bounds().Inset(s.dp2.X, s.dp2.Y)
    if s.orient == Horizontal {
        pt1 = r.RelPos(s.Factor(), 0.0)
    } else {
        pt1 = r.RelPos(0.0, 1.0-s.Factor())
    }
    if s.pushed {
        gc.SetStrokeColor(s.CtrlColor.Bright())
        gc.SetFillColor(s.CtrlColor.Bright())
    } else {
        gc.SetStrokeColor(s.CtrlColor)
        gc.SetFillColor(s.CtrlColor)
    }
    gc.SetLineWidth(pr.Size(LineWidth))
    gc.DrawCircle(pt1.X, pt1.Y, 0.5*pr.Size(SliderCtrlWidth))
    gc.FillStroke()
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
        r := s.Rect().Inset(s.dp2.X, s.dp2.Y)
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
type Circle struct {
    LeafEmbed
    StrokeColor, FillColor utils.Color
    LineWidth float64
}

func NewCircle(r float64) (*Circle) {
    c := &Circle{}
    c.Wrapper = c
    c.Init()
    c.SetMinSize(geom.Point{2*r, 2*r})
    c.StrokeColor = pr.Color(WhiteColor)
    c.FillColor   = pr.Color(FillColor)
    c.LineWidth   = 2.0
    return c
}

func (c *Circle) Paint(gc *gg.Context) {
    //log.Printf("Circle.Paint()")
    //gc.Push()
    //gc.Multiply(gg.Matrix(c.Matrix()))
    w := c.Size().X
    gc.DrawCircle(0.5*w, 0.5*w, 0.5*w)
    gc.SetLineWidth(c.LineWidth)
    gc.SetFillColor(c.FillColor)
    gc.SetStrokeColor(c.StrokeColor)
    gc.FillStroke()
    //gc.Pop()
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
    StrokeColor, FillColor utils.Color
    LineWidth float64
}

func NewEllipse(rx, ry float64) (*Ellipse) {
    e := &Ellipse{}
    e.Wrapper = e
    e.Init()
    e.SetSize(geom.Point{2*rx, 2*ry})
    e.StrokeColor = pr.Color(WhiteColor)
    e.FillColor   = pr.Color(FillColor)
    e.LineWidth   = 2.0
    return e
}

func (e *Ellipse) Paint(gc *gg.Context) {
    //log.Printf("Circle.Paint()")
    //gc.Push()
    //gc.Multiply(gg.Matrix(e.Matrix()))
    w, h := e.Size().AsCoord()
    gc.DrawEllipse(0.5*w, 0.5*h, 0.5*w, 0.5*h)
    gc.SetLineWidth(e.LineWidth)
    gc.SetFillColor(e.FillColor)
    gc.SetStrokeColor(e.StrokeColor)
    gc.FillStroke()
    //gc.Pop()
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
    StrokeColor, FillColor utils.Color
    LineWidth float64
}

func NewRectangle(w, h float64) (*Rectangle) {
    r := &Rectangle{}
    r.Wrapper = r
    r.Init()
    r.SetMinSize(geom.Point{w, h})
    r.StrokeColor = pr.Color(WhiteColor)
    r.FillColor   = pr.Color(FillColor)
    r.LineWidth   = 2.0
    return r
}

func (r *Rectangle) Paint(gc *gg.Context) {
    //log.Printf("Rectangle.Paint()")
    //gc.Push()
    //gc.Multiply(gg.Matrix(r.Matrix()))
    gc.DrawRectangle(r.Bounds().AsCoord())
    gc.SetLineWidth(r.LineWidth)
    gc.SetFillColor(r.FillColor)
    gc.SetStrokeColor(r.StrokeColor)
    gc.FillStroke()
    //gc.Pop()
}

// Wir wollen es wissen und machen einen auf Game-Entwickler
type Sprite struct {
    LeafEmbed
    imgList []image.Image
    curImg int
    ticker *time.Ticker
}

func NewSprite(imgFiles ...string) (*Sprite) {
    s := &Sprite{}
    s.Wrapper = s
    s.Init()
    s.imgList = make([]image.Image, 0)
    s.curImg = 0
    s.AddImages(imgFiles...)
    pt := geom.NewPointIMG(s.imgList[0].Bounds().Size())
    s.SetSize(pt)
    return s
}

func (s *Sprite) AddImages(imgFiles ...string) {
    for _, fileName := range imgFiles {
        fh, err := os.Open(fileName)
        utils.Check(err)
        img, _, err := image.Decode(fh)
        utils.Check(err)
        s.imgList = append(s.imgList, img)
        fh.Close()
    }
}

func (s *Sprite) Paint(gc *gg.Context) {
    s.Marks.UnmarkNeedsPaint()
    gc.Push()
    gc.Multiply(gg.Matrix(s.Matrix()))
    gc.DrawImage(s.imgList[s.curImg], int(s.Rect().Min.X), int(s.Rect().Min.Y))
    gc.Pop()
}

func (s *Sprite) StartAnim(dt time.Duration) {
    s.ticker = time.NewTicker(dt)
    go func() {
        for {
            <- s.ticker.C
            s.curImg = (s.curImg + 1) % len(s.imgList)
            s.Mark(MarkNeedsPaint)
            s.Win.Repaint()
        }
    }()
}

func (s *Sprite) StopAnim() {
    s.ticker.Stop()
}

//-----------------------------------------------------------------------------

type Canvas struct {
    LeafEmbed
    FillColor, StrokeColor utils.Color
    LineWidth float64
    Clip bool
    ObjList *list.List
    mutex *sync.Mutex
}

func NewCanvas(w, h float64) (*Canvas) {
    c := &Canvas{}
    c.Wrapper = c
    c.Init()
    c.SetSize(geom.Point{w, h})
    c.FillColor   = pr.Color(BlackColor)
    c.StrokeColor = pr.Color(WhiteColor)
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
    gc.Multiply(gg.Matrix(c.Matrix()))
    gc.SetFillColor(c.FillColor)
    gc.SetStrokeColor(c.StrokeColor)
    gc.SetLineWidth(c.LineWidth)
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

//-----------------------------------------------------------------------------

// GUI daten sowohl fuer PageButton als auch Drawer
var (
    flapWidth            = 22.0
    flapHeight           = 60.0
    flapInset            =  9.0
    flapRectRad          =  6.0
    flapFillColor        = utils.Lightgray.SetAlpha(0.6)
    flapPushedFillColor  = pr.Color(WhiteColor)
    flapArrowColor       = pr.Color(FillColor)
    flapPushedArrowColor = pr.Color(WhiteColor)
    flapArrowWidth       =  8.0
    flapSize = []geom.Point{
	    geom.Point{flapWidth, flapHeight},
	    geom.Point{flapHeight, flapWidth},
	    geom.Point{flapWidth, flapHeight},
	    geom.Point{flapHeight, flapWidth},
    }
    flapRectInsets = [][]geom.Point{
        {
            geom.Point{-flapRectRad, 0.0},
            geom.Point{},
        },
        {
            geom.Point{0.0, -flapRectRad},
            geom.Point{},
        },
        {
            geom.Point{},
            geom.Point{flapRectRad, 0.0},
        },
        {
            geom.Point{},
            geom.Point{0.0, flapRectRad},
        },
    }
    pgBtnFillColor = flapFillColor
    drwFillColor   = utils.Lightgray.SetAlpha(0.6)
)

func DrawArrow(gc *gg.Context, dst geom.Rectangle, pos Border) {
    switch pos {
    case Left:
        gc.MoveTo(dst.Max.X, dst.Min.Y)
        gc.LineTo(dst.Min.X, 0.5*(dst.Min.Y+dst.Max.Y))
        gc.LineTo(dst.Max.X, dst.Max.Y)
    case Top:
        gc.MoveTo(dst.Min.X, dst.Max.Y)
        gc.LineTo(0.5*(dst.Min.X+dst.Max.X), dst.Min.Y)
        gc.LineTo(dst.Max.X, dst.Max.Y)
    case Right:
        gc.MoveTo(dst.Min.X, dst.Min.Y)
        gc.LineTo(dst.Max.X, 0.5*(dst.Min.Y+dst.Max.Y))
        gc.LineTo(dst.Min.X, dst.Max.Y)
    case Bottom:
        gc.MoveTo(dst.Min.X, dst.Min.Y)
        gc.LineTo(0.5*(dst.Min.X+dst.Max.X), dst.Max.Y)
        gc.LineTo(dst.Max.X, dst.Min.Y)
    }
}

//-----------------------------------------------------------------------------

// PageButton dienen vorallem fuer den Wechsel zwischen den Windows, koennen
// aber auch fuer anderes verwendet werden.
type PageButton struct {
    LeafEmbed
    pos Border
    pushed bool
    ExtRect geom.Rectangle
}

func NewPageButton(pos Border) (*PageButton) {
    b := &PageButton{}
    b.Wrapper = b
    b.Init()
    b.pos = pos
    b.pushed = false
    b.SetSize(flapSize[b.pos])
    return b
}

func (b *PageButton) SetPos(pt geom.Point) {
    switch b.pos {
    case Left:
        pt.X = 0.0
        pt.Y -= 0.5*flapHeight
    case Top:
        pt.X -= 0.5*flapHeight
        pt.Y = 0.0
    case Right:
        pt.X = 298.0
        pt.Y -= 0.5*flapHeight
    case Bottom:
        pt.X -= 0.5*flapHeight
        pt.Y = 218.0
    }
    b.Wrappee().SetPos(pt)

    b.ExtRect = b.Rect().Sub(b.Rect().Min)
    b.ExtRect = geom.Rectangle{
        b.ExtRect.Min.Add(flapRectInsets[b.pos][0]),
        b.ExtRect.Max.Add(flapRectInsets[b.pos][1]),
    }
}

func (b *PageButton) Paint(gc *gg.Context) {
    //log.Printf("PageButton.Paint()")
    b.Marks.UnmarkNeedsPaint()
    gc.Push()
    gc.Translate(b.Rect().Min.AsCoord())
    gc.DrawRoundedRectangle(b.ExtRect.Min.X, b.ExtRect.Min.Y,
            b.ExtRect.Dx(), b.ExtRect.Dy(), flapRectRad)
    gc.SetFillColor(pgBtnFillColor)
    gc.Fill()

    if b.pushed {
        gc.SetStrokeColor(flapArrowColor.Bright())
    } else {
        gc.SetStrokeColor(flapArrowColor)
    }
    gc.SetLineWidth(flapArrowWidth)
    DrawArrow(gc, b.ExtRect.Inset(flapInset, flapInset), b.pos)
    gc.Stroke()
    gc.Pop()
}

func (b *PageButton) OnInputEvent(evt touch.Event) {
    //log.Printf("%T: %v", b, evt)
    switch evt.Type {
    case touch.TypePress, touch.TypeEnter:
        b.pushed = true
        b.Mark(MarkNeedsPaint)
    case touch.TypeLeave:
        b.pushed = false
        b.Mark(MarkNeedsPaint)
    case touch.TypeRelease:
        b.pushed = false
        b.Mark(MarkNeedsPaint)
    }
    b.CallTouchFunc(evt)
}

// Der Drawer (engl. Schublade) kann eine Reihe von weiteren Widgets aufnehmen
// und laesst sich bei Nichtbedarf am Rand des Bildschirms auf ein kleines
// Icon zusammenklappen. Eine Antwort auf den beschraenkten Platz des Adafruit
// TFT-Bildschirm.
var (
    drwSizeChange = [][]geom.Point{
        {
            geom.Point{},
            geom.Point{100.0, 0.0},
        },
	    {
            geom.Point{},
            geom.Point{0.0, 100.0},
        },
        {
            geom.Point{-100.0, 0.0},
            geom.Point{},
        },
        {
            geom.Point{0.0, -100.0},
            geom.Point{},
        },
    }
)

type Drawer struct {
    ContainerEmbed
    pos Border
    FillColor utils.Color
    pushed bool
    handle geom.Rectangle
    isOpen bool
    ExtRect geom.Rectangle
}

func NewDrawer(pos Border) (*Drawer) {
    d := &Drawer{}
    d.Wrapper = d
    d.Init()
    d.pos = pos
    d.FillColor = flapFillColor
    d.pushed = false
    d.isOpen = false
    d.SetSize(flapSize[d.pos])
    return d
}

func (d *Drawer) SetPos(pt geom.Point) {
    switch d.pos {
    case Left:
        pt.X = 0.0
    case Top:
        pt.Y = 0.0
    case Right:
        pt.X = 298.0
    case Bottom:
        pt.Y = 218.0
    }
    d.Wrappee().SetPos(pt)

    d.ExtRect = d.Rect().Sub(d.Rect().Min)
    d.ExtRect = geom.Rectangle{
        d.ExtRect.Min.Add(flapRectInsets[d.pos][0]),
        d.ExtRect.Max.Add(flapRectInsets[d.pos][1]),
    }
}

func (d *Drawer) Paint(gc *gg.Context) {
    //log.Printf("Drawer.Paint()")
    d.Marks.UnmarkNeedsPaint()
    gc.Push()
    gc.Translate(d.Rect().Min.AsCoord())
    gc.DrawRoundedRectangle(d.ExtRect.Min.X, d.ExtRect.Min.Y,
            d.ExtRect.Dx(), d.ExtRect.Dy(), flapRectRad)
    //log.Printf("Drawer.Paint():")
    if d.pushed {
        gc.SetFillColor(d.FillColor.Bright())
    } else {
        gc.SetFillColor(d.FillColor)
    }
    gc.Fill()

    gc.SetStrokeColor(flapArrowColor)
    gc.SetLineWidth(flapArrowWidth)
    DrawArrow(gc, d.ExtRect.Inset(flapInset, flapInset), (d.pos+2)%4)
    gc.Stroke()

    gc.DrawRectangle(d.ExtRect.AsCoord())
    gc.Clip()
    //d.ContainerEmbed.Paint(gc)
    gc.ResetClip()
    gc.Pop()
}

func (d *Drawer) OnInputEvent(evt touch.Event) {
    //log.Printf("Drawer.OnInputEvent(): %T, %v", d, evt)
    switch evt.Type {
    case touch.TypePress, touch.TypeEnter:
        d.pushed = true
        d.Mark(MarkNeedsPaint)
    case touch.TypeLeave, touch.TypeRelease:
        d.pushed = false
        d.Mark(MarkNeedsPaint)
    case touch.TypeTap:
        if d.isOpen {
            d.Close()
        } else {
            d.Open()
        }
    }
}

func (d *Drawer) IsOpen() (bool) {
    return d.isOpen
}

func (d *Drawer) Open() {
    if d.isOpen {
        return
    }
    d.isOpen = true
    d.SetPos(d.Pos().Add(drwSizeChange[d.pos][0]))
    d.SetSize(d.Size().Add(drwSizeChange[d.pos][1]))
    //d.Rect.Min = d.Rect.Min.Add(drwSizeChange[d.pos][0])
    //d.Rect.Max = d.Rect.Max.Add(drwSizeChange[d.pos][1])
    d.Mark(MarkNeedsPaint)
}


func (d *Drawer) Close() {
    if !d.isOpen {
        return
    }
    d.isOpen = false
    d.SetPos(d.Pos().Sub(drwSizeChange[d.pos][0]))
    d.SetSize(d.Size().Sub(drwSizeChange[d.pos][1]))
    //d.Rect.Min = d.Rect.Min.Sub(drwSizeChange[d.pos][0])
    //d.Rect.Max = d.Rect.Max.Sub(drwSizeChange[d.pos][1])
    d.Mark(MarkNeedsPaint)
}

