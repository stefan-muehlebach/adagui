package adagui

import (
    "container/list"
    "log"
    "github.com/stefan-muehlebach/adagui/touch"
    "github.com/stefan-muehlebach/gg"
    "github.com/stefan-muehlebach/gg/geom"
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
    PropertyEmbed
}

func (m *Embed) Init() {
    m.transl  = geom.Identity()
    m.rotate  = geom.Identity()
    m.scale   = geom.Identity()
    m.transf  = geom.Identity()
    m.visible = true
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
    //m.Mark(MarkNeedsPaint)
}
func (m *Embed) MinSize() (geom.Point) {
    return m.minSize
}
func (m *Embed) SetMinSize(s geom.Point) {
    m.minSize = s
    //m.Mark(MarkNeedsPaint)
}

func (m *Embed) LocalBounds() (geom.Rectangle) {
    return geom.Rectangle{Max: m.Size()}
}
func (m *Embed) Bounds() (geom.Rectangle) {
    return m.Wrapper.LocalBounds()
}

func (m *Embed) ParentBounds() (geom.Rectangle) {
    return geom.Rectangle{Max: m.Size()}.Add(m.Pos())
}
func (m *Embed) Rect() (geom.Rectangle) {
    return m.Wrapper.ParentBounds()
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
    Debugf(Painting, "type %T", m.Wrapper)
    m.Marks.UnmarkNeedsPaint()
    gc.Push()
    gc.Multiply(m.Matrix())
    m.Wrapper.Paint(gc)
    gc.Pop()
}

// Contains ermittelt, ob sich der Punkt pt innerhalb des Widgets befindet.
// Die Koordianten in pt muessen relativ zum Bezugssystem von m sein.
func (m *Embed) Contains(pt geom.Point) (bool) {
    Debugf(Coordinates, "type %T, pt: %v", m.Wrapper, pt)
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
    m.Mark(MarkNeedsRecalc | MarkNeedsPaint)
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
    m.Mark(MarkNeedsRecalc | MarkNeedsPaint)
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
    m.Mark(MarkNeedsRecalc | MarkNeedsPaint)
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
    Debugf(Painting, "type %T", m.Wrapper)
}

func (m *LeafEmbed) OnChildMarked(child Node, newMarks Marks) {}

func (m *LeafEmbed) SelectTarget(pt geom.Point) (Node) {
    Debugf(Coordinates, "type %T, pt: %v", m.Wrapper, pt)
    if !m.Visible() {
        return nil
    }
    if !m.Wrapper.Contains(pt) {
        return nil
    }
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

