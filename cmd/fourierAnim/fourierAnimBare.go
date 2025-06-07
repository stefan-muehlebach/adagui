//go:build ignore

package main

import (
	"flag"
	"fmt"
	"math"
	"math/cmplx"
	"os"
	"os/signal"
	"time"

	"github.com/stefan-muehlebach/adagui"
	"github.com/stefan-muehlebach/adatft"
	"github.com/stefan-muehlebach/gg"
	"github.com/stefan-muehlebach/gg/colors"
	"github.com/stefan-muehlebach/gg/geom"
)

const (
	F            float64 = 1.0
	DefMaxFreq           = 20
	Dt           float64 = 0.05 / float64(DefMaxFreq)
	DefCoeffFile         = "coeff.json"
)

var (
	disp      *adatft.Display
	gc, img   *gg.Context
	tmp       = math.Sqrt(3.0) / 2.0
	coeffFile string
	coeffList *CoeffList
)

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
	d.radius = F * cmplx.Abs(coeff.factor)
	d.initAngle = cmplx.Phase(coeff.factor)

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
	gc.Push()
	gc.Translate(d.mp.AsCoord())
	//gc.Rotate(d.angle)

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
	//gc.DrawPoint(0.0, 0.0, d.pointerSize/2)
	//gc.FillStroke()

	if d.child != nil {
		d.child.Draw(gc)
	}
	gc.Pop()
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
	p.penWidth = 1.0
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
	//p.img.DrawPoint(p.mp.X, p.mp.Y, p.penWidth)
	//p.img.FillStroke()
}

//-----------------------------------------------------------------------------

// Go routine fuer das Zeichnen der Objekte.
func paintThread(obj *FourierDisc, syncQ chan bool) {
	for {
		if _, ok := <-syncQ; !ok {
			break
		}
		t1 := time.Now()
		gc.SetFillColor(colors.DarkRed.Alpha(0.5))
		gc.Clear()
		obj.Draw(gc)
		gc.DrawImageAnchored(img.Image(), 0.0, 0.0, 0.0, 0.0)
		paintDur += time.Since(t1)
		disp.Draw(gc.Image())
	}
}

// Go routine fuer die Animation der Objekte
func animThread(obj *FourierDisc, syncQ, quitQ chan bool) {
	var ticker *time.Ticker
	var dt float64 = Dt
	var t float64 = 0.0

	ticker = time.NewTicker(step)
ForeverLoop:
	for {
		select {
		case <-ticker.C:
			if singleStep {
				ticker.Stop()
			}
			t1 := time.Now()
			obj.Animate(t)
			t += dt
			if t > 1.0 {
				t = 0.0
			}
			animDur += time.Since(t1)
			syncQ <- true
		case <-quitQ:
			break ForeverLoop
		}
	}
	close(syncQ)
	quitQ <- true
}

//-----------------------------------------------------------------------------

var (
	step              time.Duration
	singleStep        bool = false
	animDur, paintDur time.Duration
)

func main() {
	var firstDisc *FourierDisc
	var syncQ, quitQ chan bool
	var sigChan chan os.Signal
	var maxFreq int

	flag.DurationVar(&step, "step", 30*time.Millisecond,
		"time step of the animation")
	flag.IntVar(&maxFreq, "freq", DefMaxFreq, "Max. Frequence")
	flag.StringVar(&coeffFile, "in", DefCoeffFile, "Input file with coeff.")
	flag.Parse()

	adagui.StartProfiling()

	coeffList = ReadCoeffList(coeffFile)

	disp = adatft.OpenDisplay(adatft.Rotate000)

	gc = gg.NewContext(adatft.Width, adatft.Height)
	img = gg.NewContext(adatft.Width, adatft.Height)
	img.SetFillColor(colors.Transparent)
	img.Clear()

	firstDisc = NewFourierDisc(coeffList.Get(0), nil)
	disc := firstDisc
	for i := range maxFreq {
		freq := i + 1
		disc = NewFourierDisc(coeffList.Get(freq), disc)
		disc = NewFourierDisc(coeffList.Get(-freq), disc)
	}
	NewFourierPen(img, disc)
	firstDisc.SetPos(geom.Point{float64(adatft.Width) / 2.0,
		float64(adatft.Height) / 2.0})
	firstDisc.SetAngle(0.0)

	syncQ = make(chan bool)
	quitQ = make(chan bool)
	sigChan = make(chan os.Signal)

	go paintThread(firstDisc, syncQ)
	go animThread(firstDisc, syncQ, quitQ)

	signal.Notify(sigChan, os.Interrupt)
	<-sigChan

	quitQ <- true
	<-quitQ
	time.Sleep(1 * time.Second)

	disp.Close()
	adagui.StopProfiling()

	adatft.PrintStat()
	fmt.Printf("animation:\n")
	fmt.Printf("  %v total\n", animDur)
	fmt.Printf("  %.3f ms / frame\n", float64(animDur.Milliseconds())/float64(adatft.DispWatch.Num()))
}
