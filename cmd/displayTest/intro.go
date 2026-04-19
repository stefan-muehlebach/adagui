package main

import (
	"time"

    "github.com/stefan-muehlebach/adatft"

	"github.com/stefan-muehlebach/gg"
	"github.com/stefan-muehlebach/gg/colors"
	"github.com/stefan-muehlebach/gg/fonts"
	"github.com/stefan-muehlebach/gg/geom"
	"golang.org/x/image/font"
)

var (
	colorList = []colors.RGBA{
		colors.DarkGreen,
		colors.DarkRed,
		colors.DarkMagenta,
		colors.DarkCyan,
	}

	introText = "By tapping in the shown rectangles, you can switch between the different animations or quit the application."
)

//----------------------------------------------------------------------------

const (
	StrokeWidth = 8.0
)

var (
	PrevRect, QuitRect, NextRect *Rectangle
)

type Rectangle struct {
	geom.Rectangle
	color colors.RGBA
	face  font.Face
	txt   string
}

func NewRectangle(rect geom.Rectangle, color colors.RGBA,
	face font.Face, txt string) *Rectangle {
	r := &Rectangle{}
	r.Rectangle = rect
	r.color = color
	r.face = face
	r.txt = txt
	return r
}

func (r *Rectangle) Paint(gc *gg.Context) {
	gc.SetStrokeWidth(StrokeWidth)
	gc.SetStrokeColor(r.color)
	gc.SetFillColor(r.color.Alpha(0.3))
	gc.DrawRectangle(r.Inset(StrokeWidth/2, StrokeWidth/2).AsCoord())
	gc.FillStroke()
	gc.SetTextColor(r.color)
	gc.SetFontFace(r.face)
	c := r.C()
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
	a.face, _ = fonts.NewFace(fonts.LucidaBright, 20.0)
	a.txt = introText
	return a
}

func (a *IntroAnim) RefreshTime() time.Duration {
	return time.Second
}

func (a *IntroAnim) Init(gc *gg.Context) {
	a.gc = gc
	face, _ := fonts.NewFace(fonts.LucidaBrightDemiboldItalic, 26.0)
	rect := NewRectangle(prevRect, colors.DarkGreen, face, "Prev")
	a.rectList = append(a.rectList, rect)
//	rect = NewRectangle(geom.NewRectangleWH(480/5, 0, 3*480/5, 320/5),
//		colors.DarkBlue, face, "extra")
//	a.rectList = append(a.rectList, rect)
	rect = NewRectangle(quitRect, colors.DarkRed, face, "Quit")
	a.rectList = append(a.rectList, rect)
	rect = NewRectangle(nextRect, colors.DarkCyan, face, "Next")
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

func (a *IntroAnim) Handle(evt adatft.PenEvent) {}
