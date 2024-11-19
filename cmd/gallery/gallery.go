package main

// Dieses Programm dient der Demonstration der vorhandenen Widgets aus dem
// adagui-Package.
//
import (
	"flag"
	"fmt"
	"github.com/stefan-muehlebach/adagui"
	"github.com/stefan-muehlebach/adagui/binding"
	"github.com/stefan-muehlebach/adagui/touch"
	"github.com/stefan-muehlebach/adatft"
	"log"
	"os"
	"os/signal"
	//	"github.com/stefan-muehlebach/adagui/props"
	"github.com/stefan-muehlebach/gg/color"
	"github.com/stefan-muehlebach/gg/fonts"
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
	/*
	   p01 := props.ReadJSON("DefProps.json")
	   p02 := props.ReadJSON("ButtonProps.json")
	   fmt.Printf("%v\n", p01)
	   fmt.Printf("%v\n", p02)
	*/

func widgetFlex() adagui.Node {
	var iconList []*adagui.IconButton

	root := adagui.NewGroup()
	root.Layout = adagui.NewPaddedLayout()

	grpMain := adagui.NewGroupPL(root, adagui.NewVBoxLayout())

	grpLabel := adagui.NewGroupPL(grpMain, adagui.NewHBoxLayout())
	label01 := adagui.NewLabel("Hallo zusammen")
	label02 := adagui.NewLabel("Grösser")
	label02.SetFontSize(18.0)
	label03 := adagui.NewLabel("Andere Schrift")
	label03.SetFont(fonts.LucidaBright)
	grpLabel.Add(label01, label02, label03)

	grpBtn := adagui.NewGroupPL(grpMain, adagui.NewHBoxLayout())
	button01 := adagui.NewButton(16, 16)
	button02 := adagui.NewButton(32, 32)
	button03 := adagui.NewButton(48, 48)
	grpBtn.Add(button01, button02, button03)

	grpBtn = adagui.NewGroupPL(grpMain, adagui.NewHBoxLayout())
	txtBtn01 := adagui.NewTextButton("Open")
	txtBtn02 := adagui.NewTextButton("Close")
	txtBtn03 := adagui.NewTextButton("Execute")
	txtBtn03.SetFont(fonts.LucidaHandwritingItalic)
	txtBtn03.SetColor(color.Purple.Dark(0.1))
	txtBtn03.SetPushedColor(color.Purple.Bright(0.8))
	txtBtn03.SetTextColor(color.Gold)
	txtBtn03.SetPushedTextColor(color.Gold.Dark(0.8))
	grpBtn.Add(txtBtn01, txtBtn02, txtBtn03)

	grpBtn = adagui.NewGroupPL(grpMain, adagui.NewHBoxLayout())
	lst01 := []string{"Stefan", "Michael", "Luzia"}
	lst02 := []string{"A", "B", "C", "D", "E", "F"}
	lstBtn01 := adagui.NewListButton(lst01)
	lstBtn02 := adagui.NewListButton(lst02)
	grpBtn.Add(lstBtn01, lstBtn02)

	grpIcon := adagui.NewGroupPL(grpMain, adagui.NewHBoxLayout())
	iconData := binding.NewInt()
	numIcons := 13
	iconList = make([]*adagui.IconButton, numIcons)
	for i := range numIcons {
		fileName := fmt.Sprintf("icons/%d.png", i+1)
		iconList[i] = adagui.NewIconButtonWithData(fileName, i, iconData)
		grpIcon.Add(iconList[i])
	}

	grpCheck := adagui.NewGroupPL(grpMain, adagui.NewHBoxLayout())
	chk01 := adagui.NewCheckbox("Senf")
	chk02 := adagui.NewCheckbox("Mayo")
	chk03 := adagui.NewCheckbox("Ketchup")
	grpCheck.Add(chk01, chk02, chk03)

	grpRadio := adagui.NewGroupPL(grpMain, adagui.NewHBoxLayout())
	sizeVar := binding.NewInt()
	sizeVar.Set(1)
	rad01 := adagui.NewRadioButtonWithData("Klein", 1, sizeVar)
	rad02 := adagui.NewRadioButtonWithData("Mittel", 2, sizeVar)
	rad03 := adagui.NewRadioButtonWithData("Gross", 3, sizeVar)
	grpRadio.Add(rad01, rad02, rad03)

	grpMain.Add(adagui.NewSpacer())

	grpBtn = adagui.NewGroupPL(grpMain, adagui.NewHBoxLayout())
	btnQuit := adagui.NewTextButton("Quit")
	btnQuit.SetOnTap(func(evt touch.Event) {
		screen.Quit()
	})
	grpBtn.Add(adagui.NewSpacer(), btnQuit)

	return root
}
