package main

import (
	"flag"
	"fmt"
	"github.com/stefan-muehlebach/adatft"
	"github.com/stefan-muehlebach/gg"
	"github.com/stefan-muehlebach/gg/color"
	"github.com/stefan-muehlebach/gg/colornames"
	"github.com/stefan-muehlebach/gg/geom"
	//"log"
	"math"
	"math/cmplx"
	"os"
	"os/signal"
	"time"
)

const (
    F float64 = 40.0
    MaxFreq   = 19
    Dt float64 = 0.02 / float64(MaxFreq)
)

var (
	disp    *adatft.Display
	gc, img *gg.Context
	tmp     = math.Sqrt(3.0) / 2.0
)

type FourierCoeff struct {
	freq   float64
	factor complex128
}

var (
/*
    CoeffList = []FourierCoeff{
        FourierCoeff{  0, (  -0.0000  +0.6688i)},
        FourierCoeff{  1, (  +0.0077  -0.2466i)},
        FourierCoeff{ -1, (  +0.1080  +3.4708i)},
        FourierCoeff{  2, (  +0.0216  -0.3465i)},
        FourierCoeff{ -2, (  -0.0125  -0.2007i)},
        FourierCoeff{  3, (  -0.0019  +0.0203i)},
        FourierCoeff{ -3, (  -0.0238  -0.2542i)},
        FourierCoeff{  4, (  +0.0127  -0.1016i)},
        FourierCoeff{ -4, (  -0.0261  -0.2085i)},
        FourierCoeff{  5, (  -0.0013  +0.0081i)},
        FourierCoeff{ -5, (  -0.0051  -0.0323i)},
        FourierCoeff{  6, (  +0.0085  -0.0452i)},
        FourierCoeff{ -6, (  -0.0104  -0.0550i)},
        FourierCoeff{  7, (  -0.0008  +0.0034i)},
        FourierCoeff{ -7, (  -0.0014  -0.0065i)},
        FourierCoeff{  8, (  +0.0060  -0.0234i)},
        FourierCoeff{ -8, (  -0.0078  -0.0308i)},
        FourierCoeff{  9, (  -0.0014  +0.0048i)},
        FourierCoeff{ -9, (  -0.0014  -0.0048i)},
    }
*/
    CoeffList = []FourierCoeff{
        FourierCoeff{  0, (  +1.1720  -1.3875i)},
        FourierCoeff{  1, (  -0.3010  -0.5965i)},
        FourierCoeff{ -1, (  -1.4220  -2.0622i)},
        FourierCoeff{  2, (  -0.4890  +0.6058i)},
        FourierCoeff{ -2, (  +0.4454  +0.3561i)},
        FourierCoeff{  3, (  +0.0548  -0.1308i)},
        FourierCoeff{ -3, (  -0.3568  -0.5938i)},
        FourierCoeff{  4, (  -0.1810  +0.2281i)},
        FourierCoeff{ -4, (  +0.1254  -0.0810i)},
        FourierCoeff{  5, (  -0.0960  +0.1690i)},
        FourierCoeff{ -5, (  +0.1017  -0.3047i)},
        FourierCoeff{  6, (  +0.1180  +0.2466i)},
        FourierCoeff{ -6, (  +0.1081  +0.1222i)},
        FourierCoeff{  7, (  +0.0885  +0.0887i)},
        FourierCoeff{ -7, (  -0.0869  -0.0510i)},
        FourierCoeff{  8, (  +0.3233  -0.0732i)},
        FourierCoeff{ -8, (  -0.2468  -0.3610i)},
        FourierCoeff{  9, (  +0.0936  +0.0113i)},
        FourierCoeff{ -9, (  -0.0682  -0.0434i)},
        FourierCoeff{ 10, (  -0.0209  +0.1233i)},
        FourierCoeff{-10, (  +0.2035  +0.2263i)},
        FourierCoeff{ 11, (  -0.0632  -0.0998i)},
        FourierCoeff{-11, (  +0.0388  +0.0689i)},
        FourierCoeff{ 12, (  -0.0126  -0.0264i)},
        FourierCoeff{-12, (  -0.0330  +0.0401i)},
        FourierCoeff{ 13, (  -0.0257  -0.0525i)},
        FourierCoeff{-13, (  +0.0565  +0.0715i)},
        FourierCoeff{ 14, (  -0.0329  -0.0445i)},
        FourierCoeff{-14, (  -0.0345  -0.0587i)},
        FourierCoeff{ 15, (  -0.0112  +0.0007i)},
        FourierCoeff{-15, (  +0.0426  +0.0256i)},
        FourierCoeff{ 16, (  -0.0368  +0.0146i)},
        FourierCoeff{-16, (  +0.0572  -0.0071i)},
        FourierCoeff{ 17, (  +0.0306  -0.0522i)},
        FourierCoeff{-17, (  -0.0231  -0.0314i)},
        FourierCoeff{ 18, (  -0.0037  -0.0023i)},
        FourierCoeff{-18, (  -0.0255  -0.0158i)},
        FourierCoeff{ 19, (  +0.0074  +0.0306i)},
        FourierCoeff{-19, (  +0.0165  +0.0445i)},
    }

/*
	CoeffList = []FourierCoeff{
		FourierCoeff{0, 0.0 - 0.0060066i},
		FourierCoeff{1, 0.0 + 0.0390431i},
		FourierCoeff{-1, 0.0 + 0.9760778i},
		FourierCoeff{2, 0.0 - 0.1952156i},
		FourierCoeff{-2, 0.0 - 0.1952156i},
		FourierCoeff{3, 0.0 + 0.0780862i},
		FourierCoeff{-3, 0.0 - 0.2342587i},
		FourierCoeff{4, 0.0 - 0.0390431i},
		FourierCoeff{-4, 0.0 - 0.0390431i},
	}
*/
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

	fillColor, borderColor, pointerColor   color.Color
	borderWidth, pointerWidth, pointerSize float64

	freq                     float64
	sign                     bool
	radius, initAngle, angle float64
	child                    FourierObject
}

func NewFourierDisc(coeff FourierCoeff, parent *FourierDisc) *FourierDisc {
	d := &FourierDisc{}
	d.freq = coeff.freq
	d.sign = math.Signbit(d.freq)
	d.radius = F * cmplx.Abs(coeff.factor)
	d.initAngle = cmplx.Phase(coeff.factor)

	d.fillColor = colornames.LightSkyBlue.Alpha(0.1)
	d.borderColor = colornames.WhiteSmoke.Alpha(0.5)
	d.borderWidth = 1.0
	d.pointerColor = colornames.WhiteSmoke.Alpha(0.5)
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
	w := math.Cos(d.angle) * d.radius
	h := math.Sin(d.angle) * d.radius
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
	img      *gg.Context
	penColor color.Color
	penWidth float64
    firstPoint bool
}

func NewFourierPen(img *gg.Context, parent *FourierDisc) *FourierPen {
	p := &FourierPen{}
	p.img = img
	p.penColor = colornames.WhiteSmoke
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
		gc.SetFillColor(colornames.Black)
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

	flag.DurationVar(&step, "step", 30*time.Millisecond,
		"time step of the animation")
	flag.Parse()

	disp = adatft.OpenDisplay(adatft.Rotate000)

	gc = gg.NewContext(adatft.Width, adatft.Height)
	img = gg.NewContext(adatft.Width, adatft.Height)
	img.SetFillColor(color.Transparent)
	img.Clear()

	firstDisc = NewFourierDisc(CoeffList[0], nil)
	disc := firstDisc
	//log.Printf("freq : %f\n", disc.freq)
	//log.Printf("angle: %f\n", disc.angle)
	for _, coeff := range CoeffList[1:2*MaxFreq+1] {
		disc = NewFourierDisc(coeff, disc)
		//log.Printf("freq : %f\n", disc.freq)
		//log.Printf("angle: %f\n", disc.angle)
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

	adatft.PrintStat()
	fmt.Printf("animation:\n")
	fmt.Printf("  %v total\n", animDur)
	fmt.Printf("  %.3f ms / frame\n", float64(animDur.Milliseconds())/float64(adatft.DispWatch.Num()))
}
