package main

import (
    "math/rand"
    "time"
    "golang.org/x/image/font"
    "github.com/stefan-muehlebach/gg"
    "github.com/stefan-muehlebach/gg/color"
    "github.com/stefan-muehlebach/gg/fonts"
)

type TextAnim struct {
    gc *gg.Context
	textList []*TextObject
	fontList []*fonts.Font
}

func (a *TextAnim) RefreshTime() time.Duration {
    return 30 * time.Millisecond
}

func (a *TextAnim) Init(gc *gg.Context) {
    a.gc = gc
	a.textList = make([]*TextObject, numObjs)
    a.fontList = []*fonts.Font{
		fonts.LucidaBright,
		fonts.LucidaBrightItalic,
		fonts.LucidaBrightDemibold,
		fonts.LucidaBrightDemiboldItalic,
    }
	for i := 0; i < numObjs; i++ {
		a.textList[i] = NewTextObject(msg, a.fontList[i%len(a.fontList)],
			45.0+80.0*rand.Float64())
		yPos := 200.0*rand.Float64() + 20.0
		xVel := 4.0*rand.Float64() + 1.0
		if rand.Float64() < 0.5 {
			xVel *= -1.0
		}
		a.textList[i].SetAnimParam(yPos, xVel)
	}
}

func (a *TextAnim) Paint() {
	a.gc.SetFillColor(color.Black)
	a.gc.Clear()

	for _, txtObj := range a.textList {
		txtObj.Draw(gc)
	}
	for _, txtObj := range a.textList {
		if !txtObj.Animate() {
			yPos := 200.0*rand.Float64() + 20.0
			xVel := 4.0*rand.Float64() + 1.0
			if rand.Float64() < 0.5 {
				xVel *= -1.0
			}
			txtObj.SetAnimParam(yPos, xVel)
		}
	}
}

func (a *TextAnim) Clean() {}

type TextObject struct {
	x, y          float64
	txt           string
	face          font.Face
	width, height float64
	Color         color.Color
	xVel          float64
}

func NewTextObject(txt string, fnt *fonts.Font, size float64) *TextObject {
	o := &TextObject{}
	o.txt = txt
	o.face = fonts.NewFace(fnt, size)
	o.Color = color.RandColor()
	o.width = float64(font.MeasureString(o.face, o.txt)) / 64.0
	o.height = float64(o.face.Metrics().Ascent) / 64.0
	return o
}

func (o *TextObject) SetAnimParam(y, xVel float64) {
	if xVel > 0.0 {
		o.x = -o.width / 2.0
	} else {
		o.x = float64(gc.Width()) + o.width/2.0
	}
	o.y = y
	o.xVel = xVel
}

func (o *TextObject) Draw(gc *gg.Context) {
	gc.SetFontFace(o.face)
	gc.SetStrokeColor(o.Color)
	gc.DrawStringAnchored(o.txt, o.x, o.y, 0.5, 0.5)
}

func (o *TextObject) Animate() bool {
	o.x += o.xVel
	if o.x > float64(gc.Width())+o.width/2.0 || o.x < -o.width/2.0 {
		return false
	}
	return true
}

