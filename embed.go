package adagui

import (
    "container/list"
    "log"
    "github.com/stefan-muehlebach/adagui/touch"
    "github.com/stefan-muehlebach/gg/geom"
    "github.com/stefan-muehlebach/gg"
)

// Dieser Typ ist die Basis für alle graphischen Typen von AdaGui. Er kann
// jedoch nicht direkt verwendet werden, sondern dient als Basis für zwei
// weitere Embed-Typen (siehe weiter unten), welche für eigene Typen
// verwendet werden können.
type Embed struct {
    Win *Window
    Wrapper Node
    Parent *ContainerEmbed
    pos, size, minSize geom.Point
    transl, rotate, scale, transf geom.Matrix
    Marks Marks
    visible bool
    Prop *Properties
}

func (m *Embed) Init(parentProps *Properties) {
    m.transl  = geom.Identity()
    m.rotate  = geom.Identity()
    m.scale   = geom.Identity()
    m.transf  = geom.Identity()
    m.visible = true
    m.Prop    = NewProperties(parentProps)
}

func (m *Embed) Wrappee() (*Embed) {
    return m
}

func (m *Embed) IsAtFront() bool {
    var e *list.Element

    if m.Parent == nil {
        log.Fatal("node: This child is not attached")
    }
    p := m.Parent
    e = p.ChildList.Back()
    return e.Value.(*Embed) == m
}

func (m *Embed) ToBack() {
    var e *list.Element

    if m.Parent == nil {
        log.Fatal("node: This child is not attached")
    }
    p := m.Parent
    for e = p.ChildList.Front(); e != nil; e = e.Next() {
        if e.Value.(*Embed) == m {
            break
        }
    }
    if e == nil {
        return
    }
    p.ChildList.MoveToFront(e)
}

func (m *Embed) ToFront() {
    var e *list.Element

    if m.Parent == nil {
        log.Fatal("node: This child is not attached")
    }
    p := m.Parent
    for e = p.ChildList.Front(); e != nil; e = e.Next() {
        if e.Value.(*Embed) == m {
            break
        }
    }
    if e == nil {
        return
    }
    p.ChildList.MoveToBack(e)
}

func (m *Embed) Remove() {
    var e *list.Element

    if m.Parent == nil {
        log.Fatal("node: This child is not attached")
    }
    p := m.Parent
    for e = p.ChildList.Front(); e != nil; e = e.Next() {
        if e.Value.(*Embed) == m {
            break
        }
    }
    if e == nil {
        return
    }
    m.Win = nil
    p.ChildList.Remove(e)
}

func (m *Embed) Pos() (geom.Point) {
    return m.pos
}
func (m *Embed) SetPos(p geom.Point) {
    m.pos = p
    m.Translate(p)
}
func (m *Embed) Size() (geom.Point) {
    return m.size.Max(m.Wrapper.MinSize())
}
func (m *Embed) SetSize(s geom.Point) {
    m.size = s
}
func (m *Embed) MinSize() (geom.Point) {
    return m.minSize
}
func (m *Embed) SetMinSize(s geom.Point) {
    m.minSize = s
}

func (m *Embed) LocalBounds() (geom.Rectangle) {
    return geom.Rectangle{Max: m.Size()}
}
func (m *Embed) Bounds() (geom.Rectangle) {
    return geom.Rectangle{Max: m.Size()}
}

func (m *Embed) ParentBounds() (geom.Rectangle) {
    return m.LocalBounds().Add(m.Pos())
}
func (m *Embed) Rect() (geom.Rectangle) {
    return m.LocalBounds().Add(m.Pos())
}

func (m *Embed) Visible() (bool) {
    return m.visible
}
func (m *Embed) SetVisible(v bool) {
    m.visible = v
}

func (m *Embed) Mark(marks Marks) {
    oldMarks := m.Marks
    m.Marks |= marks
    changedMarks := m.Marks ^ oldMarks
    changedMarks &^= MarkNeedsRecalc
    if changedMarks != 0 && m.Parent != nil {
        m.Parent.Wrapper.OnChildMarked(m.Wrapper, changedMarks)
    }
}

func (m *Embed) Paint(gc *gg.Context) {
    //log.Printf("Embed.Paint() of %T", m.Wrapper)
    m.Marks.UnmarkNeedsPaint()
    gc.Push()
    gc.Multiply(m.Matrix())
    m.Wrapper.Paint(gc)
    gc.Pop()
}

// Contains ermittelt, ob sich der Punkt pt innerhalb des Widgets befindet.
// Die Koordianten in pt muessen relativ zum Bezugssystem von m sein.
func (m *Embed) Contains(pt geom.Point) (bool) {
    pt = m.Parent2Local(pt)
    return pt.In(m.Wrapper.Bounds())
}

// Rechnet die lokalen Koordianten pt in Bildschirmkoordinaten um.
func (m *Embed) Local2Screen(pt geom.Point) (geom.Point) {
    pt = m.Matrix().Transform(pt)
    if m.Parent == nil {
        return pt
    }
    return m.Parent.Local2Screen(pt)
}

// Rechnet die Bildschirmkoordinaten pt in lokale Koordianten um.
func (m *Embed) Screen2Local(pt geom.Point) (geom.Point) {
    if m.Parent != nil {
        pt = m.Parent.Screen2Local(pt)
    }
    return m.Matrix().Inv().Transform(pt)
}

// Rechnet die Koordinaten in pt relativ zum Parent-Node um.
func (m *Embed) Local2Parent(pt geom.Point) (geom.Point) {
    if m.Parent == nil {
        return pt
    }
    return m.Matrix().Transform(pt)
}

// Rechnet die Koordinaten in pt vom relativen Bezugsystem des Parent-Nodes
// zu lokalen Koordinaten um.
func (m *Embed) Parent2Local(pt geom.Point) (geom.Point) {
    if m.Parent == nil {
        return pt
    }
    return m.Matrix().Inv().Transform(pt)
}

// Ersetzt die aktuelle Translation des Nodes durch eine Translation um dp.
func (m *Embed) Translate(dp geom.Point) {
    m.transl = geom.Translate(dp)
    m.Mark(MarkNeedsRecalc)
}

// Ersetzt die aktuelle Rotation des Nodes durch eine Rotation um a um den
// Mittelpunkt des Nodes.
func (m *Embed) Rotate(a float64) {
    m.RotateAbout(m.Size().Mul(0.5), a)
}

// Ersetzt die aktuelle Rotation des Nodes durch eine Rotation um a um den
// angegebenen Drehpunkt.
func (m *Embed) RotateAbout(rp geom.Point, a float64) {
    m.rotate = geom.RotateAbout(rp, a)
    m.Mark(MarkNeedsRecalc)
}

// Ersetzt die aktuelle Skalierung des Nodes durch eine Skalierung um
// sx, sy. Zentrum der Skalierung ist der Mittelpunkt des Nodes.
func (m *Embed) Scale(sx, sy float64) {
    m.ScaleAbout(m.Size().Mul(0.5), sx, sy)
}

// Ersetzt die aktuelle Skalierung des Nodes durch eine Skalierung um
// sx, sy mit sp als Zentrum der Skalierung.
func (m *Embed) ScaleAbout(sp geom.Point, sx, sy float64) {
    m.scale = geom.ScaleAbout(sp, sx, sy)
    m.Mark(MarkNeedsRecalc)
}

// Liefert die aktuelle Transformationsmatrix des Nodes.
func (m *Embed) Matrix() (geom.Matrix) {
    if m.Marks.NeedsRecalc() {
        m.Marks.UnmarkNeedsRecalc()
        m.transf = m.transl.Multiply(m.scale.Multiply(m.rotate))
    }
    return m.transf
}

// Jeder GUI-Typ, der selber keine weiteren Kinder verwaltet, muss diesen
// Typ einbetten.
type LeafEmbed struct {
    Embed
    touch.TouchEmbed
}

func (m *LeafEmbed) Paint(gc *gg.Context) {
    //log.Printf("LeafEmbed.Paint() of %T", m.Wrapper)
}

func (m *LeafEmbed) OnChildMarked(child Node, newMarks Marks) {}

func (m *LeafEmbed) SelectTarget(pt geom.Point) (Node) {
    //stackLevel.Inc()
    //defer stackLevel.Dec()

    if !m.Visible() {
        return nil
    }
    //log.Printf("Leaf.SelectTarget on %T, size %v, %v", m.Wrapper, m.Wrapper.Size(), pt)
    if !m.Wrapper.Contains(pt) {
        //log.Printf("   > point is outside my rect %v", m.Bounds())
        return nil
    }
    //log.Printf("   > target found: %T!", m.Wrapper)
    return m.Wrapper
}

// Umrechnungsmethoden fuer Koordinaten.
func (m *LeafEmbed) Local2Screen(pt geom.Point) (geom.Point) {
    if m.Parent == nil {
        return pt
    }
    return m.Parent.Local2Screen(pt)
}

func (m *LeafEmbed) Screen2Local(pt geom.Point) (geom.Point) {
    if m.Parent == nil {
        return pt
    }
    return m.Parent.Screen2Local(pt)
}

func (m *LeafEmbed) Local2Parent(pt geom.Point) (geom.Point) {
    if m.Parent == nil {
        return pt
    }
    return m.Parent.Local2Parent(pt)
}

func (m *LeafEmbed) Parent2Local(pt geom.Point) (geom.Point) {
    if m.Parent == nil {
        return pt
    }
    return m.Parent.Parent2Local(pt)
}

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
    Layout LayoutManager
}

func (c *ContainerEmbed) Init(parentProps *Properties) {
    //stackLevel.Inc()
    //defer stackLevel.Dec()
    //log.Printf("ContainerEmbed.Init()")

    c.Embed.Init(parentProps)
    c.ChildList = list.New()
    c.Layout = &NullLayout{}
}

func (c *ContainerEmbed) Add(n ...Node) {
    //stackLevel.Inc()
    //defer stackLevel.Dec()
    //log.Printf("ContainerEmbed.Add()")

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
        embed.Win =  nil
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
    //stackLevel.Inc()
    //defer stackLevel.Dec()
    //log.Printf("ContainerEmbed.SetSize(%v)", s)

    c.Embed.SetSize(s)
    c.layout()
}

func (c *ContainerEmbed) MinSize() (geom.Point) {
    //stackLevel.Inc()
    //defer stackLevel.Dec()
    //log.Printf("ContainerEmbed.MinSize()")

    ms := geom.Point{}
    if c.minSize.Eq(geom.Point{0, 0}) {
        //log.Printf("  minSize is zero: calling Layout.MinSize")
        ms = c.Layout.MinSize(c.ChildList)
    } else {
        ms =  c.Embed.MinSize()
    }
    //log.Printf("  > %v", ms)
    return ms
}

func (c *ContainerEmbed) layout() {
    //stackLevel.Inc()
    //defer stackLevel.Dec()
    //log.Printf("ContainerEmbed.layout() (internal func)")

    if c.Layout == nil {
        return
    }
    c.Layout.Layout(c.ChildList, c.Wrapper.Size())
}

func (c *ContainerEmbed) Paint(gc *gg.Context) {
    //log.Printf("ContainerEmbed.Paint() of %T", c.Wrapper)
    c.Marks.UnmarkNeedsPaint()
    for elem := c.ChildList.Front(); elem != nil; elem = elem.Next() {
        child := elem.Value.(*Embed)
        //child := elem.Value.(*Embed).Wrapper
        if !child.Visible() {
            continue
        }
        child.Paint(gc)
    }
}

func (c *ContainerEmbed) OnChildMarked(child Node, newMarks Marks) {
    c.Mark(newMarks)
    if c.Parent == nil && newMarks.NeedsPaint() {
        c.Win.Repaint()
    }
}

func (c *ContainerEmbed) SelectTarget(pt geom.Point) (Node) {
    //stackLevel.Inc()
    //defer stackLevel.Dec()

    //log.Printf("Container.SelectTarget on %T, size %v, %v", c.Wrapper, c.Wrapper.Size(), pt)
    if !c.Wrapper.Contains(pt) {
        //log.Printf("   > point is outside my rect %v", c.Wrapper.LocalBounds())
        return nil
    }
    pt = c.Parent2Local(pt)
    //log.Printf("   > new local point: %v", pt)
    for elem := c.ChildList.Back(); elem != nil; elem = elem.Prev() {
        embed := elem.Value.(*Embed)
        node := embed.Wrapper.SelectTarget(pt)
        if node != nil {
            //log.Printf("   > target found: %T!", node)
            return node
        }
    }
    //log.Printf("   > no target found, sending my self back: %T", c.Wrapper)
    return c.Wrapper
}

//----------------------------------------------------------------------------

type Marks uint32

const (
    MarkNeedsMeasure = Marks(1 << 0)
    MarkNeedsLayout  = Marks(1 << 1)
    MarkNeedsPaint   = Marks(1 << 2)
    MarkNeedsRecalc  = Marks(1 << 3)
)

func (m Marks)  NeedsMeasure() (bool) { return m & MarkNeedsMeasure != 0 }
func (m Marks)  NeedsLayout() (bool)  { return m & MarkNeedsLayout  != 0 }
func (m Marks)  NeedsPaint() (bool)   { return m & MarkNeedsPaint   != 0  }
func (m Marks)  NeedsRecalc() (bool)  { return m & MarkNeedsRecalc  != 0 }

func (m *Marks) UnmarkNeedsMeasure()  { *m &^= MarkNeedsMeasure }
func (m *Marks) UnmarkNeedsLayout()   { *m &^= MarkNeedsLayout  }
func (m *Marks) UnmarkNeedsPaint()    { *m &^= MarkNeedsPaint   }
func (m *Marks) UnmarkNeedsRecalc()   { *m &^= MarkNeedsRecalc  }

