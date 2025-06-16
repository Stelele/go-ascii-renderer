package main

import (
	"fmt"
	"image"
	"image/color"
	"image/draw"
	_ "image/jpeg"
	"image/png"
	_ "image/png"
	"log"
	"os"

	"github.com/nfnt/resize"
	"golang.org/x/image/font"
	"golang.org/x/image/font/basicfont"
	"golang.org/x/image/math/fixed"
)

func getImage(loc string, w int, h int) ([]float64, error) {
	output := make([]float64, w*h)

	file, err := os.Open(loc)
	if err != nil {
		return nil, fmt.Errorf("Could not load image ", err)
	}
	defer file.Close()

	img, _, err := image.Decode(file)
	if err != nil {
		return nil, fmt.Errorf("Could not decode image", err)
	}

	img = resize.Resize(uint(w), uint(h), img, resize.Lanczos2)
	for y := range h {
		for x := range w {
			c := img.At(x, y)
			r, g, b, _ := c.RGBA()
			ra := float64(r) / float64(0xffff)
			rg := float64(g) / float64(0xffff)
			rb := float64(b) / float64(0xffff)

			output[y*w+x] = 0.2126*ra + 0.7152*rg + 0.0722*rb
		}
	}

	return output, nil
}

func printImage(imgLoc string, width int, height int) {
	img, err := getImage(imgLoc, width, height)
	if err != nil {
		log.Fatal("Failed to get image: ", err)
	}

	output := ""
	for y := range height {
		for x := range width {
			output += getBrightnessChar(img[y*width+x])
		}
		output += "\n"
	}

	fmt.Println(output)
}

func saveImage(imgLoc string, outImgLoc string, width int, height int) {
	img, err := getImage(imgLoc, width, height)
	if err != nil {
		log.Fatal("Failed to get image: ", err)
	}

	output := []string{}
	for y := range height {
		line := ""
		for x := range width {
			line += getBrightnessChar(img[y*width+x])
		}
		output = append(output, line)
	}

	charWidth, charHeight := 8, 13
	imgWidth := len(output[0]) * charWidth
	imgHeight := len(output) * charHeight

	outImg := image.NewRGBA(image.Rect(0, 0, imgWidth, imgHeight))

	draw.Draw(outImg, outImg.Bounds(), &image.Uniform{color.Black}, image.Point{}, draw.Src)

	face := basicfont.Face7x13
	d := font.Drawer{
		Dst:  outImg,
		Src:  image.NewUniform(color.White),
		Face: face,
	}

	for y, line := range output {
		d.Dot = fixed.P(0, (y+1)*charHeight)
		d.DrawString(line)
	}

	outFile, err := os.Create(outImgLoc)
	if err != nil {
		log.Fatal("Failed to create image file", err)
	}
	defer outFile.Close()

	err = png.Encode(outFile, outImg)
	if err != nil {
		log.Fatal("Failed to encode png image file", err)
	}

	fmt.Printf("Saved image to %s\n", outImgLoc)
}
