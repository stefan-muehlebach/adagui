package main

// Damit versuchen wir, die Animationen aus displayTest in einer
// AdaGUI-konformen Art umzuschreiben.
import (
	"flag"
	//"fmt"
	"github.com/stefan-muehlebach/adagui"
	"github.com/stefan-muehlebach/adagui/touch"
	"github.com/stefan-muehlebach/gg/color"
	"github.com/stefan-muehlebach/gg/colornames"
	"log"
	"math/rand"
	"time"
	//"github.com/stefan-muehlebach/adagui/value"
	"github.com/stefan-muehlebach/adatft"
	"github.com/stefan-muehlebach/gg"
	"github.com/stefan-muehlebach/gg/geom"
)

// -----------------------------------------------------------------------------
//
// Alles fuer die Animation mit Polygonen.
type Point struct {
	x, y, dx, dy float64
}

type Polygon struct {
	adagui.LeafEmbed
	ptList                 []*Point
	StrokeColor, FillColor color.Color
	LineWidth              float64
}

func NewPolygon(dispWidth, dispHeight float64, edges int) *Polygon {
	p := &Polygon{}
	p.ptList = make([]*Point, edges)
	for i := 0; i < edges; i++ {
		pt := &Point{}
		pt.x = rand.Float64() * dispWidth
		pt.y = rand.Float64() * dispHeight
		pt.dx = rand.Float64()*5.0 - 2.0
		pt.dy = rand.Float64()*5.0 - 2.0
		p.ptList[i] = pt
	}
	p.StrokeColor = colornames.White
	p.FillColor = colornames.RandColor().Alpha(0.5)
	p.LineWidth = 3.0
	return p
}

func (p *Polygon) Move(bounds geom.Rectangle) {
	for _, pt := range p.ptList {
		pt.Move(bounds)
	}
}

func (p *Polygon) Paint(gc *gg.Context) {
	gc.MoveTo(p.ptList[0].x, p.ptList[0].y)
	for _, pt := range p.ptList[1:] {
		gc.LineTo(pt.x, pt.y)
	}
	gc.ClosePath()
	gc.SetFillColor(p.FillColor)
	gc.SetStrokeColor(p.StrokeColor)
	gc.SetStrokeWidth(p.LineWidth)
	gc.FillStroke()
}

func (p *Point) Move(bounds geom.Rectangle) {
	p.x += p.dx
	p.y += p.dy
	if p.x < bounds.Min.X || p.x > bounds.Max.X {
		p.dx *= -1
		p.x += p.dx
	}
	if p.y < bounds.Min.Y || p.y > bounds.Max.Y {
		p.dy *= -1
		p.y += p.dy
	}
}

func InitAnimation() {
	polyList = make([]*Polygon, numObjs)
	for i := 0; i < numObjs; i++ {
		polyList[i] = NewPolygon(canvas.Size().X, canvas.Size().Y, numEdges)
		canvas.Add(polyList[i])
	}
}

func StartAnimation(dt time.Duration) {
	if ticker != nil {
		ticker.Reset(dt)
		return
	}

	ticker = time.NewTicker(dt)
	go func() {
		for {
			<-ticker.C
			for _, p := range polyList {
				p.Move(canvas.LocalBounds())
			}
			win.Repaint()
		}
	}()
}

func StopAnimation() {
	if ticker == nil {
		return
	}
	ticker.Stop()
}

//----------------------------------------------------------------------------

func init() {
	log.SetFlags(log.Lmicroseconds | log.Lmsgprefix)
	log.SetPrefix(": ")
}

var (
	screen            *adagui.Screen
	win               *adagui.Window
	canvas            *adagui.Canvas
	ticker            *time.Ticker
	numObjs, numEdges int
	polyList          []*Polygon
)

func main() {
	var spc float64 = 5.0
	var posPt geom.Point

	flag.IntVar(&numObjs, "objs", 5, "Number of objects")
	flag.IntVar(&numEdges, "edges", 3, "Number of edges of an object")

	//utils.StartProfiling()

	screen = adagui.NewScreen(adatft.Rotate090)
	win = screen.NewWindow()

	base := adagui.NewGroup()
	win.SetRoot(base)

	canvas = adagui.NewCanvas(320.0, 195.0)
	canvas.SetPos(geom.Point{})
	//canvasPanel.Clip = true
	//canvasPanel.FillColor = utils.Black
	base.Add(canvas)

	uiPanel := adagui.NewPanel(320.0, 45.0)
	uiPanel.SetPos(geom.Point{0.0, canvas.Size().Y})
	//uiPanel.FillColor = utils.Color{0, 0, 0, 0}
	base.Add(uiPanel)

	//canvas = adagui.NewCanvas(320.0, 240.0)
	//canvasPanel.Add(canvas)

	btn := adagui.NewTextButton("Quit")
	posPt = geom.Point{320.0 - btn.Size().X - spc, spc}
	btn.SetPos(posPt)
	uiPanel.Add(btn)
	btn.SetOnTap(func(evt touch.Event) {
		StopAnimation()
		time.Sleep(10 * time.Millisecond)
		screen.Quit()
	})

	chk := adagui.NewCheckboxWithCallback("Run Animation", func(v bool) {
		if v {
			StartAnimation(40 * time.Millisecond)
		} else {
			StopAnimation()
		}
	})
	posPt = geom.Point{spc, spc}
	chk.SetPos(posPt)
	uiPanel.Add(chk)

	sldLen := uiPanel.Size().X - btn.Size().X - chk.Size().X - 4*spc
	sld := adagui.NewSliderWithCallback(sldLen, adagui.Horizontal,
		func(f float64) {
			canvas.ScaleAbout(canvas.Bounds().Center(), f, f)
		})
	sld.SetRange(0.1, 1.9, 0.1)
	sld.SetValue(1.0)
	posPt = geom.Point{chk.Rect().X1() + spc, spc}
	sld.SetPos(posPt)
	uiPanel.Add(sld)

	InitAnimation()
	screen.SetWindow(win)
	screen.Run()

	//adatft.PrintStat()
	//utils.StopProfiling()
}
