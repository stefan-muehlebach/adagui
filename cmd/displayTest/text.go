package main

import (
	"math/rand"
	"time"

	"github.com/stefan-muehlebach/gg"
	"github.com/stefan-muehlebach/gg/colors"
	"github.com/stefan-muehlebach/gg/fonts"
	"golang.org/x/image/font"
)

var (
	numTextLines = 6
	minFontSize  = 70.0
	maxFontSize  = 130.0
	minVel       = 3.0
	maxVel       = 5.0
	fontList     = []*fonts.Font{
		fonts.LucidaBright,
		fonts.LucidaSans,
		fonts.LucidaSansTypewriter,
		fonts.LucidaFax,
		fonts.LucidaConsole,
		fonts.LucidaHandwritingItalic,
		fonts.LucidaCalligraphy,
		fonts.LucidaBlackletter,
	}

	meanFontSize   = (minFontSize + maxFontSize) / 2.0
	stddevFontSize = (maxFontSize - minFontSize) / 2.0
	meanVel        = (minVel + maxVel) / 2.0
	stddevVel      = (maxVel - minVel) / 2.0
)

func normRand(mean, stddev float64) float64 {
	return rand.NormFloat64()*stddev + mean
}
func uniRand(minVal, maxVal float64) float64 {
	return minVal + rand.Float64()*(maxVal-minVal)
}

type TextAnim struct {
	gc       *gg.Context
	textList []*TextObject
	fontList []*fonts.Font
}

func (a *TextAnim) RefreshTime() time.Duration {
	return 30 * time.Millisecond
}

func (a *TextAnim) Init(gc *gg.Context) {
	a.gc = gc
	a.textList = make([]*TextObject, numTextLines)
	a.fontList = fontList
	for i := range numTextLines {
		t := float64(i) / float64(numTextLines-1)
		a.textList[i] = NewTextObject(msg,
			a.fontList[i%len(a.fontList)],
			normRand(meanFontSize, stddevFontSize),
			colors.RandColor().Alpha(1.0-t*0.5))
	}
}

func (a *TextAnim) Animate(dt time.Duration) {
	for _, txtObj := range a.textList {
		if !txtObj.Animate(1.0) {
			yPos := uniRand(20.0, float64(gc.Height())-20.0)
			xVel := normRand(meanVel, stddevVel)
			if rand.Float64() < 0.5 {
				xVel *= -1.0
			}
			txtObj.SetAnimParam(yPos, xVel)
		}
	}
}

func (a *TextAnim) Paint() {
	a.gc.SetFillColor(colors.Black)
	a.gc.Clear()

	for _, txtObj := range a.textList {
		txtObj.Paint(gc)
	}
}

func (a *TextAnim) Clean() {}

type TextObject struct {
	x, y          float64
	txt           string
	face          font.Face
	width, height float64
	color         colors.Color
	xVel, yVel    float64
}

func NewTextObject(txt string, fnt *fonts.Font, fontSize float64,
	color colors.Color) *TextObject {
	o := &TextObject{}
	o.txt = txt
	o.face = fonts.NewFace(fnt, fontSize)
	o.color = color
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

func (o *TextObject) Animate(t float64) bool {
	if o.xVel == 0.0 {
		return false
	}
	o.x += t * o.xVel
	if o.x > float64(gc.Width())+o.width/2.0 || o.x < -o.width/2.0 {
		return false
	}
	return true
}

func (o *TextObject) Paint(gc *gg.Context) {
	gc.SetFontFace(o.face)
	gc.SetTextColor(o.color)
	gc.DrawStringAnchored(o.txt, o.x, o.y, 0.5, 0.5)
}
