package main

import (
	_ "fmt"
	_ "image"
	//"image/color"
	"github.com/stefan-muehlebach/adagui"
	"github.com/stefan-muehlebach/adagui/binding"
	"github.com/stefan-muehlebach/adagui/touch"
	"github.com/stefan-muehlebach/adatft"
	"log"
	"math/rand"
	_ "sync"
	//"github.com/stefan-muehlebach/gg/color"
	"github.com/stefan-muehlebach/gg/colornames"
	"github.com/stefan-muehlebach/gg/geom"
)

// Dieses Programm dient der Demonstration, wie mit dem Touchscreen umgegangen
// werden kann. Auf einem Panel können je nach Werkzeug Kreise oder Rechtecke
// erstellt und verschoben werden..

type ToolType int

const (
	CircleTool ToolType = iota
	RectangleTool
	EllipseTool
)

func init() {
	log.SetFlags(log.Lmicroseconds | log.Lmsgprefix)
	log.SetPrefix(": ")
}

var (
	screen *adagui.Screen
	win    *adagui.Window
	tool   ToolType
)

// Erstellt ein neues Panel der angegebenen Groesse und legt alle wichtigen
// Handler fuer das Touch-Event fest.
func NewPanel(w, h float64) *adagui.Panel {
	var circ *adagui.Circle
	var rect *adagui.Rectangle
	var elli *adagui.Ellipse

	p := adagui.NewPanel(w, h)

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
	screen = adagui.NewScreen(adatft.Rotate090)
	win = screen.NewWindow()

	group := adagui.NewGroup()
	win.SetRoot(group)

	btn1 := adagui.NewTextButton("Quit")
	w, h := btn1.Size().AsCoord()
	btn1.SetPos(geom.Point{320 - w - 5, 240 - h - 5})
	group.Add(btn1)
	btn1.SetOnTap(func(evt touch.Event) {
		screen.Quit()
	})

	btnData := binding.NewInt()
	btnData.Set(-1)
	btn2 := adagui.NewIconButtonWithData("icons/circle.png", int(CircleTool), btnData)
	btn2.SetPos(geom.Point{5, btn1.Rect().Y0()})
	group.Add(btn2)
	btn2.SetOnTap(func(evt touch.Event) {
		tool = CircleTool
	})

	btn3 := adagui.NewIconButtonWithData("icons/rectangle.png", int(RectangleTool), btnData)
	btn3.SetPos(geom.Point{btn2.Rect().X1() + 5, btn1.Rect().Y0()})
	group.Add(btn3)
	btn3.SetOnTap(func(evt touch.Event) {
		tool = RectangleTool
	})

	btn4 := adagui.NewIconButtonWithData("icons/ellipse.png", int(EllipseTool), btnData)
	btn4.SetPos(geom.Point{btn3.Rect().X1() + 5, btn1.Rect().Y0()})
	group.Add(btn4)
	btn4.SetOnTap(func(evt touch.Event) {
		tool = EllipseTool
	})

	panel := NewPanel(310, btn1.Rect().Y0()-10)
	panel.SetPos(geom.Point{5, 5})
	panel.SetColor(colornames.Gray)
	group.Add(panel)

	screen.SetWindow(win)
	screen.Run()
}
