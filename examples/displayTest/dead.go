//go:build ignore

package main

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
		Draw(gc, disp)
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
	c.strokeColor = color.RandGroupColor(color.Blues)
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
		gc.SetFillColor(color.Black)
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
		Draw(gc, disp)
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
	gc.SetStrokeColor(color.Black)
	gc.SetFillColor(color.Black.Alpha(0.5))
	gc.MoveTo(c.p1.x, c.p1.y)
	gc.CubicTo(c.c1.x, c.c1.y, c.c2.x, c.c2.y, c.p2.x, c.p2.y)
	gc.SetStrokeWidth(3)
	gc.FillStroke()

	gc.MoveTo(c.p1.x, c.p1.y)
	gc.LineTo(c.c1.x, c.c1.y)
	gc.LineTo(c.c2.x, c.c2.y)
	gc.LineTo(c.p2.x, c.p2.y)

	gc.SetStrokeColor(color.Red.Alpha(0.5))
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
		gc.SetFillColor(color.White)
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
		Draw(gc, disp)
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
		Draw(gc, disp)
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

