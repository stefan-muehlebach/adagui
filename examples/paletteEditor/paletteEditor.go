package main

import (
	"context"
	"flag"
	_ "fmt"
	"github.com/stefan-muehlebach/adagui"
	"github.com/stefan-muehlebach/adagui/touch"
	"github.com/stefan-muehlebach/adatft"
	"github.com/stefan-muehlebach/gg"
	_ "github.com/stefan-muehlebach/gg/color"
	"github.com/stefan-muehlebach/gg/colornames"
	"github.com/stefan-muehlebach/gg/geom"
	"github.com/stefan-muehlebach/mandel"
	_ "image"
	"image/color"
	"log"
	"math"
	_ "math/rand"
	_ "runtime/trace"
	_ "sync"
	_ "time"
)

//----------------------------------------------------------------------------

func init() {
	log.SetFlags(log.Lmicroseconds | log.Lmsgprefix)
	log.SetPrefix(": ")
}

//----------------------------------------------------------------------------

var (
	// Dies sind die Farben, welche fuer die Bezier-Kurven und die Kontoll-
	// Punkte zur Anpassung der Farbkurve verwendet werden. Von jeder Grund-
	// farbe (R, G, B) werden drei Helligkeitsstufen angeboten. Es hat sich
	// gezeigt, dass die Darstellung der Farben auf dem TFT von Adafruit
	// nicht homogen ist, d.h. 0xff fuer Rot ist leicht dunkler in der Wahr-
	// nehmung als 0xff fuer Gruen (was ziemlich hell erscheint). Daher koennen
	// die Farben hier einzeln definiert werden.
	BrightColors = []color.Color{
		colornames.Red,
		colornames.Lime,
		colornames.Blue,
		//color.RGBA{255, 0, 0, 255},
		//color.RGBA{0, 255, 0, 255},
		//color.RGBA{0, 0, 255, 255},
	}
	NormalColors = []color.Color{
		colornames.Red.Dark(0.37),
		colornames.Lime.Dark(0.37),
		colornames.Blue.Dark(0.37),
		//color.RGBA{159, 0, 0, 255},
		//color.RGBA{0, 159, 0, 255},
		//color.RGBA{0, 0, 159, 255},
	}
	DimmedColors = []color.Color{
		colornames.Red.Dark(0.92),
		colornames.Lime.Dark(0.92),
		colornames.Blue.Dark(0.92),
		//color.RGBA{20, 0, 0, 255},
		//color.RGBA{0, 20, 0, 255},
		//color.RGBA{0, 0, 20, 255},
	}

	BrightWhite = colornames.White
	DimmedWhite = colornames.White.Dark(0.37)
	Background  = colornames.Black
	//BrightWhite = color.RGBA{255, 255, 255, 255}
	//DimmedWhite = color.RGBA{160, 160, 160, 255}
	//Background  = color.RGBA{80, 80, 80, 255}

	palPrevWidth       = 480.0
	palPrevHeight      = 80.0
	palPrevInset       = 8.5
	palPrevLineWidth   = 2.0
	palPrevStrokeColor = BrightWhite

	gradEditWidth        = 480.0
	gradEditHeight       = 65.0
	gradEditInset        = 8.5
	gradEditLineWidth    = 2.0
	gradEditStrokeColor  = BrightWhite
	gradEditFillColor    = Background
	gradEditRectRound    = 5.0
	gradEditSplineWidth  = 4.0
	gradEditSplineColors = NormalColors

	ctrlPtRadius         = 8.0
	ctrlPtLineWidth      = 2.0
	ctrlPtActLineWidth   = 4.0
	ctrlPtStrokeColors   = BrightColors
	ctrlPtFillColors     = NormalColors
	ctrlPtSyncStrokeCols = BrightColors
	ctrlPtSyncFillCols   = DimmedColors

	procEditWidth          = 480.0
	procEditHeight         = 65.0
	procEditInset          = 8.5
	procEditLineWidth      = 2.0
	procEditStrokeColor    = BrightWhite
	procEditRectRound      = 5.0
	procEditGraphLineWidth = 4.0
	procEditGraphColors    = NormalColors
)

func UpdateVars(w *adagui.Window) {
	palPrevWidth = w.Rect.Dx()
	gradEditWidth = w.Rect.Dx()
	procEditWidth = w.Rect.Dx()
}

//----------------------------------------------------------------------------

// Mit diesem Widget wird die Palette angezeigt und Mutationen werden
// direkt nachgefuert.
type PalettePreview struct {
	adagui.LeafEmbed
	Inset   geom.Rectangle
	Palette mandel.Palette
}

// NewPalettePreview erzeugt ein neues Widget zur Anzeige einer Palette.
// Als einziges Pflichargument muss eine geladene Palette aus dem [mandel]
// Package angegeben werden.
func NewPalettePreview(palette mandel.Palette) *PalettePreview {
	n := &PalettePreview{}
	n.Wrapper = n
	n.Init()
	n.PropertyEmbed.InitByName("Default")
	n.SetSize(geom.Point{palPrevWidth, palPrevHeight})
	n.Palette = palette
	return n
}

func (n *PalettePreview) SetSize(s geom.Point) {
	n.Embed.SetSize(s)
	n.Inset = n.Bounds().Inset(palPrevInset, palPrevInset/2)
}

func (n *PalettePreview) Paint(gc *gg.Context) {
	n.Palette.SetLength(int(n.Inset.Dx()))
	pos0 := n.Inset.Min.Int()
	for col := 0; col < int(n.Inset.Dx()); col++ {
		color := n.Palette.GetColor(float64(col))
		for row := 0; row < int(n.Inset.Dy()); row++ {
			gc.SetPixel(pos0.X+col, pos0.Y+row, color)
		}
	}
	gc.DrawRectangle(n.Inset.AsCoord())
	gc.SetStrokeColor(palPrevStrokeColor)
	gc.SetStrokeWidth(palPrevLineWidth)
	gc.Stroke()
}

//----------------------------------------------------------------------------

type GradientEditor struct {
	adagui.ContainerEmbed
	Inset                   geom.Rectangle
	palette                 *mandel.GradientPalette
	color                   mandel.BaseColorType
	firstCtrlPt, lastCtrlPt *CtrlPoint
}

func NewGradientEditor(palette *mandel.GradientPalette,
	color mandel.BaseColorType) *GradientEditor {
	n := &GradientEditor{}
	n.Wrapper = n
	n.Init()
	n.PropertyEmbed.InitByName("Default")
	n.SetSize(geom.Point{gradEditWidth, gradEditHeight})
	n.palette = palette
	n.color = color
	n.CreateCtrlPoints()
	return n
}

func (n *GradientEditor) SetSize(s geom.Point) {
	n.Embed.SetSize(s)
	n.Inset = n.Bounds().Inset(gradEditInset, gradEditInset)
}

func (n *GradientEditor) Paint(gc *gg.Context) {
	gc.DrawRoundedRectangle(n.Inset.X0(), n.Inset.Y0(), n.Inset.Dx(), n.Inset.Dy(),
		gradEditRectRound)
	gc.SetStrokeWidth(gradEditLineWidth)
	gc.SetStrokeColor(gradEditStrokeColor)
	gc.Stroke()

	gc.SetStrokeWidth(gradEditSplineWidth)
	gc.SetStrokeColor(gradEditSplineColors[n.color])
	rect := n.Inset
	gpl := n.palette.GradPointList(n.color)
	for i := 0; i < len(gpl)-1; i++ {
		gp0 := gpl[i]
		gp1 := gpl[i+1]
		pt0 := rect.RelPos(gp0.Pos, 1.0-gp0.Val)
		pt1 := rect.RelPos(gp1.Pos, 1.0-gp1.Val)
		gc.MoveTo(pt0.AsCoord())
		gc.CubicTo(0.5*pt0.X+0.5*pt1.X, pt0.Y,
			0.5*pt0.X+0.5*pt1.X, pt1.Y,
			pt1.X, pt1.Y)
		gc.Stroke()
	}
	n.ContainerEmbed.Paint(gc)
	//gc.Pop()
}

func (n *GradientEditor) OnInputEvent(evt touch.Event) {
	switch evt.Type {
	case touch.TypeLongPress:
		rect := n.Inset
		pos := rect.SetInside(evt.Pos)
		fx, fy := rect.PosRel(pos)
		fy = 1 - fy
		gp := &mandel.GradPoint{fx, fy}
		n.palette.AddGradPoint(n.color, gp)
		n.palette.Update()
		cp := NewCtrlPoint(n, gp)
		n.Add(cp)
		n.Mark(adagui.MarkNeedsPaint)
	}
}

func (n *GradientEditor) CreateCtrlPoints() {
	var first, last *CtrlPoint

	n.ChildList.Init()
	gradPointList := n.palette.GradPointList(n.color)
	for i, gradPoint := range gradPointList {
		cp := NewCtrlPoint(n, gradPoint)
		n.Add(cp)
		if i == 0 {
			first = cp
		} else {
			last = cp
		}
	}
	first.sync = last
	last.sync = first
}

//----------------------------------------------------------------------------

type CtrlPoint struct {
	adagui.LeafEmbed
	color  mandel.BaseColorType
	active bool
	ge     *GradientEditor
	gp     *mandel.GradPoint
	sync   *CtrlPoint
	inSync bool
}

func NewCtrlPoint(ge *GradientEditor, gp *mandel.GradPoint) *CtrlPoint {
	n := &CtrlPoint{}
	n.Wrapper = n
	n.Init()
	n.PropertyEmbed.InitByName("Default")
	n.SetSize(geom.Point{2 * ctrlPtRadius, 2 * ctrlPtRadius})
	n.ge = ge
	n.color = ge.color
	n.active = false
	n.gp = gp
	n.sync = nil
	n.inSync = false
	n.UpdatePos()
	return n
}

func (n *CtrlPoint) SetPos(p geom.Point) {
	n.Embed.SetPos(p.Sub(geom.Point{ctrlPtRadius, ctrlPtRadius}))
}

func (n *CtrlPoint) Pos() geom.Point {
	return n.Embed.Pos().Add(geom.Point{ctrlPtRadius, ctrlPtRadius})
}

func (n *CtrlPoint) Paint(gc *gg.Context) {
	gc.DrawPoint(ctrlPtRadius, ctrlPtRadius, ctrlPtRadius)
	if n.inSync {
		gc.SetFillColor(ctrlPtSyncFillCols[n.color])
		gc.SetStrokeColor(ctrlPtSyncStrokeCols[n.color])
	} else {
		gc.SetFillColor(ctrlPtFillColors[n.color])
		gc.SetStrokeColor(ctrlPtStrokeColors[n.color])
	}
	if n.active {
		gc.SetStrokeWidth(ctrlPtActLineWidth)
	} else {
		gc.SetStrokeWidth(ctrlPtLineWidth)
	}
	gc.FillStroke()
}

func (n *CtrlPoint) OnInputEvent(evt touch.Event) {
	switch evt.Type {
	case touch.TypePress:
		n.active = true
		n.Mark(adagui.MarkNeedsPaint)
	case touch.TypeDrag:
		rect := n.ge.Inset
		pos := rect.SetInside(evt.Pos)
		if n.IsFirst() {
			pos.X = rect.Min.X
		}
		if n.IsLast() {
			pos.X = rect.Max.X
		}
		n.SetPos(pos)
		n.UpdateValue()
		if n.sync != nil && n.inSync {
			n.sync.gp.Val = n.gp.Val
			n.sync.UpdatePos()
		}
		n.ge.palette.Update()
		n.Mark(adagui.MarkNeedsPaint)
	case touch.TypeRelease:
		n.active = false
		n.Mark(adagui.MarkNeedsPaint)
	case touch.TypeLongPress:
		if !n.IsFirst() && !n.IsLast() {
			break
		}
		n.inSync = !n.inSync
		n.sync.inSync = n.inSync
		n.sync.gp.Val = n.gp.Val
		n.sync.UpdatePos()
		n.ge.palette.Update()
		n.Mark(adagui.MarkNeedsPaint)
	case touch.TypeDoubleTap:
		if n.IsFirst() || n.IsLast() {
			break
		}
		p := n.Parent
		n.Remove()
		n.ge.palette.DelGradPoint(n.color, n.gp)
		n.ge.palette.Update()
		p.Mark(adagui.MarkNeedsPaint)
	}
}

func (n *CtrlPoint) IsFirst() bool {
	return n.gp.Pos == 0.0
}

func (n *CtrlPoint) IsLast() bool {
	return n.gp.Pos == 1.0
}

func (n *CtrlPoint) UpdatePos() {
	rect := n.ge.Inset
	pos := rect.RelPos(n.gp.Pos, 1-n.gp.Val)
	n.SetPos(pos)
}

func (n *CtrlPoint) UpdateValue() {
	rect := n.ge.Inset
	fx, fy := rect.PosRel(n.Pos())
	n.gp.Pos, n.gp.Val = fx, 1-fy
}

//----------------------------------------------------------------------------

type ProcEditor struct {
	adagui.ContainerEmbed
	Inset   geom.Rectangle
	palette *mandel.ProcPalette
	color   mandel.BaseColorType
	fncPts  []*FuncPoint
}

func NewProcEditor(palette *mandel.ProcPalette,
	color mandel.BaseColorType) *ProcEditor {
	n := &ProcEditor{}
	n.Wrapper = n
	n.Init()
	n.PropertyEmbed.InitByName("Default")
	n.SetSize(geom.Point{procEditWidth, procEditHeight})
	n.palette = palette
	n.color = color
	n.fncPts = make([]*FuncPoint, 2)
	n.AddFuncPoints()
	return n
}

func (n *ProcEditor) SetSize(s geom.Point) {
	n.Embed.SetSize(s)
	n.Inset = n.Bounds().Inset(procEditInset, procEditInset)
}

func (n *ProcEditor) Paint(gc *gg.Context) {
	gc.DrawRoundedRectangle(n.Inset.X0(), n.Inset.Y0(), n.Inset.Dx(), n.Inset.Dy(),
		procEditRectRound)
	gc.SetStrokeWidth(procEditLineWidth)
	gc.SetStrokeColor(procEditStrokeColor)
	gc.Stroke()

	gc.SetStrokeWidth(procEditGraphLineWidth)
	gc.SetStrokeColor(procEditGraphColors[n.color])

	for col := n.Inset.X0(); col < n.Inset.X1(); col += 1.0 {
		fx, _ := n.Inset.PosRel(geom.Point{col, 0.0})
		fy := n.palette.Value(n.color, fx)
		pt := n.Inset.RelPos(fx, 1-fy)
		gc.LineTo(pt.AsCoord())
	}
	gc.Stroke()
	n.ContainerEmbed.Paint(gc)
}

func (n *ProcEditor) AddFuncPoints() {
	n.ChildList.Init()
	for pt := PointA; pt < NumFuncPointTypes; pt++ {
		fp := NewFuncPoint(n, pt)
		n.Add(fp)
		n.fncPts[pt] = fp
	}
}

func (n *ProcEditor) UpdateProcValues() {
	x1, y1 := n.fncPts[0].Value()
	x2, y2 := n.fncPts[1].Value()
	v0 := (y1 + y2) / 2.0
	v1 := (y1 - y2) / 2.0
	v2 := 1.0 / (2.0 * (x2 - x1))
	v3 := -x1 * v2
	n.palette.SetParamList(n.color, []float64{v0, v1, v2, v3})
	n.palette.Update()
}

//----------------------------------------------------------------------------

type FuncPointType uint8

const (
	PointA FuncPointType = iota
	PointB
	NumFuncPointTypes
)

type CoordFuncType func(l []float64) (fx, fy float64)

var (
	cfl [2]CoordFuncType
)

func init() {
	cfl[PointA] = func(l []float64) (fx, fy float64) {
		fx, fy = -l[3]/l[2], l[0]+l[1]
		if math.IsNaN(fx) {
			fx = 0.0
		}
		return fx, fy
	}
	cfl[PointB] = func(l []float64) (fx, fy float64) {
		fx, fy = 1/(2*l[2])-l[3]/l[2], l[0]-l[1]
		if math.IsNaN(fx) {
			fx = 1.0
		}
		return fx, fy
	}
}

type FuncPoint struct {
	adagui.LeafEmbed
	pe        *ProcEditor
	color     mandel.BaseColorType
	active    bool
	pointType FuncPointType
}

func NewFuncPoint(pe *ProcEditor, pointType FuncPointType) *FuncPoint {
	n := &FuncPoint{}
	n.Wrapper = n
	n.Init()
	n.PropertyEmbed.InitByName("Default")
	n.SetSize(geom.Point{2 * ctrlPtRadius, 2 * ctrlPtRadius})
	n.pe = pe
	n.color = pe.color
	n.active = false
	n.pointType = pointType
	l := n.pe.palette.ParamList(n.color)
	n.SetValue(cfl[n.pointType](l))
	return n
}

func (n *FuncPoint) SetPos(p geom.Point) {
	n.Embed.SetPos(p.Sub(geom.Point{ctrlPtRadius, ctrlPtRadius}))
}

func (n *FuncPoint) Pos() geom.Point {
	return n.Embed.Pos().Add(geom.Point{ctrlPtRadius, ctrlPtRadius})
}

func (n *FuncPoint) Paint(gc *gg.Context) {
	gc.DrawPoint(ctrlPtRadius, ctrlPtRadius, ctrlPtRadius)
	gc.SetFillColor(ctrlPtFillColors[n.color])
	gc.SetStrokeColor(ctrlPtStrokeColors[n.color])
	if n.active {
		gc.SetStrokeWidth(ctrlPtActLineWidth)
	} else {
		gc.SetStrokeWidth(ctrlPtLineWidth)
	}
	gc.FillStroke()
}

func (n *FuncPoint) OnInputEvent(evt touch.Event) {
	switch evt.Type {
	case touch.TypePress:
		n.active = true
		n.Mark(adagui.MarkNeedsPaint)
	case touch.TypeDrag:
		rect := n.pe.Inset
		pos := rect.SetInside(evt.Pos)
		n.SetPos(pos)
		n.pe.UpdateProcValues()
		n.Mark(adagui.MarkNeedsPaint)
	case touch.TypeRelease:
		n.active = false
		n.Mark(adagui.MarkNeedsPaint)
	}
}

func (n *FuncPoint) SetValue(fx, fy float64) {
	rect := n.pe.Inset
	pos := rect.RelPos(fx, 1-fy)
	n.SetPos(pos)
}

func (n *FuncPoint) Value() (float64, float64) {
	rect := n.pe.Inset
	fx, fy := rect.PosRel(n.Pos())
	return fx, 1 - fy
}

//----------------------------------------------------------------------------

var (
	screen *adagui.Screen
	win    *adagui.Window
	ctx    context.Context
)

// -----------------------------------------------------------------------------
//
// (main) --
func main() {
	var palName string
	var palette mandel.Palette
	var err error
	var rotation adatft.RotationType = adatft.Rotate090

	flag.StringVar(&palName, "palette", "Default",
		"name of the palette to edit")
	flag.Var(&rotation, "rotation", "Rotation of the Display")
	flag.Parse()

	palette, err = mandel.NewPalette(palName)
	if err != nil {
		log.Fatalf("Couldn'r read palette: %v", err)
	}

	screen = adagui.NewScreen(rotation)
	win = screen.NewWindow()
	UpdateVars(win)

	group := adagui.NewPanel(0, 0)
	group.SetColor(Background)
	win.SetRoot(group)

	// Erstellt die Vorschau auf die Palette
	palPrev := NewPalettePreview(palette)
	palPrev.SetPos(geom.Point{0.0, 0.0})
	palPrev.SetOnDoubleTap(func(evt touch.Event) {
		screen.Quit()
	})
	palPrev.SetOnTap(func(evt touch.Event) {
		screen.Save("screenshot.png")
	})
	group.Add(palPrev)

	switch palImpl := palette.(type) {

	// Erstellt alle Objekte fuer die Bearbeitung einer Gradienten-Palette.
	case *mandel.GradientPalette:
		// Editor fuer Rot
		geRed := NewGradientEditor(palImpl, mandel.Red)
		geRed.SetPos(geom.Point{0.0, palPrevHeight})
		group.Add(geRed)

		// Editor fuer Gruen
		geGreen := NewGradientEditor(palImpl, mandel.Green)
		geGreen.SetPos(geom.Point{0.0, palPrevHeight + gradEditHeight})
		group.Add(geGreen)

		// Editor fuer Blau
		geBlue := NewGradientEditor(palImpl, mandel.Blue)
		geBlue.SetPos(geom.Point{0.0, palPrevHeight + 2*gradEditHeight})
		group.Add(geBlue)

	// Analoger Aufbau, jedoch fuer eine prozedurale Palette.
	// >>> Dieser Teil ist noch im Aufbau begriffen <<<
	case *mandel.ProcPalette:

		peRed := NewProcEditor(palImpl, mandel.Red)
		peRed.SetPos(geom.Point{0.0, palPrevHeight})
		group.Add(peRed)

		peGreen := NewProcEditor(palImpl, mandel.Green)
		peGreen.SetPos(geom.Point{0.0, palPrevHeight + procEditHeight})
		group.Add(peGreen)

		peBlue := NewProcEditor(palImpl, mandel.Blue)
		peBlue.SetPos(geom.Point{0.0, palPrevHeight + 2*procEditHeight})
		group.Add(peBlue)
	}

	//    lbl := adagui.NewLabel("Test")
	//    lbl.SetPos(geom.Point{10.0, 300.0})
	//    lbl.SetTextColor(colornames.Gold)
	//    group.Add(lbl)

	// Start der Applikation
	screen.SetWindow(win)
	screen.Run()
}
