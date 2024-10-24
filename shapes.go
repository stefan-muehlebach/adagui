package adagui

import (
	"fmt"
    //"math"
    "github.com/stefan-muehlebach/adagui/touch"
    "github.com/stefan-muehlebach/gg"
    "github.com/stefan-muehlebach/gg/geom"
)

// Schoene Kreise fuer Spiele oder was auch immer lassen sich mit diesem
// Widget-Typ auf den Schirm zaubern.
var (
    fangRadius = 5.0
)

// Abstrakter, allgemeiner Typ fuer geometrische Formen
type Shape struct {
    LeafEmbed
    PushEmbed
    SelectEmbed
}

func (s *Shape) Init() {
    s.LeafEmbed.Init()
    s.PushEmbed.Init(s, nil)
    s.SelectEmbed.Init(s, nil)
}

func (s *Shape) OnInputEvent(evt touch.Event) {
    Debugf(Events, "evt: %v", evt)
    s.PushEmbed.OnInputEvent(evt)
    s.SelectEmbed.OnInputEvent(evt)
    s.CallTouchFunc(evt)
}

// Punkte
type Point struct {
    Shape
}

func NewPoint() (*Point) {
    p := &Point{}
    p.Wrapper = p
    p.Shape.Init()
    p.PropertyEmbed.InitByName("Point")
    p.SetMinSize(geom.Point{p.Width(), p.Height()})
    return p
}

func (p *Point) Paint(gc *gg.Context) {
    Debugf(Painting, "")
    mp := p.LocalBounds().Center()
    gc.DrawPoint(mp.X, mp.Y, p.Width()/2)
    if p.Pushed() {
        gc.SetFillColor(p.PushedColor())
    } else {
        gc.SetFillColor(p.Color())
    }
    gc.FillPreserve()
    if p.Pushed() || p.Selected() {
        gc.SetStrokeWidth(p.PushedBorderWidth())
        gc.SetStrokeColor(p.PushedBorderColor())
        gc.StrokePreserve()
    }
    gc.SetStrokeWidth(p.BorderWidth())
    gc.SetStrokeColor(p.BorderColor())
    gc.Stroke()
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
    p0, p1 geom.Point
}

func NewLine() (*Line) {
    l := &Line{}
    l.Wrapper = l
    l.Shape.Init()
    l.PropertyEmbed.InitByName("Shape")
    //l.SetPos(p0.Min(p1))
    l.p0 = geom.Point{}
    l.p1 = geom.Point{}
    //l.SetP0(p0)
    //l.SetP1(p1)
    return l
}

func (l *Line) P0() (geom.Point) {
    return l.Pos().Add(l.p0)
}
func (l *Line) SetP0(pt geom.Point) {
    if l.Pos().Eq(l.p0) {
        l.Wrappee().SetPos(pt)
        return
    }
    pos := pt.Min(l.Pos())
    dPos := l.Pos().Sub(pos)
    l.p1 = l.p1.Add(dPos)
    l.p0 = pt.Sub(pos)
    l.Wrappee().SetPos(pos)
    l.SetMinSize(geom.Rect(l.p0.X, l.p0.Y, l.p1.X, l.p1.Y).Size())
}

func (l *Line) P1() (geom.Point) {
    return l.Pos().Add(l.p1)
}
func (l *Line) SetP1(pt geom.Point) {
    pos := pt.Min(l.Pos())
    dPos := l.Pos().Sub(pos)
    l.p0 = l.p0.Add(dPos)
    l.p1 = pt.Sub(pos)
    l.Wrappee().SetPos(pos)
    l.SetMinSize(geom.Rect(l.p0.X, l.p0.Y, l.p1.X, l.p1.Y).Size())
}    

func (l *Line) Paint(gc *gg.Context) {
    Debugf(Painting, "")
    gc.DrawLine(l.p0.X, l.p0.Y, l.p1.X, l.p1.Y)
    if l.Pushed() || l.Selected() {
        gc.SetStrokeWidth(l.PushedBorderWidth())
        gc.SetStrokeColor(l.PushedBorderColor())
        gc.StrokePreserve()
    }
    gc.SetStrokeWidth(l.BorderWidth())
    gc.SetStrokeColor(l.BorderColor())
    gc.Stroke()
}

func (l *Line) Contains(pt geom.Point) (bool) {
    if !pt.In(l.ParentBounds()) {
        return false
    }
    fx, fy := l.ParentBounds().PosRel(pt)
    fmt.Printf("fx, fy: %f, %f\n", fx, fy)
    fmt.Printf("  p0, p1: %v, %v\n", l.p0, l.p1)
    return true
    //return math.Abs(fx-fy) <= 0.1
}

// Rechtecke
type Rectangle struct {
    Shape
}

func NewRectangle(w, h float64) (*Rectangle) {
    r := &Rectangle{}
    r.Wrapper = r
    r.Shape.Init()
    r.PropertyEmbed.InitByName("Shape")
    r.SetMinSize(geom.Point{w, h})
    return r
}

func (r *Rectangle) Paint(gc *gg.Context) {
    Debugf(Painting, "")
    gc.DrawRectangle(r.LocalBounds().AsCoord())
    if r.Pushed() {
        gc.SetFillColor(r.PushedColor())
    } else {
        gc.SetFillColor(r.Color())
    }
    gc.FillPreserve()
    if r.Pushed() || r.Selected() {
        gc.SetStrokeWidth(r.PushedBorderWidth())
        gc.SetStrokeColor(r.PushedBorderColor())
        gc.StrokePreserve()
    }
    gc.SetStrokeWidth(r.BorderWidth())
    gc.SetStrokeColor(r.BorderColor())
    gc.Stroke()
}

func (r *Rectangle) Contains(pt geom.Point) (bool) {
    outer := r.ParentBounds().Inset(-fangRadius, -fangRadius)
    return pt.In(outer)
}

// Kreise
type Circle struct {
    Shape
}

func NewCircle(r float64) (*Circle) {
    c := &Circle{}
    c.Wrapper = c
    c.Shape.Init()
    c.PropertyEmbed.InitByName("Shape")
    c.SetMinSize(geom.Point{2*r, 2*r})
    return c
}

func (c *Circle) Paint(gc *gg.Context) {
    Debugf(Painting, "")
    mp := c.LocalBounds().Center()
    r  := 0.5 * c.Size().X
    gc.DrawCircle(mp.X, mp.Y, r)
    if c.Pushed() {
        gc.SetFillColor(c.PushedColor())
    } else {
        gc.SetFillColor(c.Color())
    }
    gc.FillPreserve()
    if c.Pushed() || c.Selected() {
        gc.SetStrokeWidth(c.PushedBorderWidth())
        gc.SetStrokeColor(c.PushedBorderColor())
        gc.StrokePreserve()
    }
    gc.SetStrokeWidth(c.BorderWidth())
    gc.SetStrokeColor(c.BorderColor())
    gc.Stroke()
}

func (c *Circle) Contains(pt geom.Point) (bool) {
    outer := c.ParentBounds().Inset(-fangRadius, -fangRadius)
    if !pt.In(outer) {
        return false
    }
    return c.Pos().Distance(pt) <= c.Radius() + fangRadius
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
    e.Shape.Init()
    e.PropertyEmbed.InitByName("Shape")
    e.SetMinSize(geom.Point{2*rx, 2*ry})
    return e
}

func (e *Ellipse) Paint(gc *gg.Context) {
    Debugf(Painting, "")
    mp := e.LocalBounds().Center()
    w, h := e.Size().AsCoord()
    gc.DrawEllipse(mp.X, mp.Y, 0.5*w, 0.5*h)
    if e.Pushed() {
        gc.SetFillColor(e.PushedColor())
    } else {
        gc.SetFillColor(e.Color())
    }
    gc.FillPreserve()
    if e.Pushed() || e.Selected() {
        gc.SetStrokeWidth(e.PushedBorderWidth())
        gc.SetStrokeColor(e.PushedBorderColor())
        gc.StrokePreserve()
    }
    gc.SetStrokeWidth(e.BorderWidth())
    gc.SetStrokeColor(e.BorderColor())
    gc.Stroke()
}

func (e *Ellipse) Contains(pt geom.Point) (bool) {
    outer := e.ParentBounds().Inset(-fangRadius, -fangRadius)
    if !pt.In(outer) {
        return false
    }
    dx, dy := e.Pos().Sub(pt).AsCoord()
    rx, ry := e.Radius()
    
    return (dx*dx)/(rx*rx) + (dy*dy)/(ry*ry) <= 1.0
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

// Polygone
type Polygon struct {
    Shape
    pts []geom.Point
    Closed bool
}

func NewPolygon(p0 geom.Point) (*Polygon) {
    p := &Polygon{}
    p.Wrapper = p
    p.Shape.Init()
    p.PropertyEmbed.InitByName("Polygon")
    p.SetPos(p0)
    p.pts = make([]geom.Point, 0)
    p.Closed = false
    return p
}

func (p *Polygon) Paint(gc *gg.Context) {
    gc.SetStrokeWidth(p.BorderWidth())
    gc.SetStrokeColor(p.BorderColor())
    for _, pt := range p.pts {
        gc.LineTo(pt.X, pt.Y)
    }
    if p.Closed {
        gc.LineTo(p.pts[0].X, p.pts[0].Y)
    }
    gc.Stroke()
/*
    gc.SetFillColor(color.Black)
    for _, pt := range p.pts {
        gc.DrawPoint(pt.X, pt.Y, 2.0)
    }
    gc.Fill()
*/
}

func (p *Polygon) Contains(pt geom.Point) bool {
    return false
}

func (p *Polygon) AddPoint(pt geom.Point) {
    lpt := pt.Sub(p.Pos())
    p.pts = append(p.pts, lpt)
    min := pt.Min(p.Pos())
    max := pt.Max(p.Pos().Add(p.Size()))
    p.SetMinSize(max.Sub(min))
}

func (p *Polygon) Flatten() {
    pts := make([]geom.Point, 1)
    p0 := p.pts[0]
    pts[0] = p0
    for _, p1 := range p.pts[1:] {
        if p0.Distance(p1) < 4.0 {
            continue
        }
        pts = append(pts, p1)
        p0 = p1
    }
    p.pts = pts
}

func (p *Polygon) Points() []geom.Point {
    pts := make([]geom.Point, len(p.pts))
    for i, pt := range p.pts {
        pts[i] = pt
    }
    return pts
}

//-----------------------------------------------------------------------------

/*
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
    c.PropertyEmbed.InitByName("Default")
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
*/
