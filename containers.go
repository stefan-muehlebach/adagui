// In diesem File befinden sich alle Widgets, die im Zusammenhang mit adagui
// existieren. Aktuell sind dies:
//
// Container Widgets
// -----------------
//
//	Group
//	Panel
//	ScrollPanel
package adagui

import (
	"container/list"
	"github.com/stefan-muehlebach/adagui/binding"
	"github.com/stefan-muehlebach/adagui/touch"
	"github.com/stefan-muehlebach/gg"
//	"github.com/stefan-muehlebach/gg/color"
//	"github.com/stefan-muehlebach/gg/colornames"
	"github.com/stefan-muehlebach/gg/geom"
	"golang.org/x/image/draw"
	"image"
	"log"
)

// Alle GUI-Typen, welche weitere Nodes verwalten können (Fenster, Panels,
// etc.) müssen dagegen diesen Typ einbetten. Damit kann über die ChildList
// die angehängten Nodes verwaltet werden. Ebenso kann ein LayoutManager
// verwendet werden, der für die Platzierung der Nodes zuständig ist.
// Per Default wird das NullLayout verwendet, d.h. die Kinder müssen per
// SetPos platziert werden und bleiben an dieser Stelle.
type ContainerEmbed struct {
	Embed
	touch.TouchEmbed
	ChildList *list.List
	Layout    LayoutManager
}

func (c *ContainerEmbed) Init() {
	c.Embed.Init()
	c.ChildList = list.New()
	c.Layout = &NullLayout{}
}

func (c *ContainerEmbed) Add(n ...Node) {
	for _, node := range n {
		embed := node.Wrappee()
		if embed.Parent != nil {
			log.Fatal("Container: Add called for an attached child")
		}
		embed.Win = c.Win
		embed.Parent = c
		c.ChildList.PushBack(embed)
		c.layout()
	}
}

func (c *ContainerEmbed) Del(n Node) {
	for elem := c.ChildList.Front(); elem != nil; elem = elem.Next() {
		node := elem.Value.(Node)
		if n != node {
			continue
		}
		embed := node.Wrappee()
		embed.Win = nil
		embed.Parent = nil
		c.ChildList.Remove(elem)
		break
	}
	c.layout()
}

func (c *ContainerEmbed) DelAll() {
	for elem := c.ChildList.Front(); elem != nil; elem = elem.Next() {
		embed := elem.Value.(*Embed)
		embed.Parent = nil
		embed.Win = nil
	}
	c.ChildList.Init()
	c.layout()
}

func (c *ContainerEmbed) SetSize(s geom.Point) {
	c.Embed.SetSize(s)
	c.layout()
}

func (c *ContainerEmbed) MinSize() geom.Point {
	ms := geom.Point{}
	if c.minSize.Eq(geom.Point{0, 0}) {
		ms = c.Layout.MinSize(c.ChildList)
	} else {
		ms = c.Embed.MinSize()
	}
	return ms
}

func (c *ContainerEmbed) Paint(gc *gg.Context) {
	Debugf(Painting, "[%T], LocalBounds: %v", c.Wrapper, c.LocalBounds())
	c.Marks.UnmarkNeedsPaint()
	for elem := c.ChildList.Front(); elem != nil; elem = elem.Next() {
		child := elem.Value.(*Embed)
		if !child.Visible() {
			continue
		}
		if !c.Wrapper.LocalBounds().Overlaps(child.ParentBounds()) {
			continue
		}
		child.Paint(gc)
	}
}

func (c *ContainerEmbed) OnChildMarked(child Node, newMarks Marks) {
	c.Mark(newMarks)
}

func (c *ContainerEmbed) SelectTarget(pt geom.Point) Node {
	Debugf(Coordinates, "[%T], pt: %v", c.Wrapper, pt)
	if !c.Wrapper.Contains(pt) {
	    Debugf(Coordinates, "is not inside this container")
		return nil
	}
	pt = c.Parent2Local(pt).Add(c.Wrapper.LocalBounds().Min)
	Debugf(Coordinates, "pt after Parent2Local: %v", pt)
	for elem := c.ChildList.Back(); elem != nil; elem = elem.Prev() {
		embed := elem.Value.(*Embed)
		node := embed.Wrapper.SelectTarget(pt)
		if node != nil {
		    Debugf(Coordinates, "target found: %T", node)
			return node
		}
	}
	if !c.selectable {
    	Debugf(Coordinates, "container is not selectable")
	    return nil
	}
	Debugf(Coordinates, "no target found, returning %T", c.Wrapper)
	return c.Wrapper
}

func (c *ContainerEmbed) layout() {
	if c.Layout == nil {
		return
	}
	c.Layout.Layout(c.ChildList, c.Wrapper.Size())
}

// Eine Group ist die einfachste Form eines Containers. Sie dient
// hauptsaechlich als logisches Sammelbecken fuer Widgets auf dem Screen.
// Sie hat zwar ein eigenes Koordinatensystem und beherrscht alle Layouts,
// ist jedoch selber unsichtbar, d.h. dieses Widget hat weder Farbe, Rahmen
// noch andere optische Merkmale.
type Group struct {
	ContainerEmbed
}

func NewGroup() *Group {
	g := &Group{}
	g.Wrapper = g
	g.Init()
	g.PropertyEmbed.InitByName("Default")
	g.selectable = false
	return g
}

func NewGroupPL(parent Container, layout LayoutManager) *Group {
    g := NewGroup()
    g.Layout = layout
    if parent != nil {
        parent.Add(g)
    }
    return g
}

func (g *Group) Paint(gc *gg.Context) {
	Debugf(Painting, "[%T]", g.Wrapper)
	g.ContainerEmbed.Paint(gc)
}

// Ein Panel ist eine etwas komplexere Version eines Containers. Im Gegensatz
// zur Group ist ein Panel auf dem Bildschirm sichtbar. Ueber Properties
// laesst sich die visuelle Erscheinung beeinflussen. Panels beschneiden
// ihren Inhalt auf ihre Groesse. Sie koennen eine Hintergrundfarbe oder
// ein Hitergundbild haben.
type Panel struct {
	ContainerEmbed
	Image image.Image
	IsClipping bool
}

func NewPanel(w, h float64) *Panel {
	p := &Panel{}
	p.Wrapper = p
	p.Init()
	p.PropertyEmbed.InitByName("Panel")
	p.SetMinSize(geom.Point{w, h})
	return p
}

func (p *Panel) Paint(gc *gg.Context) {
	Debugf(Painting, "[%T], LocalBounds: %v", p.Wrapper, p.LocalBounds())

	gc.DrawRectangle(p.LocalBounds().AsCoord())
	gc.SetFillColor(p.Color())
	gc.SetStrokeColor(p.BorderColor())
	gc.SetStrokeWidth(p.BorderWidth())
	gc.FillStroke()

	if p.Image != nil {
		dst := gc.Image().(*image.RGBA)
		draw.NearestNeighbor.Transform(dst, p.Matrix().AsAff3(),
			p.Image, p.LocalBounds().Int(), draw.Over, nil)
	}

	if p.IsClipping {
		gc.DrawRectangle(p.LocalBounds().AsCoord())
		gc.Clip()
		p.ContainerEmbed.Paint(gc)
		gc.ResetClip()
	} else {
		p.ContainerEmbed.Paint(gc)
	}
}

// Komplexeres Panel mit Scrollmoeglichkeit.
type ScrollPanel struct {
	ContainerEmbed
	Image                               image.Image
	virtSize, sizeDiff, viewPort, refPt geom.Point
}

func NewScrollPanel(w, h float64) *ScrollPanel {
	p := &ScrollPanel{}
	p.Wrapper = p
	p.Init()
	p.PropertyEmbed.InitByName("Panel")

	p.SetMinSize(geom.Point{w, h})
	p.SetVirtualSize(p.Size())
	p.viewPort = geom.Point{0, 0}
	p.refPt = geom.Point{0, 0}
	return p
}

func (p *ScrollPanel) Paint(gc *gg.Context) {
	Debugf(Painting, "[%T], LocalBounds: %v, RefPt: %v", p.Wrapper,
	    p.LocalBounds(), p.refPt)
	gc.Translate(p.refPt.Neg().AsCoord())
	gc.DrawRectangle(p.LocalBounds().AsCoord())
	gc.ClipPreserve()
	gc.SetFillColor(p.Color())
	gc.SetStrokeColor(p.BorderColor())
	gc.SetStrokeWidth(p.BorderWidth())
	gc.FillStroke()
	p.ContainerEmbed.Paint(gc)
	gc.ResetClip()
}

func (p *ScrollPanel) LocalBounds() geom.Rectangle {
	return geom.Rectangle{Min: p.refPt, Max: p.refPt.Add(p.Size())}
}

func (p *ScrollPanel) Size() (geom.Point) {
    return p.size
}

func (p *ScrollPanel) SetSize(size geom.Point) {
    Debugf(Layout, "[%T], %+v", p.Wrapper, size)
	p.Embed.SetSize(size)
	p.layout()
}

func (p *ScrollPanel) MinSize() geom.Point {
	ms := geom.Point{}
	if p.minSize.Eq(geom.Point{0, 0}) {
		ms = p.Layout.MinSize(p.ChildList)
	} else {
		ms = p.Embed.MinSize()
	}
    Debugf(Layout, "[%T], %+v", p.Wrapper, ms)
	return ms
}

func (p *ScrollPanel) VisibleRange() geom.Point {
	if p.virtSize.X == 0.0 && p.virtSize.Y == 0.0 {
		return geom.Point{1, 1}
	}
	vis := p.Wrapper.Size()
	vis.X /= p.virtSize.X
	vis.Y /= p.virtSize.Y
	return vis
}

func (p *ScrollPanel) SetXView(vx float64) {
	p.refPt.X = p.sizeDiff.X * vx
}

func (p *ScrollPanel) SetYView(vy float64) {
	p.refPt.Y = p.sizeDiff.Y * vy
}

func (p *ScrollPanel) ViewPort() geom.Point {
	return p.viewPort
}

// Bestimmt die neue virtuelle Groesse des ScrolledPanels. Man kann bei
// sz keine Angaben machen, die kleiner als die eigentliche Groesse des
// Widgets ist.
func (p *ScrollPanel) VirtualSize() geom.Point {
	return p.virtSize
}
func (p *ScrollPanel) SetVirtualSize(sz geom.Point) {
	if sz.X < p.Size().X {
		sz.X = p.Size().X
	}
	if sz.Y < p.Size().Y {
		sz.Y = p.Size().Y
	}
	p.virtSize = sz
	p.sizeDiff = p.virtSize.Sub(p.Size())
}

// TabPanel und TabButton sind fuer Tabbed Windows gedacht.
/*
type TabPanel struct {
	ContainerEmbed
	menu    *TabMenu
	content *Panel
}

func NewTabPanel(w, h float64, menu *TabMenu, content *Panel) *TabPanel {
	p := &TabPanel{}
	p.Wrapper = p
	p.Init()
	p.PropertyEmbed.InitByName("TabPanel")
	p.SetMinSize(geom.Point{w, h})
	p.Layout = NewVBoxLayout(0)
	p.menu = menu
	p.menu.panel = p
	p.content = content
	p.Add(p.menu, p.content)
	return p
}
*/

type TabMenu struct {
	ContainerEmbed
	content     Container
	data        binding.Int
	contentList []Node
}

func NewTabMenu(content Container) *TabMenu {
	m := &TabMenu{}
	m.Wrapper = m
	m.Init()
	m.PropertyEmbed.InitByName("TabMenu")
    m.SetMinSize(geom.Point{m.Width(), m.Height()})
	m.Layout = NewHBoxLayout(0)
	m.data = binding.NewInt()
	m.data.Set(-1)
    m.content = content
	m.contentList = make([]Node, 0)
	m.data.AddCallback(func(d binding.DataItem) {
		idx := d.(binding.Int).Get()
		if (idx < 0) || (idx >= len(m.contentList)) ||
			(m.contentList[idx] == nil) {
			return
		}
		m.content.DelAll()
		m.content.Add(m.contentList[idx])
		m.content.layout()
	})
	return m
}

func (m *TabMenu) AddTab(label string, content Node) (int) {
	tabIndex := len(m.contentList)
	m.contentList = append(m.contentList, content)
	b := NewTabButtonWithData(label, tabIndex, m.data)
	m.Add(b)
	m.layout()
	return tabIndex
}

func (m *TabMenu) SetTab(idx int) {
	m.data.Set(idx)
}

func (m *TabMenu) Paint(gc *gg.Context) {
	Debugf(Painting, "[%T], LocalBounds: %v", m.Wrapper, m.LocalBounds())
	gc.DrawRectangle(m.LocalBounds().AsCoord())
	gc.SetFillColor(m.Color())
    gc.Fill()
//	gc.FillPreserve()
//	gc.Clip()
	m.ContainerEmbed.Paint(gc)
//	gc.ResetClip()
}

