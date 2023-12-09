
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
    // "image/color"
    "log"
    "math"
    "os"
    "sync"
    "time"
    "github.com/stefan-muehlebach/adagui/binding"
    "github.com/stefan-muehlebach/adagui/touch"
    "github.com/stefan-muehlebach/gg"
    "github.com/stefan-muehlebach/gg/color"
    //"github.com/stefan-muehlebach/gg/colornames"
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
    FillColor, StrokeColor color.Color
    LineWidth float64
    Clip bool
    virtSize, sizeDiff, viewPort, refPt geom.Point
}

func NewPanel(w, h float64) (*Panel) {
    g := &Panel{}
    g.Wrapper = g
    g.Init()
    g.SetMinSize(geom.Point{w, h})
    g.FillColor   = pr.Color(BlackColor)
    g.StrokeColor = pr.Color(BlackColor)
    g.LineWidth   = pr.Size(PanelBorderSize)
    g.Clip        = true
    g.SetVirtualSize(g.Size())
    g.viewPort    = geom.Point{0, 0}
    g.refPt       = geom.Point{0, 0}
    return g
}

func (g *Panel) Paint(gc *gg.Context) {
    //log.Printf("Panel.Paint()")
    gc.SetFillColor(g.FillColor)
    gc.SetStrokeColor(g.StrokeColor)
    gc.SetStrokeWidth(g.LineWidth)
    gc.DrawRectangle(g.Bounds().AsCoord())
    if g.Clip {
        gc.ClipPreserve()
    }
    gc.FillStroke()

/*
    if g.virtSize.X != 0.0 && g.virtSize.Y != 0.0 {
        refPt := g.Size().Sub(g.virtSize)
        refPt.X *= g.viewPort.X
        refPt.Y *= g.viewPort.Y
        gc.Translate(refPt.AsCoord())
    }
*/
    //gc.Multiply(g.Matrix())

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
    StrokeColor color.Color
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
    gc.SetStrokeWidth(s.LineWidth)
    gc.MoveTo(s.Bounds().W().AsCoord())
    gc.LineTo(s.Bounds().E().AsCoord())
    gc.Stroke()
}

// Nimmt den verfuegbaren Platz (vertikal oder horizontal) in Box-Layouts
// ein.
type Spacer struct {
    LeafEmbed
    fixHorizontal, fixVertical bool
}

func NewSpacer() (*Spacer) {
    s := &Spacer{}
    s.Wrapper = s
    s.Init()
    return s
}

func (s *Spacer) ExpandHorizontal() (bool) {
    return !s.fixHorizontal
}

func (s *Spacer) ExpandVertical() (bool) {
    return !s.fixVertical
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
type Label struct {
    LeafEmbed
    text binding.String
    fontFont *opentype.Font
    fontSize float64
    fontFace font.Face
    TextColor color.Color
    align AlignType
    rPt geom.Point
    desc float64
}

func newLabel() (*Label) {
    l := &Label{}
    l.Wrapper = l
    l.Init()
    l.fontFont  = pr.Font(RegularFont)
    l.fontSize  = pr.Size(TextSize)
    l.TextColor = pr.Color(TextColor)
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

func (l *Label) SetFont(fontFont *opentype.Font) {
    l.fontFont = fontFont
    l.updateSize()
}

func (l *Label) Font() (*opentype.Font) {
    return l.fontFont
}

func (l *Label) SetFontSize(fontSize float64) {
    l.fontSize = fontSize
    l.updateSize()
}

func (l *Label) FontSize() (float64) {
    return l.fontSize
}

func (l *Label) updateSize() {
    l.fontFace = fonts.NewFace(l.fontFont, l.fontSize)
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
    //gc.SetStrokeWidth(1.0)
    //gc.DrawRectangle(l.Bounds().AsCoord())
    //gc.Stroke()
    // Referenzpunkt fuer den Text
    //gc.SetFillColor(colornames.Lightgray)
    //gc.DrawPoint(l.rPt.X, l.rPt.Y, 5.0)
    //gc.Fill()
}

// Buttons sind neutrale Knoepfe, ohne spezifischen Inhalt, d.h. ohne Text
// oder Icons. Sie werden selten direkt verwendet, sondern dienen als
// generische Grundlage fuer die weiter unten definierten Text- oder Icon-
// Buttons.
type Button struct {
    LeafEmbed
    FillColor, FillFocusColor, BorderColor, BorderFocusColor color.Color
    LineWidth float64
    pushed bool
}

func NewButton(w, h float64) (*Button) {
    b := &Button{}
    b.Wrapper = b
    b.Init()
    b.SetMinSize(geom.Point{w, h})
    b.FillColor        = pr.Color(ButtonColor)
    b.FillFocusColor   = pr.Color(ButtonFocusColor)
    b.BorderColor      = pr.Color(ButtonBorderColor)
    b.BorderFocusColor = pr.Color(ButtonBorderFocusColor)
    b.LineWidth        = pr.Size(ButtonBorderSize)
    b.pushed           = false
    return b
}

func (b *Button) Paint(gc *gg.Context) {
    //log.Printf("Button.Paint()")
    gc.DrawRoundedRectangle(0.0, 0.0, b.Size().X, b.Size().Y,
            pr.Size(ButtonCornerRad))
    if b.pushed {
        gc.SetFillColor(b.FillFocusColor)
        gc.SetStrokeColor(b.BorderFocusColor)
    } else {
        gc.SetFillColor(b.FillColor)
        gc.SetStrokeColor(b.BorderColor)
    }
    gc.SetStrokeWidth(b.LineWidth)
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
    fontSize float64
    fontFace font.Face
    TextColor color.Color
    //align AlignType
    rPt geom.Point
    desc float64
}

func NewTextButton(label string) (*TextButton) {
    b := &TextButton{}
    b.Wrapper     = b
    b.Init()
    b.label       = label
    b.FillColor        = pr.Color(ButtonColor)
    b.FillFocusColor   = pr.Color(ButtonFocusColor)
    b.BorderColor      = pr.Color(ButtonBorderColor)
    b.BorderFocusColor = pr.Color(ButtonBorderFocusColor)
    b.LineWidth        = pr.Size(ButtonBorderSize)
    b.fontSize         = pr.Size(TextSize)
    b.fontFace         = fonts.NewFace(pr.Font(BoldFont), b.fontSize)
    b.TextColor        = pr.Color(TextColor)
    b.updateSize()
    return b
}

func (b *TextButton) SetSize(size geom.Point) {
    b.Button.SetSize(size)
    b.updateRefPoint()
}

func (b *TextButton) updateSize() {
    w := float64(font.MeasureString(b.fontFace, b.label)) / 64.0
    h := pr.Size(ButtonSize)
    //h := float64(b.fontFace.Metrics().Ascent +
    //        b.fontFace.Metrics().Descent) / 64.0
    b.desc = float64(b.fontFace.Metrics().Descent) / 64.0
    b.SetMinSize(geom.Point{w+2*pr.Size(TextButtonPaddingSize), h})
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
    fontSize float64
    fontFace font.Face
    TextColor color.Color
    Options []string
    selIdx int
    Selected string
}

func NewListButton(options []string) (*ListButton) {
    b := &ListButton{}
    //b.TextButton  = *NewTextButton(label)
    b.Wrapper     = b
    b.Init()
    b.FillColor        = pr.Color(ButtonColor)
    b.FillFocusColor   = pr.Color(ButtonFocusColor)
    b.BorderColor      = pr.Color(ButtonBorderColor)
    b.BorderFocusColor = pr.Color(ButtonBorderFocusColor)
    b.LineWidth   = pr.Size(ButtonBorderSize)
    b.fontSize         = pr.Size(TextSize)
    b.fontFace         = fonts.NewFace(pr.Font(RegularFont), b.fontSize)
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
    w := maxWidth/64.0 + 2.0*pr.Size(TextButtonPaddingSize) + pr.Size(ButtonSize)
    h := pr.Size(ButtonSize)
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
    pt := geom.Point{0.6*b.Size().Y+2*pr.Size(InnerPaddingSize), 0.5*b.Size().Y}
    gc.DrawStringAnchored(b.Selected, pt.X, pt.Y, 0.0, 0.5)

    gc.SetStrokeWidth(b.LineWidth)
    if b.pushed {
        gc.SetStrokeColor(b.BorderFocusColor)
    } else {
        gc.SetStrokeColor(b.BorderColor)
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
    gc.SetFillColor(pr.Color(ArrowColor))
    gc.SetStrokeColor(pr.Color(ArrowColor))
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
    data binding.Untyped
}

func NewIconButton(imgFile string) (*IconButton) {
    b := &IconButton{}
    b.Wrapper = b
    b.Init()
    b.pushed = false
    b.img, _ = gg.LoadPNG(imgFile)
    i := pr.Size(InnerPaddingSize)
    rect := geom.NewRectangleIMG(b.img.Bounds()).Inset(-i, -i)
    b.SetMinSize(rect.Size())
    b.FillColor        = pr.Color(IconButtonColor)
    b.FillFocusColor   = pr.Color(IconButtonFocusColor)
    b.BorderColor      = pr.Color(IconButtonBorderColor)
    b.BorderFocusColor = pr.Color(IconButtonBorderFocusColor)
    b.LineWidth        = pr.Size(ButtonBorderSize)
    b.data             = binding.NewUntyped()
    return b
}

func NewIconButtonWithCallback(imgFile string, callback func(interface {})) (*IconButton) {
    b := NewIconButton(imgFile)
    b.data.AddCallback(func (data binding.DataItem) {
        callback(data.(binding.Untyped).Get())
    })
    return b
}

func NewIconButtonWithData(imgFile string, data binding.Untyped) (*IconButton) {
    b := NewIconButton(imgFile)
    b.data = data
    b.data.AddListener(b)
    return b
}

func (b *IconButton) OnInputEvent(evt touch.Event) {
    b.Button.OnInputEvent(evt)
    if evt.Type == touch.TypeTap {
        if !b.checked {
            b.data.Set(b)
        } else {
            b.data.Set(nil)
        }
        b.Mark(MarkNeedsPaint)
    }
}

func (b *IconButton) Paint(gc *gg.Context) {
    //log.Printf("IconButton.Paint()")
    if b.checked {
        b.FillColor   = pr.Color(IconButtonSelColor)
        b.BorderColor = pr.Color(IconButtonBorderSelColor)
    } else {
        b.FillColor   = pr.Color(IconButtonColor)
        b.BorderColor = pr.Color(IconButtonBorderColor)
    }
    b.Button.Paint(gc)
    cp := b.Bounds().Center()
    gc.DrawImageAnchored(b.img, cp.X, cp.Y, 0.5, 0.5)
}

func (b *IconButton) DataChanged(data binding.DataItem) {
    value := data.(binding.Untyped).Get()
    if b == value {
        b.checked = true
    } else {
        b.checked = false
    }
}

// Buttons sind neutrale Knoepfe, ohne spezifischen Inhalt, d.h. ohne Text
// oder Icons. Sie werden selten direkt verwendet, sondern dienen als
// generische Grundlage fuer die weiter unten definierten Text- oder Icon-
// Buttons.
type TabButton struct {
    Button
    selected bool
    label string
    fontSize float64
    fontFace font.Face
    textColor color.Color
    data binding.Int
    tabIndex int
}

func NewTabButton(label string) (*TabButton) {
    b := &TabButton{}
    b.Wrapper = b
    b.Init()
    b.SetMinSize(geom.Point{pr.Size(TabButtonWidth), pr.Size(TabButtonHeight)})
    b.FillColor        = pr.Color(TabButtonColor)
    b.FillFocusColor   = pr.Color(TabButtonFocusColor)
    b.BorderColor      = pr.Color(TabButtonBorderColor)
    b.BorderFocusColor = pr.Color(TabButtonBorderFocusColor)
    b.LineWidth        = pr.Size(TabButtonBorderSize)
    b.pushed           = false
    b.selected         = false
    b.label            = label
    b.fontSize         = pr.Size(TabButtonTextSize)
    b.fontFace         = fonts.NewFace(pr.Font(BoldFont), b.fontSize)
    b.textColor        = pr.Color(TextColor)
    b.data             = binding.NewInt()
    b.tabIndex         = -1
    return b
}

func NewTabButtonWithData(label string, data binding.Int) (*TabButton) {
    b := NewTabButton(label)
    b.data = data
    b.data.AddListener(b)
    return b
}

func (b *TabButton) Paint(gc *gg.Context) {
    //log.Printf("Button.Paint()")
    gc.DrawRectangle(b.Bounds().AsCoord())
    gc.Clip()
    gc.DrawRoundedRectangle(0.0, 0.0,
            b.Size().X, b.Size().Y+pr.Size(TabButtonCornerRad),
            pr.Size(TabButtonCornerRad))
    if b.selected {
        b.FillColor   = pr.Color(TabButtonSelColor)
        b.BorderColor = pr.Color(TabButtonBorderSelColor)
    } else {
        b.FillColor   = pr.Color(TabButtonColor)
        b.BorderColor = pr.Color(TabButtonBorderColor)
    }
    if b.pushed {
        gc.SetFillColor(b.FillFocusColor)
        gc.SetStrokeColor(b.BorderFocusColor)
    } else {
        gc.SetFillColor(b.FillColor)
        gc.SetStrokeColor(b.BorderColor)
    }
    gc.SetStrokeWidth(b.LineWidth)
    gc.FillStroke()
    gc.ResetClip()

    mp := b.Bounds().Center()
    if b.selected || b.pushed {
        gc.SetStrokeColor(b.textColor)
    } else {
        gc.SetStrokeColor(pr.Color(TextDimColor))
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
        if !b.selected {
            b.data.Set(b.tabIndex)
        } else {
            b.data.Set(-1)
        }
        b.Mark(MarkNeedsPaint)
    }
    b.CallTouchFunc(evt)
}

func (b *TabButton) SetTabIndex(idx int) {
    b.tabIndex = idx
}

func (b *TabButton) TabIndex() (int) {
    return b.tabIndex
}

func (b *TabButton) DataChanged(data binding.DataItem) {
    value := data.(binding.Int).Get()
    if b.tabIndex == value {
        b.selected = true
    } else {
        b.selected = false
    }
}

// Checkboxen verhalten sich sehr aehnlich zu RadioButtons, sind jeoch eigen-
// staendig und nicht Teil einer Gruppe.
type Checkbox struct {
    Button
    label string
    fontSize float64
    fontFace font.Face
    TextColor color.Color
    value binding.Bool
}

func NewCheckbox(label string) (*Checkbox) {
    c := &Checkbox{}
    c.Wrapper = c
    c.Init()
    c.FillColor        = pr.Color(ButtonColor)
    c.FillFocusColor   = pr.Color(ButtonFocusColor)
    c.BorderColor      = pr.Color(ButtonBorderColor)
    c.BorderFocusColor = pr.Color(ButtonBorderFocusColor)
    c.LineWidth        = pr.Size(ButtonBorderSize)
    c.label       = label
    c.fontSize         = pr.Size(TextSize)
    c.fontFace         = fonts.NewFace(pr.Font(RegularFont), c.fontSize)
    c.TextColor   = pr.Color(TextColor)
    w := float64(font.MeasureString(c.fontFace, label))/64.0
    c.SetMinSize(geom.Point{pr.Size(CheckSize)+pr.Size(InnerPaddingSize)+w, pr.Size(CheckSize)})
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
    gc.DrawRoundedRectangle(0.0, 0.0, pr.Size(CheckSize), pr.Size(CheckSize),
            pr.Size(CheckCornerRad))
    if c.pushed {
        gc.SetStrokeColor(c.BorderFocusColor)
        gc.SetFillColor(c.FillFocusColor)
    } else {
        gc.SetStrokeColor(c.BorderColor)
        gc.SetFillColor(c.FillColor)
    }
    gc.SetStrokeWidth(c.LineWidth)
    gc.FillStroke()
    if c.Checked() {
        gc.SetStrokeWidth(pr.Size(CheckLineSize))
        gc.SetStrokeColor(pr.Color(ArrowColor))
        gc.MoveTo(4, 9)
        gc.LineTo(8, 14)
        gc.LineTo(14, 5)
        gc.Stroke()
    }
    x := pr.Size(CheckSize) + pr.Size(InnerPaddingSize)
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
    fontSize float64
    fontFace font.Face
    TextColor color.Color
    checked bool
    value int
    data binding.Int
}

func NewRadioButtonWithData(label string, value int, data binding.Int) (*RadioButton) {
    b := &RadioButton{}
    b.Wrapper = b
    b.Init()
    b.FillColor        = pr.Color(ButtonColor)
    b.FillFocusColor   = pr.Color(ButtonFocusColor)
    b.BorderColor      = pr.Color(ButtonBorderColor)
    b.BorderFocusColor = pr.Color(ButtonBorderFocusColor)
    b.LineWidth        = pr.Size(ButtonBorderSize)
    b.label       = label
    b.fontSize         = pr.Size(TextSize)
    b.fontFace         = fonts.NewFace(pr.Font(RegularFont), b.fontSize)
    b.TextColor   = pr.Color(TextColor)
    w := float64(font.MeasureString(b.fontFace, label))/64.0
    b.SetMinSize(geom.Point{pr.Size(RadioSize)+pr.Size(InnerPaddingSize)+w, pr.Size(RadioSize)})
    b.value = value
    b.data = data
    b.data.AddListener(b)
    return b
}

func (b *RadioButton) Paint(gc *gg.Context) {
    //log.Printf("RadioButton.Paint()")
    mp := geom.Point{0.5*pr.Size(RadioSize), 0.5*pr.Size(RadioSize)}
    gc.DrawCircle(mp.X, mp.Y, 0.5*pr.Size(RadioSize))
    if b.pushed {
        gc.SetStrokeColor(b.BorderFocusColor)
        gc.SetFillColor(b.FillFocusColor)
    } else {
        gc.SetStrokeColor(b.BorderColor)
        gc.SetFillColor(b.FillColor)
    }
    gc.SetStrokeWidth(b.LineWidth)
    gc.FillStroke()
    if b.checked {
        gc.DrawCircle(mp.X, mp.Y, 0.5*pr.Size(RadioDotSize))
	gc.SetFillColor(pr.Color(ArrowColor))
	gc.Fill()
    }
    x := pr.Size(RadioSize)+pr.Size(InnerPaddingSize)
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
type Scrollbar struct {
    LeafEmbed
    BarColor, BarFocusColor, CtrlColor, CtrlFocusColor color.Color
    orient Orientation
    initValue, visiRange float64
    //len float64
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
    s.orient = orient
    if s.orient == Horizontal {
        s.SetMinSize(geom.Point{len, pr.Size(ScrollSize)})
        s.barStart = geom.Point{0.5*pr.Size(ScrollBarSize), 0.5*pr.Size(ScrollSize)}
        s.ctrlStart = geom.Point{0.5*max(pr.Size(ScrollCtrlSize), pr.Size(ScrollBarSize)),
            0.5*pr.Size(ScrollSize)}
    } else {
        s.SetMinSize(geom.Point{pr.Size(ScrollSize), len})
        s.barStart = geom.Point{0.5*pr.Size(ScrollSize), 0.5*pr.Size(ScrollBarSize)}
        s.ctrlStart = geom.Point{0.5*pr.Size(ScrollSize), 0.5*max(pr.Size(ScrollCtrlSize),
            pr.Size(ScrollBarSize))}
    }
    s.initValue = 0.0
    s.visiRange = 0.1
    s.value     = binding.NewFloat()
    s.BarColor       = pr.Color(ScrollBarColor)
    s.BarFocusColor  = pr.Color(ScrollBarFocusColor)
    s.CtrlColor      = pr.Color(ScrollCtrlColor)
    s.CtrlFocusColor = pr.Color(ScrollCtrlFocusColor)
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
        gc.SetStrokeColor(s.BarFocusColor)
    } else {
        gc.SetStrokeColor(s.BarColor)
    }
    gc.SetStrokeWidth(pr.Size(ScrollBarSize))
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
        gc.SetStrokeColor(s.CtrlFocusColor)
    } else {
        gc.SetStrokeColor(s.CtrlColor)
    }
    gc.SetStrokeWidth(pr.Size(ScrollCtrlSize))
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
type Slider struct {
    LeafEmbed
    BarColor, BarFocusColor, CtrlColor, CtrlFocusColor color.Color
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
    s.orient = orient
    if s.orient == Horizontal {
        s.SetMinSize(geom.Point{len, pr.Size(SliderSize)})
        s.barStart = geom.Point{0.5*pr.Size(SliderBarSize), 0.5*pr.Size(SliderSize)}
        s.ctrlStart = geom.Point{0.5*max(pr.Size(SliderCtrlSize), pr.Size(SliderBarSize)),
            0.5*pr.Size(SliderSize)}
    } else {
        s.SetMinSize(geom.Point{pr.Size(SliderSize), len})
        s.barStart = geom.Point{0.5*pr.Size(SliderSize), 0.5*pr.Size(SliderBarSize)}
        s.ctrlStart = geom.Point{0.5*pr.Size(SliderSize), 0.5*max(pr.Size(SliderCtrlSize),
            pr.Size(SliderBarSize))}
    }
    s.initValue = 0.0
    s.minValue  = 0.0
    s.maxValue  = 1.0
    s.stepSize  = 0.1
    s.value     = binding.NewFloat()
    s.BarColor       = pr.Color(SliderBarColor)
    s.BarFocusColor  = pr.Color(SliderBarFocusColor)
    s.CtrlColor      = pr.Color(SliderCtrlColor)
    s.CtrlFocusColor = pr.Color(SliderCtrlFocusColor)
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
        gc.SetStrokeColor(s.BarFocusColor)
    } else {
        gc.SetStrokeColor(s.BarColor)
    }
    gc.SetStrokeWidth(pr.Size(SliderBarSize))
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
        gc.SetStrokeColor(s.CtrlFocusColor)
    } else {
        gc.SetStrokeColor(s.CtrlColor)
    }
    gc.SetStrokeWidth(pr.Size(SliderCtrlSize))
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
type Circle struct {
    LeafEmbed
    StrokeColor, FillColor color.Color
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
    gc.SetStrokeWidth(c.LineWidth)
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
    StrokeColor, FillColor color.Color
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
    gc.SetStrokeWidth(e.LineWidth)
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
    StrokeColor, FillColor color.Color
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
    gc.SetStrokeWidth(r.LineWidth)
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
        check(err)
        img, _, err := image.Decode(fh)
        check(err)
        s.imgList = append(s.imgList, img)
        fh.Close()
    }
}

func (s *Sprite) Paint(gc *gg.Context) {
    s.Marks.UnmarkNeedsPaint()
    gc.Push()
    gc.Multiply(s.Matrix())
    gc.DrawImage(s.imgList[s.curImg], s.Rect().Min.X, s.Rect().Min.Y)
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
    FillColor, StrokeColor color.Color
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
    gc.Multiply(c.Matrix())
    gc.SetFillColor(c.FillColor)
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

