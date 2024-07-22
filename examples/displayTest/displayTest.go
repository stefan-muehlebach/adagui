package main

import (
	"flag"
	"fmt"
	"github.com/stefan-muehlebach/adatft"
	"github.com/stefan-muehlebach/gg"
	"log"
	"os"
	"os/signal"
	"path"
	"runtime"
	"runtime/pprof"
	"runtime/trace"
	"syscall"
	"time"
)

//-----------------------------------------------------------------------------

var (
	doCpuProf, doMemProf, doTrace       bool
	cpuProfFile, memProfFile, traceFile string
	fhCpu, fhMem, fhTrace               *os.File
)

func init() {
	cpuProfFile = fmt.Sprintf("%s.cpuprof", path.Base(os.Args[0]))
	memProfFile = fmt.Sprintf("%s.memprof", path.Base(os.Args[0]))
	traceFile = fmt.Sprintf("%s.trace", path.Base(os.Args[0]))

	flag.BoolVar(&doCpuProf, "cpuprof", false,
		"write cpu profile to "+cpuProfFile)
	flag.BoolVar(&doMemProf, "memprof", false,
		"write memory profile to "+memProfFile)
	flag.BoolVar(&doTrace, "trace", false,
		"write trace data to "+traceFile)
}

func StartProfiling() {
	var err error

	if doCpuProf {
		fhCpu, err = os.Create(cpuProfFile)
		if err != nil {
			log.Fatal("couldn't create cpu profile: ", err)
		}
		err = pprof.StartCPUProfile(fhCpu)
		if err != nil {
			log.Fatal("couldn't start cpu profiling: ", err)
		}
	}

	if doMemProf {
		fhMem, err = os.Create(memProfFile)
		if err != nil {
			log.Fatal("couldn't create memory profile: ", err)
		}
	}

	if doTrace {
		fhTrace, err = os.Create(traceFile)
		if err != nil {
			log.Fatal("couldn't create tracefile: ", err)
		}
		trace.Start(fhTrace)
	}
}

func StopProfiling() {
	if fhCpu != nil {
		pprof.StopCPUProfile()
		fhCpu.Close()
	}

	if fhMem != nil {
		runtime.GC()
		err := pprof.WriteHeapProfile(fhMem)
		if err != nil {
			log.Fatal("couldn't write memory profile: ", err)
		}
		fhMem.Close()
	}

	if fhTrace != nil {
		trace.Stop()
		fhTrace.Close()
	}
}

//-----------------------------------------------------------------------------

type DrawFunc func(gc *gg.Context, disp *adatft.Display)

var Draw DrawFunc = DrawNormal

func DrawNormal(gc *gg.Context, disp *adatft.Display) {
	disp.Draw(gc.Image())
}

func DrawScreenshot(gc *gg.Context, disp *adatft.Display) {
	gc.SavePNG("images/screenshot.png")
	disp.Draw(gc.Image())
	Draw = DrawNormal
}

func DrawMovie(gc *gg.Context, disp *adatft.Display) {
	fileName := fmt.Sprintf("images/movie.%03d.png", movieCurrentFrame)
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
		if penEvent.Type == adatft.PenRelease {
            if penEvent.X > float64(gc.Width())/2.0 {
		        animNum += 1
                if animNum >= len(AnimationList) {
                    animNum %= len(AnimationList)
                }
            } else {
		        animNum -= 1
                if animNum < 0 {
                    animNum += len(AnimationList)
                }
            }
			runFlag = false
		}
	}
}

//-----------------------------------------------------------------------------

type Animation interface {
    RefreshTime() time.Duration
	Init(gc *gg.Context)
	Paint()
    Clean()
}

func ShowAnimation(a Animation) {
	a.Init(gc)
	ticker := time.NewTicker(a.RefreshTime())
	defer ticker.Stop()
	for range ticker.C {
		if !runFlag {
			break
		}
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
		{"Dancing Polygons", &PolygonAnim{}},
		{"Rotating Cube (3D)", &Cube3DAnim{}},
		{"Text on the run", &TextAnim{}},
		{"Plasma... some sort of", &PlasmaAnim{}},
		{"SBB (are you Swiss?)", &SBBAnim{}},
        {"Scrolling Text", &ScrollAnim{}},
	}
)

//-----------------------------------------------------------------------------

/*
type displayPageType struct {
	description string
	pageFunc    func()
}

var (
	displayPageList = []displayPageType{
		{"Dancing Polygons", PolygonAnimation},
		{"Rotating Cube (3D)", Cube3DAnimation},
		{"Text on the run", TextAnimation},
		// {"Beziers wherever you look", BezierAnimation},
		// {"Let's fade the colors", FadingColors},
		{"Plasma (dont burn yourself!)", PlasmaAnimation},
		{"Fading Circles", CircleAnimation},
		{"SBB (are you Swiss?)", SBBAnimation},
		{"Scrolling Text", ScrollingText},
		// {"Matrix Tests", MatrixTest},
	}
)
*/

var (
	IntroText                           string = "Im Folgenden habe ich einige kleine Beispiele, Animationen oder Interaktionen zusammengestellt, um die Möglichkeiten des TFT-Displays mit Go zu demonstrieren Sämtliche Animationen werden direkt gerechnet. Die Beispiele laufen jeweils unbegrenzt, für den Wechsel zwischen den Beispielen, verwende die Pfeil-Buttons unten links und rechts."
	disp                                *adatft.Display
	touch                               *adatft.Touch
	gc                                  *gg.Context
	pageNum                             int
	animNum                             int
	numObjs, numEdges                   int
	blurFactor                          float64
	msg                                 string
	rotation                            adatft.RotationType = adatft.Rotate090
	runFlag, quitFlag                   bool
	movieTotalFrames, movieCurrentFrame int
)

func main() {
	flag.IntVar(&animNum, "anim", 0, "Start with a given animation")
	flag.IntVar(&numObjs, "objs", 5, "Number of objects")
	flag.IntVar(&numEdges, "edges", 3, "Number of edges of an object")
	flag.Float64Var(&blurFactor, "blur", 1.0, "(Only for Anim 1) Blur factor [0,1] (1: no blur, 0: max blur).\nIn order to see something, choose a value < 0.1")
	flag.StringVar(&msg, "text", "Hello, world!", "The text that will be displayed in animation 3")
	flag.Var(&rotation, "rotation", "Display rotation")
	flag.Parse()

	StartProfiling()

	disp = adatft.OpenDisplay(rotation)
	touch = adatft.OpenTouch(rotation)
	gc = gg.NewContext(adatft.Width, adatft.Height)

	go SignalHandler()
	go TouchHandler()

	quitFlag = false
	for !quitFlag {
		runFlag = true
		log.Printf("[%d] %s", animNum, AnimationList[animNum].description)
		ShowAnimation(AnimationList[animNum].animation)
	}

	disp.Close()
	touch.Close()
	StopProfiling()

	adatft.PrintStat()
}
