// Einfachstes Beispiel eines AdaGui-Programmes. Erzeugt nur ein leeres,
// farbiges Feld.
package main

import (
	"flag"
	"fmt"
	"log"
	"math"
	"math/rand"
	"os"
	"os/signal"
	//"sort"
	"github.com/stefan-muehlebach/adagui"
	"github.com/stefan-muehlebach/adagui/binding"
	"github.com/stefan-muehlebach/adagui/props"
	"github.com/stefan-muehlebach/adagui/touch"
	"github.com/stefan-muehlebach/adatft"
	"github.com/stefan-muehlebach/gg/color"
	"github.com/stefan-muehlebach/gg/fonts"
	"github.com/stefan-muehlebach/gg/geom"
)

var (
	scr *adagui.Screen
	win *adagui.Window
)

func SignalHandler() {
	sigChan := make(chan os.Signal)
	signal.Notify(sigChan, os.Interrupt)
    <-sigChan
	scr.Quit()
}

type panelInitFunc func(size geom.Point) adagui.Node

type panelType struct {
    name string
	initFunc panelInitFunc
	panel    adagui.Node
}

var (
	panelList = []panelType{
		{"Overview", WidgetGallery, nil},
		{"Fonts", ScrolledFontPanel, nil},
		{"Colors", ScrolledColorPanel, nil},
		{"Draw", NestedTransformations, nil},
//		{"Type", Keyboard, nil},
		{"Sliders", SliderAndScrollbar, nil},
        {"Text", TextAlignment, nil},
		//{"7", BorderLayout, nil},
		//{"8", EmptyScrollPanel, nil},
	}
)

// WidgetGallery zeigt die meisten der GUI-Elemente in Aktion.
func WidgetGallery(size geom.Point) adagui.Node {
//	root := adagui.NewGroup()
//	root.Layout = adagui.NewPaddedLayout()
//	root.SetMinSize(size)

	grpMain := adagui.NewGroup()
	grpMain.Layout = adagui.NewVBoxLayout()
//	root.Add(grpMain)

	grpTop := adagui.NewGroup()
	grpTop.Layout = adagui.NewHBoxLayout()
	grpMain.Add(grpTop)

	grpChk := adagui.NewGroupPL(grpTop, adagui.NewVBoxLayout())
	chk1 := adagui.NewCheckbox("Extra Pommes")
	chk2 := adagui.NewCheckbox("Big Coke")
	chk3 := adagui.NewCheckbox("Zum Mitnehmen")
	chk4 := adagui.NewCheckbox("Zum Fortschmeissen")
	grpChk.Add(chk1, chk2, chk3, chk4)

	grpChk = adagui.NewGroupPL(grpTop, adagui.NewVBoxLayout())
	chk1 = adagui.NewCheckbox("Schweiz")
	chk2 = adagui.NewCheckbox("Europa")
	chk3 = adagui.NewCheckbox("Planet Erde")
	grpChk.Add(chk1, chk2, chk3)

	grpTop.Add(adagui.NewSpacer())

	radVal := binding.NewInt()
	radVal.Set(1)

	grpRad := adagui.NewGroupPL(grpTop, adagui.NewVBoxLayout())
	rad1 := adagui.NewRadioButtonWithData("Anfänger", 1, radVal)
	rad2 := adagui.NewRadioButtonWithData("Fortgeschritten", 2, radVal)
	rad3 := adagui.NewRadioButtonWithData("Profi", 3, radVal)
	grpRad.Add(rad1, rad2, rad3)

	grpIcon := adagui.NewGroupPL(grpMain, adagui.NewHBoxLayout())
	iconData := binding.NewInt()
	for i := 1; i <= 10; i++ {
		fileName := fmt.Sprintf("32x32/%02d.png", i)
		icn := adagui.NewIconButtonWithData(fileName, i, iconData)
		grpIcon.Add(icn)
	}

	grpSld1 := adagui.NewGroupPL(grpMain, adagui.NewHBoxLayout())

	val01 := binding.NewFloat()
	str01 := binding.FloatToStringWithFormat(val01, "%.1f")

	sld01 := adagui.NewSliderWithData(200, adagui.Horizontal, val01)
	sld01.SetRange(0.0, 1.0, 0.2)
	lbl8 := adagui.NewLabelWithData(str01)
	lbl8.SetSize(geom.Point{42, sld01.Size().Y})
	grpSld1.Add(sld01)
	grpSld1.Add(lbl8)

	grpSld2 := adagui.NewGroupPL(grpMain, adagui.NewHBoxLayout())

	val02 := binding.NewFloat()
	str02 := binding.FloatToStringWithFormat(val02, "%.3f")

	sld02 := adagui.NewSliderWithData(200, adagui.Horizontal, val02)
	sld02.SetRange(0.0, 2*math.Pi, math.Pi/36.0)
	lbl9 := adagui.NewLabelWithData(str02)
	lbl9.SetSize(geom.Point{42, sld02.Size().Y})
	grpSld2.Add(sld02)
	grpSld2.Add(lbl9)

	grpMain.Add(adagui.NewSpacer())

	grpBtn := adagui.NewGroup()
	grpBtn.Layout = adagui.NewHBoxLayout()
	grpMain.Add(grpBtn)

	btn1 := adagui.NewTextButton("Open")
	btn1.SetOnTap(func(evt touch.Event) {
		log.Printf("What do you want to open?")
	})
	btn2 := adagui.NewListButton([]string{"Sinus", "Triangle", "Sawtooth",
		"Random", "White Noise"})
	btn2.SetOnTap(func(evt touch.Event) {
		log.Printf("Set sound generator to '%s'", btn2.Selected)
	})
	btn3 := adagui.NewTextButton("Quit")
	btn3.SetOnTap(func(evt touch.Event) {
		scr.Quit()
	})
	grpBtn.Add(btn1, btn2, adagui.NewSpacer(), btn3)

	return grpMain
}

// ScrolledFontPanel zeigt erstens die Moeglichkeiten, Text in ansprechenden
// Fonts darzustellen und den Einsatz eines ScrolledPanels.
func ScrolledFontPanel(size geom.Point) adagui.Node {
	var i int
	var fontName string
    var scrHori, scrVert *adagui.Scrollbar

    log.Printf("ScrolledFontPanel(size): %+v", size)

	fontSize := 18.0
	textColor := color.WhiteSmoke
	fontList := fonts.Names

	virtualWidth := 2048.0
	virtualHeight := 3000.0

	panel := adagui.NewScrollPanel(0, 0)
	panel.Layout = adagui.NewVBoxLayout(10)
	panel.SetVirtualSize(geom.Point{virtualWidth, virtualHeight})

	for i, fontName = range fontList[:8] {
		if fontName == "Elegante" {
			fontList = fontList[i:]
			break
		}
		lbl := adagui.NewLabel(fontName)
		lbl.SetTextColor(textColor)
		lbl.SetFont(fonts.Map[fontName])
		lbl.SetFontSize(fontSize)
		panel.Add(lbl)
	}
/*
	fontSize *= 2
	for i, fontName = range fontList {
		if fontName == "Elzevier" {
			fontList = fontList[i:]
			break
		}
		lbl := adagui.NewLabel(fontName)
		lbl.SetTextColor(textColor)
		lbl.SetFont(fonts.Map[fontName])
		lbl.SetFontSize(fontSize)
		panel.Add(lbl)
	}
	fontSize *= 3
	for _, fontName = range fontList {
		lbl := adagui.NewLabel(fontName)
		lbl.SetTextColor(textColor)
		lbl.SetFont(fonts.Map[fontName])
		lbl.SetFontSize(fontSize)
		panel.Add(lbl)
	}
*/

	scrVert = adagui.NewScrollbarWithCallback(160, adagui.Vertical,
		func(f float64) {
			panel.SetYView(f)
		})

	scrHori = adagui.NewScrollbarWithCallback(240, adagui.Horizontal,
		func(f float64) {
			panel.SetXView(f)
		})

	visRange := panel.VisibleRange()
	scrHori.SetVisiRange(visRange.X)
	scrVert.SetVisiRange(visRange.Y)

    main := adagui.NewPanel(0, 0)
    layout := adagui.NewBorderLayout(nil, scrHori, nil, scrVert)
    //layout.(*adagui.BorderLayout).Padding = 5.0
    main.Layout = layout
    //main.SetColor(color.DarkRed.Alpha(0.5))
    main.Add(scrHori, scrVert, panel)

	return main
}

// ScrolledColorPanel
type ColorInfo struct {
	name  string
	color color.Color
}

func ScrolledColorPanel(size geom.Point) adagui.Node {
	root := adagui.NewGroup()
	root.SetMinSize(size)

	s := props.PropsMap["Scrollbar"].Size(props.Width)
	w, h := size.X-s-2, size.Y

	panel := adagui.NewScrollPanel(w, h)
	panel.SetVirtualSize(geom.Point{w, 870})
	panel.SetColor(color.Transparent)
	panel.SetBorderColor(color.Transparent)
	root.Add(panel)

	scrV := adagui.NewScrollbarWithCallback(h, adagui.Vertical,
		func(f float64) {
			panel.SetYView(f)
		})
	scrV.SetPos(panel.Rect().NE().AddXY(1, 0))
	root.Add(scrV)

	visRange := panel.VisibleRange()
	scrV.SetVisiRange(visRange.Y)

	colorList := make([]ColorInfo, 0)
	for _, nameList := range color.Groups {
		for _, name := range nameList {
			colorList = append(colorList, ColorInfo{name, color.Map[name]})
		}
	}

	numColumns := 5
	numRows := 29
	tileWidth := w / float64(numColumns)
	tileHeight := 870 / float64(numRows)

	for col := 0; col < numColumns; col++ {
		for row := 0; row < numRows; row++ {
			idx := col*numRows + row
			if idx >= len(colorList) {
				continue
			}
			x0, y0 := float64(col)*tileWidth, float64(row)*tileHeight
			colorInfo := colorList[idx]
			tile := adagui.NewRectangle(tileWidth-4.0, tileHeight-4.0)
			tile.SetPos(geom.Point{x0, y0})

			tile.SetColor(colorInfo.color)
			tile.SetPushedColor(colorInfo.color.Bright(0.5))
			tile.SetSelectedColor(colorInfo.color)
			tile.SetBorderColor(colorInfo.color)
			tile.SetPushedBorderColor(colorInfo.color.Bright(0.5))
			tile.SetSelectedBorderColor(colorInfo.color)

			tile.SetBorderWidth(0.0)
			tile.SetPushedBorderWidth(0.0)
			tile.SetSelectedBorderWidth(0.0)

			tile.SetOnDoubleTap(func(evt touch.Event) {
				log.Printf("remove tile")
				tile.Remove()
				panel.Mark(adagui.MarkNeedsPaint)
			})
			panel.Add(tile)
		}
	}
	return root
}

// NestedTransformations
//
func NestedTransformations(size geom.Point) adagui.Node {
	var root *adagui.Group
	var panel01, panel02, panel03 *adagui.Panel
	var color02, color03 color.Color
	var rotPt1, rotPt2 geom.Point
    var colorFactor float64 = 0.5

	root = adagui.NewGroup()
	root.Layout = adagui.NewPaddedLayout(0)
	root.SetMinSize(size)

	panel01 = NewPanel(0, 0)
	panel01.SetColor(color.RandColor().Dark(0.8))
	root.Add(panel01)

	grp0 := adagui.NewGroup()
	grp1 := adagui.NewGroup()
	grp0.SetPos(geom.Point{0, 0})
	grp1.SetPos(geom.Point{0, panel01.Size().Y - 26})
	panel01.Add(grp0, grp1)

	obj := NewCircle(60.0)
	obj.SetPos(geom.Point{15, panel01.Rect().Dy() / 2})

	rotVal1 := binding.NewFloat()
	scaleVal1 := binding.NewFloat()

    interGap := 10.0
    intraGap := 5.0
	len := (panel01.Size().X - interGap) / 2

    refPt := geom.Point{0, 3}
    lbl1 := adagui.NewLabel("Rotate:")
    lbl1.SetFont(fonts.GoBold)
    lbl1.SetPos(refPt)
    lbl1.SetTextColor(lbl1.TextColor().Alpha(0.7))
	rotSld1 := adagui.NewSliderWithData(len - lbl1.Size().X - intraGap,
        adagui.Horizontal, rotVal1)
	rotSld1.SetPos(refPt.AddXY(lbl1.Size().X + intraGap, 0))
	rotSld1.SetRange(-math.Pi/3, math.Pi/3, math.Pi/72.0)
	rotSld1.SetInitValue(0.0)

    refPt = refPt.AddXY(len+interGap, 0)
    lbl2 := adagui.NewLabel("Scale:")
    lbl2.SetFont(fonts.GoBold)
    lbl2.SetPos(refPt)
    lbl2.SetTextColor(lbl2.TextColor().Alpha(0.7))
	scaleSld1 := adagui.NewSliderWithData(len - lbl2.Size().X - intraGap,
        adagui.Horizontal, scaleVal1)
	scaleSld1.SetPos(refPt.AddXY(lbl2.Size().X + intraGap, 0))
	scaleSld1.SetRange(0.2, 1.8, 0.05)
	scaleSld1.SetInitValue(1.0)

	hSpc := 30.0
	vSpc := 10.0
	w, h := panel01.Size().AsCoord()
	w, h = w-2*hSpc, h-2*vSpc-scaleSld1.Size().Y

	color02 = color.RandColor().Dark(colorFactor)
	panel02 = NewPanel(w, h)
	panel02.SetPos(geom.Point{hSpc, vSpc})
	panel02.SetColor(color02)

	grp0.Add(obj, panel02)
	grp1.Add(rotSld1, scaleSld1, lbl1, lbl2)

	grp0 = adagui.NewGroup()
	grp1 = adagui.NewGroup()
	grp0.SetPos(geom.Point{0, 0})
	grp1.SetPos(geom.Point{0, panel02.Size().Y - 26})
	panel02.Add(grp0, grp1)

	obj = NewCircle(30.0)
	obj.SetPos(panel02.Size())
	grp0.Add(obj)

	rotPt1 = panel02.Size()
	rotPt1.X = 0.0
	rotPt1.Y /= 2.0

	rotVal1.AddCallback(func(data binding.DataItem) {
		f := data.(binding.Float).Get()
		panel02.RotateAbout(rotPt1, f)
	})
	scaleVal1.AddCallback(func(data binding.DataItem) {
		f := data.(binding.Float).Get()
		panel02.ScaleAbout(rotPt1, f, f)
	})

	bool01 := binding.NewBool()
	bool01.AddCallback(func(data binding.DataItem) {
		if data.(binding.Bool).Get() {
			panel02.SetColor(color02.Bright(colorFactor))
		} else {
			panel02.SetColor(color02)
		}
	})

	chk := adagui.NewCheckboxWithData("Background Bright", bool01)
	chk.SetPos(geom.Point{5, 5})
	grp0.Add(chk)

	rotVal2 := binding.NewFloat()
	scaleVal2 := binding.NewFloat()

	len = (panel02.Size().X - 9) / 2
    lbl1 = adagui.NewLabel("Rotate")
    lbl1.SetFont(fonts.GoBold)
    lbl1.SetPos(geom.Point{3, 3})
    lbl1.SetTextColor(lbl1.TextColor().Alpha(0.7))
	rotSld2 := adagui.NewSliderWithData(len, adagui.Horizontal, rotVal2)
	rotSld2.SetPos(geom.Point{3, 3})
	rotSld2.SetRange(-math.Pi/3, math.Pi/3, math.Pi/72.0)
	rotSld2.SetInitValue(0.0)

    lbl2 = adagui.NewLabel("Scale")
    lbl2.SetFont(fonts.GoBold)
    lbl2.SetPos(geom.Point{len + 6, 3})
    lbl2.SetTextColor(lbl2.TextColor().Alpha(0.7))
	scaleSld2 := adagui.NewSliderWithData(len, adagui.Horizontal, scaleVal2)
	scaleSld2.SetPos(geom.Point{len + 6, 3})
	scaleSld2.SetRange(0.2, 1.8, 0.05)
	scaleSld2.SetInitValue(1.0)

	grp1.Add(rotSld2, scaleSld2, lbl1, lbl2)

	hSpc = 5.0
	vSpc = 30.0
	w, h = panel02.Size().AsCoord()
	w, h = w-hSpc-5, h-vSpc-25.0-5

	color03 = color.RandColor().Dark(colorFactor)
	panel03 = NewPanel(w, h)
	panel03.SetPos(geom.Point{hSpc, vSpc})
	panel03.SetColor(color03)
	grp0.Add(panel03)

	rotPt2 = panel03.Size()
	rotPt2.Y /= 2.0

	rotVal2.AddCallback(func(data binding.DataItem) {
		f := data.(binding.Float).Get()
		panel03.RotateAbout(rotPt2, f)
	})
	scaleVal2.AddCallback(func(data binding.DataItem) {
		f := data.(binding.Float).Get()
		panel03.ScaleAbout(rotPt2, f, f)
	})

	bool02 := binding.NewBool()
	bool02.AddCallback(func(data binding.DataItem) {
		if data.(binding.Bool).Get() {
			panel03.SetColor(color03.Bright(colorFactor))
		} else {
			panel03.SetColor(color03)
		}
	})

	chk = adagui.NewCheckboxWithData("Background Bright", bool02)
	chk.SetPos(geom.Point{5, 5})
	panel03.Add(chk)

	return root
}

func NewPanel(w, h float64) *adagui.Panel {
	var c *adagui.Circle

	p := adagui.NewPanel(w, h)

	p.SetOnPress(func(evt touch.Event) {
		// log.Printf("Press on Panel: %v", evt)
	})
	p.SetOnLongPress(func(evt touch.Event) {
		c = NewCircle(1.0)
		c.SetPos(evt.Pos)
		p.Add(c)
		p.Mark(adagui.MarkNeedsPaint)
	})
	p.SetOnDrag(func(evt touch.Event) {
		if !evt.LongPressed {
			return
		}
		r := evt.Pos.Distance(evt.InitPos)
		c.SetRadius(r)
		p.Mark(adagui.MarkNeedsPaint)
	})
	p.SetOnTap(func(evt touch.Event) {
		r := 30.0 + 10.0*rand.Float64()
		c = NewCircle(r)
		c.SetPos(evt.Pos)
		p.Add(c)
		p.Mark(adagui.MarkNeedsPaint)
	})

	return p
}

func NewCircle(r float64) *adagui.Circle {
	var dp geom.Point

	c := adagui.NewCircle(r)
	col := color.RandColor()

	c.SetColor(col)
	c.SetPushedColor(col.Alpha(0.5))
	c.SetSelectedColor(col.Alpha(0.5))

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
		if !c.IsAtFront() {
			c.ToFront()
		} else {
			c.ToBack()
		}
		c.Mark(adagui.MarkNeedsPaint)
	})
	c.SetOnDoubleTap(func(evt touch.Event) {
		p := c.Wrappee().Parent
		c.Remove()
		p.Mark(adagui.MarkNeedsPaint)
	})

	return c
}

// SliderAndScrollbar
func SliderAndScrollbar(size geom.Point) adagui.Node {
	root := adagui.NewGroup()
	root.Layout = adagui.NewPaddedLayout()
	root.SetMinSize(size)

	panel := adagui.NewGroupPL(root, adagui.NewVBoxLayout())

	scr := adagui.NewScrollbar(160.0, adagui.Horizontal)
	scr.SetVisiRange(0.1)
	panel.Add(scr)

	scr = adagui.NewScrollbar(160.0, adagui.Horizontal)
	scr.SetVisiRange(0.3)
	panel.Add(scr)

	scr = adagui.NewScrollbar(160.0, adagui.Horizontal)
	scr.SetVisiRange(0.7)
	panel.Add(scr)

	scr = adagui.NewScrollbar(160.0, adagui.Horizontal)
	scr.SetVisiRange(1.0)
	panel.Add(scr)

	val1 := binding.NewFloat()
	str1 := binding.FloatToStringWithFormat(val1, "%.1f")
	val2 := binding.NewFloat()
	str2 := binding.FloatToStringWithFormat(val2, "%-1.0f")
	val3 := binding.NewFloat()
	str3 := binding.FloatToStringWithFormat(val3, "%.3f")

	sld1 := adagui.NewSliderWithData(160.0, adagui.Horizontal, val1)
	panel.Add(sld1)

	sld2 := adagui.NewSliderWithData(160.0, adagui.Horizontal, val2)
	sld2.SetRange(-5, 5, 1)
	sld2.SetValue(0)
	panel.Add(sld2)

	sld3 := adagui.NewSliderWithData(160.0, adagui.Horizontal, val3)
	sld3.SetRange(0.0, 2*math.Pi, 0.01)
	sld3.SetValue(0.0)
	panel.Add(sld3)

	row := adagui.NewGroupPL(panel, adagui.NewHBoxLayout())

	lbl1 := adagui.NewLabelWithData(str1)
	lbl2 := adagui.NewLabelWithData(str2)
	lbl3 := adagui.NewLabelWithData(str3)
	row.Add(lbl1, adagui.NewSpacer(), lbl2, adagui.NewSpacer(), lbl3)

	return root
}

func TextAlignment(size geom.Point) adagui.Node {
	root := adagui.NewGroup()
	root.Layout = adagui.NewPaddedLayout()
	root.SetMinSize(size)

	panel := adagui.NewGroup()
    root.Add(panel)

    r := panel.Bounds()

    lbl := adagui.NewLabel("AlignLeft | AlignTop")
    lbl.SetAlign(adagui.AlignLeft | adagui.AlignTop)
    lbl.SetFont(fonts.GoBold)
    lbl.SetFontSize(18.0)
    lbl.SetPos(r.NW())
	panel.Add(lbl)

    lbl = adagui.NewLabel("AlignRight | AlignTop")
    lbl.SetAlign(adagui.AlignRight | adagui.AlignTop)
    lbl.SetFont(fonts.GoBold)
    lbl.SetFontSize(18.0)
    lbl.SetPos(r.NE())
	panel.Add(lbl)

    lbl = adagui.NewLabel("AlignLeft  | AlignBottom")
    lbl.SetAlign(adagui.AlignLeft | adagui.AlignBottom)
    lbl.SetFont(fonts.GoBold)
    lbl.SetFontSize(18.0)
    lbl.SetPos(r.SW())
	panel.Add(lbl)

    lbl = adagui.NewLabel("AlignRight | AlignBottom")
    lbl.SetAlign(adagui.AlignRight | adagui.AlignBottom)
    lbl.SetFont(fonts.GoBold)
    lbl.SetFontSize(18.0)
    lbl.SetPos(r.SE())
	panel.Add(lbl)

    lbl = adagui.NewLabel("AlignCenter | AlignMiddle")
    lbl.SetAlign(adagui.AlignCenter | adagui.AlignMiddle)
    lbl.SetFont(fonts.GoBold)
    lbl.SetFontSize(18.0)
    lbl.SetPos(r.C())
	panel.Add(lbl)

	return root
}
// Keyboard als Miniatur-Tastatur
var (
	KeyboardProps              = props.NewProperties(props.PropsMap["Button"])
	KeyboardShift binding.Bool = binding.NewBool()
)

func init() {
	KeyboardProps.SetSize(props.BorderWidth, 0.0)
	KeyboardProps.SetSize(props.CornerRadius, 4.0)
	KeyboardProps.SetSize(props.Width, 30.0)
	KeyboardProps.SetSize(props.Height, 30.0)

	KeyboardShift.Set(false)
}

func NewKeyboard(text binding.String) adagui.Node {
	var keySpace float64 = 2.0

	panel := adagui.NewGroup()
	panel.Layout = adagui.NewVBoxLayout(keySpace)

	row0 := adagui.NewGroup()
	row0.Layout = adagui.NewHBoxLayout(keySpace)
	for _, ch := range "`1234567890-=" {
		key := NewKey(ch, text)
		row0.Add(key)
	}

	row1 := adagui.NewGroup()
	row1.Layout = adagui.NewHBoxLayout(keySpace)
	spc := adagui.NewSpacer()
	spc.FixHorizontal = true
	spc.SetMinSize(geom.Point{32 + 14, 0})
	row1.Add(spc)
	for _, ch := range "qwertyuiop[]\\" {
		key := NewKey(ch, text)
		row1.Add(key)
	}

	row2 := adagui.NewGroup()
	row2.Layout = adagui.NewHBoxLayout(keySpace)
	spc = adagui.NewSpacer()
	spc.FixHorizontal = true
	spc.SetMinSize(geom.Point{32 + 16 + 8, 0})
	row2.Add(spc)
	for _, ch := range "asdfghjkl;'" {
		key := NewKey(ch, text)
		row2.Add(key)
	}

	row3 := adagui.NewGroup()
	row3.Layout = adagui.NewHBoxLayout(keySpace)
	shift := NewShiftKey()
	row3.Add(shift)
	for _, ch := range "zxcvbnm,./" {
		key := NewKey(ch, text)
		row3.Add(key)
	}
	shift = NewShiftKey()
	row3.Add(shift)

	row4 := adagui.NewGroup()
	row4.Layout = adagui.NewHBoxLayout(keySpace)
	space := NewKey(' ', text)
	space.SetMinSize(geom.Point{4*31 + 29, 29})
	row4.Add(space)
	del := NewDelKey(text)
	row4.Add(del)

	panel.Add(row0, row1, row2, row3, row4)

	return panel
}

func NewKey(ch rune, text binding.String) *adagui.TextButton {
	var keySize geom.Point = geom.Point{
		KeyboardProps.Size(props.Width),
		KeyboardProps.Size(props.Height),
	}

	btn := adagui.NewTextButton(string(ch))
	btn.LeafEmbed.Init()
	btn.PushEmbed.Init(btn, nil)
	btn.PropertyEmbed.Init(KeyboardProps)
	btn.SetMinSize(keySize)
	btn.SetOnTap(func(evt touch.Event) {
		var str string
		if ch >= 'A' && ch <= 'Z' {
			if KeyboardShift.Get() {
				str = string(ch)
				KeyboardShift.Set(false)
			} else {
				str = string(ch + ('a' - 'A'))
			}
		} else {
			str = string(ch)
		}
		text.Set(text.Get() + str)
	})
	return btn
}

func NewShiftKey() *adagui.IconButton {
	var keySize geom.Point = geom.Point{60, 29}

	btn := adagui.NewIconButton("icons/1.png")
	btn.LeafEmbed.Init()
	btn.PushEmbed.Init(btn, KeyboardShift)
	btn.PropertyEmbed.Init(KeyboardProps)
	btn.SetMinSize(keySize)
	return btn
}

func NewDelKey(text binding.String) *adagui.TextButton {
	var keySize geom.Point = geom.Point{29, 29}

	btn := adagui.NewTextButton("\u25c4")
	btn.LeafEmbed.Init()
	btn.PushEmbed.Init(btn, nil)
	btn.PropertyEmbed.Init(KeyboardProps)
	btn.SetMinSize(keySize)
	btn.SetOnTap(func(evt touch.Event) {
		str := text.Get()
		if len(str) < 1 {
			return
		}
		text.Set(str[:len(str)-1])
	})
	return btn
}

func Keyboard(size geom.Point) adagui.Node {
	var keySpace float64 = 2.0
	var text binding.String = binding.NewString()

	root := adagui.NewGroup()
	root.Layout = adagui.NewPaddedLayout(5)
	root.SetMinSize(size)

	panel := adagui.NewGroup()
	panel.Layout = adagui.NewVBoxLayout(keySpace)

	lbl := adagui.NewLabelWithData(text)
	lbl.SetFontSize(24.0)
	lbl.SetFont(fonts.GoRegular)
	lbl.SetBorderColor(color.Silver)
	lbl.SetBorderWidth(1.0)

	keyboard := NewKeyboard(text)
	panel.Add(lbl, adagui.NewSpacer(), keyboard)

	root.Add(panel)
	return root
}

// BorderLayout example
func BorderLayout(size geom.Point) adagui.Node {
	topWidget := adagui.NewPanel(0, 0)
	topWidget.Layout = adagui.NewHBoxLayout()

	lbl := adagui.NewLabel("Euclid on a RaspberryPi")
	lbl.SetTextColor(color.WhiteSmoke)
	lbl.SetFont(fonts.LucidaBrightDemibold)
	lbl.SetFontSize(18.0)
	topWidget.Add(lbl)

	bottomWidget := adagui.NewPanel(0, 0)
	bottomWidget.Layout = adagui.NewHBoxLayout()

	toolVal := binding.NewInt()
	toolVal.Set(0)
	for i := 1; i <= 5; i++ {
		fileName := fmt.Sprintf("icons/%d.png", i)
		icn := adagui.NewIconButtonWithData(fileName, i, toolVal)
		bottomWidget.Add(icn)
	}
	bottomWidget.Add(adagui.NewSpacer())

	btnQuit := adagui.NewTextButton("Quit")
	btnQuit.SetOnTap(func(evt touch.Event) {
		scr.Quit()
	})
	bottomWidget.Add(btnQuit)

	/*
	   icn := adagui.NewIconButtonWithData("icons/handInv.png", toolVal)
	   bottomWidget.Add(icn)
	   icn  = adagui.NewIconButtonWithData("icons/point.png", toolVal)
	   bottomWidget.Add(icn)
	   icn  = adagui.NewIconButtonWithData("icons/segment.png", toolVal)
	   bottomWidget.Add(icn)
	*/

	leftWidget := adagui.NewPanel(20, 20)
	leftWidget.Layout = adagui.NewVBoxLayout()
	//leftWidget.FillColor = color.Gold

	rightWidget := adagui.NewPanel(50, 50)
	rightWidget.Layout = adagui.NewVBoxLayout()
	//rightWidget.FillColor = color.DarkGreen

	centerWidget := adagui.NewPanel(0, 0)
	centerWidget.SetColor(color.MidnightBlue)

	root := adagui.NewGroup()
	root.Layout = adagui.NewPaddedLayout()
	root.SetMinSize(size)

	panel := adagui.NewPanel(0, 0)
	panel.Layout = adagui.NewBorderLayout(topWidget, bottomWidget, nil, nil)
	//panel.FillColor = color.DimGray

	panel.Add(topWidget, bottomWidget, centerWidget)
	root.Add(panel)

	return root
}

// Ein leeres ScrollPanel, damit die Koordinaten-Transformationen besser
// studiert werden koennen.
func EmptyScrollPanel(size geom.Point) adagui.Node {
	virtualWidth := 512.0
	virtualHeight := 512.0

	root := adagui.NewGroup()
	root.SetMinSize(size)

	s := props.PropsMap["Scrollbar"].Size(props.Width)
	w, h := size.X-s-2, size.Y-s-2

	panel := adagui.NewScrollPanel(w, h)
	panel.SetVirtualSize(geom.Point{virtualWidth, virtualHeight})
	panel.SetColor(color.AntiqueWhite)
	panel.SetBorderColor(color.AntiqueWhite)
	//panel.SetOnPress(func(evt touch.Event) {
	//	log.Printf("evt: %+v", evt)
	//	log.Printf("LocalBounds  : %v", panel.LocalBounds())
	//})
	//panel.SetOnDrag(func(evt touch.Event) {
	//	log.Printf("evt: %+v", evt)
	//})
	root.Add(panel)

	c := adagui.NewCircle(20.0)
	c.SetPos(geom.Point{128, 128})
	//c.SetOnPress(func(evt touch.Event) {
	//    log.Printf("press on circle 1: event data: %+v", evt)
	//})
	//c.SetOnDrag(func(evt touch.Event) {
	//    log.Printf("drag on circle 1: event data: %+v", evt)
	//})
	//c.SetOnLeave(func(evt touch.Event) {
	//    log.Printf("leave on circle 1: event data: %+v", evt)
	//})
	panel.Add(c)
	c = adagui.NewCircle(20.0)
	c.SetPos(geom.Point{256, 128})
	panel.Add(c)
	c = adagui.NewCircle(20.0)
	c.SetPos(geom.Point{384, 128})
	panel.Add(c)

	scrV := adagui.NewScrollbarWithCallback(h, adagui.Vertical,
		func(f float64) {
			panel.SetYView(f)
		})
	scrV.SetPos(panel.Rect().NE().AddXY(1, 0))

	scrH := adagui.NewScrollbarWithCallback(w, adagui.Horizontal,
		func(f float64) {
			panel.SetXView(f)
		})
	scrH.SetPos(panel.Rect().SW().AddXY(0, 1))

	visRange := panel.VisibleRange()
	scrH.SetVisiRange(visRange.X)
	scrV.SetVisiRange(visRange.Y)

	root.Add(scrV, scrH)

	return root
}

//----------------------------------------------------------------------------

func main() {
	var tabPanel *adagui.Panel
	var tabMenu *adagui.TabMenu
	var tabContent *adagui.Panel

	flag.Parse()
	adagui.StartProfiling()

	scr = adagui.NewScreen(adatft.Rotate090)
	win = scr.NewWindow()

	tabContent = adagui.NewPanel(0, 0)
    tabContent.SetColor(color.Teal.Dark(0.8))
	tabContent.Layout = adagui.NewMaxLayout()
	tabMenu = adagui.NewTabMenu(tabContent)
    tabPanel = adagui.NewPanel(0, 0)
 //   tabPanel.SetColor(color.DarkViolet.Alpha(0.5))
    tabPanel.Layout = adagui.NewBorderLayout(tabMenu, nil, nil, nil)
    tabPanel.Add(tabMenu, tabContent)
	//tabPanel = adagui.NewTabPanel(float64(adatft.Width),
	//	float64(adatft.Height), tabMenu, tabContent)
	win.SetRoot(tabPanel)

    //log.Printf("tabContent.Size(): %v", tabContent.Rect())

	for i, panelInfo := range panelList {
		if panelInfo.initFunc == nil {
			continue
		}
		panelList[i].panel = panelInfo.initFunc(geom.Point{
			float64(adatft.Width), float64(adatft.Height) - tabMenu.Height()})
	}

	go SignalHandler()

	for _, panelInfo := range panelList {
		tabMenu.AddTab(panelInfo.name, panelInfo.panel)
	}
    tabMenu.SetTab(0)
	scr.SetWindow(win)
	scr.Run()

	adagui.StopProfiling()
}
