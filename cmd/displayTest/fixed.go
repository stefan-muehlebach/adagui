package main

import (
	"image"
	"image/color"
	"time"

	"github.com/stefan-muehlebach/gg"
	"golang.org/x/image/draw"
	"golang.org/x/image/font"
	"golang.org/x/image/font/basicfont"
	"golang.org/x/image/math/fixed"
)

//-----------------------------------------------------------------------------

type FixedFontAnim struct {
	gc       *gg.Context
	drawer   *font.Drawer
	faceList []*basicfont.Face
	text     string
	col      color.Color
	pt       fixed.Point26_6
	lh       fixed.Int26_6
}

func NewFixedFontAnim(pt image.Point, text string) *FixedFontAnim {
	a := &FixedFontAnim{}
	a.drawer = &font.Drawer{}
	a.faceList = []*basicfont.Face{Pico3x5, Pico6x10, Pico9x15, Pico12x20}
	a.text = text
	a.col = color.White
	a.pt = fixed.P(pt.X, pt.Y)
	return a
}

func (a *FixedFontAnim) RefreshTime() time.Duration {
	return time.Second
}

func (a *FixedFontAnim) Init(gc *gg.Context) {
	a.gc = gc
	a.drawer.Dst = gc.Image().(draw.Image)
	a.drawer.Src = image.NewUniform(a.col)
	a.gc.SetFillColor(color.Black)
	a.gc.Clear()
}

func (a *FixedFontAnim) Animate(dt time.Duration) {}

func (a *FixedFontAnim) Paint() {
	a.drawer.Dot = a.pt
	for _, face := range a.faceList {
		a.drawer.Face = face
		a.lh = fixed.I(face.Height).Mul(fixed.I(2))
		a.drawer.DrawString(a.text)
		a.drawer.Dot.X = a.pt.X
		a.drawer.Dot.Y += a.lh
	}
}

func (a *FixedFontAnim) Clean() {
}
