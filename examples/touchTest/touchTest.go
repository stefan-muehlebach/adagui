package main

import (
	"flag"
	"fmt"
	"github.com/stefan-muehlebach/adatft"
	"github.com/stefan-muehlebach/gg"
	"github.com/stefan-muehlebach/gg/color"
	"github.com/stefan-muehlebach/gg/color"
	"github.com/stefan-muehlebach/gg/geom"
	"image"
	"image/draw"
	"time"
)

//----------------------------------------------------------------------------

const ()

var (
	pointColor  = color.WhiteSmoke.Alpha(0.5)
	pointRadius = 1.0

	crossColor = color.White
	crossSize  = 40.0

	gridColor     = color.DarkGreen
	gridMargin    = 20.0
	gridSpace     = 40.0
	gridPointSize = 5.0
	gridWallSize  = 3.0

	touch          *adatft.Touch
	disp           *adatft.Display
	verbose, nogui bool

	wallConfig = [][][]int{
		{{0, 0}, {2, 0}, {2, 1}, {5, 1}},
		{{0, 1}, {0, 3}, {2, 3}},
	}
)

func printEvent(event adatft.PenEvent) {
	if !verbose {
		return
	}
	fmt.Printf("[%d]: %10s: %v => %v (%d)\n",
		event.Time.UnixMilli(), event.Type, event.TouchRawPos,
		event.TouchPos, event.FifoSize)
}

func drawPoint(gc *gg.Context, x, y float64) {
	gc.DrawPoint(x, y, pointRadius)
	gc.SetFillColor(pointColor)
	gc.SetStrokeColor(pointColor)
	gc.FillStroke()
}

func drawCross(gc *gg.Context, x, y float64) {
	gc.Clear()
	gc.DrawLine(x-crossSize/2, y, x+crossSize/2, y)
	gc.DrawLine(x, y-crossSize/2, x, y+crossSize/2)
	gc.Stroke()
}

func setupGrid(gc *gg.Context, actCol, actRow int) {
	gc.SetFillColor(color.Transparent)
	gc.Clear()
	gc.SetFillColor(gridColor)
	for y := gridMargin; y < float64(adatft.Height); y += gridSpace {
		for x := gridMargin; x < float64(adatft.Width); x += gridSpace {
			gc.DrawPoint(x, y, gridPointSize)
			gc.Fill()
		}
	}

	gc.SetStrokeColor(gridColor)
	gc.SetStrokeWidth(gridWallSize)
	for _, wall := range wallConfig {
		for _, pos := range wall {
			x := 20 + 40*float64(pos[0])
			y := 20 + 40*float64(pos[1])
			gc.LineTo(x, y)
		}
		gc.Stroke()
	}

	gc.SetFillColor(color.SteelBlue)
	gc.MoveTo(0, 0)
	gc.LineTo(20, 0)
	gc.LineTo(0, 20)
	gc.LineTo(0, 0)
	gc.Fill()
}

func composeScreen(out *image.RGBA, grid, trace, cross *gg.Context) {
	draw.Draw(out, out.Bounds(), image.Black, image.Point{0, 0}, draw.Src)
	draw.Draw(out, out.Bounds(), trace.Image(), image.Point{0, 0}, draw.Over)
	draw.Draw(out, out.Bounds(), cross.Image(), image.Point{0, 0}, draw.Over)
	draw.Draw(out, out.Bounds(), grid.Image(), image.Point{0, 0}, draw.Over)
}

func main() {
	var rotation adatft.RotationType = adatft.Rotate000
	var grid, trace, cross *gg.Context

	flag.BoolVar(&verbose, "verbose", false, "write events to stdout")
	flag.BoolVar(&nogui, "nogui", false, "dont paint on the screen")
	flag.Var(&rotation, "rotation", "display rotation")
	flag.Parse()

	//adatft.Init()
	disp = adatft.OpenDisplay(rotation)

	touch = adatft.OpenTouch()
	touch.ReadConfig()

	fmt.Printf("Done opening TFT\n")

	grid = gg.NewContext(adatft.Width, adatft.Height)
	setupGrid(grid, 0, 0)

	trace = gg.NewContext(adatft.Width, adatft.Height)
	trace.SetFillColor(color.Transparent)
	trace.Clear()
	trace.SetStrokeWidth(1.0)
	trace.SetStrokeColor(pointColor)
	trace.SetFillColor(pointColor)

	cross = gg.NewContext(adatft.Width, adatft.Height)
	cross.SetFillColor(color.Transparent)
	cross.Clear()
	cross.SetStrokeWidth(2.0)
	cross.SetStrokeColor(crossColor)

	out := image.NewRGBA(image.Rect(0, 0, adatft.Width, adatft.Height))

	done := make(chan bool)
	ticker := time.NewTicker(30 * time.Millisecond)

	// Draw oder Paint Thread. Via Ticker zeitgesteuert (alle 30 ms)
	go func() {
		for {
			select {
			case <-ticker.C:
				composeScreen(out, grid, trace, cross)
				disp.Draw(out)
			case <-done:
				return
			}
		}
	}()

	quitPt := geom.NewPoint(20.0, 20.0)
	for event := range touch.EventQ {
		if quitPt.Distance(geom.Point{event.X, event.Y}) <= 10.0 {
			break
		}
		printEvent(event)
		if nogui {
			continue
		}
		switch event.Type {
		case adatft.PenPress, adatft.PenDrag:
			drawPoint(trace, event.X, event.Y)
			drawCross(cross, event.X, event.Y)
		}
	}

	done <- true

	grid.SetFillColor(color.Black)
	grid.Clear()
	trace.SetFillColor(color.Black)
	trace.Clear()
	cross.SetFillColor(color.Black)
	cross.Clear()
	composeScreen(out, grid, trace, cross)
	disp.Draw(out)

	disp.Close()
	touch.Close()

	adatft.PrintStat()
}
