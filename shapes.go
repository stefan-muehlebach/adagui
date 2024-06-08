package adagui

import (
//    "container/list"
    "math"
//    "sync"
    "github.com/stefan-muehlebach/adagui/touch"
    "github.com/stefan-muehlebach/gg"
//    "github.com/stefan-muehlebach/gg/colornames"
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

// Kreis
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
    if c.Pushed() || c.Selected() {
        gc.SetStrokeColor(c.PushedBorderColor())
        gc.SetStrokeWidth(c.PushedBorderWidth())
        gc.StrokePreserve()
    }
    gc.SetFillColor(c.Color())
    gc.SetStrokeColor(c.BorderColor())
    gc.SetStrokeWidth(c.BorderWidth())
    gc.FillStroke()
}

func (c *Circle) Contains(pt geom.Point) (bool) {
    outer := c.ParentBounds().Inset(-fangRadius, -fangRadius)
    if !pt.In(outer) {
        return false
    }
    return c.Pos().Distance(pt) <= c.Radius()
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
    p.pts = make([]geom.Point, 1)
    p.pts[0] = geom.Point{}
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
    gc.SetFillColor(colornames.Black)
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
        //pts[i] = pt.Add(p.Pos())
    }
    return pts
}

// Geraden
type Line struct {
    Shape
    p0, p1 geom.Point
}

func NewLine(p0, p1 geom.Point) (*Line) {
    l := &Line{}
    l.Wrapper = l
    l.Shape.Init()
    l.PropertyEmbed.InitByName("Shape")
    l.SetPos(p0.Min(p1))
    l.p0 = geom.Point{}
    l.p1 = geom.Point{}
    l.SetP0(p0)
    l.SetP1(p1)
    return l
}

func (l *Line) P0() (geom.Point) {
    return l.Pos().Add(l.p0)
}
func (l *Line) SetP0(pt geom.Point) {
    pos := l.Pos()
    if pt.X < pos.X {
        pos.X = pt.X
    }
    if pt.Y < pos.Y {
        pos.Y = pt.Y
    }
    l.p1 = l.p1.Add(l.Pos().Sub(pos))
    l.p0 = pt.Sub(pos)
    l.SetPos(pos)
    l.SetMinSize(geom.Rect(l.p0.X, l.p0.Y, l.p1.X, l.p1.Y).Size())
}

func (l *Line) P1() (geom.Point) {
    return l.Pos().Add(l.p1)
}
func (l *Line) SetP1(pt geom.Point) {
    pos := l.Pos()
    if pt.X < pos.X {
        pos.X = pt.X
    }
    if pt.Y < pos.Y {
        pos.Y = pt.Y
    }
    l.p0 = l.p0.Add(l.Pos().Sub(pos))
    l.p1 = pt.Sub(pos)
    l.SetPos(pos)
    l.SetMinSize(geom.Rect(l.p0.X, l.p0.Y, l.p1.X, l.p1.Y).Size())
}    

func (l *Line) Paint(gc *gg.Context) {
    Debugf(Painting, "")
    gc.MoveTo(l.p0.AsCoord())
    gc.LineTo(l.p1.AsCoord())
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
    outer := l.ParentBounds().Inset(-fangRadius, -fangRadius)
    if !pt.In(outer) {
        return false
    }
    fx, fy := l.ParentBounds().PosRel(pt)
    Debugf(Coordinates, "fx, fy: %f, %f", fx, fy)
    return math.Abs(fx-fy) <= 0.1
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
    if e.Pushed() || e.Selected() {
        gc.SetStrokeWidth(e.PushedBorderWidth())
        gc.SetStrokeColor(e.PushedBorderColor())
        gc.StrokePreserve()
    }
    gc.SetFillColor(e.Color())
    gc.SetStrokeWidth(e.BorderWidth())
    gc.SetStrokeColor(e.BorderColor())
    gc.FillStroke()
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

// Und wo es Kreise gibt, da sind auch die Rechtecke nicht weit.
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

func (r *Rectangle) Contains(pt geom.Point) (bool) {
    outer := r.ParentBounds().Inset(-5.0, -5.0)
    //inner := r.ParentBounds().Inset(+5.0, +5.0)
    return pt.In(outer)
}

func (r *Rectangle) Paint(gc *gg.Context) {
    Debugf(Painting, "")
    gc.DrawRectangle(r.LocalBounds().AsCoord())
    if r.Pushed() || r.Selected() {
        gc.SetStrokeWidth(r.PushedBorderWidth())
        gc.SetStrokeColor(r.PushedBorderColor())
        gc.StrokePreserve()
    }
    gc.SetFillColor(r.Color())
    gc.SetStrokeWidth(r.BorderWidth())
    gc.SetStrokeColor(r.BorderColor())
    gc.FillStroke()
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
