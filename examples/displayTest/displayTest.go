package main

import (
	"flag"
	"fmt"
	"github.com/stefan-muehlebach/adatft"
	"github.com/stefan-muehlebach/gg"
	"github.com/stefan-muehlebach/gg/color"
	"github.com/stefan-muehlebach/gg/colornames"
	"github.com/stefan-muehlebach/gg/fonts"
	"golang.org/x/image/font"
	//"golang.org/x/image/font/opentype"
	"image"
	"log"
	"math"
	"math/rand"
	"os"
	"os/signal"
	"path"
	"runtime"
	"runtime/pprof"
	"runtime/trace"
	"strings"
	"time"
)

// PolygonAnimation --
//
// Animation von halbtransparenten Polygonen. Die Anzahl der Polygone kann
// ueber das Flag 'objs' und die Anzahl Ecken der Polygone ueber das Flag
// 'edges' gesteuert werden.
func PolygonAnimation() {
	var polyList []*Polygon

	polyList = make([]*Polygon, numObjs)
	for i := 0; i < numObjs; i++ {
		polyList[i] = NewPolygon(adatft.Width, adatft.Height, numEdges)
	}

	gc.SetStrokeWidth(3)
	gc.SetLineCapRound()
	gc.SetLineJoinRound()
	gc.SetFillColor(colornames.Black)
	gc.Clear()

	ticker := time.NewTicker(30 * time.Millisecond)
	for range ticker.C {
		if !runFlag {
			break
		}
		adatft.PaintWatch.Start()
		gc.SetFillColor(color.RGBAF{0, 0, 0, blurFactor})
		gc.DrawRectangle(0, 0, float64(adatft.Width), float64(adatft.Height))
		gc.Fill()

		for _, p := range polyList {
			p.Draw(gc)
			p.Move(0.0, float64(adatft.Width-1), 0.0, float64(adatft.Height-1))
		}
		adatft.PaintWatch.Stop()
		disp.Draw(gc.Image())
	}
}

type Polygon struct {
	p                      []*Point
	strokeColor, fillColor color.Color
}

func NewPolygon(dispWidth, dispHeight, edges int) *Polygon {
	p := &Polygon{}
	p.p = make([]*Point, edges)
	for i := 0; i < edges; i++ {
		pt := &Point{}
		pt.x = rand.Float64() * float64(dispWidth)
		pt.y = rand.Float64() * float64(dispHeight)
		pt.dx = rand.Float64()*5.0 - 2.0
		pt.dy = rand.Float64()*5.0 - 2.0
		p.p[i] = pt
	}
	p.strokeColor = colornames.White
	p.fillColor = colornames.RandColor().Alpha(0.5)
	return p
}

func (p *Polygon) Move(xmin, xmax, ymin, ymax float64) {
	for _, p := range p.p {
		p.Move(xmin, xmax, ymin, ymax)
	}
}

func (p *Polygon) Draw(gc *gg.Context) {
	gc.MoveTo(p.p[0].x, p.p[0].y)
	for _, p := range p.p[1:] {
		gc.LineTo(p.x, p.y)
	}
	gc.ClosePath()
	gc.SetStrokeStyle(gg.NewSolidPattern(p.strokeColor))
	gc.SetFillStyle(gg.NewSolidPattern(p.fillColor))
	gc.FillStroke()
}

type Point struct {
	x, y, dx, dy float64
}

func (p *Point) Move(xmin, xmax, ymin, ymax float64) {
	p.x += p.dx
	p.y += p.dy
	if p.x < xmin || p.x > xmax {
		p.dx *= -1
		p.x += p.dx
	}
	if p.y < ymin || p.y > ymax {
		p.dy *= -1
		p.y += p.dy
	}
}

// 3D-Animation ---------------------------------------------------------------
func Cube3DAnimation() {
	var mBase, m Matrix
	var cube, cubeT *Cube
	var cloud, cloudT *Cloud
	var zero, xAxis, yAxis, zAxis Vector
	var zeroT, xAxisT, yAxisT, zAxisT Vector
	var alpha, dAlpha float64
	var beta, dBeta float64

	cube = NewCube(70.0)
	cubeT = &Cube{}

	cloud = NewCloud(0, 0, 0, 70.0, 5*numObjs, 4.0)
	cloudT = &Cloud{}

	zero = NewVector(0.0, 0.0, 0.0)
	xAxis = NewVector(70.0, 0.0, 0.0)
	yAxis = NewVector(0.0, 70.0, 0.0)
	zAxis = NewVector(0.0, 0.0, 70.0)

	alpha = math.Pi / 12.0
	dAlpha = math.Pi / 162.0
	beta = math.Pi / 18.0
	dBeta = math.Pi / 126.0

	mBase = Identity().Multiply(Scale(NewVector(1.0, -1.0, 1.0))).Multiply(Translate(NewVector(240.0, -160.0, 0.0)))

	ticker := time.NewTicker(30 * time.Millisecond)
	for range ticker.C {
		if !runFlag {
			break
		}
		adatft.PaintWatch.Start()

		gc.SetFillColor(colornames.Black)
		gc.Clear()

		m = mBase.Multiply(RotateX(alpha)).Multiply(RotateY(beta))
		cube.Transform(m, cubeT)
		cubeT.Draw(gc)

		cloud.Transform(m, cloudT)
		cloudT.Draw(gc)

		zeroT = m.Transform(zero)
		xAxisT = m.Transform(xAxis)
		yAxisT = m.Transform(yAxis)
		zAxisT = m.Transform(zAxis)

		gc.SetStrokeWidth(4.0)
		gc.SetStrokeColor(colornames.DarkRed)
		gc.DrawLine(zeroT.X, zeroT.Y, xAxisT.X, xAxisT.Y)
		gc.Stroke()

		gc.SetStrokeColor(colornames.DarkGreen)
		gc.DrawLine(zeroT.X, zeroT.Y, yAxisT.X, yAxisT.Y)
		gc.Stroke()

		gc.SetStrokeColor(colornames.DarkBlue)
		gc.DrawLine(zeroT.X, zeroT.Y, zAxisT.X, zAxisT.Y)
		gc.Stroke()

		alpha += dAlpha
		beta += dBeta

		adatft.PaintWatch.Stop()

		disp.Draw(gc.Image())
	}
}

type Cube struct {
	Pts         []Vector
	LineWidth   float64
	StrokeColor color.Color
}

func NewCube(s float64) *Cube {
	c := &Cube{LineWidth: 3.0, StrokeColor: colornames.Silver}
	c.Pts = []Vector{
		NewVector(-s, -s, -s),
		NewVector(-s, s, -s),
		NewVector(s, s, -s),
		NewVector(s, -s, -s),
		NewVector(-s, -s, s),
		NewVector(-s, s, s),
		NewVector(s, s, s),
		NewVector(s, -s, s),
	}
	return c
}

func (c *Cube) Transform(a Matrix, d *Cube) {
	d.LineWidth = c.LineWidth
	d.StrokeColor = c.StrokeColor
	d.Pts = make([]Vector, len(c.Pts))
	for i, pt := range c.Pts {
		d.Pts[i] = a.Transform(pt)
	}
}

func (c *Cube) Draw(gc *gg.Context) {
	gc.SetStrokeWidth(c.LineWidth)
	gc.SetStrokeColor(c.StrokeColor)
	gc.MoveTo(c.Pts[0].X, c.Pts[0].Y)
	gc.LineTo(c.Pts[1].X, c.Pts[1].Y)
	gc.LineTo(c.Pts[2].X, c.Pts[2].Y)
	gc.LineTo(c.Pts[3].X, c.Pts[3].Y)
	gc.ClosePath()
	gc.MoveTo(c.Pts[4].X, c.Pts[4].Y)
	gc.LineTo(c.Pts[5].X, c.Pts[5].Y)
	gc.LineTo(c.Pts[6].X, c.Pts[6].Y)
	gc.LineTo(c.Pts[7].X, c.Pts[7].Y)
	gc.ClosePath()
	gc.DrawLine(c.Pts[0].X, c.Pts[0].Y, c.Pts[4].X, c.Pts[4].Y)
	gc.DrawLine(c.Pts[1].X, c.Pts[1].Y, c.Pts[5].X, c.Pts[5].Y)
	gc.DrawLine(c.Pts[2].X, c.Pts[2].Y, c.Pts[6].X, c.Pts[6].Y)
	gc.DrawLine(c.Pts[3].X, c.Pts[3].Y, c.Pts[7].X, c.Pts[7].Y)
	gc.Stroke()
}

type Cloud struct {
	Pts   []Vector
	Color color.Color
	Size  float64
}

func NewCloud(x, y, z, w float64, numObjs int, size float64) *Cloud {
	c := &Cloud{}
	c.Pts = make([]Vector, numObjs)
	for i := 0; i < numObjs; i++ {
		px := rand.NormFloat64()*w + x
		py := rand.NormFloat64()*w + y
		pz := rand.NormFloat64()*w + z
		c.Pts[i] = NewVector(px, py, pz)
	}
	c.Color = colornames.YellowGreen
	c.Size = size
	return c
}

func (c *Cloud) Transform(m Matrix, d *Cloud) {
	d.Color = c.Color
	d.Size = c.Size
	d.Pts = make([]Vector, len(c.Pts))
	for i, pt := range c.Pts {
		d.Pts[i] = m.Transform(pt)
	}
}

func (c *Cloud) Draw(gc *gg.Context) {
	gc.SetFillColor(c.Color)
	for _, pt := range c.Pts {
		gc.DrawPoint(pt.X, pt.Y, c.Size)
	}
	gc.Fill()
}

type Matrix struct {
	M11, M12, M13, M14 float64
	M21, M22, M23, M24 float64
	M31, M32, M33, M34 float64
}

type Matrix3 struct {
	M11, M12, M13 float64
	M21, M22, M23 float64
}

type Vector struct {
	X, Y, Z float64
}

func NewVector(x, y, z float64) Vector {
	return Vector{x, y, z}
}

func Identity() Matrix {
	return Matrix{
		1, 0, 0, 0,
		0, 1, 0, 0,
		0, 0, 1, 0,
	}
}

func Translate(v Vector) Matrix {
	return Matrix{
		1, 0, 0, v.X,
		0, 1, 0, v.Y,
		0, 0, 1, v.Z,
	}
}

func Scale(sv Vector) Matrix {
	return Matrix{
		sv.X, 0, 0, 0,
		0, sv.Y, 0, 0,
		0, 0, sv.Z, 0,
	}
}

func RotateX(angle float64) Matrix {
	c := math.Cos(angle)
	s := math.Sin(angle)
	return Matrix{
		1, 0, 0, 0,
		0, c, -s, 0,
		0, s, c, 0,
	}
}

func RotateY(angle float64) Matrix {
	c := math.Cos(angle)
	s := math.Sin(angle)
	return Matrix{
		c, 0, -s, 0,
		0, 1, 0, 0,
		s, 0, c, 0,
	}
}

func RotateZ(angle float64) Matrix {
	c := math.Cos(angle)
	s := math.Sin(angle)
	return Matrix{
		c, -s, 0, 0,
		s, c, 0, 0,
		0, 0, 1, 0,
	}
}

func (a Matrix) Multiply(b Matrix) Matrix {
	return Matrix{
		a.M11*b.M11 + a.M12*b.M21 + a.M13*b.M31,
		a.M11*b.M12 + a.M12*b.M22 + a.M13*b.M32,
		a.M11*b.M13 + a.M12*b.M23 + a.M13*b.M33,
		a.M11*b.M14 + a.M12*b.M24 + a.M13*b.M34 + a.M14,
		a.M21*b.M11 + a.M22*b.M21 + a.M23*b.M31,
		a.M21*b.M12 + a.M22*b.M22 + a.M23*b.M32,
		a.M21*b.M13 + a.M22*b.M23 + a.M23*b.M33,
		a.M21*b.M14 + a.M22*b.M24 + a.M23*b.M34 + a.M24,
		a.M31*b.M11 + a.M32*b.M21 + a.M33*b.M31,
		a.M31*b.M12 + a.M32*b.M22 + a.M33*b.M32,
		a.M31*b.M13 + a.M32*b.M23 + a.M33*b.M33,
		a.M31*b.M14 + a.M32*b.M24 + a.M33*b.M34 + a.M34,
	}
}

func (a Matrix) Det() float64 {
	return a.M11*a.M22*a.M33 + a.M12*a.M23*a.M31 + a.M13*a.M21*a.M32 -
		a.M11*a.M23*a.M32 - a.M12*a.M21*a.M33 - a.M13*a.M22*a.M31
}

func (a Matrix3) Det() float64 {
	return a.M11*a.M22 - a.M12*a.M21
}

//func (a Matrix) Kof(i, j int) (float64)

func (a Matrix) Inv() Matrix {
	det := a.Det()
	adj := Matrix{}
	adj.M11 = Matrix3{
		a.M22, a.M23, a.M24,
		a.M32, a.M33, a.M34,
	}.Det() / det
	adj.M12 = -Matrix3{
		a.M21, a.M23, a.M24,
		a.M31, a.M33, a.M34,
	}.Det() / det
	adj.M13 = Matrix3{
		a.M21, a.M22, a.M24,
		a.M31, a.M32, a.M34,
	}.Det() / det
	adj.M14 = -Matrix3{
		a.M21, a.M22, a.M23,
		a.M31, a.M32, a.M33,
	}.Det() / det

	adj.M21 = -Matrix3{
		a.M12, a.M13, a.M14,
		a.M32, a.M33, a.M34,
	}.Det() / det
	adj.M22 = Matrix3{
		a.M11, a.M13, a.M14,
		a.M31, a.M33, a.M34,
	}.Det() / det
	adj.M23 = -Matrix3{
		a.M11, a.M12, a.M14,
		a.M31, a.M32, a.M34,
	}.Det() / det
	adj.M24 = Matrix3{
		a.M11, a.M12, a.M13,
		a.M31, a.M32, a.M33,
	}.Det() / det

	adj.M31 = Matrix3{
		a.M12, a.M13, a.M14,
		a.M22, a.M23, a.M24,
	}.Det() / det
	adj.M32 = -Matrix3{
		a.M11, a.M13, a.M14,
		a.M21, a.M23, a.M24,
	}.Det() / det
	adj.M33 = Matrix3{
		a.M11, a.M12, a.M14,
		a.M21, a.M22, a.M24,
	}.Det() / det
	adj.M34 = -Matrix3{
		a.M11, a.M12, a.M13,
		a.M21, a.M22, a.M23,
	}.Det() / det

	return Matrix{adj.M11, adj.M21, adj.M31, adj.M14,
		adj.M12, adj.M22, adj.M32, adj.M24,
		adj.M13, adj.M23, adj.M33, adj.M34,
	}
}

func (a Matrix) Transform(v Vector) Vector {
	return NewVector(
		(a.M11*v.X + a.M12*v.Y + a.M13*v.Z + a.M14),
		(a.M21*v.X + a.M22*v.Y + a.M23*v.Z + a.M24),
		(a.M31*v.X + a.M32*v.Y + a.M33*v.Z + a.M34))
}

func (a Matrix) String() string {
	return fmt.Sprintf("[%.4v %.4v %.4v %.4v]\n[%.4v %.4v %.4v %.4v]\n[%.4v %.4v %.4v %.4v]",
		a.M11, a.M12, a.M13, a.M14,
		a.M21, a.M22, a.M23, a.M24,
		a.M31, a.M32, a.M33, a.M34)
}

var (
	left, right = -5.0, 5.0
	top, bottom = 3.75, -3.75
	near, far   = -5.0, 30.0
)

func ProjectionMatrix(xRot, yRot, dist float64) Matrix {
	scale := NewVector(320.0/(right-left), 240.0/(bottom-top),
		120.0/(far-near))
	trans := NewVector((0.0 - left), -(top - 0.0), (far+near)/2)

	rotMat := RotateY(yRot).Multiply(RotateX(xRot))
	transl := rotMat.Transform(NewVector(0, 0, -dist))
	camMat := Translate(NewVector(-transl.X, -transl.Y, -transl.Z)).Multiply(rotMat)
	ndcMat := Scale(scale).Multiply(Translate(trans))
	ndcMat = ndcMat.Multiply(Matrix{1, 0, 0, 0,
		0, 1, 0, 0,
		0, 0, (near + far) / near, -far})
	return ndcMat.Multiply(camMat.Inv())
}

// Animation of Text ----------------------------------------------------------
func TextAnimation() {
	var textList []*TextObject
	var fontList = []*fonts.Font{
		fonts.LucidaBright,
		fonts.LucidaBrightItalic,
		fonts.LucidaBrightDemibold,
		fonts.LucidaBrightDemiboldItalic,
	}

	textList = make([]*TextObject, numObjs)
	for i := 0; i < numObjs; i++ {
		textList[i] = NewTextObject(msg, fontList[i%len(fontList)],
			45.0+80.0*rand.Float64())
		yPos := 200.0*rand.Float64() + 20.0
		xVel := 4.0*rand.Float64() + 1.0
		if rand.Float64() < 0.5 {
			xVel *= -1.0
		}
		textList[i].SetAnimParam(yPos, xVel)
	}

	ticker := time.NewTicker(30 * time.Millisecond)
	for range ticker.C {
		if !runFlag {
			break
		}
		adatft.PaintWatch.Start()
		gc.SetFillColor(colornames.Black)
		gc.Clear()

		for _, txtObj := range textList {
			txtObj.Draw(gc)
		}
		for _, txtObj := range textList {
			if !txtObj.Animate() {
				yPos := 200.0*rand.Float64() + 20.0
				xVel := 4.0*rand.Float64() + 1.0
				if rand.Float64() < 0.5 {
					xVel *= -1.0
				}
				txtObj.SetAnimParam(yPos, xVel)
			}
		}
		adatft.PaintWatch.Stop()

		disp.Draw(gc.Image())
	}
}

type TextObject struct {
	x, y          float64
	txt           string
	face          font.Face
	width, height float64
	Color         color.Color
	xVel          float64
}

func NewTextObject(txt string, fnt *fonts.Font, size float64) *TextObject {
	o := &TextObject{}
	o.txt = txt
	o.face = fonts.NewFace(fnt, size)
	o.Color = colornames.RandColor()
	o.width = float64(font.MeasureString(o.face, o.txt)) / 64.0
	o.height = float64(o.face.Metrics().Ascent) / 64.0
	return o
}

func (o *TextObject) SetAnimParam(y, xVel float64) {
	if xVel > 0.0 {
		o.x = -o.width / 2.0
	} else {
		o.x = float64(adatft.Width) + o.width/2.0
	}
	o.y = y
	o.xVel = xVel
}

func (o *TextObject) Draw(gc *gg.Context) {
	gc.SetFontFace(o.face)
	gc.SetStrokeColor(o.Color)
	gc.DrawStringAnchored(o.txt, o.x, o.y, 0.5, 0.5)
}

func (o *TextObject) Animate() bool {
	o.x += o.xVel
	if o.x > float64(adatft.Width)+o.width/2.0 || o.x < -o.width/2.0 {
		return false
	}
	return true
}

// Manually show color fadings ------------------------------------------------
func FadingColors() {
	var r, g, b uint8
	var idx int
	var f1, f2 float64
	var cmd, mode byte

	img := gc.Image().(*image.RGBA)

MainLoop:
	for {
		fmt.Printf("Waehle:\n")
		fmt.Printf("  1: Rot\n")
		fmt.Printf("  2: Grün\n")
		fmt.Printf("  3: Blau\n")
		fmt.Printf("  4: Gelb\n")
		fmt.Printf("  5: Cyan\n")
		fmt.Printf("  6: Magenta\n")
		fmt.Printf("  7: Weiss\n\n")
		fmt.Printf("  a: Rot (Spalten) + Grün (Zeilen)\n")
		fmt.Printf("  b: Grün (Spalten) + Blau (Zeilen)\n")
		fmt.Printf("  c: Blau (Spalten) + Rot (Zeilen)\n\n")
		fmt.Printf("  x: Programm beenden\n")
		fmt.Scanf("%c", &cmd)

		switch cmd {
		case 'x':
			break MainLoop
		case '1', '2', '3', '4', '5', '6', '7':
			mode = cmd - '0'
		case 'a', 'b', 'c':
			mode = 10 + (cmd - 'a')
		default:
			continue MainLoop
		}

		idx = 0
		for row := 0; row < adatft.Height; row++ {
			f1 = float64(row) / float64(adatft.Height-1)
			for col := 0; col < adatft.Width; col++ {
				f2 = float64(col) / float64(adatft.Width-1)
				r, g, b = 0x00, 0x00, 0x00
				switch mode {
				case 1:
					r = uint8(255 * f2)
				case 2:
					g = uint8(255 * f2)
				case 3:
					b = uint8(255 * f2)
				case 4:
					r = uint8(255 * f2)
					g = uint8(255 * f2)
				case 5:
					g = uint8(255 * f2)
					b = uint8(255 * f2)
				case 6:
					r = uint8(255 * f2)
					b = uint8(255 * f2)
				case 7:
					r = uint8(255 * f2)
					g = uint8(255 * f2)
					b = uint8(255 * f2)
				case 10:
					r = uint8(255 * f2)
					g = uint8(255 * f1)
				case 11:
					g = uint8(255 * f2)
					b = uint8(255 * f1)
				case 12:
					b = uint8(255 * f2)
					r = uint8(255 * f1)
				}
				img.Pix[idx+0] = r
				img.Pix[idx+1] = g
				img.Pix[idx+2] = b
				idx += 4
			}
		}
		disp.Draw(gc.Image())
	}
}

// The famous plasma animation ------------------------------------------------
var (
	ColorFuncList = []ColorFuncType{
		ColorFunc01,
		ColorFunc02,
		ColorFunc03,
	}
)

const (
	numThreads = 3
)

func PlasmaAnimation() {
	var v, v1, v2, v3 float64
	var t float64
	var col, row, pixIdx, valIdx int
	var c color.Color
	var pal *Palette
	var orderQ [numThreads]chan float64
	var doneQ [numThreads]chan bool
	var valFlds [numThreads]*ValFieldType

	for i := 0; i < numThreads; i++ {
		orderQ[i] = make(chan float64)
		doneQ[i] = make(chan bool)
		valFlds[i] = NewValField(adatft.Width, adatft.Height,
			-dispWidth/2.0, dispWidth/2.0,
			dispHeight/2.0, -dispHeight/2.0,
			ColorFuncList[i])
		go UpdateThread(valFlds[i], orderQ[i], doneQ[i])
	}

	pal = NewPalette("Simpel",
		colornames.RandColor(),
		colornames.RandColor(),
		colornames.RandColor(),
		colornames.RandColor(),
		colornames.RandColor(),
	)

	img := gc.Image().(*image.RGBA)
	t = 0.0
	ticker := time.NewTicker(30 * time.Millisecond)
	for range ticker.C {
		if !runFlag {
			break
		}
		adatft.PaintWatch.Start()
		for i := range numThreads {
			orderQ[i] <- t
		}
		for i := range numThreads {
			<-doneQ[i]
		}

		pixIdx = 0
		valIdx = 0
		for row = 0; row < adatft.Height; row++ {
			for col = 0; col < adatft.Width; col++ {
				v1 = valFlds[0].Vals[valIdx]
				v2 = valFlds[1].Vals[valIdx]
				v3 = valFlds[2].Vals[valIdx]
				v = (v1 + v2 + v3 + 3.0) / 6.0
				c = pal.GetColor(v)
				r, g, b, _ := c.RGBA()
				img.Pix[pixIdx+0] = uint8(r >> 8)
				img.Pix[pixIdx+1] = uint8(g >> 8)
				img.Pix[pixIdx+2] = uint8(b >> 8)
				pixIdx += 4
				valIdx += 1
			}
		}
		t += dt
		adatft.PaintWatch.Stop()
		disp.Draw(gc.Image())
	}
	for i := range numThreads {
		close(orderQ[i])
	}
}

func UpdateThread(valFld *ValFieldType, orderQ chan float64, doneQ chan bool) {
	var t float64
	var ok bool

	for {
		if t, ok = <-orderQ; !ok {
			break
		}
		valFld.Update(t)
		doneQ <- true
	}
}

type ValFieldType struct {
	Vals               []float64
	Cols, Rows         int
	Xmin, Ymax, Dx, Dy float64
	Fnc                ColorFuncType
}

func NewValField(cols, rows int, xmin, xmax, ymin, ymax float64,
	fnc ColorFuncType) *ValFieldType {
	v := &ValFieldType{}
	v.Vals = make([]float64, cols*rows)
	v.Cols = cols
	v.Rows = rows
	v.Xmin = xmin
	v.Ymax = ymax
	v.Dx = (xmax - xmin) / float64(cols)
	v.Dy = (ymax - ymin) / float64(rows)
	v.Fnc = fnc
	return v
}

func (v *ValFieldType) Update(t float64) {
	var x, y float64
	var col, row, idx int

	y = v.Ymax
	idx = 0
	for row = 0; row < v.Rows; row++ {
		x = v.Xmin
		for col = 0; col < v.Cols; col++ {
			v.Vals[idx] = v.Fnc(x, y, t)
			x += v.Dx
			idx++
		}
		y -= v.Dy
	}
}

type ColorFuncType func(x, y, t float64) float64

func ColorFunc01(x, y, t float64) float64 {
	return math.Sin(x*f1p1 + t)
}

func ColorFunc02(x, y, t float64) float64 {
	return math.Sin(f2p1*(x*math.Sin(t/f2p2)+y*math.Cos(t/f2p3)) + t)
}

func ColorFunc03(x, y, t float64) float64 {
	cx := x + 0.5*math.Sin(t/f3p1)
	cy := y + 0.5*math.Cos(t/f3p2)
	return math.Sin(math.Sqrt(100.0*(cx*cx+cy*cy)+1.0) + t)
}

// -----------------------------------------------------------------------------
//
// Scroll-Text
const (
	textMargin  = 10.0
	textWidth   = 300.0
	fontSize    = 18.0
	lineSpacing = 1.3
)

var (
	fontList = []*fonts.Font{
		fonts.GoRegular,
		fonts.LucidaBright,
		fonts.LucidaSans,
		fonts.Seaford,
		fonts.Garamond,
		fonts.Elegante,
	}
)

func AnimatedText() {
	var face font.Face

	for _, font := range fontList {
		if quitFlag {
			break
		}
		face = fonts.NewFace(font, fontSize)
		ScrollText(BlindText, face, textMargin, 0.0, lineSpacing, true)
		runFlag = true
	}
}

func FadeText(gc *gg.Context, dsp *adatft.Display,
	txt string, face font.Face, x, y, lineSpacing float64, fadeIn bool) {
}

func ScrollText(txt string, face font.Face, x, y, lineSpacing float64,
	scrollUp bool) {
	var textList []string
	var textWidth float64
	var h1, h2, h float64
	var ticker *time.Ticker

	gc.SetFontFace(face)
	textWidth = float64(gc.Width()) - 2*x
	textList = gc.WordWrap(txt, textWidth)
	txt = strings.Join(textList, "\n")

	h1 = float64(gc.Height())
	_, h2 = gc.MeasureMultilineString(txt, lineSpacing)
	h = h1 + h2

	ticker = time.NewTicker(30 * time.Millisecond)
	for range ticker.C {
		if !runFlag {
			break
		}
		gc.SetFillColor(colornames.Black)
		gc.Clear()
		gc.SetStrokeColor(colornames.White)
		if scrollUp {
			y = h - h2
		} else {
			y = -(h - h1)
		}
		gc.DrawStringWrapped(txt, x, y, 0, 0, textWidth,
			lineSpacing, gg.AlignLeft)
		disp.Draw(gc.Image())
		if h -= 1.0; h < 0.0 {
			break
		}
	}
}

// -----------------------------------------------------------------------------

type Circle struct {
	mx, my, rx, ry         float64
	dx, dy, drx, dry       float64
	age, dage              float64
	strokeColor, fillColor color.Color
	borderWidth            float64
	waveSpace              float64
	//startTime time.Time
}

var (
	aspectRatio = 0.4
)

func NewCircle(w, h int) *Circle {
	c := &Circle{}
	c.strokeColor = colornames.RandGroupColor(colornames.Blues)
	c.borderWidth = 3.0
	c.Init()
	return c
}

func (c *Circle) Init() {
	c.mx = rand.Float64() * float64(adatft.Width)
	c.my = rand.Float64() * float64(adatft.Height)
	c.dx = 0.0
	c.dy = 0.0
	c.rx = 0.0
	c.ry = 0.0
	c.drx = 0.1*rand.NormFloat64() + 0.5
	c.dry = c.drx * aspectRatio
	c.age = 1.0
	c.dage = 0.001*rand.NormFloat64() + 0.004
	c.waveSpace = 50.0 * c.drx
}

func (c *Circle) Animate() bool {
	//if time.Now().Before(c.startTime) {
	//    return true
	//}
	c.mx += c.dx
	c.my += c.dy
	c.rx += c.drx
	c.ry += c.dry
	c.age -= c.dage
	if c.age <= 0.0 {
		return false
	}
	return true
}

func (c *Circle) Draw(gc *gg.Context) {
	gc.SetStrokeWidth(c.borderWidth)
	gc.SetStrokeColor(c.strokeColor.Alpha(c.age))
	for rx := c.rx; rx > 0 && rx > c.rx-4*c.waveSpace; rx -= c.waveSpace {
		gc.DrawEllipse(c.mx, c.my, rx, rx*aspectRatio)
		gc.Stroke()
	}
}

func CircleAnimation() {
	var circleList []*Circle

	circleList = make([]*Circle, 100)

	go func() {
		i := 0
		for runFlag {
			time.Sleep(time.Duration(rand.Intn(3000)) * time.Millisecond)
			for circleList[i] != nil {
				i = (i + 1) % len(circleList)
			}
			circleList[i] = NewCircle(adatft.Width, adatft.Height)
			i = (i + 1) % len(circleList)
		}
	}()

	gc.SetLineCapRound()
	gc.SetLineJoinRound()

	ticker := time.NewTicker(30 * time.Millisecond)
	for range ticker.C {
		if !runFlag {
			break
		}
		adatft.PaintWatch.Start()
		gc.SetFillColor(colornames.Black)
		gc.Clear()

		for i, c := range circleList {
			if c == nil {
				continue
			}
			c.Draw(gc)
			if !c.Animate() {
				circleList[i] = nil
			}
		}
		adatft.PaintWatch.Stop()
		disp.Draw(gc.Image())
	}
}

//----------------------------------------------------------------------------

func MatrixTest() {
	//m := Identity()
	//m := Matrix{1, 0, 0, 0,
	//            0, 2, 0, 0,
	//            0, 0, 3, 0}
	m := Matrix{4, 5, 6, 4,
		5, 1, 8, 1,
		6, 2, 9, 4}
	fmt.Printf("m:\n%v\n", m)
	mi := m.Inv()
	fmt.Printf("mi:\n%v\n", mi)
	p := ProjectionMatrix(10.0/180.0*math.Pi, 10.0/180.0*math.Pi, 10.0)
	fmt.Printf("p:\n%v\n", p)
}

// -----------------------------------------------------------------------------
//
// Text-Animation
// -----------------------------------------------------------------------------
//
// Bezier-Animation
func random() float64 {
	return rand.Float64()*2 - 1
}

func point() (x, y float64) {
	return random(), random()
}

type Curve struct {
	p1, c1, c2, p2 *Point
}

func NewCurve() (c *Curve) {
	c = &Curve{}
	c.p1 = &Point{random(), random(), 0.0, 0.0}
	c.c1 = &Point{random(), random(), 0.05 * random(), 0.05 * random()}
	c.c2 = &Point{random(), random(), 0.05 * random(), 0.05 * random()}
	c.p2 = &Point{random(), random(), 0.0, 0.0}
	return c
}

func (c *Curve) Move(xmin, xmax, ymin, ymax float64) {
	c.c1.Move(xmin, xmax, ymin, ymax)
	c.c2.Move(xmin, xmax, ymin, ymax)
}

func (c *Curve) Draw(gc *gg.Context) {
	gc.SetStrokeColor(colornames.Black)
	gc.SetFillColor(colornames.Black.Alpha(0.5))
	gc.MoveTo(c.p1.x, c.p1.y)
	gc.CubicTo(c.c1.x, c.c1.y, c.c2.x, c.c2.y, c.p2.x, c.p2.y)
	gc.SetStrokeWidth(3)
	gc.FillStroke()

	gc.MoveTo(c.p1.x, c.p1.y)
	gc.LineTo(c.c1.x, c.c1.y)
	gc.LineTo(c.c2.x, c.c2.y)
	gc.LineTo(c.p2.x, c.p2.y)

	gc.SetStrokeColor(colornames.Red.Alpha(0.5))
	gc.SetStrokeWidth(2)
	gc.Stroke()
}

const (
	S = 160
	W = 3
	H = 2
)

func BezierAnimation() {
	var curves []*Curve

	curves = make([]*Curve, W*H)
	for i := 0; i < W*H; i++ {
		curves[i] = NewCurve()
	}
	ticker := time.NewTicker(30 * time.Millisecond)
	for range ticker.C {
		if !runFlag {
			break
		}
		adatft.PaintWatch.Start()
		gc.SetFillColor(colornames.White)
		gc.Clear()
		for j := 0; j < H; j++ {
			for i := 0; i < W; i++ {
				x := float64(i)*S + S/2
				y := float64(j)*S + S/2
				gc.Push()
				gc.Translate(x, y)
				gc.Scale(S/2, S/2)
				curves[W*j+i].Move(-1, 1, -1, 1)
				curves[W*j+i].Draw(gc)
				gc.Pop()
			}
		}
		adatft.PaintWatch.Stop()
		disp.Draw(gc.Image())
	}
}

//-----------------------------------------------------------------------------

func Animation05(gc *gg.Context, disp *adatft.Display, d time.Duration) {
	var r, g, b uint8
	var idx int

	img := gc.Image().(*image.RGBA)
	for runFlag {
		adatft.PaintWatch.Start()
		idx = 0
		for y := 0; y < adatft.Height; y++ {
			for x := 0; x < adatft.Width; x++ {
				r, g, b = 0, 0, 0
				if y < 80 {
					r = uint8((x + adatft.PaintWatch.Num()) % 256)
				} else if y < 160 {
					g = uint8((x + adatft.PaintWatch.Num()) % 256)
				} else {
					b = uint8((x + adatft.PaintWatch.Num()) % 256)
				}
				img.Pix[idx+0] = r
				img.Pix[idx+1] = g
				img.Pix[idx+2] = b
				idx += 4
			}
		}
		adatft.PaintWatch.Stop()
		disp.Draw(gc.Image())
	}
}

func Animation051(gc *gg.Context, disp *adatft.Display, d time.Duration) {
	var r, g, b uint8
	var idx int
	var f, dColor float64

	dColor = 255.0 / 319.0

	img := gc.Image().(*image.RGBA)
	for mode := 0; mode < 6; mode++ {
		idx = 0
		for row := 0; row < adatft.Height; row++ {
			for col := 0; col < adatft.Width; col++ {
				f = float64(col)
				r, g, b = 0x00, 0x00, 0x00
				switch mode {
				case 0:
					r = uint8(f * dColor)
				case 1:
					g = uint8(f * dColor)
				case 2:
					b = uint8(f * dColor)
				case 3:
					r = uint8(f * dColor)
					g = uint8(f * dColor)
				case 4:
					g = uint8(f * dColor)
					b = uint8(f * dColor)
				case 5:
					r = uint8(f * dColor)
					b = uint8(f * dColor)
				}
				img.Pix[idx+0] = r
				img.Pix[idx+1] = g
				img.Pix[idx+2] = b
				idx += 4
			}
		}
		disp.DrawSync(gc.Image())
		time.Sleep(2 * time.Second)
	}
}

//-----------------------------------------------------------------------------

type Palette struct {
	name      string
	colorList []color.Color
	shadeList []color.Color
}

func NewPalette(name string, colors ...color.Color) *Palette {
	p := &Palette{}
	p.name = name
	for _, color := range colors {
		p.colorList = append(p.colorList, color)
	}
	p.CalcLinearShades()
	return p
}

func (p *Palette) GetColor(t float64) color.Color {
	idx := int(t * float64(numShades*(len(p.colorList)-1)))
	return p.shadeList[idx]
}

// Berechnet die Abstufung zwischen den Stuetzfarben einfach, d.h.
// linear.
func (p *Palette) CalcLinearShades() {
	var i, j, k int
	var color1, color2 color.Color
	var t float64

	p.shadeList = make([]color.Color, numShades*(len(p.colorList)-1))
	for i = 0; i < len(p.colorList)-1; i++ {
		for j = 0; j < numShades; j++ {
			color1 = p.colorList[i]
			color2 = p.colorList[i+1]
			t = float64(j) / float64(numShades)
			k = i*numShades + j
			p.shadeList[k] = color1.Interpolate(color2, t)
		}
	}
}

// Die folgende Funktion berechnet nach dem Verfahren der dividierten
// Differenzen die Koeffizienten fuer ein Newton-Polynom. Diese Funktion
// geht davon aus, dass die Stuetzstellen aequidistant sind. Der Array
// f enthaelt folglich nur die Funktionswerte.
func DivDiff(f []float64) []float64 {
	var table []float64
	var res []float64
	var n int = len(f)
	var i1, i2 int

	table = make([]float64, n*(n+1)/2)
	res = make([]float64, n)

	for i, v := range f {
		table[i] = v
	}
	i1 = 0
	i2 = n
	res[0] = f[0]

	for i := 1; i < n; i++ {
		for k := 0; k < n-i; k++ {
			table[i2+k] = (table[i1+k+1] - table[i1+k]) / float64(i)
		}
		i1 = i2
		i2 += (n - i)
		res[i] = table[i1]
	}
	return res
}

// Und mit dieser Funktion wird das Newton-Polynom mit den Koeffizienten a
// an der Stelle x ausgwertet.
//
// func NewtonPoly(x float64, a []float64) (float64) {
//     var res float64
//     var n int = len(a)
//
//     res = a[n-1]
//     for i:=n-2; i>=0; i-- {
//         res = res*(x-float64(i)) + a[i]
//     }
//     return res
// }

// Berechnet die Abstufung zwischen den Stuetzfarben mit der Approximation
// durch Newton-Polynome
//
// func (p *Palette) CalcNewtonShades() {
//     var i, j, k int
//     var x float64
//     var fRed, aRed, fGreen, aGreen, fBlue, aBlue []float64
//     var n int = len(p.colorList)/3
//
//     fRed   = make([]float64, n)
//     fGreen = make([]float64, n)
//     fBlue  = make([]float64, n)
//     for i=0; i<n; i++ {
//         fRed[i]   = float64(p.colorList[3*i])
//         fGreen[i] = float64(p.colorList[3*i+1])
//         fBlue[i]  = float64(p.colorList[3*i+2])
//     }
//     aRed   = DivDiff(fRed)
//     aGreen = DivDiff(fGreen)
//     aBlue  = DivDiff(fBlue)
//
//     p.shadeList = make([]uint8, 3 * (numShades * (n-1)))
//     for i=0; i<(n-1); i++ {
//         for j=0; j<numShades; j++ {
//             x = float64(i) + float64(j)/float64(numShades)
//             k = 3 * (i*numShades + j)
//             p.shadeList[k+0] = uint8(NewtonPoly(x, aRed))
//             p.shadeList[k+1] = uint8(NewtonPoly(x, aGreen))
//             p.shadeList[k+2] = uint8(NewtonPoly(x, aBlue))
//         }
//     }
// }

const (
	numShades  = 256
	dispWidth  = 1.6
	dispHeight = 1.2
	dt         = 0.05

	f1p1 = 10.0

	f2p1 = 10.0
	f2p2 = 2.0
	f2p3 = 3.0

	f3p1 = 5.0
	f3p2 = 3.0
)

/*
func ColorFunc01(dat animData) (float64) {
    return math.Sin(dat.x * f1p1 + dat.t)
}

func ColorFunc02(dat animData) (float64) {
    return math.Sin(f2p1*(dat.x*math.Sin(dat.t/f2p2)+dat.y*math.Cos(dat.t/f2p3))+dat.t)
}

func ColorFunc03(dat animData) (float64) {
    cx := dat.x + 0.5 * math.Sin(dat.t/f3p1)
    cy := dat.y + 0.5 * math.Cos(dat.t/f3p2)
    return math.Sin(math.Sqrt(100.0*(cx*cx + cy*cy)+1.0)+dat.t)
}
*/

//-----------------------------------------------------------------------------

func SBBAnimation() {
	var imgList []image.Image
	var imgFileList []string = []string{
		"sbbWatch/dial.png",
		"sbbWatch/hour.png",
		"sbbWatch/minute.png",
		"sbbWatch/second.png",
	}
	var seconds, minutes, hours float64
	var rotationList []float64
	var xm, ym float64
	var err error

	c1 := 1.0 / 60.0
	c2 := 1.0 / 12.0
	c3 := 2.0 * math.Pi

	imgList = make([]image.Image, len(imgFileList))
	for i, fileName := range imgFileList {
		imgList[i], err = gg.LoadPNG(fileName)
		if err != nil {
			log.Fatalf("error loading image: %v", err)
		}
	}
	rotationList = make([]float64, len(imgFileList))
	xm = float64(adatft.Width) / 2.0
	ym = float64(adatft.Height) / 2.0

	gc.SetFillColor(colornames.DeepSkyBlue)
	gc.Clear()
	disp.Draw(gc.Image())

	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	for t := range ticker.C {
		if !runFlag {
			break
		}
		adatft.PaintWatch.Start()
		seconds = float64(t.Second()) * c1
		minutes = (float64(t.Minute()) + seconds) * c1
		hours = (float64(t.Hour()) + minutes) * c2
		rotationList[0] = 0.0
		rotationList[1] = c3 * hours
		rotationList[2] = c3 * minutes
		rotationList[3] = c3 * seconds

		gc.SetFillColor(colornames.DeepSkyBlue)
		gc.Clear()
		for i, r := range rotationList {
			gc.Push()
			gc.RotateAbout(r, float64(xm), float64(ym))
			gc.DrawImageAnchored(imgList[i], xm, ym, 0.5, 0.5)
			gc.Pop()
		}
		adatft.PaintWatch.Stop()
		disp.Draw(gc.Image())
	}
}

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

func SignalHandler() {
	sigChan := make(chan os.Signal)
	signal.Notify(sigChan, os.Interrupt)
	<-sigChan
	runFlag = false
	quitFlag = true
}

func TouchHandler() {
	for penEvent := range touch.EventQ {
		if penEvent.Type == adatft.PenRelease {
			runFlag = false
		}
	}
}

type displayPageType struct {
	description string
	pageFunc    func()
}

var (
	displayPageList = []displayPageType{
		{"Dancing Polygons", PolygonAnimation},
		{"Rotating Cube (3D)", Cube3DAnimation},
		{"Text on the run", TextAnimation},
		{"Beziers wherever you look", BezierAnimation},
		{"Let's fade the colors", FadingColors},
		{"Plasma (dont burn yourself!)", PlasmaAnimation},
		{"Fading Circles", CircleAnimation},
		{"SBB (are you Swiss?)", SBBAnimation},
		{"Scrolling Text", AnimatedText},
		//		{"Matrix Tests", MatrixTest},
	}
)

var (
	IntroText         string = "Im Folgenden habe ich einige kleine Beispiele, Animationen oder Interaktionen zusammengestellt, um die Möglichkeiten des TFT-Displays mit Go zu demonstrieren Sämtliche Animationen werden direkt gerechnet. Die Beispiele laufen jeweils unbegrenzt, für den Wechsel zwischen den Beispielen, verwende die Pfeil-Buttons unten links und rechts."
	BlindText         string = "Lorem ipsum dolor sit amet, consetetur sadipscing elitr, sed diam nonumy eirmod tempor invidunt ut labore et dolore magna aliquyam erat, sed diam voluptua. At vero eos et accusam et justo duo dolores et ea rebum."
	cpuprofile        string
	memprofile        string
	disp              *adatft.Display
	touch             *adatft.Touch
	gc                *gg.Context
	pageNum           int
	numObjs, numEdges int
	blurFactor        float64
	msg               string
	rotation          adatft.RotationType = adatft.Rotate090
	runFlag, quitFlag bool
)

func main() {
	flag.IntVar(&pageNum, "page", 0, "Start with a given Page")
	flag.IntVar(&numObjs, "objs", 5, "Number of objects")
	flag.IntVar(&numEdges, "edges", 3, "Number of edges of an object")
	flag.Float64Var(&blurFactor, "blur", 1.0, "(Only for Anim 1) Blur factor [0,1] (1: no blur, 0: max blur).\nIn order to see something, choose a value < 0.1")
	flag.StringVar(&msg, "text", "Hello, world!", "The text that will be displayed in animation 3")
	flag.Var(&rotation, "rotation", "Display rotation")
	flag.Parse()

	StartProfiling()
	disp = adatft.OpenDisplay(rotation)
	touch = adatft.OpenTouch()
	gc = gg.NewContext(adatft.Width, adatft.Height)

	go SignalHandler()
	go TouchHandler()

	quitFlag = false
	for !quitFlag {
		runFlag = true
		log.Printf("[%d] %s", pageNum, displayPageList[pageNum].description)
		displayPageList[pageNum].pageFunc()
		pageNum = (pageNum + 1) % len(displayPageList)
	}
	disp.Close()
	touch.Close()
	StopProfiling()

	adatft.PrintStat()
}
