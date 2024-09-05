package main

import (
	"flag"
	"fmt"
	"github.com/stefan-muehlebach/adatft"
	"github.com/stefan-muehlebach/gg"
	"github.com/stefan-muehlebach/gg/color"
	"github.com/stefan-muehlebach/gg/color"
	"github.com/stefan-muehlebach/gg/fonts"
	"github.com/stefan-muehlebach/gg/geom"
	"golang.org/x/image/font"
	"math"
	"time"
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

type Text struct {
	txt   string
	pos   geom.Point
	width float64
	face  font.Face
	align gg.Align
	color color.Color
	vis   bool
}

func NewText(txt string, pos geom.Point, width float64, face font.Face) *Text {
	t := &Text{}
	t.txt = txt
	t.pos = pos
	t.width = width
	t.face = face
	t.align = gg.AlignLeft
	t.color = color.WhiteSmoke
	t.vis = true
	return t
}

func (t *Text) IsVisible() bool {
	return t.vis
}

func (t *Text) Draw(gc *gg.Context) {
	gc.SetStrokeColor(t.color)
	gc.SetFontFace(t.face)
	gc.DrawStringWrapped(t.txt, t.pos.X, t.pos.Y, 0.0, 0.0, t.width, 1.5,
		t.align)
}

//-----------------------------------------------------------------------------

const (
	outerRadius = 18.0
	innerRadius = 9.0
)

var (
	fillColorActive   = color.ForestGreen
	fillColorDeactive = fillColorActive.Dark(0.6)
	lineColorActive   = color.White
	lineColorDeactive = lineColorActive.Dark(0.6)

	numSamples int
)

type Target struct {
	refPt     adatft.RefPointType
	pos       geom.Point
	samples   []adatft.TouchRawPos
	curSample int
	avg       adatft.TouchRawPos
	vis       bool
}

func NewTarget(refPt adatft.RefPointType, pos geom.Point) *Target {
	t := &Target{}
	t.refPt = refPt
	t.pos = pos
	t.samples = make([]adatft.TouchRawPos, numSamples)
	t.curSample = numSamples
	t.vis = true
	return t
}

func (t *Target) Reset() {
	t.curSample = 0
}

func (t *Target) AddSample(pos adatft.TouchRawPos) bool {
	var sumX, sumY int

	if t.curSample >= len(t.samples) {
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
		return false
	}
	return true
}

func (t *Target) IsVisible() bool {
	return t.vis
}

func (t *Target) Draw(gc *gg.Context) {
	tmp := outerSep + outerRadius
	active := t.curSample < len(t.samples)
	pct := float64(t.curSample) / float64(len(t.samples))

	gc.Push()
	gc.Translate(t.pos.X+0.5, t.pos.Y+0.5)
	if active {
		gc.SetFillColor(fillColorActive)
	} else {
		gc.SetFillColor(fillColorDeactive)
	}
	gc.DrawArc(0.0, 0.0, outerRadius, -math.Pi/2.0, pct*2.0*math.Pi-math.Pi/2.0)
	gc.LineTo(0.0, 0.0)
	gc.ClosePath()
	gc.Fill()

	gc.SetStrokeWidth(1.0)
	gc.DrawCircle(0.0, 0.0, innerRadius)
	if active {
		gc.SetStrokeColor(lineColorActive)
	} else {
		gc.SetStrokeColor(lineColorDeactive)
	}
	gc.SetFillColor(color.Black)
	gc.FillStroke()

	gc.DrawCircle(0.0, 0.0, outerRadius)
	gc.DrawLine(-tmp, 0.0, tmp, 0.0)
	gc.DrawLine(0.0, -tmp, 0.0, tmp)
	gc.Stroke()
	gc.Pop()
}

//-----------------------------------------------------------------------------

const (
	outerSep        = 2.0
	margin          = outerSep + outerRadius
	defNumSamples   = 256
	defDispRotation = adatft.Rotate000
)

var (
	textList = []string{
		"Mit diesem Programm wird die Verzerrung zwischen TouchScreen und Display gemessen und ein Parameterfile erstellt, mit welchem die Verzerrung ausgeglichen werden kann.",
		"Dazu muss die Mitte der 4 Kontrollpunkte möglichst genau angetippt und während einer gewissen Zeit gehalten werden.",
		"Im Anschluss werden die berechneten Parameter in die Datei 'RotateXXX.json' geschrieben. Um diese Parameter zu aktivieren, muss diese Datei unter '~/.config/adatft' abgelegt werden.",
	}
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
	gc.SetFillColor(color.Black)
	gc.Clear()
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

func DrawDisplay(gc, gcText *gg.Context) {
	gc.DrawImage(gcText.Image(), 0, 0)
}

//-----------------------------------------------------------------------------

func main() {
	flag.IntVar(&numSamples, "samples", defNumSamples,
		"number of samples to collect")
	flag.Var(&dispRotation, "rotation", "display rotation")
	flag.Parse()

	dsp = adatft.OpenDisplay(dispRotation)
	tch = adatft.OpenTouch()
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

	//gc.SetFillColor(color.Black)
	//gc.Clear()

	targetList = make([]*Target, adatft.NumRefPoints)
	for refPt := adatft.RefTopLeft; refPt < adatft.NumRefPoints; refPt++ {
		target := NewTarget(refPt, geom.Point(penPosList[refPt]))
		targetList[refPt] = target
		graphObjList = append(graphObjList, target)
	}

	infoText := NewText("", geom.Point{margin, 2.5 * margin}, width-2*margin,
		fonts.NewFace(fonts.GoRegular, 15.0))
	statusText := NewText("Tap für weiter...",
		geom.Point{margin, height - 3.5*margin}, width-2*margin,
		fonts.NewFace(fonts.GoRegular, 15.0))
	statusText.align = gg.AlignRight

	graphObjList = append(graphObjList, infoText, statusText)
	for _, txt := range textList {
		infoText.txt = txt
		UpdateDisplay()
		waitForEvent(tch, adatft.PenRelease)
	}

	infoText.txt = "Bereit?"
	infoText.face = fonts.NewFace(fonts.GoBold, 48.0)
	infoText.align = gg.AlignCenter
	UpdateDisplay()
	waitForEvent(tch, adatft.PenRelease)

	infoText.txt = "Los!"
	statusText.vis = false
	for _, target := range targetList {
		target.Reset()
	}
	UpdateDisplay()

	for refPt := adatft.RefTopLeft; refPt < adatft.NumRefPoints; refPt++ {
		CollectData(targetList[refPt])
	}

	infoText.txt = "Fertig!"
	infoText.face = fonts.NewFace(fonts.GoBold, 63.0)
	UpdateDisplay()
	time.Sleep(5 * time.Second)

	distortedPlane := &adatft.DistortedPlane{}
	for i, target := range targetList {
		fmt.Printf("%d: %v, %v\n", i, target.pos, target.avg)
		distortedPlane.SetRefPoint(target.refPt, target.avg, penPosList[i])
	}
	distortedPlane.WriteConfigFile("Calib.json")

	tch.Close()
	dsp.Close()
}
