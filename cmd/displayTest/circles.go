package main

import (
	"math/rand"
	"time"

	"github.com/stefan-muehlebach/adatft"
	"github.com/stefan-muehlebach/gg"
	"github.com/stefan-muehlebach/gg/colors"
	"github.com/stefan-muehlebach/gg/geom"
)

// CircleAnim --
//
const (
	numCircles = 5
)

type CircleAnim struct {
	gc       *gg.Context
	rect geom.Rectangle
	circleList []*Circle
}

func (a *CircleAnim) RefreshTime() time.Duration {
	return 70 * time.Millisecond
}

func (a *CircleAnim) Init(gc *gg.Context, rect geom.Rectangle) {
	a.gc = gc
	a.rect = rect
	a.circleList = make([]*Circle, numCircles)
	for i := range numCircles {
		a.circleList[i] = NewCircle(a)
		a.circleList[i].age = rand.Float64()
	}
	a.gc.SetLineWidth(3.0)
	a.gc.SetLineCapRound()
	a.gc.SetLineJoinRound()
	a.gc.SetFillColor(colors.Black)
	a.gc.Clear()
}

func (a *CircleAnim) Animate(dt time.Duration) {
	for _, c := range a.circleList {
		if ! c.Animate(dt) {
			c.Init()
		}
	}
}

func (a *CircleAnim) Paint() {
	a.gc.SetFillColor(colors.Black)
	a.gc.Clear()

	for _, c := range a.circleList {
		c.Paint()
	}
}

func (a *CircleAnim) Clean() {}

func (a *CircleAnim) Handle(evt adatft.PointerEvent) {
	switch evt.Type {
	case adatft.PointerPress, adatft.PointerDrag:
	case adatft.PointerRelease:
		a.circleList[0].Init()
		a.circleList[0].mx = float64(evt.Pos.X)
		a.circleList[0].my = float64(evt.Pos.Y)
	}
}

// Circle --
//
const (
	circleRatio = 0.35
	numWaves    = 5
	waveStep    = 10.0
)

type Circle struct {
	anim *CircleAnim
	mx, my, rx, ry, drx, dry, age, dAge, t float64
	color colors.RGBA
}

func NewCircle(anim *CircleAnim) *Circle {
	c := &Circle{}
	c.anim = anim
	c.Init()
	return c
}

func (c *Circle) Init() {
	c.mx = rand.Float64() * float64(c.anim.rect.Dx())
	c.my = float64(c.anim.rect.Dy()/4.0) + rand.Float64() * float64(c.anim.rect.Dy()/2.0)
	c.rx, c.ry = 0.0, 0.0
	c.drx = 0.1 * rand.NormFloat64() + 0.5
	c.dry = c.drx * circleRatio
	c.age = 1.0
	c.dAge = 0.0005 * rand.NormFloat64() + 0.002
	c.t = 0.0
	c.color = colors.RandColorByGroup(colors.Blues)
}

func (c *Circle) Animate(dt time.Duration) bool {
	c.age -= c.dAge
	if c.age <= 0.0 {
		c.age = 0.0
		return false
	}
	c.rx += c.drx
	c.ry += c.dry
	c.t = 1.0 - c.age
	return true
}

func (c *Circle) Paint() {
	c.anim.gc.SetLineColor(c.color.Alpha(1.0 - c.t*c.t))
	rx, ry := c.rx, c.ry
	for range numWaves {
		if rx <= 0.0 {
			break
		}
	    c.anim.gc.DrawEllipse(c.mx, c.my, rx, ry)
	    c.anim.gc.Stroke()
		rx = rx - waveStep
		ry = ry - circleRatio*waveStep
	}
}

