package main

import (
	"flag"
	"fmt"
	"math"
	"time"

	"github.com/stefan-muehlebach/adatft"
	"github.com/stefan-muehlebach/gg"
	"github.com/stefan-muehlebach/gg/colors"
	"github.com/stefan-muehlebach/gg/fonts"
	"github.com/stefan-muehlebach/gg/geom"
	"golang.org/x/image/font"
)

const (
	outerSep        = 5.0
	outerRadius     = 18.0
	innerRadius     = 9.0
	margin          = outerSep + outerRadius
	defNumSamples   = 128
	defDispRotation = adatft.Rotate000
)

var (
	fillColorActive   = colors.YellowGreen
	fillColorDeactive = fillColorActive.Alpha(0.4)
	lineColorActive   = colors.White
	lineColorDeactive = lineColorActive.Alpha(0.4)
	lineWidth         = 1.0

	backgroundColor = colors.DarkBlue.Alpha(0.5)
	numSamples      int
)

//-----------------------------------------------------------------------------

type GraphicObject interface {
	IsVisible() bool
	Draw(gc *gg.Context)
}

var (
	graphObjList []GraphicObject
)

func init() {
	graphObjList = make([]GraphicObject, 0)
}

//-----------------------------------------------------------------------------

type Arrow struct {
	p1, c, p2 geom.Point
	vis       bool
}

func NewArrow(p1, c, p2 geom.Point) *Arrow {
	a := &Arrow{}
	a.p1 = p1
	a.c = c
	a.p2 = p2
	a.vis = true
	return a
}

func (a *Arrow) IsVisible() bool {
	return a.vis
}

func (a *Arrow) Draw(gc *gg.Context) {
	gc.MoveTo(a.p1.AsCoord())
	gc.QuadraticTo(a.c.X, a.c.Y, a.p2.X, a.p2.Y)
	gc.SetStrokeWidth(5.0)
	gc.SetStrokeColor(colors.WhiteSmoke)
	gc.Stroke()

	v := a.p2.Sub(a.c)
	fmt.Printf("v: %+v\n", v)
	angle := v.Angle()
	fmt.Printf("Angle: %f\n", angle)
	rotMat := geom.Rotate(angle)
	v1 := rotMat.Transform(geom.Point{-4, 9}).Add(a.p2)
	v2 := rotMat.Transform(geom.Point{4, 9}).Add(a.p2)

	gc.MoveTo(v1.AsCoord())
	gc.LineTo(a.p2.AsCoord())
	gc.LineTo(v2.AsCoord())
	gc.Stroke()
}

//-----------------------------------------------------------------------------

type Text struct {
	txt   string
	pos   geom.Point
	width float64
	face  font.Face
	align gg.Align
	color colors.Color
	vis   bool
}

func NewText(txt string, pos geom.Point, width float64, face font.Face) *Text {
	t := &Text{}
	t.txt = txt
	t.pos = pos
	t.width = width
	t.face = face
	t.align = gg.AlignLeft
	t.color = colors.WhiteSmoke
	t.vis = true
	return t
}

func (t *Text) IsVisible() bool {
	return t.vis
}

func (t *Text) Draw(gc *gg.Context) {
	gc.SetTextColor(t.color)
	gc.SetFontFace(t.face)
	gc.DrawStringWrapped(t.txt, t.pos.X, t.pos.Y, 0.0, 0.0, t.width, 1.5,
		t.align)
}

//-----------------------------------------------------------------------------

type Target struct {
	refPt     adatft.RefPointType
	pos       geom.Point
	samples   []adatft.TouchRawPos
	curSample int
	avg       adatft.TouchRawPos
	vis, act  bool
}

func NewTarget(refPt adatft.RefPointType, pos geom.Point) *Target {
	t := &Target{}
	t.refPt = refPt
	t.pos = pos
	t.samples = make([]adatft.TouchRawPos, numSamples)
	t.curSample = 0
	t.vis = true
	t.act = false
	return t
}

func (t *Target) Reset() {
	t.curSample = 0
	t.vis = true
	t.act = true
}

func (t *Target) IsActive() bool {
	return t.act
}

func (t *Target) IsVisible() bool {
	return t.vis
}

func (t *Target) AddSample(pos adatft.TouchRawPos) bool {
	var sumX, sumY int

	if !t.IsActive() {
		return false
	}
	t.samples[t.curSample] = pos
	t.curSample++
	if t.curSample == len(t.samples) {
		sumX, sumY = 0, 0
		for _, rawPos := range t.samples {
			sumX += int(rawPos.RawX)
			sumY += int(rawPos.RawY)
		}
		t.avg.RawX = uint16(sumX / len(t.samples))
		t.avg.RawY = uint16(sumY / len(t.samples))
		t.act = false
		return false
	}
	return true
}

func (t *Target) Draw(gc *gg.Context) {
	tmp := outerSep + outerRadius
	pct := float64(t.curSample) / float64(len(t.samples))

	gc.Push()
	gc.Translate(t.pos.X+0.5, t.pos.Y+0.5)
	if t.act {
		gc.SetFillColor(fillColorActive)
	} else {
		gc.SetFillColor(fillColorDeactive)
	}
	gc.DrawArc(0.0, 0.0, outerRadius, -math.Pi/2.0, pct*2.0*math.Pi-math.Pi/2.0)
	gc.LineTo(0.0, 0.0)
	gc.ClosePath()
	gc.Fill()

	gc.DrawCircle(0.0, 0.0, innerRadius)
	gc.SetFillColor(backgroundColor)
	gc.Fill()

	gc.SetStrokeWidth(lineWidth)
	if t.act {
		gc.SetStrokeColor(lineColorActive)
	} else {
		gc.SetStrokeColor(lineColorDeactive)
	}
	gc.DrawCircle(0.0, 0.0, innerRadius)
	gc.DrawCircle(0.0, 0.0, outerRadius)
	gc.DrawLine(-tmp, 0.0, tmp, 0.0)
	gc.DrawLine(0.0, -tmp, 0.0, tmp)
	gc.Stroke()
	gc.Pop()
}

//-----------------------------------------------------------------------------

type Iterator[T any] struct {
	lst   []T
	idx   int
	cycle bool
}

func NewIterator[T any](lst []T, cycle bool) *Iterator[T] {
	i := &Iterator[T]{}
	i.lst = lst
	i.idx = 0
	i.cycle = cycle
	return i
}

func (i *Iterator[T]) Has() bool {
	if i.cycle {
		return true
	} else {
		return i.idx < len(i.lst)
	}
}

func (i *Iterator[T]) Next() T {
	j := i.idx
	if i.cycle {
		i.idx = (i.idx + 1) % len(i.lst)
	} else {
		if i.idx < len(i.lst)-1 {
			i.idx += 1
		}
	}
	return i.lst[j]
}

var (
	textList = []string{
		`Mit diesem Programm wird der Versatz zwischen TouchScreens und Displays gemessen. Das daraus resultierende Parameterfile wird von AdaTFT verwendet, um den Versatz zu begleichen.`,
		`Damit dieses Wunder vollbracht werden kann, ist Deine Mithilfe notwendig! Im Folgenden werden dir vier Kontrollpunkte in den Ecken des Bildschirms gezeigt. Diese gilt es so präzis wie möglich mit einem Stift anzutippen und über eine bestimmte Zeit zu halten.`,
		`Im Anschluss wird im aktuellen Verzeichnis eine Datei namens «TouchCalib.json» erstellt, welche die berechneten Werte enthält. Um diese Datei zu aktivieren, verschiebt man sie einfach in das Verzeichnis «~/.config/adatft».`,
		`Die Werte wurden in die Datei «TouchCalib.json» geschrieben.`,
		`Drück mit dem Stift in das Zentrum der Markierung und halte die Position...`,
		`Heb den Stift kurz hoch, damit das nächste Ziel aktiv wird.`,
	}
	textIter *Iterator[string] = NewIterator(textList, false)

	backColorList = []colors.Color{
		colors.DarkMagenta.Alpha(0.3),
		colors.DarkBlue.Alpha(0.3),
		colors.DarkCyan.Alpha(0.3),
		colors.DarkGreen.Alpha(0.3),
		colors.Gold.Alpha(0.3),
		colors.DarkRed.Alpha(0.3),
	}
	colorIter *Iterator[colors.Color] = NewIterator(backColorList, true)
)

var (
	dispRotation = defDispRotation

	dsp        *adatft.Display
	tch        *adatft.Touch
	gc, gcText *gg.Context

	penPosList [4]adatft.TouchPos

	targetList  []*Target
	curRefPoint adatft.RefPointType

	width, height float64
)

func UpdateDisplay() {
	gc.SetFillColor(backgroundColor)
	gc.Clear()
	//gc.SetFillColor(backgroundColor)
	//gc.DrawRectangle(0, 0, width, height)
	//gc.Fill()
	for _, graphObj := range graphObjList {
		if !graphObj.IsVisible() {
			continue
		}
		graphObj.Draw(gc)
	}
	dsp.Draw(gc.Image())
}

func CollectData(target *Target) {
	var tick *time.Ticker
	var curPos adatft.TouchRawPos
	var isCollecting, isDone bool
	var evt adatft.PenEvent

	tick = time.NewTicker(30 * time.Millisecond)
	isCollecting = false
	isDone = false

	for !isDone {
		select {
		case evt = <-tch.EventQ:
			switch evt.Type {
			case adatft.PenPress:
			case adatft.PenDrag:
				isCollecting = true
				curPos = evt.TouchRawPos
			case adatft.PenRelease:
				isCollecting = false
			}
		case <-tick.C:
			if !isCollecting {
				continue
			}
			if !target.AddSample(curPos) {
				isDone = true
			}
		}
		UpdateDisplay()
	}
	tick.Stop()
}

func waitForEvent(tch *adatft.Touch, typ adatft.PenEventType) {
	for evt := range tch.EventQ {
		if evt.Type == typ {
			return
		}
	}
}

//-----------------------------------------------------------------------------

func main() {
	flag.IntVar(&numSamples, "samples", defNumSamples,
		"number of samples to collect")
	dispRotation = defDispRotation
	flag.Parse()

	dsp = adatft.OpenDisplay(dispRotation)
	tch = adatft.OpenTouch(dispRotation)
	gc = gg.NewContext(adatft.Width, adatft.Height)

	width = float64(adatft.Width)
	height = float64(adatft.Height)

	penPosList[0] = adatft.TouchPos{X: margin, Y: margin}
	penPosList[1] = adatft.TouchPos{X: width - margin, Y: margin}
	penPosList[2] = adatft.TouchPos{X: width - margin, Y: height - margin}
	penPosList[3] = adatft.TouchPos{X: margin, Y: height - margin}

	curRefPoint = adatft.NumRefPoints
	//pctData = 0.0
	//penDataList = make([]adatft.TouchRawPos, numSamples)
	//curSample = 0
	//collecting = false

	//doneQ = make(chan bool)

	//go displayThread(doneQ)

	//gc.SetFillColor(colors.Black)
	//gc.Clear()

	targetList = make([]*Target, adatft.NumRefPoints)
	for refPt := adatft.RefTopLeft; refPt < adatft.NumRefPoints; refPt++ {
		target := NewTarget(refPt, geom.Point{penPosList[refPt].X, penPosList[refPt].Y})
		targetList[refPt] = target
		graphObjList = append(graphObjList, target)
	}

	infoText := NewText("", geom.Point{margin, 3.0 * margin}, width-2*margin,
		fonts.NewFace(fonts.GoRegular, 18.0))
	statusText := NewText("Tap für weiter...",
		geom.Point{margin, height - 4*margin}, width-2*margin,
		fonts.NewFace(fonts.GoItalic, 18.0))
	statusText.align = gg.AlignRight
	noteText := NewText(textList[4], geom.Point{2.1 * margin, 2.1 * margin}, width-4.2*margin,
		fonts.NewFace(fonts.LucidaHandwritingItalic, 14.0))
	noteText.vis = false

	graphObjList = append(graphObjList, infoText, statusText, noteText)
	for i := 0; i < 3; i++ {
		infoText.txt = textIter.Next()
		backgroundColor = colorIter.Next()
		UpdateDisplay()
		waitForEvent(tch, adatft.PenRelease)
	}

	infoText.pos = infoText.pos.AddXY(0, 100)
	infoText.txt = "Bereit?"
	infoText.face = fonts.NewFace(fonts.GoBold, 48.0)
	infoText.align = gg.AlignCenter
	backgroundColor = colorIter.Next()
	UpdateDisplay()
	waitForEvent(tch, adatft.PenRelease)

	infoText.txt = "Los!"
	statusText.vis = false
	backgroundColor = colorIter.Next()

	noteText.vis = true

	for refPt := adatft.RefTopLeft; refPt < adatft.NumRefPoints; refPt++ {
		targetList[refPt].Reset()
		UpdateDisplay()
		CollectData(targetList[refPt])
		if refPt == adatft.RefTopLeft {
			noteText.txt = textList[5]
			UpdateDisplay()
			noteText.vis = false
		}
		waitForEvent(tch, adatft.PenRelease)
	}

	infoText.txt = "Fertig!"
	UpdateDisplay()
	time.Sleep(2 * time.Second)

	distortedPlane := &adatft.DistortedPlane{}
	for i, target := range targetList {
		fmt.Printf("%d: %v, %v\n", i, target.pos, target.avg)
		distortedPlane.SetRefPoint(target.refPt, target.avg, penPosList[i])
	}
	distortedPlane.WriteConfigFile("TouchCalib.json")

	infoText.txt = textList[3]
	infoText.face = fonts.NewFace(fonts.GoRegular, 18.0)
	infoText.align = gg.AlignLeft
	statusText.vis = true
	backgroundColor = colorIter.Next()
	UpdateDisplay()
	waitForEvent(tch, adatft.PenRelease)

	tch.Close()
	dsp.Close()
}
