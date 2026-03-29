package main

import (
	"math/rand"
	"time"

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
		a.polyList[i] = NewPolygon(gc.Width(), gc.Height(), numEdges)
	}

	a.gc.SetStrokeWidth(3)
	a.gc.SetLineCapRound()
	a.gc.SetLineJoinRound()
	a.gc.SetFillColor(colors.Black)
	a.gc.Clear()
}

func (a *PolygonAnim) Animate(dt time.Duration) {}

func (a *PolygonAnim) Paint() {
	a.gc.SetFillColor(colors.RGBAF{0, 0, 0, blurFactor})
	a.gc.DrawRectangle(a.gc.Bounds().AsCoord())
	a.gc.Fill()

	for _, p := range a.polyList {
		p.Draw(a.gc)
		p.Move(a.gc.Bounds().AsCoord())
	}
}

func (a *PolygonAnim) Clean() {}

type Polygon struct {
	p                      []*Point
	strokeColor, fillColor colors.Color
}

func NewPolygon(dispWidth, dispHeight, edges int) *Polygon {
	p := &Polygon{}
	p.p = make([]*Point, edges)
	for i := 0; i < edges; i++ {
		pt := &Point{}
		pt.x = rand.Float64() * float64(dispWidth)
		pt.y = rand.Float64() * float64(dispHeight)
		pt.dx = rand.Float64()*5.0 - 2.0
		pt.dy = rand.Float64()*5.0 - 2.0
		p.p[i] = pt
	}
	p.strokeColor = colors.White
	p.fillColor = colors.RandColor().Alpha(0.5)
	return p
}

func (p *Polygon) Move(xmin, ymin, xmax, ymax float64) {
	for _, p := range p.p {
		p.Move(xmin, xmax, ymin, ymax)
	}
}

func (p *Polygon) Draw(gc *gg.Context) {
	gc.MoveTo(p.p[0].x, p.p[0].y)
	for _, p := range p.p[1:] {
		gc.LineTo(p.x, p.y)
	}
	gc.ClosePath()
	gc.SetStrokeStyle(gg.NewSolidPattern(p.strokeColor))
	gc.SetFillStyle(gg.NewSolidPattern(p.fillColor))
	gc.FillStroke()
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
