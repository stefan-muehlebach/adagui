package main

import (
	"flag"
	"fmt"
	"image"
	"image/draw"
	"log"
	"time"

	"github.com/stefan-muehlebach/adatft"
	"github.com/stefan-muehlebach/gg"
	"github.com/stefan-muehlebach/gg/colors"
	"github.com/stefan-muehlebach/gg/geom"
)

//----------------------------------------------------------------------------

const ()

var (
	pointColor  = colors.WhiteSmoke.Alpha(0.3)
	pointRadius = 10.0

	crossColor = colors.White
	crossSize  = 40.0

	gridColor     = colors.DarkGreen
	gridMargin    = 20.0
	gridSpace     = 40.0
	gridPointSize = 5.0
	gridWallSize  = 3.0

	touch        *adatft.Touch
	disp         *adatft.Display
	debug, nogui bool

	wallConfig = [][][]int{
		{{1, 0}, {6, 0}, {6, 1}, {7, 1}, {7, 10},
			{2, 10}, {2, 9}, {1, 9},
			{1, 2}, {4, 2}, {4, 3}, {5, 3},
			{5, 8}, {4, 8}, {4, 7}, {3, 7}, {3, 5}},
		{{6, 11}, {1, 11}, {1, 10}, {0, 10},
			{0, 1}, {5, 1}, {5, 2}, {6, 2},
			{6, 9}, {3, 9}, {3, 8}, {2, 8}, {2, 3},
			{3, 3}, {3, 4}, {4, 4}, {4, 6}},
	}
)

func printEvent(event adatft.PenEvent) {
	if !debug {
		return
	}
	fmt.Printf("[%d]: %10s: %v => %v\n",
		event.Time.UnixMilli(), event.Type, event.TouchRawPos,
		event.TouchPos)
}

func drawPoint(gc *gg.Context, x, y float64) {
	gc.SetFillColor(pointColor)
	gc.DrawPoint(x, y, pointRadius)
	gc.Fill()
}

func drawCross(gc *gg.Context, x, y float64) {
	gc.Clear()
	gc.DrawLine(x-crossSize/2, y, x+crossSize/2, y)
	gc.DrawLine(x, y-crossSize/2, x, y+crossSize/2)
	gc.Stroke()
}

func setupGrid(gc *gg.Context, actCol, actRow int) {
	// Clear the context and draw the grid first
	gc.SetFillColor(colors.Transparent)
	gc.Clear()
	gc.SetFillColor(gridColor)
	for y := gridMargin; y < float64(adatft.Height); y += gridSpace {
		for x := gridMargin; x < float64(adatft.Width); x += gridSpace {
			gc.DrawPoint(x, y, gridPointSize)
			gc.Fill()
		}
	}

	// The draw the walls in order to get a fine Irrgarten.
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

	// The 'Buttons' for Clear and Quit are drawn as the last part.
	gc.SetFillColor(colors.SteelBlue)
	gc.DrawCircle(0, 0, 35)
	gc.Fill()
	gc.SetFillColor(colors.GreenYellow)
	gc.DrawCircle(float64(adatft.Width), 0, 35)
	gc.Fill()
	gc.SetFillColor(colors.Gold)
	gc.DrawCircle(float64(adatft.Width), float64(adatft.Height), 35)
	gc.Fill()
	gc.SetFillColor(colors.OrangeRed)
	gc.DrawCircle(0, float64(adatft.Height), 35)
	gc.Fill()
}

func composeScreen(out *image.RGBA, grid, trace, cross *gg.Context) {
	draw.Draw(out, out.Bounds(), image.Black, image.Point{}, draw.Src)
	draw.Draw(out, out.Bounds(), grid.Image(), image.Point{}, draw.Over)
	draw.Draw(out, out.Bounds(), trace.Image(), image.Point{}, draw.Over)
	draw.Draw(out, out.Bounds(), cross.Image(), image.Point{}, draw.Over)
}

func main() {
	var rotation adatft.RotationType = adatft.Rotate090
	var grid, trace, cross *gg.Context

	flag.BoolVar(&debug, "debug", false, "write events to stdout")
	flag.BoolVar(&nogui, "nogui", false, "dont paint on the screen")
	flag.Var(&rotation, "rotation", "display rotation")
	flag.Parse()

	//adatft.Init()
	log.Printf("> OpenDisplay()\n")
	disp = adatft.OpenDisplay(rotation)
	log.Printf("> OpenTouch()\n")
	touch = adatft.OpenTouch(rotation)

	log.Printf("> NewContext()\n")
	grid = gg.NewContext(adatft.Width, adatft.Height)
	setupGrid(grid, 0, 0)

	trace = gg.NewContext(adatft.Width, adatft.Height)
	trace.SetFillColor(colors.Transparent)
	trace.Clear()
	trace.SetStrokeWidth(pointRadius)
	trace.SetStrokeColor(pointColor)
	trace.SetFillColor(pointColor)

	cross = gg.NewContext(adatft.Width, adatft.Height)
	cross.SetFillColor(colors.Transparent)
	cross.Clear()
	cross.SetStrokeWidth(2.0)
	cross.SetStrokeColor(crossColor)

	out := image.NewRGBA(disp.Bounds())

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

	quitPt := geom.NewPoint(0.0, 0.0)
	clearPt := geom.NewPoint(float64(adatft.Width), 0.0)
	color1Pt := geom.NewPoint(0.0, float64(adatft.Height))
	color2Pt := geom.NewPoint(float64(adatft.Width), float64(adatft.Height))

	for event := range touch.EventQ {
		printEvent(event)

		pt := geom.Point{event.X, event.Y}
		if quitPt.Distance(pt) <= 35.0 {
			break
		}
		if clearPt.Distance(pt) <= 35.0 {
			trace.SetFillColor(colors.Transparent)
			trace.Clear()
			trace.SetFillColor(pointColor)
			continue
		}
		if color1Pt.Distance(pt) <= 35.0 {
			pointColor = colors.OrangeRed.Alpha(0.2)
			pointRadius = 20
			continue
		}
		if color2Pt.Distance(pt) <= 35.0 {
			pointColor = colors.Gold.Alpha(0.7)
			pointRadius = 5
			continue
		}

		if nogui {
			continue
		}
		switch event.Type {
		case adatft.PenPress, adatft.PenDrag, adatft.PenRelease:
			drawPoint(trace, event.X, event.Y)
			drawCross(cross, event.X, event.Y)
		}
	}

	done <- true

	grid.SetFillColor(colors.Black)
	grid.Clear()
	trace.SetFillColor(colors.Black)
	trace.Clear()
	cross.SetFillColor(colors.Black)
	cross.Clear()
	composeScreen(out, grid, trace, cross)
	disp.Draw(out)

	disp.Close()
	touch.Close()

	adatft.PrintStat()
}
