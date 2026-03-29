package main

import (
	"strings"
	"time"

	"github.com/stefan-muehlebach/gg"
	"github.com/stefan-muehlebach/gg/colors"
	"github.com/stefan-muehlebach/gg/fonts"
	"golang.org/x/image/font"
)

// Scroll-Text
const (
	textMargin  = 25.0
	fontSize    = 28.0
	lineSpacing = 1.5
)

var (
	fontIdx          = 0
	blindText string = "Lorem ipsum dolor sit amet, consetetur sadipscing elitr, sed diam nonumy eirmod tempor invidunt ut labore et dolore magna aliquyam erat, sed diam voluptua."
)

type ScrollAnim struct {
	gc        *gg.Context
	face      font.Face
	x, y      float64
	text      string
	textWidth float64
	h1, h2, h float64
	scrollUp  bool
}

func (a *ScrollAnim) RefreshTime() time.Duration {
	return 30 * time.Millisecond
}

func (a *ScrollAnim) Init(gc *gg.Context) {
	a.gc = gc
	a.Setup()
	a.gc.SetFillColor(colors.Black)
	a.gc.SetTextColor(colors.White)
}

func (a *ScrollAnim) Animate(dt time.Duration) {
	if a.h -= 1.5; a.h < 0.0 {
		a.Setup()
	}
	if a.scrollUp {
		a.y = a.h - a.h2
	} else {
		a.y = -(a.h - a.h1)
	}
}

func (a *ScrollAnim) Paint() {
	a.gc.Clear()
	a.gc.DrawStringWrapped(a.text, a.x, a.y, 0, 0, a.textWidth, lineSpacing,
		gg.AlignLeft)
}

func (a *ScrollAnim) Clean() {}

func (a *ScrollAnim) Setup() {
	a.face = fonts.NewFace(fontList[fontIdx], fontSize)
	a.gc.SetFontFace(a.face)
	a.x, a.y = textMargin, textMargin
	a.textWidth = float64(a.gc.Width()) - 2*textMargin
	textList := gc.WordWrap(blindText, a.textWidth)
	a.text = strings.Join(textList, "\n")
	a.h1 = float64(a.gc.Height())
	_, a.h2 = a.gc.MeasureMultilineString(a.text, lineSpacing)
	a.h = a.h1 + a.h2 + textMargin
	a.scrollUp = !a.scrollUp
	fontIdx = (fontIdx + 1) % len(fontList)
}
