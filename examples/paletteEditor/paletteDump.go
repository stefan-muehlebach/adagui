//go:build ignore

package main

import (
	_ "flag"
	"fmt"
	"image"
	"image/png"
	"log"
	"os"

	"github.com/stefan-muehlebach/gg"
	"github.com/stefan-muehlebach/gg/colornames"
	"github.com/stefan-muehlebach/gg/fonts"
	"github.com/stefan-muehlebach/mandel"
)

const (
	numColumns     = 2
	colorBarWidth  = 1024
	colorBarHeight = 100
	textHeight     = 30
	padding        = 10
	stripeHeight   = colorBarHeight + textHeight + 2*padding
	stripeWidth    = colorBarWidth + 2*padding
)

func main() {
	var palType string

	palNameList, err := mandel.PaletteNames()
	if err != nil {
		log.Fatalf("couldn't read palette names: %v", err)
	}
	numPals := len(palNameList)
	numRows := numPals / numColumns
	if numPals%numColumns != 0 {
		numRows += 1
	}
	imgWidth := numColumns * stripeWidth
	imgHeight := numRows * stripeHeight
	img := image.NewRGBA(image.Rect(0, 0, imgWidth, imgHeight))
	gc := gg.NewContextForRGBA(img)

	gc.SetFillColor(colornames.WhiteSmoke)
	gc.Clear()
	gc.SetFontFace(fonts.NewFace(fonts.GoRegular, 18.0))
	gc.SetStrokeColor(colornames.Black)
	for i, palName := range palNameList {
		col := i / numRows
		row := i % numRows
		x0 := float64(col * stripeWidth)
		y0 := float64(row * stripeHeight)

		gc.SetStrokeWidth(1.0)
		gc.DrawRectangle(x0, y0, stripeWidth, stripeHeight)
		gc.FillStroke()

		fmt.Printf("  [%2d]: %s\n", i, palName)
		pal, err := mandel.NewPalette(palName)
		if err != nil {
			log.Fatalf("couldn't create palette: %v", err)
		}
		switch pal.(type) {
		case *mandel.GradientPalette:
			palType = "Gradient Palette"
		case *mandel.ProcPalette:
			palType = "Procedure Palette"
		}
		pal.SetLength(colorBarWidth)
		pal.LenIsMaxIter()
		pal.SetOffset(0.0)

		for x := 0; x < colorBarWidth; x++ {
			color := pal.GetColor(float64(x))
			for y := 0; y < colorBarHeight; y++ {
				img.Set(int(x0+padding)+x, int(y0+textHeight+padding)+y, color)
			}
		}
		gc.SetStrokeColor(colornames.DarkSlateGrey)
		gc.SetStrokeWidth(3.0)
		gc.DrawRectangle(x0+padding, y0+textHeight+padding, colorBarWidth, colorBarHeight)
		gc.Stroke()

		gc.DrawStringAnchored(palName, x0+padding, y0+padding+textHeight-padding, 0.0, 0.0)
		gc.DrawStringAnchored(palType, x0+stripeWidth-padding,
			y0+padding+textHeight-padding, 1.0, 0.0)
	}
	fileName := "palette.png"
	fh, err := os.Create(fileName)
	if err != nil {
		log.Fatalf("couldn't create file: %v", err)
	}
	png.Encode(fh, img)
	fh.Close()
}
