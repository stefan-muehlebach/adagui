package main

import (
	"strings"
	"time"

    "github.com/stefan-muehlebach/adatft"

	"github.com/stefan-muehlebach/gg"
	"github.com/stefan-muehlebach/gg/colors"
	"github.com/stefan-muehlebach/gg/fonts"
	"github.com/stefan-muehlebach/gg/geom"
	"golang.org/x/image/font"
)

// Scroll-Text
const (
	textMargin  = 10.0
	fontSize    = 22.0
	lineSpacing = 1.3
)

var (
	fontIdx          = 0
	blindText string = "Lorem ipsum dolor sit amet, consetetur sadipscing elitr, sed diam nonumy eirmod tempor invidunt ut labore et dolore magna aliquyam erat, sed diam voluptua."
)

type ScrollAnim struct {
	gc        *gg.Context
	rect geom.Rectangle
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

func (a *ScrollAnim) Init(gc *gg.Context, rect geom.Rectangle) {
	a.gc = gc
	a.rect = rect
	a.Setup(a.rect)
	a.gc.SetFillColor(colors.Black)
	a.gc.SetTextColor(colors.White)
}

func (a *ScrollAnim) Animate(dt time.Duration) {
	if a.h -= 1.5; a.h < 0.0 {
		a.Setup(a.rect)
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

func (a *ScrollAnim) Handle(evt adatft.PointerEvent) {}

func (a *ScrollAnim) Setup(rect geom.Rectangle) {
	a.face, _ = fonts.NewFace(fontList[fontIdx], fontSize)
	a.gc.SetFontFace(a.face)
	a.x, a.y = textMargin, textMargin
	a.textWidth = rect.Dx() - 2*textMargin
	textList := gc.WordWrap(blindText, a.textWidth)
	a.text = strings.Join(textList, "\n")
	a.h1 = rect.Dy()
	_, a.h2 = a.gc.MeasureMultilineString(a.text, lineSpacing)
	a.h = a.h1 + a.h2 + textMargin
	a.scrollUp = !a.scrollUp
	fontIdx = (fontIdx + 1) % len(fontList)
}
