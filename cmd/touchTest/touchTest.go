package main

import (
	"flag"
	"image"
	"image/draw"
	"log"
	"math"
	"time"

	"github.com/stefan-muehlebach/adatft"
	"github.com/stefan-muehlebach/adagui"
	"github.com/stefan-muehlebach/gg"
	"github.com/stefan-muehlebach/gg/colors"
	"github.com/stefan-muehlebach/gg/fonts"
	"github.com/stefan-muehlebach/gg/geom"
)

//----------------------------------------------------------------------------

const ()

var (
	pointColors = []colors.RGBA{
		colors.WhiteSmoke.Alpha(0.3),
		colors.Gold.Alpha(0.7),
		colors.OrangeRed.Alpha(0.2),
	}
	pointRadiae = []float64{
		10.0,
		5.0,
		20.0,
	}
	pointIdx    = 0

	crossColor     = colors.White
	crossLineWidth = 2.0
	crossSize      = 40.0

	gridColor     = colors.WhiteSmoke
	gridMargin    = 20.0
	gridSpace     = 40.0
	gridPointSize = 5.0
	gridWallSize  = 3.0

	font          = fonts.LucidaBrightDemibold
	fontColor     = colors.WhiteSmoke
	fontSize      = 22.0
	fontFace, _   = fonts.NewFace(font, fontSize)

	btnWidth       = 80.0
	btnHeight      = 40.0
	btnBorderWidth = 6.0
	btnFont        = fonts.SeafordBold
	btnFontSize    = 16.0
	btnColors      = []colors.RGBA{
		colors.DeepSkyBlue,
		colors.LawnGreen,
		colors.Gold,
		colors.Red,
	}
	btnTitles = []string{
		"Exit",
		"Clear",
		"Gold",
		"Red",
	}
	btnRects       = []geom.Rectangle{
		geom.NewRectangleWH(0, 0, btnWidth, btnHeight),
		geom.NewRectangleWH(-btnWidth, 0, btnWidth, btnHeight),
		geom.NewRectangleWH(0, -btnHeight, btnWidth, btnHeight),
		geom.NewRectangleWH(-btnWidth, -btnHeight, btnWidth, btnHeight),
	}

	pointer        *adagui.Pointer
	disp         *adatft.Display
	debug, nogui bool
	W, H         float64

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

func printEvent(event adatft.PointerEvent) {

	if !debug {
		return
	}
	log.Printf("[%d]: %10s: %v",
		event.Time.UnixMilli(), event.Type, event.Pos)
}

func initGrid(gc *gg.Context, actCol, actRow int) {
	// Clear the context and draw the grid first
	gc.SetFillColor(colors.Black)
	gc.Clear()
	gc.SetFillColor(gridColor)
	rows := int(math.Round((H-2*gridMargin)/gridSpace))
	cols := int(math.Round((W-2*gridMargin)/gridSpace))
	for row, y := 0, gridMargin; row <= rows; row, y = row+1, y+gridSpace {
		for col, x := 0, gridMargin; col <= cols; col, x = col+1, x+gridSpace {
			if (row==0 || row==rows) && (col<2 || col>cols-2) {
				continue
			}
			gc.DrawPoint(x, y, gridPointSize)
			gc.Fill()
		}
	}

	// The draw the walls in order to get a fine Irrgarten.
	/*
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
	*/

	gc.SetTextColor(fontColor)
	gc.SetFontFace(fontFace)
	gc.DrawStringAnchored("Oben", float64(adatft.Width)/2.0, fontSize, 0.5, 0.5)
}

func initTrace(gc *gg.Context) {
	gc.SetFillColor(colors.Transparent)
	gc.Clear()
	pointIdx = 0
	//gc.SetStrokeWidth(pointRadius)
	//gc.SetStrokeColor(pointColor)
	//gc.SetFillColor(pointColor)
}

func updateTrace(gc *gg.Context, pt geom.Point) {
	gc.SetFillColor(pointColors[pointIdx])
	gc.DrawPoint(pt.X, pt.Y, pointRadiae[pointIdx])
	gc.Fill()
}

func initCross(gc *gg.Context) {
	gc.SetFillColor(colors.Transparent)
	gc.Clear()
	gc.SetStrokeWidth(crossLineWidth)
	gc.SetStrokeColor(crossColor)
}

func drawCross(gc *gg.Context, pt geom.Point) {
	gc.Clear()
	gc.DrawLine(pt.X-crossSize/2, pt.Y, pt.X+crossSize/2, pt.Y)
	gc.DrawLine(pt.X, pt.Y-crossSize/2, pt.X, pt.Y+crossSize/2)
	gc.Stroke()
}

func initCtrls(gc *gg.Context) {
	for i, rect := range btnRects {
		switch i {
		case 0:
		case 1:
			btnRects[i] = rect.Add(geom.NewPoint(W, 0))
		case 2:
			btnRects[i] = rect.Add(geom.NewPoint(0, H))
		case 3:
			btnRects[i] = rect.Add(geom.NewPoint(W, H))
		}
	}
	face, _ := fonts.NewFace(btnFont, btnFontSize)

    gc.SetFillColor(colors.Transparent)
    gc.Clear()
    gc.SetStrokeWidth(btnBorderWidth)
	gc.SetFontFace(face)
	for i, rect := range btnRects {
    	gc.SetFillColor(btnColors[i].Alpha(0.5))
    	gc.SetStrokeColor(btnColors[i])
		gc.SetTextColor(btnColors[i])
    	gc.DrawRectangle(rect.Inset(btnBorderWidth/2, btnBorderWidth/2).AsCoord())
    	gc.FillStroke()
		mp := rect.C()
		gc.DrawStringAnchored(btnTitles[i], mp.X, mp.Y, 0.5, 0.5)
	}
}

func composeScreen(out *image.RGBA, grid, trace, cross, ctrls *gg.Context) {
	draw.Draw(out, out.Bounds(), grid.Image(), image.Point{}, draw.Over)
	draw.Draw(out, out.Bounds(), trace.Image(), image.Point{}, draw.Over)
	draw.Draw(out, out.Bounds(), cross.Image(), image.Point{}, draw.Over)
	draw.Draw(out, out.Bounds(), ctrls.Image(), image.Point{}, draw.Over)
}

func main() {
	var rotation adatft.RotationType = adatft.Rotate000
	var grid, trace, cross, ctrls *gg.Context
	var out *image.RGBA

	flag.BoolVar(&debug, "debug", false, "write events to stdout")
	flag.BoolVar(&nogui, "nogui", false, "dont paint on the screen")
	flag.Var(&rotation, "rotate", "display rotation")
	flag.Parse()

	//adatft.Init()
	log.Printf("> OpenDisplay()\n")
	disp = adatft.OpenDisplay(rotation)
	//m := disp.Matrix()

	log.Printf("> OpenPointer()\n")
	pointer = adagui.OpenPointer()
	r := geom.Rectangle{Max: disp.DrawBounds().Size()}
	pointer.SetPosRange(r)
	pointer.SetWheelRange(0, 0, 50)
	pointer.StartEvents()

	W, H = float64(adatft.Width), float64(adatft.Height)
	//sz := disp.DispBounds().Size().Int()

	log.Printf("> NewContext() for Grid\n")
	grid = disp.Canvas()
	//grid.SetMatrix(m)
	initGrid(grid, 0, 0)

	log.Printf("> NewContext() for Trace\n")
	trace = disp.Canvas()
	//trace.SetMatrix(m)
	initTrace(trace)

	log.Printf("> NewContext() for Cross\n")
	cross = disp.Canvas()
	//cross.SetMatrix(m)
	initCross(cross)

	log.Printf("> NewContext() for Controls\n")
	ctrls = disp.Canvas()
	//ctrls.SetMatrix(m)
	initCtrls(ctrls)

	out = image.NewRGBA(disp.DispBounds().Int())

	done := make(chan bool)
	ticker := time.NewTicker(30 * time.Millisecond)

	// Draw oder Paint Thread. Via Ticker zeitgesteuert (alle 30 ms)
	go func() {
		for {
			select {
			case <-ticker.C:
				adatft.PaintWatch.Start()
				composeScreen(out, grid, trace, cross, ctrls)
				adatft.PaintWatch.Stop()
				disp.Draw(out)
			case <-done:
				return
			}
		}
	}()

EVENT_LOOP:
	for event := range pointer.EventQ {
		printEvent(event)

		pt := event.Pos

		switch event.Type {
		case adatft.PointerMove:
			drawCross(cross, event.Pos)
		case adatft.PointerPress, adatft.PointerDrag:
			updateTrace(trace, event.Pos)
			drawCross(cross, event.Pos)
		case adatft.PointerRelease:
			switch {
			case pt.In(btnRects[0]):
				break EVENT_LOOP
			case pt.In(btnRects[1]):
				initTrace(trace)
				continue
			case pt.In(btnRects[2]):
				pointIdx = 1
				continue
			case pt.In(btnRects[3]):
				pointIdx = 2
				continue
			}
		}
	}

	done <- true

	grid.SetFillColor(colors.Black)
	grid.Clear()
	trace.SetFillColor(colors.Black)
	trace.Clear()
	cross.SetFillColor(colors.Black)
	cross.Clear()
	ctrls.SetFillColor(colors.Black)
	ctrls.Clear()
	composeScreen(out, grid, trace, cross, ctrls)
	disp.Draw(out)

	disp.Close()
	pointer.Close()

	adatft.PrintStat()
}
