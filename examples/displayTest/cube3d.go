package main

import (
    "fmt"
    "math"
    "math/rand"
    "time"
    "github.com/stefan-muehlebach/gg"
    "github.com/stefan-muehlebach/gg/color"
)

type Cube3DAnim struct {
    gc *gg.Context
	mBase, m Matrix
	cube, cubeT *Cube
	cloud, cloudT *Cloud
	zero, xAxis, yAxis, zAxis Vector
	zeroT, xAxisT, yAxisT, zAxisT Vector
	alpha, dAlpha float64
	beta, dBeta float64
}

func (a *Cube3DAnim) RefreshTime() time.Duration {
    return 30 * time.Millisecond
}

func (a *Cube3DAnim) Init(gc *gg.Context) {
    a.gc = gc
   
	a.cube = NewCube(70.0)
	a.cubeT = &Cube{}

	a.cloud = NewCloud(0, 0, 0, 70.0, 5*numObjs, 4.0)
	a.cloudT = &Cloud{}

	a.zero = NewVector(0.0, 0.0, 0.0)
	a.xAxis = NewVector(70.0, 0.0, 0.0)
	a.yAxis = NewVector(0.0, 70.0, 0.0)
	a.zAxis = NewVector(0.0, 0.0, 70.0)

	a.alpha = math.Pi / 12.0
	a.dAlpha = math.Pi / 162.0
	a.beta = math.Pi / 18.0
	a.dBeta = math.Pi / 126.0

	a.mBase = Identity().Multiply(Scale(NewVector(1.0, -1.0, 1.0))).Multiply(Translate(NewVector(240.0, -160.0, 0.0)))
}

func (a *Cube3DAnim) Paint() {
	a.gc.SetFillColor(color.Black)
	a.gc.Clear()

	a.m = a.mBase.Multiply(RotateX(a.alpha)).Multiply(RotateY(a.beta))
	a.cube.Transform(a.m, a.cubeT)
	a.cubeT.Draw(a.gc)

	a.cloud.Transform(a.m, a.cloudT)
	a.cloudT.Draw(a.gc)

	a.zeroT = a.m.Transform(a.zero)
	a.xAxisT = a.m.Transform(a.xAxis)
	a.yAxisT = a.m.Transform(a.yAxis)
	a.zAxisT = a.m.Transform(a.zAxis)

	a.gc.SetStrokeWidth(4.0)
	a.gc.SetStrokeColor(color.DarkRed)
	a.gc.DrawLine(a.zeroT.X, a.zeroT.Y, a.xAxisT.X, a.xAxisT.Y)
	a.gc.Stroke()

	a.gc.SetStrokeColor(color.DarkGreen)
	a.gc.DrawLine(a.zeroT.X, a.zeroT.Y, a.yAxisT.X, a.yAxisT.Y)
	a.gc.Stroke()

	a.gc.SetStrokeColor(color.DarkBlue)
	a.gc.DrawLine(a.zeroT.X, a.zeroT.Y, a.zAxisT.X, a.zAxisT.Y)
	a.gc.Stroke()

	a.alpha += a.dAlpha
	a.beta += a.dBeta
}

func (a *Cube3DAnim) Clean() {}

// 3D-Animation ---------------------------------------------------------------
/*
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

		gc.SetFillColor(color.Black)
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
		gc.SetStrokeColor(color.DarkRed)
		gc.DrawLine(zeroT.X, zeroT.Y, xAxisT.X, xAxisT.Y)
		gc.Stroke()

		gc.SetStrokeColor(color.DarkGreen)
		gc.DrawLine(zeroT.X, zeroT.Y, yAxisT.X, yAxisT.Y)
		gc.Stroke()

		gc.SetStrokeColor(color.DarkBlue)
		gc.DrawLine(zeroT.X, zeroT.Y, zAxisT.X, zAxisT.Y)
		gc.Stroke()

		alpha += dAlpha
		beta += dBeta

		adatft.PaintWatch.Stop()

		Draw(gc, disp)
	}
}
*/

type Cube struct {
	Pts         []Vector
	LineWidth   float64
	StrokeColor color.Color
}

func NewCube(s float64) *Cube {
	c := &Cube{LineWidth: 3.0, StrokeColor: color.Silver}
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
	c.Color = color.YellowGreen
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

