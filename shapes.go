// In diesem File befinden sich alle Widgets, die im Zusammenhang mit adagui
// existieren. Aktuell sind dies:
//
// Leaf Widgets (Graphik bezogen)
// ------------------------------
//   Circle
//   Rectangle
//   Line       (Geplant)
//
package adagui

import (
    "container/list"
    //"log"
    "math"
    "sync"
    "github.com/stefan-muehlebach/adagui/touch"
    "github.com/stefan-muehlebach/gg"
    "github.com/stefan-muehlebach/gg/color"
    "github.com/stefan-muehlebach/gg/geom"
)

// Schoene Kreise fuer Spiele oder was auch immer lassen sich mit diesem
// Widget-Typ auf den Schirm zaubern.
var (
    fangRadius = 5.0

    ShapeProps = NewProps(DefProps,
        map[ColorPropertyName]color.Color{
            Color:               color.Transparent,
            PressedColor:        color.Transparent,
            SelectedColor:       color.Transparent,
            BorderColor:         DefProps.Color(WhiteColor),
            PressedBorderColor:  DefProps.Color(WhiteColor).Alpha(0.5),
            SelectedBorderColor: DefProps.Color(RedColor),
        },
        nil,
        map[SizePropertyName]float64{
            BorderWidth:         2.0,
            PressedBorderWidth:  1.0,
            SelectedBorderWidth: 2.0,
        })

    PointProps = NewProps(ShapeProps,
        map[ColorPropertyName]color.Color{
            Color:               DefProps.Color(WhiteColor),
            PressedColor:        DefProps.Color(WhiteColor).Alpha(0.5),
            SelectedColor:       DefProps.Color(RedColor),
        },
        nil,
        map[SizePropertyName]float64{
            Width:               8.0,
            Height:              8.0,
            BorderWidth:         0.0,
            PressedBorderWidth:  0.0,
            SelectedBorderWidth: 0.0,
        })
)

// Abstrakter, allgemeiner Typ fuer geometrische Formen
type Shape struct {
    LeafEmbed
    pushed bool
    selected bool
}

func (s *Shape) OnInputEvent(evt touch.Event) {
    switch evt.Type {
    case touch.TypePress:
        s.pushed = true
        s.Mark(MarkNeedsPaint)
    case touch.TypeRelease:
        s.pushed = false
        s.Mark(MarkNeedsPaint)
    case touch.TypeTap:
        s.selected = !s.selected
        s.Mark(MarkNeedsPaint)
    }
    s.CallTouchFunc(evt)
}

func (s *Shape) Paint(gc *gg.Context) {
    if s.pushed {
        Debugf("paint pushed")
        gc.SetFillColor(s.PressedColor())
        gc.SetStrokeWidth(s.PressedBorderWidth())
        gc.SetStrokeColor(s.PressedBorderColor())
    } else {
        if s.selected {
            Debugf("paint selected")
            gc.SetFillColor(s.SelectedColor())
            gc.SetStrokeWidth(s.SelectedBorderWidth())
            gc.SetStrokeColor(s.SelectedBorderColor())
        } else {
            Debugf("paint normally")
            gc.SetFillColor(s.Color())
            gc.SetStrokeWidth(s.BorderWidth())
            gc.SetStrokeColor(s.BorderColor())
        }
    }
}

// Punkte
type Point struct {
    Shape
}

func NewPoint() (*Point) {
    p := &Point{}
    p.Wrapper = p
    p.Init()
    p.PropertyEmbed.Init(PointProps)
    p.SetMinSize(geom.Point{p.Width(), p.Height()})
    return p
}

func (p *Point) Paint(gc *gg.Context) {
    Debugf("-")
    p.Shape.Paint(gc)
    mp := p.LocalBounds().Center()
    gc.DrawPoint(mp.X, mp.Y, p.Width()/2)
    gc.FillStroke()
}

func (p *Point) Contains(pt geom.Point) (bool) {
    outer := p.ParentBounds().Inset(-fangRadius, -fangRadius)
    if !pt.In(outer) {
        return false
    }
    return p.Pos().Distance(pt) <= 0.5*p.Width() + fangRadius
}

func (p *Point) Pos() (geom.Point) {
    return p.ParentBounds().Center()
}
func (p *Point) SetPos(mp geom.Point) {
    p.Wrappee().SetPos(mp.Sub(p.Size().Mul(0.5)))
}

// Geraden
type Line struct {
    Shape
    isRaising bool
}

func NewLine(dx, dy float64, isRaising bool) (*Line) {
    l := &Line{}
    l.Wrapper = l
    l.Init()
    l.PropertyEmbed.Init(ShapeProps)
    l.SetMinSize(geom.Point{dx, dy})
    l.isRaising = isRaising
    return l
}

func (l *Line) Paint(gc *gg.Context) {
    Debugf("-")
    l.Shape.Paint(gc)
    if l.isRaising {
        gc.MoveTo(l.Bounds().SW().AsCoord())
        gc.LineTo(l.Bounds().NE().AsCoord())
    } else {
        gc.MoveTo(l.Bounds().Min.AsCoord())
        gc.LineTo(l.Bounds().Max.AsCoord())
    }
    gc.Stroke()
}

func (l *Line) Contains(pt geom.Point) (bool) {
    outer := l.ParentBounds().Inset(-fangRadius, -fangRadius)
    if !pt.In(outer) {
        return false
    }
    fx, fy := l.ParentBounds().PosRel(pt)
    if l.isRaising {
        fy = 1.0 - fy
    }
    Debugf("fx, fy: %f, %f", fx, fy)
    return math.Abs(fx-fy) <= 0.1
}

// Kreis
type Circle struct {
    Shape
}

func NewCircle(r float64) (*Circle) {
    c := &Circle{}
    c.Wrapper = c
    c.Init()
    c.PropertyEmbed.Init(ShapeProps)
    c.SetMinSize(geom.Point{2*r, 2*r})
    return c
}

func (c *Circle) Paint(gc *gg.Context) {
    Debugf("-")
    c.Shape.Paint(gc)
    mp := c.LocalBounds().Center()
    r  := 0.5 * c.Size().X
    gc.DrawCircle(mp.X, mp.Y, r)
    gc.FillStroke()
}

func (c *Circle) Contains(pt geom.Point) (bool) {
    outer := c.ParentBounds().Inset(-fangRadius, -fangRadius)
    if !pt.In(outer) {
        return false
    }
    return math.Abs(c.Radius() - c.Pos().Distance(pt)) <= fangRadius
}

func (c *Circle) Pos() (geom.Point) {
    return c.ParentBounds().Center()
}
func (c *Circle) SetPos(mp geom.Point) {
    c.Wrappee().SetPos(mp.Sub(c.Size().Mul(0.5)))
}

func (c *Circle) Radius() (float64) {
    return 0.5 * c.Size().X
}
func (c *Circle) SetRadius(r float64) {
    mp := c.Pos()
    c.SetMinSize(geom.Point{2*r, 2*r})
    c.Wrappee().SetPos(mp.Sub(c.Size().Mul(0.5)))
}

// Ein allgemeinerer Widget Typ ist die Ellipse.
type Ellipse struct {
    Shape
}

func NewEllipse(rx, ry float64) (*Ellipse) {
    e := &Ellipse{}
    e.Wrapper = e
    e.Init()
    e.PropertyEmbed.Init(ShapeProps)
    e.SetMinSize(geom.Point{2*rx, 2*ry})
    return e
}

func (e *Ellipse) Paint(gc *gg.Context) {
    Debugf("-")
    e.Shape.Paint(gc)
    mp := e.LocalBounds().Center()
    w, h := e.Size().AsCoord()
    gc.DrawEllipse(mp.X, mp.Y, 0.5*w, 0.5*h)
    gc.FillStroke()
}

func (e *Ellipse) Contains(pt geom.Point) (bool) {
    outer := e.ParentBounds().Inset(-fangRadius, -fangRadius)
    if !pt.In(outer) {
        return false
    }
    Debugf("pt    : %v", pt)
    mp := e.Pos()
    Debugf("mp    : %v", mp)
    rx, ry := e.Size().Mul(0.5).AsCoord()
    Debugf("rx, ry: %f, %f", rx, ry)
    dp := mp.Sub(pt)
    Debugf("dp    : %v", dp)
    dpOut := geom.Point{dp.X/(rx+fangRadius), dp.Y/(ry+fangRadius)}
    dpIn  := geom.Point{dp.X/(rx-fangRadius), dp.Y/(ry-fangRadius)}
    return dpOut.Abs() < 1.0 && dpIn.Abs() > 1.0
}

func (e *Ellipse) Pos() (geom.Point) {
    return e.ParentBounds().Center()
}
func (e *Ellipse) SetPos(mp geom.Point) {
    e.Wrappee().SetPos(mp.Sub(e.Size().Mul(0.5)))
}

func (e *Ellipse) Radius() (float64, float64) {
    return e.Size().Mul(0.5).AsCoord()
}
func (e *Ellipse) SetRadius(rx, ry float64) {
    mp := e.Pos()
    e.SetMinSize(geom.Point{2*rx, 2*ry})
    e.Wrappee().SetPos(mp.Sub(e.Size().Mul(0.5)))
}

// Und wo es Kreise gibt, da sind auch die Rechtecke nicht weit.
type Rectangle struct {
    Shape
}

func NewRectangle(w, h float64) (*Rectangle) {
    r := &Rectangle{}
    r.Wrapper = r
    r.Init()
    r.PropertyEmbed.Init(ShapeProps)
    r.SetMinSize(geom.Point{w, h})
    return r
}

func (r *Rectangle) Contains(pt geom.Point) (bool) {
    outer := r.ParentBounds().Inset(-5.0, -5.0)
    inner := r.ParentBounds().Inset(+5.0, +5.0)
    return pt.In(outer) && !pt.In(inner)
}

func (r *Rectangle) Paint(gc *gg.Context) {
    Debugf("-")
    r.Shape.Paint(gc)
    gc.DrawRectangle(r.LocalBounds().AsCoord())
    gc.FillStroke()
}

//-----------------------------------------------------------------------------

type Canvas struct {
    LeafEmbed
    Clip bool
    ObjList *list.List
    mutex *sync.Mutex
}

func NewCanvas(w, h float64) (*Canvas) {
    c := &Canvas{}
    c.Wrapper = c
    c.Init()
    c.PropertyEmbed.Init(DefProps)
    c.SetSize(geom.Point{w, h})
    c.Clip        = false
    c.ObjList     = list.New()
    c.mutex       = &sync.Mutex{}
    return c
}

func (c *Canvas) OnInputEvent(evt touch.Event) {
    c.CallTouchFunc(evt)
}

func (c *Canvas) Paint(gc *gg.Context) {
    c.Marks.UnmarkNeedsPaint()
    gc.Push()
    gc.Multiply(c.Matrix())
    gc.SetFillColor(c.Color())
    gc.SetStrokeColor(c.BorderColor())
    gc.SetStrokeWidth(c.BorderWidth())
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
