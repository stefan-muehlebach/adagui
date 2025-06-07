package main

// Dieses Programm dient der Demonstration der vorhandenen Widgets aus dem
// adagui-Package.
//
import (
	"flag"
	"fmt"
	"log"
	"math"
	"math/rand"
	"os"
	"os/signal"

	"github.com/stefan-muehlebach/adagui"
	"github.com/stefan-muehlebach/adagui/binding"
	"github.com/stefan-muehlebach/adagui/props"
	"github.com/stefan-muehlebach/adagui/touch"
	"github.com/stefan-muehlebach/adatft"
	"github.com/stefan-muehlebach/gg/colors"
	"github.com/stefan-muehlebach/gg/fonts"
	"github.com/stefan-muehlebach/gg/geom"
)

//-----------------------------------------------------------------------

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
	var rotation adatft.RotationType = adatft.Rotate090

	flag.Var(&rotation, "rotation", "display rotation")
	flag.Parse()

	screen = adagui.NewScreen(rotation)
	win = screen.NewWindow()

	widgetFlex := widgetFlex()

	win.SetRoot(widgetFlex)
	screen.SetWindow(win)
	screen.Run()
}

// ScrolledFontPanel zeigt erstens die Moeglichkeiten, Text in ansprechenden
// Fonts darzustellen und den Einsatz eines ScrolledPanels.
func ScrolledFontPanel() adagui.Node {
	var fontName string
	var scrHori, scrVert *adagui.Scrollbar
	var size geom.Point = geom.Point{float64(adatft.Width),
		float64(adatft.Height - 30)}

	//log.Printf("ScrolledFontPanel(size): %+v", size)

	fontSize := 24.0
	textColor := colors.WhiteSmoke
	fontList := fonts.Names

	virtualWidth := 1024.0
	virtualHeight := 4000.0

	panel := adagui.NewScrollPanel(0, 0)
	panel.Layout = adagui.NewVBoxLayout(10)
	panel.SetSize(size)
	panel.SetVirtualSize(geom.Point{virtualWidth, virtualHeight})

	for _, fontName = range fontList {
		//		if fontName == "Elegante" {
		//			fontList = fontList[i:]
		//			break
		//		}
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

	scrVert = adagui.NewScrollbarWithCallback(200, adagui.Vertical,
		func(f float64) {
			panel.SetYView(f)
		})

	scrHori = adagui.NewScrollbarWithCallback(400, adagui.Horizontal,
		func(f float64) {
			panel.SetXView(f)
		})

	visRange := panel.VisibleRange()
	scrHori.SetVisiRange(visRange.X)
	scrVert.SetVisiRange(visRange.Y)

	main := adagui.NewGroup()
	main.SetMinSize(size)
	main.Layout = adagui.NewBorderLayout(nil, scrHori, nil, scrVert)
	main.Add(scrHori, scrVert, panel)

	return main
}

// ScrolledColorPanel
type ColorInfo struct {
	name  string
	color colors.Color
}

func ScrolledColorPanel() adagui.Node {
	var size geom.Point = geom.Point{float64(adatft.Width),
		float64(adatft.Height - 30)}

	s := props.PropsMap["Scrollbar"].Size(props.Width)
	w, h := size.X-s-2, size.Y

	panel := adagui.NewScrollPanel(0, 0)
	panel.SetSize(size)
	panel.SetVirtualSize(geom.Point{w, 870})
	panel.SetColor(colors.Transparent)
	panel.SetBorderColor(colors.Transparent)

	scrVert := adagui.NewScrollbarWithCallback(h, adagui.Vertical,
		func(f float64) {
			panel.SetYView(f)
		})
	scrVert.SetPos(panel.Rect().NE().AddXY(1, 0))

	visRange := panel.VisibleRange()
	scrVert.SetVisiRange(visRange.Y)

	colorList := make([]ColorInfo, 0)
	for _, nameList := range colors.Groups {
		for _, name := range nameList {
			colorList = append(colorList, ColorInfo{name, colors.Map[name]})
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

	main := adagui.NewGroup()
	main.SetMinSize(size)
	main.Layout = adagui.NewBorderLayout(nil, nil, nil, scrVert)
	main.Add(panel, scrVert)

	return main
}

func WidgetPanel01() adagui.Node {
	var iconList []*adagui.IconButton

	grpMain := adagui.NewGroup()
	grpMain.Layout = adagui.NewVBoxLayout()

	grpOptions := adagui.NewGroupPL(grpMain, adagui.NewHBoxLayout())

	grpCheck := adagui.NewGroupPL(grpOptions, adagui.NewVBoxLayout())
	chk01 := adagui.NewCheckbox("Senf")
	chk02 := adagui.NewCheckbox("Mayo")
	chk03 := adagui.NewCheckbox("Ketchup")
	chk04 := adagui.NewCheckbox("Knoblauchsauce")
	grpCheck.Add(chk01, chk02, chk03, chk04)

	grpRadio := adagui.NewGroupPL(grpOptions, adagui.NewVBoxLayout())
	sizeVar := binding.NewInt()
	sizeVar.Set(1)
	rad01 := adagui.NewRadioButtonWithData("Klein", 1, sizeVar)
	rad02 := adagui.NewRadioButtonWithData("Mittel", 2, sizeVar)
	rad03 := adagui.NewRadioButtonWithData("Gross", 3, sizeVar)
	grpRadio.Add(rad01, rad02, rad03)

	grpIcon := adagui.NewGroupPL(grpMain, adagui.NewHBoxLayout())
	iconData := binding.NewInt()
	numIcons := 12
	iconList = make([]*adagui.IconButton, numIcons)
	for i := range numIcons {
		fileName := fmt.Sprintf("icons/%d.png", i+1)
		iconList[i] = adagui.NewIconButtonWithData(fileName, i, iconData)
		grpIcon.Add(iconList[i])
	}

	grpSlider := adagui.NewGroupPL(grpMain, adagui.NewHBoxLayout())
	val := binding.NewFloat()
	str := binding.FloatToStringWithFormat(val, "%.1f")
	sld := adagui.NewSliderWithData(200, adagui.Horizontal, val)
	sld.SetRange(0.0, 1.0, 0.2)
	lbl := adagui.NewLabelWithData(str)
	grpSlider.Add(sld, lbl)

	grpSlider = adagui.NewGroupPL(grpMain, adagui.NewHBoxLayout())
	val = binding.NewFloat()
	str = binding.FloatToStringWithFormat(val, "%.3f")
	sld = adagui.NewSliderWithData(200, adagui.Horizontal, val)
	sld.SetRange(0.0, 2*math.Pi, math.Pi/36.0)
	lbl = adagui.NewLabelWithData(str)
	grpSlider.Add(sld, lbl)

	grpLstBtn := adagui.NewGroupPL(grpMain, adagui.NewHBoxLayout())
	lst01 := []string{"Stefan", "Michael", "Luzia"}
	lst02 := []string{"A", "B", "C", "D", "E", "F"}
	lstBtn01 := adagui.NewListButton(lst01)
	lstBtn02 := adagui.NewListButton(lst02)
	grpLstBtn.Add(lstBtn01, lstBtn02)

	grpMain.Add(adagui.NewSpacer())

	grpBtn := adagui.NewGroupPL(grpMain, adagui.NewHBoxLayout())
	txtBtn01 := adagui.NewTextButton("Open")
	txtBtn02 := adagui.NewTextButton("Execute")
	txtBtn03 := adagui.NewTextButton("Quit")
	txtBtn02.SetFont(fonts.LucidaHandwritingItalic)
	txtBtn02.SetColor(colors.Purple.Dark(0.1))
	txtBtn02.SetPushedColor(colors.Purple.Bright(0.8))
	txtBtn02.SetTextColor(colors.Gold)
	txtBtn02.SetPushedTextColor(colors.Gold.Dark(0.8))
	txtBtn02.SetPushedBorderWidth(5.0)
	txtBtn02.SetPushedBorderColor(colors.Gold)
	grpBtn.Add(txtBtn01, txtBtn02, adagui.NewSpacer(), txtBtn03)

	return grpMain
}

func WidgetPanel02() adagui.Node {
	var iconList []*adagui.IconButton

	grpMain := adagui.NewGroup()
	grpMain.Layout = adagui.NewVBoxLayout()

	grpLabel := adagui.NewGroupPL(grpMain, adagui.NewHBoxLayout())
	label01 := adagui.NewLabel("Hallo zusammen")
	label02 := adagui.NewLabel("Grösser")
	label02.SetFontSize(28.0)
	label03 := adagui.NewLabel("Andere Schrift")
	label03.SetFont(fonts.LucidaHandwritingItalic)
	label03.SetFontSize(22.0)
	grpLabel.Add(label03, label02, label01)

	grpBtn := adagui.NewGroupPL(grpMain, adagui.NewHBoxLayout())
	button01 := adagui.NewButton(16, 16)
	button02 := adagui.NewButton(32, 32)
	button03 := adagui.NewButton(48, 48)
	grpBtn.Add(button03, button02, button01)

	grpBtn = adagui.NewGroupPL(grpMain, adagui.NewHBoxLayout())
	lst01 := []string{"Stefan", "Michael", "Luzia"}
	lst02 := []string{"A", "B", "C", "D", "E", "F"}
	lstBtn01 := adagui.NewListButton(lst01)
	lstBtn02 := adagui.NewListButton(lst02)
	grpBtn.Add(lstBtn02, lstBtn01, adagui.NewSpacer())

	txtBtn01 := adagui.NewTextButton("Open")
	txtBtn02 := adagui.NewTextButton("Close")
	txtBtn03 := adagui.NewTextButton("Execute")
	txtBtn03.SetFont(fonts.LucidaHandwritingItalic)
	txtBtn03.SetColor(colors.Purple.Dark(0.1))
	txtBtn03.SetPushedColor(colors.Purple.Bright(0.8))
	txtBtn03.SetTextColor(colors.Gold)
	txtBtn03.SetPushedTextColor(colors.Gold.Dark(0.8))
	txtBtn03.SetPushedBorderWidth(5.0)
	txtBtn03.SetPushedBorderColor(colors.Gold)
	grpBtn.Add(txtBtn03, txtBtn02, txtBtn01)

	grpIcon := adagui.NewGroupPL(grpMain, adagui.NewHBoxLayout())
	iconData := binding.NewInt()
	numIcons := 12
	iconList = make([]*adagui.IconButton, numIcons)
	for i := range numIcons {
		fileName := fmt.Sprintf("icons/%d.png", i+1)
		iconList[i] = adagui.NewIconButtonWithData(fileName, i, iconData)
		grpIcon.Add(iconList[i])
	}

	grpOptions := adagui.NewGroupPL(grpMain, adagui.NewHBoxLayout())

	grpCheck := adagui.NewGroupPL(grpOptions, adagui.NewVBoxLayout())
	chk01 := adagui.NewCheckbox("Senf")
	chk02 := adagui.NewCheckbox("Mayo")
	chk03 := adagui.NewCheckbox("Ketchup")
	chk04 := adagui.NewCheckbox("Knoblauchsauce")
	grpCheck.Add(chk04, chk03, chk02, chk01)

	grpRadio := adagui.NewGroupPL(grpOptions, adagui.NewVBoxLayout())
	sizeVar := binding.NewInt()
	sizeVar.Set(1)
	rad01 := adagui.NewRadioButtonWithData("Klein", 1, sizeVar)
	rad02 := adagui.NewRadioButtonWithData("Mittel", 2, sizeVar)
	rad03 := adagui.NewRadioButtonWithData("Gross", 3, sizeVar)
	grpRadio.Add(rad03, rad02, rad01)

	grpOptions.Add(adagui.NewSpacer())
	grpSlider := adagui.NewGroupPL(grpOptions, adagui.NewVBoxLayout())
	sld01 := adagui.NewSlider(220, adagui.Horizontal)
	sld01.SetRange(0.0, 1.0, 0.01)
	sld02 := adagui.NewSlider(220, adagui.Horizontal)
	sld02.SetRange(0.0, 1.0, 0.1)
	scr01 := adagui.NewScrollbar(220, adagui.Horizontal)
	scr01.SetVisiRange(0.45)
	scr02 := adagui.NewScrollbar(220, adagui.Horizontal)
	scr02.SetVisiRange(0.25)
	grpSlider.Add(scr01, scr02, sld01, sld02)

	return grpMain
}

func NestedTransformations() adagui.Node {
	var root *adagui.Group
	var panel01, panel02, panel03 *adagui.Panel
	var color02, color03 colors.Color
	var rotPt1, rotPt2 geom.Point
	var colorFactor float64 = 0.5
	var size geom.Point = geom.Point{float64(adatft.Width),
		float64(adatft.Height - 30)}

	root = adagui.NewGroup()
	root.Layout = adagui.NewPaddedLayout(0)
	root.SetMinSize(size)

	panel01 = NewPanel(0, 0)
	panel01.SetColor(colors.RandColor().Dark(0.8))
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
	rotSld1 := adagui.NewSliderWithData(len-lbl1.Size().X-intraGap,
		adagui.Horizontal, rotVal1)
	rotSld1.SetPos(refPt.AddXY(lbl1.Size().X+intraGap, 0))
	rotSld1.SetRange(-math.Pi/3, math.Pi/3, math.Pi/72.0)
	rotSld1.SetInitValue(0.0)

	refPt = refPt.AddXY(len+interGap, 0)
	lbl2 := adagui.NewLabel("Scale:")
	lbl2.SetFont(fonts.GoBold)
	lbl2.SetPos(refPt)
	lbl2.SetTextColor(lbl2.TextColor().Alpha(0.7))
	scaleSld1 := adagui.NewSliderWithData(len-lbl2.Size().X-intraGap,
		adagui.Horizontal, scaleVal1)
	scaleSld1.SetPos(refPt.AddXY(lbl2.Size().X+intraGap, 0))
	scaleSld1.SetRange(0.2, 1.8, 0.05)
	scaleSld1.SetInitValue(1.0)

	hSpc := 30.0
	vSpc := 10.0
	w, h := panel01.Size().AsCoord()
	w, h = w-2*hSpc, h-2*vSpc-scaleSld1.Size().Y

	color02 = colors.RandColor().Dark(colorFactor)
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
	obj.SetPos(panel02.Size().SubXY(30, 30))
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

	len = (panel02.Size().X - interGap - 6) / 2

	refPt = geom.Point{3, 3}
	lbl1 = adagui.NewLabel("Rotate:")
	lbl1.SetFont(fonts.GoBold)
	lbl1.SetPos(refPt)
	lbl1.SetTextColor(lbl1.TextColor().Alpha(0.7))
	rotSld2 := adagui.NewSliderWithData(len-lbl1.Size().X-intraGap,
		adagui.Horizontal, rotVal2)
	rotSld2.SetPos(refPt.AddXY(lbl1.Size().X+intraGap, 0))
	rotSld2.SetRange(-math.Pi/3, math.Pi/3, math.Pi/72.0)
	rotSld2.SetInitValue(0.0)

	refPt = refPt.AddXY(len+interGap, 0)
	lbl2 = adagui.NewLabel("Scale:")
	lbl2.SetFont(fonts.GoBold)
	lbl2.SetPos(refPt)
	lbl2.SetTextColor(lbl2.TextColor().Alpha(0.7))
	scaleSld2 := adagui.NewSliderWithData(len-lbl2.Size().X-intraGap,
		adagui.Horizontal, scaleVal2)
	scaleSld2.SetPos(refPt.AddXY(lbl2.Size().X+intraGap, 0))
	scaleSld2.SetRange(0.2, 1.8, 0.05)
	scaleSld2.SetInitValue(1.0)

	grp1.Add(rotSld2, scaleSld2, lbl1, lbl2)

	hSpc = 5.0
	vSpc = 30.0
	w, h = panel02.Size().AsCoord()
	w, h = w-hSpc-5, h-vSpc-25.0-5

	color03 = colors.RandColor().Dark(colorFactor)
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
	p.IsClipping = true

	//p.SetOnPress(func(evt touch.Event) {
	// log.Printf("Press on Panel: %v", evt)
	//})
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
	col := colors.RandColor()

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

func widgetFlex() adagui.Node {
	var main, cont *adagui.Group
	var menu *adagui.TabMenu

	main = adagui.NewGroup()
	cont = adagui.NewGroup()
	menu = adagui.NewTabMenu(cont)
	main.Layout = adagui.NewBorderLayout(menu, nil, nil, nil)
	cont.Layout = adagui.NewMaxLayout()
	menu.Layout = adagui.NewHBoxLayout()
	main.Add(menu, cont)

	menu.AddTab("Widgets 1", WidgetPanel01())
	menu.AddTab("Widgets 2", WidgetPanel02())
	menu.AddTab("Fonts", ScrolledFontPanel())
	menu.AddTab("Colors", ScrolledColorPanel())
	menu.AddTab("Draw", NestedTransformations())
	menu.SetTab(0)

	return main
}
