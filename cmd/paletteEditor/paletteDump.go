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
	"github.com/stefan-muehlebach/gg/colors"
	"github.com/stefan-muehlebach/gg/fonts"
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
	palNames, palMap, err := colors.ReadPaletteFile("palette.json")
	if err != nil {
		log.Fatalf("couldn't read palette names: %v", err)
	}
	numPals := len(palNames)
	numRows := numPals / numColumns
	if numPals%numColumns != 0 {
		numRows += 1
	}
	imgWidth := numColumns * stripeWidth
	imgHeight := numRows * stripeHeight
	img := image.NewRGBA(image.Rect(0, 0, imgWidth, imgHeight))
	gc := gg.NewContextForRGBA(img)

	gc.SetFillColor(colors.WhiteSmoke)
	gc.Clear()
	face, _ := fonts.NewFace(fonts.GoRegular, 18.0)
	gc.SetFontFace(face)
	gc.SetStrokeColor(colors.Black)
	for i, palName := range palNames {
		col := i / numRows
		row := i % numRows
		x0 := float64(col * stripeWidth)
		y0 := float64(row * stripeHeight)

		gc.SetStrokeWidth(1.0)
		gc.DrawRectangle(x0, y0, stripeWidth, stripeHeight)
		gc.FillStroke()

		fmt.Printf("  [%2d]: %s\n", i, palName)
		pal := palMap[palName]
		//pal.SetLength(colorBarWidth)
		//pal.LenIsMaxIter()
		//pal.SetOffset(0.0)

		for x := 0; x < colorBarWidth; x++ {
			t := float64(x)/float64(colorBarWidth-1)
			color := pal.Color(t)
			for y := 0; y < colorBarHeight; y++ {
				img.Set(int(x0+padding)+x, int(y0+textHeight+padding)+y, color)
			}
		}
		gc.SetStrokeColor(colors.DarkSlateGrey)
		gc.SetStrokeWidth(3.0)
		gc.DrawRectangle(x0+padding, y0+textHeight+padding, colorBarWidth, colorBarHeight)
		gc.Stroke()

		gc.DrawStringAnchored(palName, x0+padding, y0+padding+textHeight-padding, 0.0, 0.0)
		gc.DrawStringAnchored(pal.Type().String(), x0+stripeWidth-padding,
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
