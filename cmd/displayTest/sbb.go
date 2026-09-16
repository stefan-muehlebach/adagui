package main

import (
	"image"
	"log"
	"math"
	"time"

    "github.com/stefan-muehlebach/adatft"

	"github.com/stefan-muehlebach/gg"
	"github.com/stefan-muehlebach/gg/colors"
	"github.com/stefan-muehlebach/gg/geom"
)

type SBBAnim struct {
	gc           *gg.Context
	rect	geom.Rectangle
	imgList      []image.Image
	imgFileList  []string
	rotationList []float64
	xm, ym       float64
	c1, c2, c3   float64
}

func (a *SBBAnim) RefreshTime() time.Duration {
	return 30 * time.Millisecond
}

func (a *SBBAnim) Init(gc *gg.Context, rect geom.Rectangle) {
	var err error

	a.gc = gc
	a.rect = rect
	a.c1 = 1.0 / 60.0
	a.c2 = 1.0 / 12.0
	a.c3 = 2.0 * math.Pi

	a.imgFileList = []string{
		"sbb/dial.png",
		"sbb/hour.png",
		"sbb/minute.png",
		"sbb/second.png",
	}

	a.imgList = make([]image.Image, len(a.imgFileList))
	for i, fileName := range a.imgFileList {
		a.imgList[i], err = gg.LoadPNG(fileName)
		if err != nil {
			log.Fatalf("error loading image: %v", err)
		}
	}
	a.rotationList = make([]float64, len(a.imgFileList))
	a.xm = a.rect.Dx() / 2.0
	a.ym = a.rect.Dy() / 2.0

	a.gc.SetFillColor(colors.DeepSkyBlue)
	a.gc.Clear()
}

func (a *SBBAnim) Animate(dt time.Duration) {
	t := time.Now()
	seconds := float64(t.Second()) * a.c1
	minutes := (float64(t.Minute()) + seconds) * a.c1
	hours := (float64(t.Hour()) + minutes) * a.c2
	a.rotationList[0] = 0.0
	a.rotationList[1] = a.c3 * hours
	a.rotationList[2] = a.c3 * minutes
	a.rotationList[3] = a.c3 * seconds
}

func (a *SBBAnim) Paint() {
	a.gc.SetFillColor(colors.DeepSkyBlue)
	a.gc.Clear()
	for i, r := range a.rotationList {
		a.gc.Push()
		a.gc.RotateAbout(r, a.xm, a.ym)
		a.gc.DrawImageAnchored(a.imgList[i], a.xm, a.ym, 0.5, 0.5)
		a.gc.Pop()
	}
}

func (a *SBBAnim) Clean() {}

func (a *SBBAnim) Handle(evt adatft.PointerEvent) {}
