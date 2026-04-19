package main

import (
	"math/rand"
	"time"

    "github.com/stefan-muehlebach/adatft"

	"github.com/stefan-muehlebach/gg"
	"github.com/stefan-muehlebach/gg/colors"
)

// PolygonAnimation --
//
// Animation von halbtransparenten Polygonen. Die Anzahl der Polygone kann
// ueber das Flag 'objs' und die Anzahl Ecken der Polygone ueber das Flag
// 'edges' gesteuert werden.

type PolygonAnim struct {
	gc       *gg.Context
	polyList []*Polygon
}

func (a *PolygonAnim) RefreshTime() time.Duration {
	return 30 * time.Millisecond
}

func (a *PolygonAnim) Init(gc *gg.Context) {
	a.gc = gc

	a.polyList = make([]*Polygon, numObjs)
	for i := 0; i < numObjs; i++ {
		a.polyList[i] = NewPolygon(gc, numEdges)
	}

	a.gc.SetStrokeWidth(3)
	a.gc.SetLineCapRound()
	a.gc.SetLineJoinRound()
	a.gc.SetFillColor(colors.Black)
	a.gc.Clear()
}

func (a *PolygonAnim) Animate(dt time.Duration) {
	for _, p := range a.polyList {
		p.Animate(dt)
	}
}

func (a *PolygonAnim) Paint() {
	a.gc.SetFillColor(colors.RGBAF{0, 0, 0, blurFactor})
	a.gc.DrawRectangle(a.gc.Bounds().AsCoord())
	a.gc.Fill()

	for _, p := range a.polyList {
		p.Paint()
		//p.Move(gc.Bounds().AsCoord())
	}
}

func (a *PolygonAnim) Clean() {}

func (a *PolygonAnim) Handle(evt adatft.PenEvent) {}

type Polygon struct {
	gc *gg.Context
	xmin, ymin, xmax, ymax float64
	pts                      []*Point
	strokeColor, fillColor colors.RGBA
}

func NewPolygon(gc *gg.Context, edges int) *Polygon {
	p := &Polygon{}
	p.gc = gc
	p.xmin, p.ymin, p.xmax, p.ymax = gc.Bounds().AsCoord()
	p.pts = make([]*Point, edges)
	for i := 0; i < edges; i++ {
		pt := &Point{}
		pt.x = rand.Float64() * p.xmax
		pt.y = rand.Float64() * p.ymax
		pt.dx = rand.Float64()*5.0 - 2.0
		pt.dy = rand.Float64()*5.0 - 2.0
		p.pts[i] = pt
	}
	p.strokeColor = colors.White
	p.fillColor = colors.RandColor().Alpha(0.5)
	return p
}

func (p *Polygon) Animate(dt time.Duration) {
	for _, pt := range p.pts {
		pt.Move(p.xmin, p.xmax, p.ymin, p.ymax)
	}
}

func (p *Polygon) Paint() {
	p.gc.MoveTo(p.pts[0].x, p.pts[0].y)
	for _, pt := range p.pts[1:] {
		p.gc.LineTo(pt.x, pt.y)
	}
	p.gc.ClosePath()
	p.gc.SetStrokeStyle(gg.NewSolidPattern(p.strokeColor))
	p.gc.SetFillStyle(gg.NewSolidPattern(p.fillColor))
	p.gc.FillStroke()
}

type Point struct {
	x, y, dx, dy float64
}

func (p *Point) Move(xmin, xmax, ymin, ymax float64) {
	p.x += p.dx
	p.y += p.dy
	if p.x < xmin || p.x > xmax {
		p.dx *= -1
		p.x += p.dx
	}
	if p.y < ymin || p.y > ymax {
		p.dy *= -1
		p.y += p.dy
	}
}
