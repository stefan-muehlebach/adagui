package main

import (
	"flag"
	"image"
	"image/draw"
	"image/png"
	"log"
	"os"
	"os/signal"
	"time"

	"github.com/stefan-muehlebach/adagui"
	"github.com/stefan-muehlebach/adagui/binding"
	"github.com/stefan-muehlebach/adagui/props"
	"github.com/stefan-muehlebach/adagui/touch"
	"github.com/stefan-muehlebach/adatft"
	"github.com/stefan-muehlebach/gg"
	"github.com/stefan-muehlebach/gg/colors"
	"github.com/stefan-muehlebach/gg/geom"
	"github.com/stefan-muehlebach/ledgrid"
	ledcolor "github.com/stefan-muehlebach/ledgrid/colors"
)

//-----------------------------------------------------------------------

var (
	PickerProps = props.NewProperties(props.PropsMap["Default"])
)

func init() {
	PickerProps.SetColor(props.Color, colors.Black)
	PickerProps.SetColor(props.BorderColor, colors.WhiteSmoke)
	PickerProps.SetSize(props.BorderWidth, 2.0)
	PickerProps.SetSize(props.FieldSize, 20.0)
}

type ColorPicker struct {
	adagui.ContainerEmbed
	orient           adagui.Orientation
	colorList        [][]ledcolor.LedColor
	fieldSize        float64
	numCols, numRows int
	selIdx0, selIdx1 int
	value            binding.Untyped
	alpha            float64
}

func NewColorPicker(orient adagui.Orientation, numColsRows int) *ColorPicker {
	c := &ColorPicker{}
	c.Wrapper = c
	c.Init()
	c.PropertyEmbed.Init(PickerProps)
	c.orient = orient
	c.colorList = make([][]ledcolor.LedColor, numColsRows)
	c.fieldSize = c.FieldSize()
	switch orient {
	case adagui.Horizontal:
		c.numCols = 1
		c.numRows = numColsRows
	case adagui.Vertical:
		c.numCols = numColsRows
		c.numRows = 1
	}
	c.SetMinSize(geom.Point{c.fieldSize * float64(c.numCols),
		c.fieldSize * float64(c.numRows)})
	c.selIdx0, c.selIdx1 = 0, 0
	c.value = binding.NewUntyped()
	c.alpha = 0.5
	return c
}

func NewColorPickerWithCallback(orient adagui.Orientation,
	numColsRows int, callback func(colors.Color)) *ColorPicker {
	c := NewColorPicker(orient, numColsRows)
	c.value.AddCallback(func(data binding.DataItem) {
		callback(data.(binding.Untyped).Get().(colors.Color))
	})
	return c
}

func NewColorPickerWithData(orient adagui.Orientation,
	numColsRows int, data binding.Untyped) *ColorPicker {
	c := NewColorPicker(orient, numColsRows)
	c.value = data
	return c
}

func (c *ColorPicker) SetColors(colRowIdx int, cl []ledcolor.LedColor) {
	c.colorList[colRowIdx] = make([]ledcolor.LedColor, len(cl))
	copy(c.colorList[colRowIdx], cl)
	switch c.orient {
	case adagui.Horizontal:
		c.numCols = max(c.numCols, len(cl))
	case adagui.Vertical:
		c.numRows = max(c.numRows, len(cl))
	}
	c.SetSize(geom.Point{c.fieldSize * float64(c.numCols),
		c.fieldSize * float64(c.numRows)})
}

func (c *ColorPicker) Color() ledcolor.LedColor {
	return c.value.Get().(ledcolor.LedColor)
}

func (c *ColorPicker) SetColor(col ledcolor.LedColor) {
	ledColor := col
	for c.selIdx0 = range c.colorList {
		for c.selIdx1 = range c.colorList[c.selIdx0] {
			if ledColor == c.colorList[c.selIdx0][c.selIdx1] {
				c.value.Set(ledColor.Alpha(c.alpha))
				return
			}
		}
	}
}

func (c *ColorPicker) OnInputEvent(evt touch.Event) {
	//log.Printf("ColorPicker: %v", evt.Pos)
	switch evt.Type {
	case touch.TypeRelease:
		selPt := c.Bounds().SetInside(evt.Pos)
		dx, dy := selPt.AsCoord()
		col := int(dx / c.fieldSize)
		row := int(dy / c.fieldSize)
		switch c.orient {
		case adagui.Horizontal:
			c.selIdx0 = row
			c.selIdx1 = col
		case adagui.Vertical:
			c.selIdx0 = col
			c.selIdx1 = row
		}
		if c.selIdx0 >= len(c.colorList) ||
			c.selIdx1 >= len(c.colorList[c.selIdx0]) {
			return
		}
		c.value.Set(c.colorList[c.selIdx0][c.selIdx1].Alpha(c.alpha))
		c.Mark(adagui.MarkNeedsPaint)
	}
}

func (c *ColorPicker) SetSize(size geom.Point) {
	switch c.orient {
	case adagui.Horizontal:
		c.fieldSize = min(c.fieldSize, size.X/float64(c.numCols))
		size.Y = c.fieldSize * float64(c.numRows)
	case adagui.Vertical:
		c.fieldSize = min(c.fieldSize, size.Y/float64(c.numRows))
		size.X = c.fieldSize * float64(c.numCols)
	}
	c.ContainerEmbed.SetMinSize(size)
}

func (c *ColorPicker) Paint(gc *gg.Context) {
	var col, row float64
	var r geom.Rectangle

	gc.SetStrokeWidth(0.0)
	for i, colorSlice := range c.colorList {
		for j, color := range colorSlice {
			gc.SetFillColor(color)
			switch c.orient {
			case adagui.Horizontal:
				col = float64(j) * c.fieldSize
				row = float64(i) * c.fieldSize
			case adagui.Vertical:
				col = float64(i) * c.fieldSize
				row = float64(j) * c.fieldSize
			}
			gc.DrawRectangle(col, row, c.fieldSize, c.fieldSize)
			gc.FillStroke()
		}
	}

	gc.SetFillColor(c.Color())
	gc.SetStrokeColor(c.BorderColor())
	gc.SetStrokeWidth(c.BorderWidth())
	gc.DrawRectangle(c.Bounds().AsCoord())
	gc.Stroke()

	switch c.orient {
	case adagui.Horizontal:
		r = geom.NewRectangleWH(float64(c.selIdx1)*c.fieldSize,
			float64(c.selIdx0)*c.fieldSize, c.fieldSize, c.fieldSize)
	case adagui.Vertical:
		r = geom.NewRectangleWH(float64(c.selIdx0)*c.fieldSize,
			float64(c.selIdx1)*c.fieldSize, c.fieldSize, c.fieldSize)
	}
	r = r.Inset(-3, -3)
	color := c.value.Get().(ledcolor.LedColor)
	gc.SetFillColor(color)
	gc.SetStrokeWidth(2.0)
	gc.SetStrokeColor(c.BorderColor())
	gc.DrawRectangle(r.AsCoord())
	gc.FillStroke()
}

//-----------------------------------------------------------------------

var (
	GridProps = props.NewProperties(props.PropsMap["Default"])
)

func init() {
	GridProps.SetColor(props.Color, colors.Black)
	GridProps.SetColor(props.BorderColor, colors.WhiteSmoke)
	GridProps.SetColor(props.LineColor, colors.WhiteSmoke)
	GridProps.SetSize(props.BorderWidth, 2.0)
	GridProps.SetSize(props.LineWidth, 1.0)
	GridProps.SetSize(props.FieldSize, 30.0)
}

type LedGrid struct {
	adagui.ContainerEmbed
	fieldSize float64
	DrawColor ledcolor.LedColor
	quitQ     chan bool
	client    ledgrid.GridClient
	grid      *ledgrid.LedGrid
}

func NewLedGrid(size image.Point, host string, port uint) *LedGrid {
	g := &LedGrid{}
	g.Wrapper = g
	g.Init()
	g.PropertyEmbed.Init(GridProps)
	g.fieldSize = g.FieldSize()
	g.SetMinSize(geom.Point{g.fieldSize, g.fieldSize})
	g.DrawColor = ledcolor.LedColor{0x00, 0x00, 0x00, 0xFF}
	g.quitQ = make(chan bool)
	g.client = ledgrid.NewNetGridClient("raspi-3", "udp", ledgrid.DefDataPort,
		ledgrid.DefRPCPort)
	modConf := g.client.ModuleConfig()
	g.grid = ledgrid.NewLedGrid(g.client, modConf)
	//	g.grid = ledgrid.NewLedGridBySize(host, port, size)
	//	g.ctrl = ledgrid.NewNetPixelClient(host, port)
	return g
}

func (g *LedGrid) OnInputEvent(evt touch.Event) {
	//log.Printf("LedGrid: %v", evt.Pos)
	//log.Printf("    Pos: %v", g.Pos())
	//log.Printf("    Pos: %v", g.Pos())
	switch evt.Type {
	case touch.TypeDrag:
		//if !evt.Pos.In(g.Rect()) {
		//	break
		//}
		fx, fy := g.Bounds().PosRel(evt.Pos)
		col := int(fx * float64(g.grid.Rect.Dx()))
		row := int(fy * float64(g.grid.Rect.Dy()))
		oldColor := g.grid.LedColorAt(col, row)
		newColor := g.DrawColor.Mix(oldColor, ledcolor.Blend)
		g.grid.SetLedColor(col, row, newColor)
		g.Mark(adagui.MarkNeedsPaint)
	}
}

func (g *LedGrid) SetSize(size geom.Point) {
	if size.X > size.Y {
		g.fieldSize = size.X / float64(g.grid.Rect.Dx())
		size.Y = float64(g.grid.Rect.Dy()) * g.fieldSize
	} else {
		g.fieldSize = size.Y / float64(g.grid.Rect.Dy())
		size.X = float64(g.grid.Rect.Dx()) * g.fieldSize
	}
	g.ContainerEmbed.SetMinSize(size)
}

func (g *LedGrid) Clear(c ledcolor.LedColor) {
	for idx := 0; idx < len(g.grid.Pix); idx += 3 {
		g.grid.Pix[idx+0] = c.R
		g.grid.Pix[idx+1] = c.G
		g.grid.Pix[idx+2] = c.B
	}
	g.Mark(adagui.MarkNeedsPaint)
}

func (g *LedGrid) Paint(gc *gg.Context) {
	for row := 0; row < g.grid.Rect.Dy(); row++ {
		y0 := float64(row) * g.fieldSize
		for col := 0; col < g.grid.Rect.Dx(); col++ {
			x0 := float64(col) * g.fieldSize
			c := g.grid.At(col, row)
			gc.SetFillColor(c)
			gc.DrawRectangle(x0, y0, g.fieldSize, g.fieldSize)
			gc.Fill()
		}
	}

	gc.SetStrokeWidth(g.LineWidth() / 2.0)
	gc.SetStrokeColor(g.BorderColor())
	for t := g.fieldSize; t < g.Size().Y; t += g.fieldSize {
		gc.DrawLine(0.0, t, g.Size().X, t)
	}
	for t := g.fieldSize; t < g.Size().X; t += g.fieldSize {
		gc.DrawLine(t, 0.0, t, g.Size().Y)
	}
	gc.Stroke()

	gc.SetStrokeWidth(g.BorderWidth())
	gc.SetStrokeColor(g.BorderColor())
	gc.DrawRectangle(g.Bounds().AsCoord())
	gc.Stroke()

	g.grid.Show()
}

func (g *LedGrid) Save(fileName string) {
	fh, err := os.Create(fileName)
	if err != nil {
		log.Fatal(err)
	}
	if err := png.Encode(fh, g.grid); err != nil {
		fh.Close()
		log.Fatal(err)
	}
	if err := fh.Close(); err != nil {
		log.Fatal(err)
	}
}

func (g *LedGrid) Load(fileName string) {
	fh, err := os.Open(fileName)
	if err != nil {
		log.Fatal(err)
	}
	img, err := png.Decode(fh)
	if err != nil {
		fh.Close()
		log.Fatal(err)
	}
	draw.Draw(g.grid, img.Bounds(), img, image.Point{}, draw.Over)
	g.Mark(adagui.MarkNeedsPaint)
}

//-----------------------------------------------------------------------

const (
	host    = "raspi-3"
	port    = 5333
	iconDir = "icons"
)

func f1(t float64) float64 {
	return 3*t*t - 2*t*t*t
}

func init() {
	log.SetFlags(log.Lmicroseconds | log.Lmsgprefix)
	log.SetPrefix(": ")
}

func SignalHandler() {
	sigChan := make(chan os.Signal)
	signal.Notify(sigChan, os.Interrupt)
	<-sigChan
	screen.Quit()
}

var (
	screen *adagui.Screen
	win    *adagui.Window
)

func main() {
	var iconFiles []string
	//var iconIdx int

	flag.Parse()
	adagui.StartProfiling()

	files, err := os.ReadDir(iconDir)
	if err != nil {
		log.Fatal(err)
	}
	iconFiles = make([]string, 0)
	for _, file := range files {
		iconFiles = append(iconFiles, file.Name())
	}

	screen = adagui.NewScreen(adatft.Rotate090)
	win = screen.NewWindow()

	root := adagui.NewGroup()
	root.Layout = adagui.NewPaddedLayout()
	win.SetRoot(root)

	mainGrp := adagui.NewGroup()
	mainGrp.Layout = adagui.NewVBoxLayout(5)
	root.Add(mainGrp)

	drawColor := ledcolor.LedColor{}
	colorValue := binding.BindUntyped(&drawColor)
	colorPicker := NewColorPickerWithData(adagui.Horizontal, 2, colorValue)

	ledPanel := NewLedGrid(image.Point{40, 10}, host, port)

	picoPal := ledgrid.PaletteMap["Pico08"].(*ledgrid.SlicePalette)
	colorPicker.SetColors(0, picoPal.Colors[:16])
	colorPicker.SetColors(1, picoPal.Colors[16:32])
	colorValue.AddCallback(func(data binding.DataItem) {
		if data == nil {
			return
		}
		ledPanel.DrawColor = data.(binding.Untyped).Get().(ledcolor.LedColor)
	})
	colorPicker.SetColor(picoPal.Colors[0])

	//    alphaGroup := adagui.NewGroup()
	//    alphaGroup.Layout = adagui.NewHBoxLayout()

	//    alphaLabel := adagui.NewLabel("Alpha:")
	alphaSlider := adagui.NewSliderWithCallback(10.0, adagui.Horizontal,
		func(v float64) {
			colorPicker.alpha = v
			color := colorPicker.Color()
			colorPicker.value.Set(color.Alpha(colorPicker.alpha))
		})
	alphaSlider.SetInitValue(1.0)
	alphaSlider.SetRange(0.0, 1.0, 0.1)
	//    alphaGroup.Add(alphaLabel, alphaSlider)

	mainGrp.Add(ledPanel, colorPicker, alphaSlider)

	/*
	   	toolGrp := adagui.NewGroupPL(mainGrp, adagui.NewHBoxLayout())
	   	btnPrev := adagui.NewTextButton("Prev")
	   	btnPrev.SetOnTap(func(evt touch.Event) {
	           if iconIdx <= 0 {
	               return
	           }
	           iconIdx -= 1
	           fileName := filepath.Join(iconDir, iconFiles[iconIdx])
	           ledPanel.Load(fileName)
	   	})
	   	btnNext := adagui.NewTextButton("Next")
	   	btnNext.SetOnTap(func(evt touch.Event) {
	           if iconIdx >= len(iconFiles)-1 {
	               return
	           }
	           iconIdx += 1
	           fileName := filepath.Join(iconDir, iconFiles[iconIdx])
	           ledPanel.Load(fileName)
	   	})
	   	toolGrp.Add(btnPrev, adagui.NewSpacer(), btnNext)
	*/

	//mainGrp.Add(adagui.NewSpacer())

	buttonGrp := adagui.NewGroup()
	buttonGrp.Layout = adagui.NewHBoxLayout()
	mainGrp.Add(buttonGrp)

	btnQuit := adagui.NewTextButton("Quit")
	btnQuit.SetOnTap(func(evt touch.Event) {
		ledPanel.Clear(ledcolor.LedColor{0, 0, 0, 0xff})
		time.Sleep(100 * time.Millisecond)
		screen.Quit()
	})
	buttonGrp.Add(adagui.NewSpacer(), btnQuit)

	screen.SetWindow(win)
	screen.Run()

	adagui.StopProfiling()
	time.Sleep(100 * time.Millisecond)
}
