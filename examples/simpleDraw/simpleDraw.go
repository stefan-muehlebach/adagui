package main

import (
	"flag"
    "os"
    _ "image/png"
	//"image/color"
	"github.com/stefan-muehlebach/adagui"
	"github.com/stefan-muehlebach/adagui/binding"
	"github.com/stefan-muehlebach/adagui/touch"
	"github.com/stefan-muehlebach/adatft"
	"log"
	"math/rand"
	_ "sync"
    "encoding/json"
	//"github.com/stefan-muehlebach/gg/color"
	"github.com/stefan-muehlebach/gg/colornames"
	"github.com/stefan-muehlebach/gg/geom"
    "github.com/cpmech/gosl/fun/fftw"
)

// Dieses Programm dient der Demonstration, wie mit dem Touchscreen umgegangen
// werden kann. Auf einem Panel können je nach Werkzeug Kreise oder Rechtecke
// erstellt und verschoben werden..

type ToolType int

const (
	CircleTool ToolType = iota
	RectangleTool
	EllipseTool
    PolygonTool
)

func init() {
	log.SetFlags(log.Lmicroseconds | log.Lmsgprefix)
	log.SetPrefix(": ")
}

var (
	screen *adagui.Screen
	win    *adagui.Window
	tool   ToolType
    poly *adagui.Polygon
    outFile string
)

//-----------------------------------------------------------------------------

type Complex struct {
    Re, Im float64
}

func NewComplex(c complex128) Complex {
    return Complex{real(c), imag(c)}
}

func (c *Complex) AsComplex() complex128 {
    return complex(c.Re, c.Im)
}

//-----------------------------------------------------------------------------

// Erstellt ein neues Panel der angegebenen Groesse und legt alle wichtigen
// Handler fuer das Touch-Event fest.
func NewPanel(w, h float64) *adagui.Panel {
	var circ *adagui.Circle
	var rect *adagui.Rectangle
	var elli *adagui.Ellipse

	p := adagui.NewPanel(w, h)
    p.SetColor(colornames.DarkSlateGrey)
//    fh, _ := os.Open("taube.png")
//    p.Image, _ = png.Decode(fh)

	p.SetOnTap(func(evt touch.Event) {
		switch tool {
		case CircleTool:
			r := 30.0 + 10.0*rand.Float64()
			circ = NewCircle(r)
			circ.SetPos(evt.Pos)
			p.Add(circ)
		case RectangleTool:
			w, h := 60.0+20.0*rand.Float64(), 45.0+15.0*rand.Float64()
			rect = NewRectangle(w, h)
			rect.SetPos(evt.Pos.Sub(geom.Point{w / 2, h / 2}))
			p.Add(rect)
		case EllipseTool:
			rx, ry := 30.0+20.0*rand.Float64(), 20.0+10.0*rand.Float64()
			elli = NewEllipse(rx, ry)
			elli.SetPos(evt.Pos)
			p.Add(elli)
		}
		p.Mark(adagui.MarkNeedsPaint)
	})

	p.SetOnLongPress(func(evt touch.Event) {
		switch tool {
		case CircleTool:
			circ = NewCircle(1.0)
			circ.SetPos(evt.Pos)
			p.Add(circ)
		case RectangleTool:
			rect = NewRectangle(1.0, 1.0)
			rect.SetPos(evt.Pos)
			p.Add(rect)
		case EllipseTool:
			elli = NewEllipse(1.0, 1.0)
			elli.SetPos(evt.Pos)
			p.Add(elli)
        case PolygonTool:
            if poly != nil {
                poly.Remove()
            }
            poly = adagui.NewPolygon(evt.Pos)
            poly.Closed = true
            p.Add(poly)
		}
		p.Mark(adagui.MarkNeedsPaint)
	})

	p.SetOnDrag(func(evt touch.Event) {
		if !evt.LongPressed {
			return
		}
		switch tool {
		case CircleTool:
			r := evt.Pos.Distance(evt.InitPos)
			circ.SetRadius(r)
		case RectangleTool:
			d := evt.Pos.Sub(evt.InitPos)
			rect.SetSize(d)
		case EllipseTool:
			rx, ry := evt.Pos.Sub(evt.InitPos).AsCoord()
			elli.SetRadius(rx, ry)
        case PolygonTool:
            poly.AddPoint(evt.Pos)
		}
		p.Mark(adagui.MarkNeedsPaint)
	})

	p.SetOnRelease(func(evt touch.Event) {
		if !evt.LongPressed {
			return
		}
		switch tool {
		case RectangleTool:
			r := rect.Rect().Canon()
			rect.SetPos(r.Min)
			rect.SetSize(r.Size())
			rect.Mark(adagui.MarkNeedsPaint)
		case EllipseTool:
			r := elli.Rect().Canon()
			elli.SetPos(r.Min)
			elli.SetSize(r.Size())
			elli.Mark(adagui.MarkNeedsPaint)
        case PolygonTool:
            poly.Flatten()
			poly.Mark(adagui.MarkNeedsPaint)
		}
	})

	return p
}

// Erstellt einen neuen Kreis und definiert alle Handler, welche fuer diese
// Objekte spezifisch sind.
func NewCircle(r float64) *adagui.Circle {
	var dp geom.Point

	c := adagui.NewCircle(r)
	col := colornames.RandGroupColor(colornames.Blues)
	c.SetColor(col)
	c.SetPushedColor(col.Alpha(0.5))

	c.SetOnPress(func(evt touch.Event) {
		dp = evt.Pos.Sub(c.Pos())
		c.Mark(adagui.MarkNeedsPaint)
	})

	c.SetOnDrag(func(evt touch.Event) {
		c.SetPos(evt.Pos.Sub(dp))
		c.Mark(adagui.MarkNeedsPaint)
	})

	c.SetOnRelease(func(evt touch.Event) {
		c.Mark(adagui.MarkNeedsPaint)
	})

	c.SetOnLongPress(func(evt touch.Event) {
		c.ToBack()
		c.Mark(adagui.MarkNeedsPaint)
	})

	c.SetOnTap(func(evt touch.Event) {
		if evt.LongPressed {
			return
		}
		c.ToFront()
		c.Mark(adagui.MarkNeedsPaint)
	})

	c.SetOnDoubleTap(func(evt touch.Event) {
		p := c.Wrappee().Parent
		c.Remove()
		p.Mark(adagui.MarkNeedsPaint)
	})

	return c
}

// Erstellt eine neue Ellipse.
func NewEllipse(rx, ry float64) *adagui.Ellipse {
	var dp geom.Point

	e := adagui.NewEllipse(rx, ry)
	col := colornames.RandGroupColor(colornames.Greens)
	e.SetColor(col)
	e.SetPushedColor(col.Alpha(0.5))

	e.SetOnPress(func(evt touch.Event) {
		dp = evt.Pos.Sub(e.Pos())
		e.Mark(adagui.MarkNeedsPaint)
	})

	e.SetOnDrag(func(evt touch.Event) {
		e.SetPos(evt.Pos.Sub(dp))
		e.Mark(adagui.MarkNeedsPaint)
	})

	e.SetOnRelease(func(evt touch.Event) {
		e.Mark(adagui.MarkNeedsPaint)
	})

	e.SetOnLongPress(func(evt touch.Event) {
		e.ToBack()
		e.Mark(adagui.MarkNeedsPaint)
	})

	e.SetOnTap(func(evt touch.Event) {
		if evt.LongPressed {
			return
		}
		e.ToFront()
		e.Mark(adagui.MarkNeedsPaint)
	})

	e.SetOnDoubleTap(func(evt touch.Event) {
		p := e.Wrappee().Parent
		e.Remove()
		p.Mark(adagui.MarkNeedsPaint)
	})

	return e
}

// Erstellt ein neues Rechteck und definiert alle Handler, welche fuer diese
// Objekte spezifisch sind.
func NewRectangle(w, h float64) *adagui.Rectangle {
	var dp geom.Point

	r := adagui.NewRectangle(w, h)
	col := colornames.RandGroupColor(colornames.Reds)
	r.SetColor(col)
	r.SetPushedColor(col.Alpha(0.5))

	r.SetOnPress(func(evt touch.Event) {
		dp = evt.Pos.Sub(r.Pos())
		r.Mark(adagui.MarkNeedsPaint)
	})

	r.SetOnDrag(func(evt touch.Event) {
		r.SetPos(evt.Pos.Sub(dp))
		r.Mark(adagui.MarkNeedsPaint)
	})

	r.SetOnRelease(func(evt touch.Event) {
		r.Mark(adagui.MarkNeedsPaint)
	})

	r.SetOnLongPress(func(evt touch.Event) {
		r.ToBack()
		r.Mark(adagui.MarkNeedsPaint)
	})

	r.SetOnTap(func(evt touch.Event) {
		if evt.LongPressed {
			return
		}
		r.ToFront()
		r.Mark(adagui.MarkNeedsPaint)
	})

	r.SetOnDoubleTap(func(evt touch.Event) {
		p := r.Wrappee().Parent
		r.Remove()
		p.Mark(adagui.MarkNeedsPaint)
	})

	return r
}

// Hauptprogramm.
func main() {
    flag.StringVar(&outFile, "out", "coeff.json", "Output File")
    flag.Parse()

	screen = adagui.NewScreen(adatft.Rotate090)
	win = screen.NewWindow()

	group := adagui.NewGroupPL(nil, adagui.NewVBoxLayout())
	win.SetRoot(group)

	panel := NewPanel(float64(adatft.Width-10), float64(adatft.Height-10-40))
	//panel.SetColor(colornames.Gray)
	group.Add(panel)

    btnGrp := adagui.NewGroupPL(group, adagui.NewHBoxLayout())

	btnData := binding.NewInt()
	btnData.Set(-1)

	btn5 := adagui.NewIconButtonWithData("icons/line.png", int(PolygonTool), btnData)
	btn5.SetOnTap(func(evt touch.Event) {
		tool = PolygonTool
	})

	btn2 := adagui.NewIconButtonWithData("icons/circle.png", int(CircleTool), btnData)
	btn2.SetOnTap(func(evt touch.Event) {
		tool = CircleTool
	})

	btn3 := adagui.NewIconButtonWithData("icons/rectangle.png", int(RectangleTool), btnData)
	btn3.SetOnTap(func(evt touch.Event) {
		tool = RectangleTool
	})

	btn4 := adagui.NewIconButtonWithData("icons/ellipse.png", int(EllipseTool), btnData)
	btn4.SetOnTap(func(evt touch.Event) {
		tool = EllipseTool
	})

    btnFFT := adagui.NewTextButton("FFT")
    btnFFT.SetOnTap(func(evt touch.Event) {
        if poly == nil {
            return
        }
        pts := poly.Points()
        data := make([]complex128, len(pts))
        out  := make([]Complex, len(pts))
        fftPlan := fftw.NewPlan1d(data, false, true)
        for i, pt := range pts {
            data[i] = complex(pt.X, pt.Y)
        }
        fftPlan.Execute()
        n := complex(float64(len(data)), 0.0)
        for i, dat := range data {
            out[i] = NewComplex(dat / n)
        }

        fh, err := os.Create(outFile)
        if err != nil {
            log.Fatal(err)
        }
        b, err := json.Marshal(out)
        if err != nil {
            log.Fatal(err)
        }
        fh.Write(b)
        fh.Close()
/*
        fmt.Printf("package main\n")
        fmt.Printf("var (\n")
        fmt.Printf("    CoeffList = []FourierCoeff{\n")
        for i, dat := range data[:len(data)/2] {
            fmt.Printf("        FourierCoeff{%3d, %+9.4f},\n", i, dat/n)
            if i > 0 {
                fmt.Printf("        FourierCoeff{%3d, %+9.4f},\n", -i, data[len(data)-i]/n)
            }
        }
        fmt.Printf("    }\n")
        fmt.Printf(")\n")
*/
        fftPlan.Free()
    })

	btnGrp.Add(btn5, btn2, btn3, btn4, btnFFT, adagui.NewSpacer())

	btnQuit := adagui.NewTextButton("Quit")
	btnGrp.Add(btnQuit)
	btnQuit.SetOnTap(func(evt touch.Event) {
		screen.Quit()
	})

	screen.SetWindow(win)
	screen.Run()
}
