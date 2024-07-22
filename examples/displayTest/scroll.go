package main

import (
    "strings"
    "time"
    "golang.org/x/image/font"
    "github.com/stefan-muehlebach/gg"
    "github.com/stefan-muehlebach/gg/color"
    "github.com/stefan-muehlebach/gg/fonts"
)

// Scroll-Text
const (
	textMargin  = 20.0
	fontSize    = 28.0
	lineSpacing = 1.3
)

var (
	fontList = []*fonts.Font{
		fonts.GoRegular,
		fonts.Seaford,
		fonts.LucidaBright,
		fonts.LucidaSans,
        fonts.LucidaHandwritingItalic,
        fonts.LucidaCalligraphy,
		fonts.Garamond,
	}
    fontIdx = 0
    blindText string = "Lorem ipsum dolor sit amet, consetetur sadipscing elitr, sed diam nonumy eirmod tempor invidunt ut labore et dolore magna aliquyam erat, sed diam voluptua. At vero eos et accusam et justo duo dolores et ea rebum."
)

type ScrollAnim struct {
    gc   *gg.Context
    face font.Face
    x, y float64
    text string
    textWidth float64
    h1, h2, h float64
    scrollUp bool
}

func (a *ScrollAnim) RefreshTime() time.Duration {
    return 30 * time.Millisecond
}

func (a *ScrollAnim) Init(gc *gg.Context) {
    a.gc = gc
    a.Setup()
    a.gc.SetFillColor(color.Black)
    a.gc.SetStrokeColor(color.White)
}

func (a *ScrollAnim) Paint() {
    a.gc.Clear()
    if a.scrollUp {
        a.y = a.h - a.h2
    } else {
        a.y = -(a.h - a.h1)
    }
    a.gc.DrawStringWrapped(a.text, a.x, a.y, 0, 0, a.textWidth, lineSpacing,
        gg.AlignLeft)
    if a.h -= 1.5; a.h < 0.0 {
        a.Setup()
    }
}

func (a *ScrollAnim) Clean() {}

func (a *ScrollAnim) Setup() {
    a.face = fonts.NewFace(fontList[fontIdx], fontSize)
    a.gc.SetFontFace(a.face)
    a.x, a.y = textMargin, 0.0
    a.textWidth = float64(a.gc.Width()) - 2*a.x
    textList := gc.WordWrap(blindText, a.textWidth)
    a.text = strings.Join(textList, "\n")
    a.h1 = float64(a.gc.Height())
    _, a.h2 = a.gc.MeasureMultilineString(a.text, lineSpacing)
    a.h = a.h1 + a.h2
    a.scrollUp = !a.scrollUp
    fontIdx = (fontIdx + 1) % len(fontList)
}


/*
func ScrollingText() {
	var face font.Face

	for _, font := range fontList {
		if quitFlag {
			break
		}
		face = fonts.NewFace(font, fontSize)
		ScrollText(BlindText, face, textMargin, 0.0, lineSpacing, true)
		runFlag = true
	}
}

func FadeText(gc *gg.Context, dsp *adatft.Display,
	txt string, face font.Face, x, y, lineSpacing float64, fadeIn bool) {
}

func ScrollText(txt string, face font.Face, x, y, lineSpacing float64,
	scrollUp bool) {
	var textList []string
	var textWidth float64
	var h1, h2, h float64
	var ticker *time.Ticker

	gc.SetFontFace(face)
	textWidth = float64(gc.Width()) - 2*x
	textList = gc.WordWrap(txt, textWidth)
	txt = strings.Join(textList, "\n")

	h1 = float64(gc.Height())
	_, h2 = gc.MeasureMultilineString(txt, lineSpacing)
	h = h1 + h2

	ticker = time.NewTicker(30 * time.Millisecond)
	for range ticker.C {
		if !runFlag {
			break
		}
		gc.SetFillColor(color.Black)
		gc.Clear()
		gc.SetStrokeColor(color.White)
		if scrollUp {
			y = h - h2
		} else {
			y = -(h - h1)
		}
		gc.DrawStringWrapped(txt, x, y, 0, 0, textWidth,
			lineSpacing, gg.AlignLeft)
		Draw(gc, disp)
		if h -= 1.0; h < 0.0 {
			break
		}
	}
}
*/
