package main

import (
	"time"

	"github.com/stefan-muehlebach/gg"
	"github.com/stefan-muehlebach/gg/colors"
	"github.com/stefan-muehlebach/gg/fonts"
	"github.com/stefan-muehlebach/gg/geom"
	"golang.org/x/image/font"
)

var (
	colorList = []colors.Color{
		colors.DarkGreen,
		colors.DarkRed,
		colors.DarkMagenta,
		colors.DarkCyan,
	}

	introText = "By tapping in the shown rectangles, you can switch between the different animations or quit the application."
)

//----------------------------------------------------------------------------

type Rectangle struct {
	geom.Rectangle
	color colors.Color
	face  font.Face
	txt   string
}

func NewRectangle(rect geom.Rectangle, color colors.Color,
	face font.Face, txt string) *Rectangle {
	r := &Rectangle{}
	r.Rectangle = rect
	r.color = color
	r.face = face
	r.txt = txt
	return r
}

func (r *Rectangle) Paint(gc *gg.Context) {
	gc.SetStrokeWidth(5.0)
	gc.SetStrokeColor(r.color)
	gc.SetFillColor(r.color.Alpha(0.3))
	gc.DrawRectangle(r.Inset(3, 3).AsCoord())
	gc.FillStroke()
	gc.SetTextColor(r.color)
	gc.SetFontFace(r.face)
	c := r.Center()
	gc.DrawStringAnchored(r.txt, c.X, c.Y, 0.5, 0.5)
}

//----------------------------------------------------------------------------

type IntroAnim struct {
	gc            *gg.Context
	rectList      []*Rectangle
	face          font.Face
	txt           string
	x, y          float64
	width, height float64
	cnt           int
}

func NewIntroAnim() *IntroAnim {
	a := &IntroAnim{}
	a.rectList = make([]*Rectangle, 0)
	a.face = fonts.NewFace(fonts.LucidaBright, 20.)
	a.txt = introText
	return a
}

func (a *IntroAnim) RefreshTime() time.Duration {
	return time.Second
}

func (a *IntroAnim) Init(gc *gg.Context) {
	a.gc = gc
	rect := NewRectangle(geom.Rect(0, 0, 480/3, 320), colors.DarkGreen,
		fonts.NewFace(fonts.LucidaBrightDemibold, 40.0), "prev")
	a.rectList = append(a.rectList, rect)
	rect = NewRectangle(geom.Rect(480/3, 0, 2*480/3, 320/2), colors.DarkBlue,
		fonts.NewFace(fonts.LucidaBrightDemibold, 40.0), "extra")
	a.rectList = append(a.rectList, rect)
	rect = NewRectangle(geom.Rect(480/3, 320/2, 2*480/3, 320), colors.DarkRed,
		fonts.NewFace(fonts.LucidaBrightDemibold, 40.0), "quit")
	a.rectList = append(a.rectList, rect)
	rect = NewRectangle(geom.Rect(2*480/3, 0, 480, 320), colors.DarkCyan,
		fonts.NewFace(fonts.LucidaBrightDemibold, 40.0), "next")
	a.rectList = append(a.rectList, rect)
}

func (a *IntroAnim) Animate(dt time.Duration) {}

func (a *IntroAnim) Paint() {
	a.gc.SetFillColor(colors.Black)
	a.gc.Clear()
	for _, rect := range a.rectList {
		rect.Paint(a.gc)
	}
}

func (a *IntroAnim) Clean() {}
