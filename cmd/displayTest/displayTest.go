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
	"github.com/stefan-muehlebach/adagui"
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

func PointerHandler() {
	for pointerEvent := range pointer.EventQ {
		// log.Printf("pointerEvent: %#v", pointerEvent)
		pt := pointerEvent.Pos
		switch {
		case pt.In(prevRect):
			switch pointerEvent.Type {
			case adatft.PointerRelease:
				animNum -= 1
				if animNum < 0 {
					animNum += len(AnimationList)
				}
			default:
				continue
			}
		case pt.In(quitRect):
			switch pointerEvent.Type {
			case adatft.PointerRelease:
				quitFlag = true
			default:
				continue
			}
		case pt.In(nextRect):
			switch pointerEvent.Type {
			case adatft.PointerRelease:
				animNum += 1
				if animNum >= len(AnimationList) {
					animNum %= len(AnimationList)
				}
			default:
				continue
			}
		default:
			AnimationList[animNum].animation.Handle(pointerEvent)
			continue
		}
		runFlag = false
	}
}

//-----------------------------------------------------------------------------

type Animation interface {
	RefreshTime() time.Duration
	Init(gc *gg.Context, rect geom.Rectangle)
	Animate(dt time.Duration)
	Paint()
	Clean()
	Handle(ev adatft.PointerEvent)
}

func ShowAnimation(gc *gg.Context, rect geom.Rectangle, a Animation) {
	dt := a.RefreshTime()

	a.Init(gc, rect)
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
		pointer.Draw(gc)
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
	pointer                             *adagui.Pointer
	gc                                  *gg.Context
	pageNum                             int
	animNum                             int
	numObjs                             = 10
	numEdges                            = 3
	blurFactor                          float64
	msg                                 string
	rotation                            adatft.RotationType = adatft.Rotate000
	runFlag, quitFlag                   bool

	prevRect, quitRect, nextRect geom.Rectangle
)

func main() {
	flag.IntVar(&animNum, "anim", 0, "Start with a given animation")
	flag.Float64Var(&blurFactor, "blur", 1.0, "(Only for Anim 1) Blur factor [0,1] (1: no blur, 0: max blur).\nIn order to see something, choose a value < 0.1")
	flag.StringVar(&msg, "text", "Hello, world!", "Sample text")
	flag.Var(&rotation, "rotate", "Display rotation")
	flag.Parse()

	adagui.StartProfiling()

	log.Printf("> OpenDisplay()")
	disp = adatft.OpenDisplay(rotation)
//	trans := disp.Matrix()
	log.Printf(" > done")

	log.Printf("> OpenPointer()")
	pointer = adagui.OpenPointer()
	r := geom.Rectangle{Max: disp.DrawBounds().Size()}
	pointer.SetPosRange(r)
	pointer.SetWheelRange(0, 0, 30)
	pointer.StartEvents()
	log.Printf(" > done")

	log.Printf("> NewContext()")
//	sz := disp.Bounds().Size().Int()
	gc = disp.Canvas()
//	gc.SetMatrix(trans)
	log.Printf(" > done")

	dispB := disp.DispBounds()
	drawB := disp.DrawBounds()

	log.Printf("dispB: %v", dispB)
	log.Printf("drawB: %v", drawB)

	w := float64(adatft.Width)/3.0
	h := 64.0
	ypos := float64(adatft.Height)-h

	prevRect = geom.NewRectangleWH(0.0, ypos, w, h)
	quitRect = geom.NewRectangleWH(w, ypos, w, h)
	nextRect = geom.NewRectangleWH(2*w, ypos, w, h)

	go SignalHandler()
	go PointerHandler()

	quitFlag = false
	for !quitFlag {
		runFlag = true
		log.Printf("[%d] %s", animNum, AnimationList[animNum].description)
		ShowAnimation(gc, drawB, AnimationList[animNum].animation)
		adatft.PrintStat()
		adatft.ResetStat()
	}

	disp.Close()
	pointer.Close()
	adagui.StopProfiling()
}
