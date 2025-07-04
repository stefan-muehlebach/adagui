package main

import (
	"flag"
	"fmt"
	"log"
	"math"
	"math/cmplx"
	"os"
	"os/signal"
	"time"

	"github.com/cpmech/gosl/fun/fftw"
	"github.com/stefan-muehlebach/adagui"
	"github.com/stefan-muehlebach/adagui/touch"
	"github.com/stefan-muehlebach/adatft"
	"github.com/stefan-muehlebach/gg"
	"github.com/stefan-muehlebach/gg/colors"
	"github.com/stefan-muehlebach/gg/geom"
)

func init() {
	log.SetFlags(log.Lmicroseconds | log.Lmsgprefix)
	log.SetPrefix(": ")
}

func SignalHandler() {
	sigChan := make(chan os.Signal)
	signal.Notify(sigChan, os.Interrupt)
	<-sigChan
	fourObj.StopAnim()
	screen.Quit()
}

const (
	F          float64 = 1.0
	DefMaxFreq         = 20
	Dt         float64 = 0.05 / float64(DefMaxFreq)
)

var (
	screen                    *adagui.Screen
	drawWin, calcWin, animWin *adagui.Window
	overlayFile               string
	maxFreq                   int
	poly                      *adagui.Polygon
	fourObj                   *FourierThingy
)

//-----------------------------------------------------------------------------

type FourierThingy struct {
	adagui.LeafEmbed
	firstDisc    *FourierDisc
	cfList       *CoeffList
	maxFreq      int
	layer        *gg.Context
	syncQ, quitQ chan bool
	ticker       *time.Ticker
}

func NewFourierThingy(cfList *CoeffList, maxFreq int) *FourierThingy {
	f := &FourierThingy{}
	f.Wrapper = f
	f.LeafEmbed.Init()
	f.layer = gg.NewContext(adatft.Width, adatft.Height)
	f.layer.SetFillColor(colors.Transparent)
	f.layer.Clear()
	//f.layer.Translate(f.layer.Bounds().Size().Mul(0.5).AsCoord())
	f.cfList = cfList
	f.maxFreq = maxFreq
	f.syncQ = make(chan bool)
	f.quitQ = make(chan bool)

	f.firstDisc = NewFourierDisc(f.cfList.Get(0), nil)
	f.SetMinSize(geom.Point{2 * f.firstDisc.radius, 2 * f.firstDisc.radius})
	disc := f.firstDisc
	for i := range f.maxFreq {
		freq := i + 1
		disc = NewFourierDisc(f.cfList.Get(freq), disc)
		disc = NewFourierDisc(f.cfList.Get(-freq), disc)
	}
	NewFourierPen(f.layer, disc)
	f.firstDisc.SetPos(geom.Point{0, 0})
	f.firstDisc.SetAngle(0.0)
	f.Mark(adagui.MarkNeedsPaint)
	return f
}

func (f *FourierThingy) Draw(gc *gg.Context) {
	fmt.Printf("Draw of FourierThingy\n")
	f.firstDisc.Draw(gc)
	gc.DrawImageAnchored(f.layer.Image(), 0.0, 0.0, 0.5, 0.5)
	fmt.Printf("   done\n")
}

func (f *FourierThingy) StartAnim() {
	// go fourObj.paintFunc()
	go f.animateFunc()
}

func (f *FourierThingy) StopAnim() {
	f.quitQ <- true
	<-f.quitQ
}

func (f *FourierThingy) animateFunc() {
	var dt float64 = 0.05 / float64(DefMaxFreq)
	var t float64 = 0.0
	var step time.Duration = 30 * time.Millisecond

	f.ticker = time.NewTicker(step)

MainLoop:
	for {
		select {
		case <-f.ticker.C:
			f.firstDisc.Animate(t)
			f.Mark(adagui.MarkNeedsPaint)
			t += dt
			if t > 1.0 {
				t = 0.0
			}
			screen.Repaint()
		case <-f.quitQ:
			break MainLoop
		}
	}
	fmt.Printf("Quitting main loop\n")
	close(f.syncQ)
}

//-----------------------------------------------------------------------------

type FourierObject interface {
	SetPos(mp geom.Point)
	SetAngle(angle float64)
	Animate(t float64)
	Draw(gc *gg.Context)
}

//-----------------------------------------------------------------------------

type FourierDisc struct {
	mp, nextMp geom.Point

	fillColor, borderColor, pointerColor   colors.Color
	borderWidth, pointerWidth, pointerSize float64

	freq                     float64
	sign                     bool
	radius, initAngle, angle float64
	child                    FourierObject
}

func NewFourierDisc(coeff FourierCoeff, parent *FourierDisc) *FourierDisc {
	d := &FourierDisc{}
	d.freq = float64(coeff.freq)
	d.sign = math.Signbit(d.freq)
	d.radius = cmplx.Abs(coeff.factor)
	//d.initAngle = cmplx.Phase(coeff.factor)
	d.initAngle = cmplx.Phase(coeff.factor) + math.Pi/2.0

	d.fillColor = colors.LightSkyBlue.Alpha(0.1)
	d.borderColor = colors.WhiteSmoke.Alpha(0.5)
	d.borderWidth = 1.0
	d.pointerColor = colors.WhiteSmoke.Alpha(0.5)
	d.pointerWidth = 1.0
	d.pointerSize = 6.0
	d.SetAngle(0.0)

	if parent != nil {
		parent.child = d
	}

	return d
}

func (d *FourierDisc) SetPos(mp geom.Point) {
	d.mp = mp
}

func (d *FourierDisc) SetAngle(angle float64) {
	d.angle = math.Mod(d.initAngle+angle, 2.0*math.Pi)
	w := math.Sin(d.angle) * d.radius
	h := math.Cos(d.angle) * d.radius
	//d.nextMp = d.mp.Add(geom.Point{w, -h})
	d.nextMp = geom.Point{w, -h}
	if d.child != nil {
		d.child.SetPos(d.nextMp)
	}
}

func (d *FourierDisc) Animate(t float64) {
	if d.freq != 0.0 {
		d.SetAngle((t * d.freq) * 2.0 * math.Pi)
	}
	if d.child != nil {
		d.child.Animate(t)
	}
}

func (d *FourierDisc) Draw(gc *gg.Context) {
	fmt.Printf("Draw of FourierDisc\n")
	gc.Push()
	gc.Translate(d.mp.AsCoord())

	gc.SetStrokeWidth(d.borderWidth)
	gc.SetStrokeColor(d.borderColor)
	gc.SetFillColor(d.fillColor)
	gc.DrawCircle(0.0, 0.0, d.radius)
	gc.FillStroke()

	gc.SetStrokeWidth(d.pointerWidth)
	gc.SetStrokeColor(d.pointerColor)
	gc.SetFillColor(d.pointerColor)
	gc.DrawLine(0.0, 0.0, d.nextMp.X, d.nextMp.Y)
	gc.Stroke()

	if d.child != nil {
		d.child.Draw(gc)
	}
	gc.Pop()
	fmt.Printf("   done\n")
}

//-----------------------------------------------------------------------------

type FourierPen struct {
	mp, prevPt geom.Point
	img        *gg.Context
	penColor   colors.Color
	penWidth   float64
	firstPoint bool
}

func NewFourierPen(img *gg.Context, parent *FourierDisc) *FourierPen {
	p := &FourierPen{}
	p.img = img
	p.penColor = colors.WhiteSmoke
	p.penWidth = 2.0
	p.firstPoint = true
	parent.child = p
	return p
}

func (p *FourierPen) SetPos(mp geom.Point) {
	p.mp = mp
}

func (p *FourierPen) SetAngle(angle float64) {}

func (p *FourierPen) Animate(t float64) {}

func (p *FourierPen) Draw(gc *gg.Context) {
	fmt.Printf("Draw of FourierPen\n")
	pt := gc.Matrix().Transform(p.mp)
	if !p.firstPoint {
		p.img.SetFillColor(p.penColor)
		p.img.SetStrokeColor(p.penColor)
		p.img.SetStrokeWidth(p.penWidth)
		p.img.DrawLine(p.prevPt.X, p.prevPt.Y, pt.X, pt.Y)
		p.img.Stroke()
	} else {
		p.firstPoint = false
	}
	p.prevPt = pt
	fmt.Printf("   done\n")
}

//-----------------------------------------------------------------------------

// Erstellt ein neues Panel, in welchem der Benutzer mit den Stift einen
// Umriss zeichnen kann, der dann als Ausgangsfunktion fuer die Fourier-
// Transformation verwendet wird.
func NewDrawPanel(w, h float64) *adagui.Panel {
	var err error

	p := adagui.NewPanel(w, h)
	if overlayFile != "" {
		p.Image, err = gg.LoadPNG(overlayFile)
		if err != nil {
			log.Fatal(err)
		}
	}

	p.SetOnPress(func(evt touch.Event) {
		if poly != nil {
			poly.Remove()
		}
		poly = adagui.NewPolygon(p.Bounds().C())
		poly.Closed = true
		poly.AddPoint(evt.Pos)
		p.Add(poly)
		p.Mark(adagui.MarkNeedsPaint)
	})

	p.SetOnLongPress(func(evt touch.Event) {
		log.Printf("LongPressed!")
		//p.Win.Stop()
	})

	p.SetOnDrag(func(evt touch.Event) {
		poly.AddPoint(evt.Pos)
		p.Mark(adagui.MarkNeedsPaint)
	})

	p.SetOnRelease(func(evt touch.Event) {
		if poly == nil {
			return
		}
		poly.Flatten()
		poly.Mark(adagui.MarkNeedsPaint)

		screen.SetWindow(calcWin)
		time.Sleep(2 * time.Second)

		// Create the window for the animation
		animWin = screen.NewWindow()
		root := adagui.NewGroupPL(nil, adagui.NewPaddedLayout())
		animWin.SetRoot(root)

		panel := NewAnimPanel(10, 10)
		panel.SetColor(colors.DarkRed.Alpha(0.5))
		panel.SetBorderWidth(1.0)
		panel.SetBorderColor(colors.DarkRed.Bright(0.5))
		root.Add(panel)

		label := adagui.NewLabel("Tap somewhre on the screen to start.")
		label.SetTextColor(colors.DarkRed.Bright(0.7))
		label.SetAlign(adagui.AlignCenter | adagui.AlignBottom)
		label.SetPos(panel.Bounds().S().AddXY(0, -10))
		panel.Add(label)

		screen.SetWindow(animWin)
	})

	return p
}

func NewAnimPanel(w, h float64) *adagui.Panel {
	var cfList *CoeffList

	p := adagui.NewPanel(w, h)

	p.SetOnTap(func(evt touch.Event) {
		if fourObj != nil {
			return
		}
		// screen.StopPaint()
		fmt.Printf("len(cfList): %d\n", len(cfList.data))
		fourObj = NewFourierThingy(cfList, maxFreq)
		fourObj.SetPos(p.Bounds().C())
		p.Add(fourObj)
		fourObj.Mark(adagui.MarkNeedsPaint)
		fourObj.StartAnim()
		fmt.Printf("custom animation started\n")
	})

	pts := poly.Points()
	data := make([]complex128, len(pts))
	fftPlan := fftw.NewPlan1d(data, false, true)
	for i, pt := range pts {
		data[i] = complex(pt.X, pt.Y)
	}
	fftPlan.Execute()
	n := complex(float64(len(data)), 0.0)
	for i, dat := range data {
		data[i] = dat / n
	}
	fftPlan.Free()
	cfList = NewCoeffList(data)

	return p
}

// Hauptprogramm.
func main() {
	flag.StringVar(&overlayFile, "overlay", "", "Optional overlay graphic file")
	flag.IntVar(&maxFreq, "freq", DefMaxFreq, "Maximal frequency")
	flag.Parse()
	adagui.StartProfiling()

	go SignalHandler()

	screen = adagui.NewScreen(adatft.Rotate090)

	// Create the windows for the sketching part
	drawWin = screen.NewWindow()
	root := adagui.NewGroupPL(nil, adagui.NewPaddedLayout())
	drawWin.SetRoot(root)

	panel := NewDrawPanel(10, 10)
	panel.SetColor(colors.DarkGreen.Alpha(0.5))
	panel.SetBorderWidth(1.0)
	panel.SetBorderColor(colors.DarkGreen.Bright(0.5))
	root.Add(panel)

	label := adagui.NewLabel("Draw something, but use only one stroke!")
	label.SetTextColor(colors.DarkGreen.Bright(0.7))
	label.SetPos(panel.Bounds().S().AddXY(0, -10))
	label.SetAlign(adagui.AlignCenter | adagui.AlignBottom)
	panel.Add(label)

	// Create the window for the calculation
	calcWin = screen.NewWindow()
	root = adagui.NewGroupPL(nil, adagui.NewPaddedLayout())
	calcWin.SetRoot(root)

	panel = adagui.NewPanel(10, 10)
	panel.SetColor(colors.DarkBlue.Alpha(0.5))
	panel.SetBorderWidth(1.0)
	panel.SetBorderColor(colors.DarkBlue.Bright(0.5))
	root.Add(panel)

	label = adagui.NewLabel("Calculating... Please wait!")
	label.SetFontSize(32.0)
	label.SetTextColor(colors.DarkBlue.Bright(0.7))
	label.SetPos(panel.Bounds().C())
	label.SetAlign(adagui.AlignCenter | adagui.AlignMiddle)
	panel.Add(label)

	screen.SetWindow(drawWin)
	screen.Run()

	adagui.StopProfiling()
}
