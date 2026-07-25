package main

import (
	"flag"
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
//rotation = 90.0
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
	rotation                adatft.RotationType = adatft.Rotate270
	alpha                   float64
	fadeIn = time.Duration(3 * time.Second)
	hold   = time.Duration(3 * time.Second)
	fadeOut = time.Duration(3 * time.Second)
)

func main() {
	flag.Var(&rotation, "rotation", "rotation of the main display")
	flag.DurationVar(&fadeIn, "fadein", fadeIn, "Duration for fade in")
	flag.DurationVar(&hold, "hold", hold, "Duration for hold")
	flag.DurationVar(&fadeOut, "fadeout", fadeOut, "Duration for fade out")
	flag.Parse()

	if len(flag.Args()) <= 0 {
		fmt.Printf("usage: %s [-rotation <rot>] <file> [...]\n", os.Args[0])
		os.Exit(1)
	}

	alpha = float64(rotation) / 180.0 * math.Pi
	disp = adatft.OpenDisplay(rotation)
	pngImg = image.NewRGBA(image.Rect(0, 0, adatft.Width, adatft.Height))
	tftImg = image.NewRGBA(image.Rect(0, 0, adatft.Width, adatft.Height))
	backColor = image.NewUniform(color.Black)
	alphaMask = image.NewUniform(color.Alpha{128})
	drawOpts = draw.Options{
		alphaMask, image.Point{},
		alphaMask, image.Point{},
	}

	for _, imgFile = range flag.Args() {
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

		fadeInSteps := int(fadeIn / (70 * time.Millisecond))
		fadeOutSteps := int(fadeOut / (70 * time.Millisecond))
		dAlpha := 256 / fadeInSteps

		//t0 := time.Now()
		for alpha := 0; alpha < 256; alpha += dAlpha {
			alphaMask.C = color.Alpha{uint8(alpha)}
			draw.Copy(tftImg, image.Point{}, backColor, tftImg.Bounds(),
				draw.Src, nil)
			draw.Copy(tftImg, image.Point{}, pngImg, tftImg.Bounds(),
				draw.Over, &drawOpts)
			disp.Draw(tftImg)
			time.Sleep(20 * time.Millisecond)
		}
		//d0 := time.Since(t0)
		//log.Printf("fadeIn in %v, using %d steps\n", d0, fadeInSteps)

		//t0 = time.Now()
		time.Sleep(hold)
		//d0 = time.Since(t0)
		//log.Printf("holding for %v\n", d0)

		dAlpha = 256 / fadeOutSteps

		//t0 = time.Now()
		for alpha := 255; alpha >= 0; alpha -= dAlpha {
			alphaMask.C = color.Alpha{uint8(alpha)}
			draw.Copy(tftImg, image.Point{}, backColor, tftImg.Bounds(),
				draw.Src, nil)
			draw.Copy(tftImg, image.Point{}, pngImg, tftImg.Bounds(),
				draw.Over, &drawOpts)
			disp.Draw(tftImg)
			time.Sleep(20 * time.Millisecond)
		}
		//d0 = time.Since(t0)
		//log.Printf("fadeOut in %v, using %d steps\n", d0, fadeOutSteps)
	}
	disp.Close()
}
