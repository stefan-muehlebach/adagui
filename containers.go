
// In diesem File befinden sich alle Widgets, die im Zusammenhang mit adagui
// existieren. Aktuell sind dies:
//
// Container Widgets
// -----------------
//   Group
//   Panel
//   ScrollPanel
//
package adagui

import (
    "image"
    //"image/draw"
    "golang.org/x/image/draw"
    "github.com/stefan-muehlebach/adagui/binding"
    "github.com/stefan-muehlebach/gg"
    "github.com/stefan-muehlebach/gg/color"
    "github.com/stefan-muehlebach/gg/colornames"
    "github.com/stefan-muehlebach/gg/geom"
)

// Eine Group ist die einfachste Form eines Containers. Sie dient
// hauptsaechlich als logisches Sammelbecken fuer Widgets auf dem Screen.
// Sie hat zwar ein eigenes Koordinatensystem und beherrscht alle Layouts,
// ist jedoch selber unsichtbar, d.h. dieses Widget hat weder Farbe, Rahmen
// noch andere optische Merkmale.
type Group struct {
    ContainerEmbed
}

func NewGroup() (*Group) {
    g := &Group{}
    g.Wrapper = g
    g.Init()
    g.PropertyEmbed.Init(DefProps)
    return g
}

func (g *Group) Paint(gc *gg.Context) {
    Debugf("type %T", g.Wrapper)
    g.ContainerEmbed.Paint(gc)
}

// Ein Panel ist eine etwas komplexere Version eines Containers. Im Gegensatz
// zur Group ist ein Panel auf dem Bildschirm sichtbar. Ueber Properties
// laesst sich die visuelle Erscheinung beeinflussen. Panels beschneiden
// ihren Inhalt auf ihre Groesse. Sie koennen eine Hintergrundfarbe oder
// ein Hitergundbild haben.
var (
    PanelProps = NewProps(DefProps,
        map[ColorPropertyName]color.Color{
            Color:        colornames.Black,
            BorderColor:  colornames.Black,
        },
        nil,
        map[SizePropertyName]float64{
            BorderWidth:  1.0,
        })
)

// Einfaches Panel, welches seine Objekte 
type Panel struct {
    ContainerEmbed
    Image image.Image
}

func NewPanel(w, h float64) (*Panel) {
    p := &Panel{}
    p.Wrapper = p
    p.Init()
    p.PropertyEmbed.Init(PanelProps)
    p.SetMinSize(geom.Point{w, h})
    return p
}

func (p *Panel) Paint(gc *gg.Context) {
    Debugf("type %T", p.Wrapper)
    Debugf("LocalBounds: %v", p.LocalBounds())
    gc.DrawRectangle(p.LocalBounds().AsCoord())
    if p.Image != nil {
        dst := gc.Image().(*image.RGBA)
        draw.NearestNeighbor.Transform(dst, p.Matrix().AsAff3(),
                p.Image, p.LocalBounds().Int(), draw.Src, nil)
        //draw.Draw(dst, p.LocalBounds().Int(), p.Image,
        //        image.Point{0, 0}, draw.Src)
    } else {
        gc.SetFillColor(p.Color())
        gc.FillPreserve()
    }
    gc.ClipPreserve()
    gc.SetStrokeColor(p.BorderColor())
    gc.SetStrokeWidth(p.BorderWidth())
    gc.Stroke()
    p.ContainerEmbed.Paint(gc)
    gc.ResetClip()
}

// Komplexeres Panel mit Scrollmoeglichkeit.
type ScrollPanel struct {
    ContainerEmbed
    Image image.Image
    virtSize, sizeDiff, viewPort, refPt geom.Point
}

func NewScrollPanel(w, h float64) (*ScrollPanel) {
    p := &ScrollPanel{}
    p.Wrapper = p
    p.Init()
    p.PropertyEmbed.Init(PanelProps)

    p.SetMinSize(geom.Point{w, h})
    p.SetVirtualSize(p.Size())
    p.viewPort    = geom.Point{0, 0}
    p.refPt       = geom.Point{0, 0}
    return p
}

func (p *ScrollPanel) Paint(gc *gg.Context) {
    Debugf("type %T", p.Wrapper)
    Debugf("LocalBounds: %v", p.LocalBounds())
    Debugf("RefPt      : %v", p.refPt)
    gc.Translate(p.refPt.AsCoord())
    gc.DrawRectangle(p.LocalBounds().AsCoord())
    gc.ClipPreserve()
    gc.SetFillColor(p.Color())
    gc.SetStrokeColor(p.BorderColor())
    gc.SetStrokeWidth(p.BorderWidth())
    gc.FillStroke()

    p.ContainerEmbed.Paint(gc)
    gc.ResetClip()
}

func (p *ScrollPanel) LocalBounds() (geom.Rectangle) {
    return geom.Rectangle{Max: p.Size()}.Sub(p.refPt)
}

func (p *ScrollPanel) VisibleRange() (geom.Point) {
    if p.virtSize.X == 0.0 && p.virtSize.Y == 0.0 {
        return geom.Point{1, 1}
    }
    vis := p.Wrapper.Size()
    vis.X /= p.virtSize.X
    vis.Y /= p.virtSize.Y
    return vis
}

func (p *ScrollPanel) SetXView(vx float64) {
    p.viewPort.X = vx
    p.refPt.X = p.sizeDiff.X * p.viewPort.X
}

func (p *ScrollPanel) SetYView(vy float64) {
    p.viewPort.Y = vy
    p.refPt.Y = p.sizeDiff.Y * p.viewPort.Y
}

func (p *ScrollPanel) ViewPort() (geom.Point) {
    return p.viewPort
}

func (p *ScrollPanel) SetVirtualSize(sz geom.Point) {
    p.virtSize = sz
    p.sizeDiff = p.Size().Sub(p.virtSize)
}

func (p *ScrollPanel) VirtualSize() (geom.Point) {
    return p.virtSize
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
    p.Init()
    p.PropertyEmbed.Init(DefProps)
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

