package main

import (
	"flag"
	"fmt"
	"github.com/stefan-muehlebach/adagui"
	"github.com/stefan-muehlebach/adagui/binding"
	"github.com/stefan-muehlebach/adagui/props"
	"github.com/stefan-muehlebach/adagui/touch"
	"github.com/stefan-muehlebach/adatft"
	"github.com/stefan-muehlebach/gg/color"
	"github.com/stefan-muehlebach/gg/geom"
	"os"
	"os/signal"
)

type Tool int

const (
	HandTool Tool = iota
	ErasureTool
	PointTool
	SegmentTool
	RayTool
	LineTool
	CircleTool
	PolygonTool
	MidpointTool
	PerpendBisecTool
	AngleBisectTool
	PerpendTool
	ParallelTool
	MoveTool
)

type ToolData struct {
	Idx      Tool
	IconFile string
}

var (
	ToolList = []ToolData{
		{HandTool, "90.png"},
		{ErasureTool, "91.png"},
		{PointTool, "01.png"},
		{SegmentTool, "02.png"},
		{RayTool, "03.png"},
		{LineTool, "04.png"},
		{CircleTool, "05.png"},
		{PolygonTool, "06.png"},
		{MidpointTool, "07.png"},
		{PerpendBisecTool, "08.png"},
		{AngleBisectTool, "09.png"},
		{PerpendTool, "10.png"},
		{ParallelTool, "11.png"},
		{MoveTool, "12.png"},
	}
)

var (
	scr        *adagui.Screen
	win        *adagui.Window
	root       *adagui.Panel
	canvas     *adagui.Panel
	p0, p1, p2 *adagui.Point
	l0, l1, l2 *adagui.Line
	c0, c1, c2 *adagui.Circle
	e0, e1, e2 *adagui.Ellipse
	r0, r1, r2 *adagui.Rectangle
	dbgDom     adagui.DebugDomain

	paleBlue = color.RGBAF{0.447, 0.706, 0.839, 1.0}
)

func SignalHandler() {
	sigChan := make(chan os.Signal)
	signal.Notify(sigChan, os.Interrupt)
	<-sigChan
	scr.Quit()
}

func NewPoint() *adagui.Point {
	var dp geom.Point

	p := adagui.NewPoint()

	p.SetOnPress(func(evt touch.Event) {
		dp = evt.Pos.Sub(p.Pos())
		adagui.Debugf(dbgDom, "dp: %v", dp)
		p.Mark(adagui.MarkNeedsPaint)
	})
	p.SetOnDrag(func(evt touch.Event) {
		adagui.Debugf(dbgDom, "new pos: %v", evt.Pos)
		p.SetPos(evt.Pos.Sub(dp))
		p.Mark(adagui.MarkNeedsPaint)
	})
	p.SetOnRelease(func(evt touch.Event) {
		p.Mark(adagui.MarkNeedsPaint)
	})

	return p
}

func NewLine() *adagui.Line {
	var dp geom.Point

	l := adagui.NewLine()

	l.SetOnPress(func(evt touch.Event) {
		dp = evt.Pos.Sub(l.Pos())
		adagui.Debugf(dbgDom, "dp: %v", dp)
        l.SetP0(evt.Pos)
		l.Mark(adagui.MarkNeedsPaint)
	})
	l.SetOnDrag(func(evt touch.Event) {
		adagui.Debugf(dbgDom, "new pos: %v", evt.Pos)
		//l.SetPos(evt.Pos.Sub(dp))
        l.SetP0(evt.Pos)
		l.Mark(adagui.MarkNeedsPaint)
	})
	l.SetOnRelease(func(evt touch.Event) {
		l.Mark(adagui.MarkNeedsPaint)
	})

	return l
}

func NewCircle(r float64) *adagui.Circle {
	var dp geom.Point

	c := adagui.NewCircle(r)

	c.SetOnPress(func(evt touch.Event) {
		dp = evt.Pos.Sub(c.Pos())
		adagui.Debugf(dbgDom, "dp: %v", dp)
		c.Mark(adagui.MarkNeedsPaint)
	})
	c.SetOnDrag(func(evt touch.Event) {
		adagui.Debugf(dbgDom, "new pos: %v", evt.Pos)
		c.SetPos(evt.Pos.Sub(dp))
		c.Mark(adagui.MarkNeedsPaint)
	})
	c.SetOnRelease(func(evt touch.Event) {
		c.Mark(adagui.MarkNeedsPaint)
	})

	return c
}

func NewEllipse(rx, ry float64) *adagui.Ellipse {
	var dp geom.Point

	e := adagui.NewEllipse(rx, ry)

	e.SetOnPress(func(evt touch.Event) {
		dp = evt.Pos.Sub(e.Pos())
		adagui.Debugf(dbgDom, "dp: %v", dp)
		e.Mark(adagui.MarkNeedsPaint)
	})
	e.SetOnDrag(func(evt touch.Event) {
		adagui.Debugf(dbgDom, "new pos: %v", evt.Pos)
		e.SetPos(evt.Pos.Sub(dp))
		e.Mark(adagui.MarkNeedsPaint)
	})
	e.SetOnRelease(func(evt touch.Event) {
		e.Mark(adagui.MarkNeedsPaint)
	})

	return e
}

func NewRectangle(w, h float64) *adagui.Rectangle {
	var dp geom.Point

	r := adagui.NewRectangle(w, h)

	r.SetOnPress(func(evt touch.Event) {
		dp = evt.Pos.Sub(r.Pos())
		adagui.Debugf(dbgDom, "dp: %v", dp)
		r.Mark(adagui.MarkNeedsPaint)
	})
	r.SetOnDrag(func(evt touch.Event) {
		adagui.Debugf(dbgDom, "new pos: %v", evt.Pos)
		r.SetPos(evt.Pos.Sub(dp))
		r.Mark(adagui.MarkNeedsPaint)
	})
	r.SetOnRelease(func(evt touch.Event) {
		r.Mark(adagui.MarkNeedsPaint)
	})

	return r
}

func NewCanvas(w, h float64) *adagui.Panel {
	var l *adagui.Line

	p := adagui.NewPanel(w, h)

	p.SetOnPress(func(evt touch.Event) {
		l = NewLine()
        l.SetP0(evt.Pos)
		p.Add(l)
		p.Mark(adagui.MarkNeedsPaint)
	})
	p.SetOnDrag(func(evt touch.Event) {
		l.SetP1(evt.Pos)
		p.Mark(adagui.MarkNeedsPaint)
	})
	p.SetOnRelease(func(evt touch.Event) {
		p.Mark(adagui.MarkNeedsPaint)
	})
	p.SetOnTap(func(evt touch.Event) {
		pt := NewPoint()
		pt.SetPos(evt.Pos)
		p.Add(pt)
		p.Mark(adagui.MarkNeedsPaint)
	})

	return p
}

func main() {
	//var img image.Image
	//var err error
	var icon *adagui.IconButton

	prop := props.NewPropsFromUser(props.PropsMap["Button"],
		"IconButtonProps.json")
	props.PropsMap["IconButton"] = prop

	flag.Parse()

	dbgDom = adagui.NewDebugDomain()
	adagui.AddDebugDomain(dbgDom)

	scr = adagui.NewScreen(adatft.Rotate090)
	win = scr.NewWindow()

	root = adagui.NewPanel(10, 10)
	root.Layout = adagui.NewPaddedLayout()
	root.SetColor(paleBlue.Bright(0.7))
	win.SetRoot(root)

	main := adagui.NewGroupPL(root, adagui.NewVBoxLayout())

	canvas = NewCanvas(0.0, win.Rect.Dy()-88.0)
	canvas.SetColor(color.Snow)
	canvas.SetBorderWidth(0.0)

	iconBar01 := adagui.NewGroup()
	iconBar01.Layout = adagui.NewHBoxLayout()
	iconBar02 := adagui.NewGroup()
	iconBar02.Layout = adagui.NewHBoxLayout()

	toolIdx := binding.NewInt()
	handIcon := adagui.NewIconButtonWithData("24x24/90.png", 0, toolIdx)
	erasureIcon := adagui.NewIconButtonWithData("24x24/91.png", 1, toolIdx)
	iconBar01.Add(handIcon, adagui.NewSpacer())
	iconBar02.Add(erasureIcon, adagui.NewSpacer())

	for n := range 18 {
		fileName := fmt.Sprintf("24x24/%02d.png", n+1)
		icon = adagui.NewIconButtonWithData(fileName, 2+n, toolIdx)
		if n < 9 {
			iconBar01.Add(icon)
		} else {
			iconBar02.Add(icon)
		}
	}

	main.Add(canvas, iconBar01, iconBar02)

	//fmt.Printf("size icon bar: %v\n", iconBar.Size())

	c0 = NewCircle(30)
	c0.SetPos(geom.Point{80, 80})

	c1 = NewCircle(30)
	c1.SetPos(geom.Point{160, 40})

	l1 = NewLine()
    l1.SetP0(geom.Point{20, 20})
    l1.SetP1(geom.Point{40, 40})


	//l2 = NewLine(geom.Point{10, 100}, geom.Point{130.0, 220.0})

	r1 = NewRectangle(80, 40)
	r1.SetPos(geom.Point{20, 100})

	p0 = NewPoint()
	p0.SetPos(geom.Point{50, 20})
	p1 = NewPoint()
	p1.SetPos(geom.Point{140, 15})

	e0 = NewEllipse(20.0, 60.0)
	e0.SetPos(geom.Point{240.0, 160.0})

	canvas.Add(c0, c1, l1, r1, p0, p1, e0)

	go SignalHandler()

	//    fmt.Printf("size of icon bar: %v\n", iconBar.Size())
	//    fmt.Printf("size of icon    : %v\n", icon.Size())

	win.SetRoot(root)
	scr.SetWindow(win)

	//    fmt.Printf("bounds of root: %v\n", root.Rect())

	scr.Run()
}
