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
    ShapeProps = NewProps(DefProps,
        map[ColorPropertyName]color.Color{
            Color:              color.Transparent,
            PressedColor:       DefProps.Color(Color).Alpha(0.5),
            BorderColor:        DefProps.Color(WhiteColor),
            PressedBorderColor: DefProps.Color(WhiteColor).Alpha(0.5),
        },
        nil,
        map[SizePropertyName]float64{
            BorderWidth:        1.0,
            PressedBorderWidth: 1.0,
        })
)

// Abstrakter, allgemeiner Typ fuer geometrische Formen
type Shape struct {
    LeafEmbed
    pushed bool
}

func (s *Shape) OnInputEvent(evt touch.Event) {
    switch evt.Type {
    case touch.TypePress:
        s.pushed = true
        s.Mark(MarkNeedsPaint)
    case touch.TypeRelease:
        s.pushed = false
        s.Mark(MarkNeedsPaint)
    }
    s.CallTouchFunc(evt)
}

func (s *Shape) Paint(gc *gg.Context) {
    Debugf("-")
    if s.pushed {
        Debugf("paint pushed")
        gc.SetStrokeWidth(s.PressedBorderWidth())
        gc.SetFillColor(s.PressedColor())
        gc.SetStrokeColor(s.PressedBorderColor())
    } else {
        Debugf("paint not pushed")
        gc.SetStrokeWidth(s.BorderWidth())
        gc.SetFillColor(s.Color())
        gc.SetStrokeColor(s.BorderColor())
    }
}

// Geraden
type Line struct {
    Shape
    slope float64
}

func NewLine(dx, dy, slope float64) (*Line) {
    l := &Line{}
    l.Wrapper = l
    l.Init()
    l.PropertyEmbed.Init(ShapeProps)
    l.SetMinSize(geom.Point{dx, dy})
    l.slope = slope
    return l
}

func (l *Line) Paint(gc *gg.Context) {
    Debugf("-")
    l.Shape.Paint(gc)
    if l.slope > 0.0 {
        gc.MoveTo(l.Bounds().SW().AsCoord())
        gc.LineTo(l.Bounds().NE().AsCoord())
    } else {
        gc.MoveTo(l.Bounds().Min.AsCoord())
        gc.LineTo(l.Bounds().Max.AsCoord())
    }
    gc.Stroke()
}

func (l *Line) Contains(pt geom.Point) (bool) {
    if !l.Embed.Contains(pt) {
        return false
    }
    fx, fy := l.Rect().PosRel(pt)
    //log.Printf("fx, fy: %f, %f", fx, fy)
    if l.slope > 0.0 && math.Abs((fx+fy) - 1.0) < 0.05 {
        return true
    }
    if l.slope < 0.0 && math.Abs(fx-fy) < 0.05 {
        return true
    }
    return false
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
    w := c.Size().X
    gc.DrawCircle(0.5*w, 0.5*w, 0.5*w)
    gc.FillStroke()
}

func (c *Circle) Contains(pt geom.Point) (bool) {
    if !c.Embed.Contains(pt) {
        return false
    }
    r := c.Radius()
    return pt.Dist2(c.Center()) <= r*r
}

func (c *Circle) Center() (geom.Point) {
    return c.Rect().Center()
}

func (c *Circle) SetCenter(p geom.Point) {
    c.SetPos(p.Sub(c.Size().Mul(0.5)))
}

func (c *Circle) Radius() (float64) {
    return 0.5 * c.Size().X
}

func (c *Circle) SetRadius(r float64) {
    mp := c.Center()
    c.SetMinSize(geom.Point{2*r, 2*r})
    c.SetPos(mp.Sub(c.Size().Mul(0.5)))
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
    w, h := e.Size().AsCoord()
    gc.DrawEllipse(0.5*w, 0.5*h, 0.5*w, 0.5*h)
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

func (e *Ellipse) Center() (geom.Point) {
    return e.Rect().Center()
}

func (e *Ellipse) SetCenter(p geom.Point) {
    e.SetPos(p.Sub(e.Size().Mul(0.5)))
}

func (e *Ellipse) Radius() (float64, float64) {
    return e.Size().Mul(0.5).AsCoord()
}

func (e *Ellipse) SetRadius(rx, ry float64) {
    mp := e.Center()
    e.SetMinSize(geom.Point{2*rx, 2*ry})
    e.SetPos(mp.Sub(e.Size().Mul(0.5)))
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

func (r *Rectangle) Paint(gc *gg.Context) {
    Debugf("-")
    r.Shape.Paint(gc)
    gc.DrawRectangle(r.Bounds().AsCoord())
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

