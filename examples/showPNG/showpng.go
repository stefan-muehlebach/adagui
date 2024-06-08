package main

import (
	"fmt"
	"github.com/stefan-muehlebach/adatft"
	"github.com/stefan-muehlebach/gg"
	"golang.org/x/image/draw"
	"golang.org/x/image/math/f64"
	"image"
	"image/color"
	"log"
	"math"
	"os"
	"time"
)

func check(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

//----------------------------------------------------------------------------

const (
	rotation = 90.0
	alpha    = rotation / 180.0 * math.Pi
)

var (
	disp                    *adatft.Display
	png                     image.Image
	backColor               *image.Uniform
	alphaMask               *image.Uniform
	pngImg, tftImg          *image.RGBA
	err                     error
	imgFile                 string
	sinAlpha                = math.Sin(alpha)
	cosAlpha                = math.Cos(alpha)
	t                       f64.Aff3
	drawOpts                draw.Options
	scale, offsetX, offsetY float64
	imgSize, pngSize        image.Point
	dstRect                 image.Rectangle
)

func main() {
	if len(os.Args) < 2 {
		fmt.Printf("usage: %s <file> [...]\n", os.Args[0])
		os.Exit(1)
	}

	disp = adatft.OpenDisplay(adatft.Rotate090)
	pngImg = image.NewRGBA(image.Rect(0, 0, adatft.Width, adatft.Height))
	tftImg = image.NewRGBA(image.Rect(0, 0, adatft.Width, adatft.Height))
	backColor = image.NewUniform(color.Black)
	alphaMask = image.NewUniform(color.Alpha{128})
	drawOpts = draw.Options{
		alphaMask, image.Point{},
		alphaMask, image.Point{},
	}

	for _, imgFile = range os.Args[1:] {
		png, err = gg.LoadPNG(imgFile)
		check(err)
		imgSize = pngImg.Bounds().Size()
		pngSize = png.Bounds().Size()

		scale = min(1.0, min(float64(imgSize.X)/float64(pngSize.X),
			float64(imgSize.Y)/float64(pngSize.Y)))

		offsetX = 0.5 * (float64(imgSize.X) - scale*float64(pngSize.X))
		offsetY = 0.5 * (float64(imgSize.Y) - scale*float64(pngSize.Y))

		t = f64.Aff3{
			scale, 0.0, offsetX,
			0.0, scale, offsetY,
		}
        draw.Copy(pngImg, image.Point{}, backColor, pngImg.Bounds(),
            draw.Src, nil)
		draw.BiLinear.Transform(pngImg, t, png, png.Bounds(), draw.Src, nil)

		for alpha := 0; alpha < 256; alpha += 4 {
			alphaMask.C = color.Alpha{uint8(alpha)}
			draw.Copy(tftImg, image.Point{}, backColor, tftImg.Bounds(),
				draw.Src, nil)
			draw.Copy(tftImg, image.Point{}, pngImg, tftImg.Bounds(),
				draw.Over, &drawOpts)
			disp.Draw(tftImg)
			time.Sleep(20 * time.Millisecond)
		}
		time.Sleep(5 * time.Second)
		for alpha := 255; alpha >= 0; alpha -= 4 {
			alphaMask.C = color.Alpha{uint8(alpha)}
			draw.Copy(tftImg, image.Point{}, backColor, tftImg.Bounds(),
				draw.Src, nil)
			draw.Copy(tftImg, image.Point{}, pngImg, tftImg.Bounds(),
				draw.Over, &drawOpts)
			disp.Draw(tftImg)
			time.Sleep(20 * time.Millisecond)
		}
	}
	disp.Close()
}
