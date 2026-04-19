package main

import (
	"flag"
	"fmt"
	"image"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/stefan-muehlebach/adatft"
	"github.com/stefan-muehlebach/gg"
	"github.com/stefan-muehlebach/gg/geom"
)

//-----------------------------------------------------------------------------

var (
	Draw DrawFunc = DrawNormal
	movieTotalFrames, movieCurrentFrame int
)

type DrawFunc func(gc *gg.Context, disp *adatft.Display)

func DrawNormal(gc *gg.Context, disp *adatft.Display) {
	disp.Draw(gc.Image())
}

func DrawScreenshot(gc *gg.Context, disp *adatft.Display) {
	gc.SavePNG("images/screenshot.png")
	disp.Draw(gc.Image())
	Draw = DrawNormal
}

func DrawMovie(gc *gg.Context, disp *adatft.Display) {
	fileName := fmt.Sprintf("images/movie.%04d.png", movieCurrentFrame)
	gc.SavePNG(fileName)
	disp.Draw(gc.Image())
	movieCurrentFrame++
	if movieCurrentFrame >= movieTotalFrames {
		Draw = DrawNormal
	}
}

//-----------------------------------------------------------------------------

func SignalHandler() {
	sigChan := make(chan os.Signal)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGUSR1, syscall.SIGUSR2)
	for sig := range sigChan {
		switch sig {
		case os.Interrupt:
			runFlag = false
			quitFlag = true
			break
		case syscall.SIGUSR1:
			Draw = DrawScreenshot
		case syscall.SIGUSR2:
			movieTotalFrames = 150
			movieCurrentFrame = 0
			Draw = DrawMovie
		}
	}
}

func TouchHandler() {
	for penEvent := range touch.EventQ {
		//log.Printf("penEvent: %#v", penEvent)
		pt := geom.Point{penEvent.X, penEvent.Y}
		switch {
		case pt.In(prevRect):
			switch penEvent.Type {
			case adatft.PenPress, adatft.PenDrag:
				continue
			case adatft.PenRelease:
				animNum -= 1
				if animNum < 0 {
					animNum += len(AnimationList)
				}
			}
		case pt.In(quitRect):
			switch penEvent.Type {
			case adatft.PenPress, adatft.PenDrag:
				continue
			case adatft.PenRelease:
				quitFlag = true
			}
		case pt.In(nextRect):
			switch penEvent.Type {
			case adatft.PenPress, adatft.PenDrag:
				continue
			case adatft.PenRelease:
				animNum += 1
				if animNum >= len(AnimationList) {
					animNum %= len(AnimationList)
				}
			}
		default:
			AnimationList[animNum].animation.Handle(penEvent)
			continue
		}
		runFlag = false
	}
}

//-----------------------------------------------------------------------------

type Animation interface {
	RefreshTime() time.Duration
	Init(gc *gg.Context)
	Animate(dt time.Duration)
	Paint()
	Clean()
	Handle(ev adatft.PenEvent)
}

func ShowAnimation(gc *gg.Context, a Animation) {
	dt := a.RefreshTime()

	a.Init(gc)
	ticker := time.NewTicker(dt)
	defer ticker.Stop()
	for range ticker.C {
		if !runFlag {
			break
		}
		adatft.AnimWatch.Start()
		a.Animate(dt)
		adatft.AnimWatch.Stop()
		adatft.PaintWatch.Start()
		a.Paint()
		adatft.PaintWatch.Stop()
		Draw(gc, disp)
	}
	a.Clean()
}

type AnimationListType struct {
	description string
	animation   Animation
}

var (
	AnimationList = []AnimationListType{
		{"Introduction", NewIntroAnim()},
		{"Circle", &CircleAnim{}},
		{"Dancing Polygons", &PolygonAnim{}},
		{"Rotating Cube (3D)", &Cube3DAnim{}},
		{"Text on the run", &TextAnim{}},
		{"Plasma... some sort of", &PlasmaAnim{}},
		{"SBB (are you Swiss?)", &SBBAnim{}},
		{"Scrolling Text", &ScrollAnim{}},
		{"Using Pico-8 font",
			NewFixedFontAnim(image.Point{20, 100}, "Hello Pico-8 | HELLO PICO-8")},
	}
)

//-----------------------------------------------------------------------------

var (
	IntroText                           string = "Im Folgenden habe ich einige kleine Beispiele, Animationen oder Interaktionen zusammengestellt, um die Möglichkeiten des TFT-Displays mit Go zu demonstrieren Sämtliche Animationen werden direkt gerechnet. Die Beispiele laufen jeweils unbegrenzt, für den Wechsel zwischen den Beispielen, verwende die Pfeil-Buttons unten links und rechts."
	disp                                *adatft.Display
	touch                               *adatft.Touch
	gc                                  *gg.Context
	pageNum                             int
	animNum                             int
	numObjs                             = 10
	numEdges                            = 3
	blurFactor                          float64
	msg                                 string
	rotation                            adatft.RotationType = adatft.Rotate270
	runFlag, quitFlag                   bool

	prevRect = geom.NewRectangleWH(0, 4*320/5, 480/3, 320/5)
	quitRect = geom.NewRectangleWH(480/3, 4*320/5, 480/3, 320/5)
	nextRect = geom.NewRectangleWH(2*480/3, 4*320/5, 480/3, 320/5)
)

func main() {
	InitProfiling()

	flag.IntVar(&animNum, "anim", 0, "Start with a given animation")
	flag.Float64Var(&blurFactor, "blur", 1.0, "(Only for Anim 1) Blur factor [0,1] (1: no blur, 0: max blur).\nIn order to see something, choose a value < 0.1")
	flag.StringVar(&msg, "text", "Hello, world!", "Sample text")
	flag.Var(&rotation, "rotation", "Display rotation")
	flag.Parse()

	StartProfiling()

	log.Printf("> OpenDisplay()\n")
	disp = adatft.OpenDisplay(rotation)
	log.Printf(" > done\n")

	log.Printf("> OpenTouch()\n")
	touch = adatft.OpenTouch(rotation)
	log.Printf(" > done\n")

	log.Printf("> NewContext()\n")
	gc = gg.NewContext(adatft.Width, adatft.Height)
	log.Printf(" > done\n")

	go SignalHandler()
	go TouchHandler()

	quitFlag = false
	for !quitFlag {
		runFlag = true
		log.Printf("[%d] %s", animNum, AnimationList[animNum].description)
		ShowAnimation(gc, AnimationList[animNum].animation)
		adatft.PrintStat()
		adatft.ResetStat()
	}

	disp.Close()
	touch.Close()
	StopProfiling()
}
