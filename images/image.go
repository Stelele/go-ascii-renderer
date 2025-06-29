package main

import (
	"fmt"
	"image"
	"image/color"
	"image/draw"
	_ "image/jpeg"
	"image/png"
	"log"
	"os"

	"github.com/disintegration/imaging"
	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"
	"golang.org/x/image/math/fixed"
)

type Grad struct {
	Dx    float64
	Dy    float64
	Mag   float64
	Theta float64
}

func getImage(loc string, w int, h int) ([]float64, error) {
	img, err := getRawImage(loc, w, h)
	if err != nil {
		return nil, fmt.Errorf("could not get raw image: %v", err)
	}

	output := getLuminence(img, w, h)
	return output, nil
}

func getRawImage(loc string, w int, h int) (*image.NRGBA, error) {
	file, err := imaging.Open(loc)
	if err != nil {
		return nil, fmt.Errorf("could not load image: %v", err)
	}

	img := imaging.Resize(file, w, h, imaging.Lanczos)

	return img, nil
}

func getLuminence(img *image.NRGBA, w int, h int) []float64 {
	output := make([]float64, w*h)
	for y := range h {
		for x := range w {
			output[y*w+x] = getSingleLuminence(img.At(x, y))
		}
	}

	return output
}

func getSingleLuminence(c color.Color) float64 {
	r, g, b, _ := c.RGBA()
	nr := float64(r) / float64(0xffff)
	ng := float64(g) / float64(0xffff)
	nb := float64(b) / float64(0xffff)

	return 0.2126*nr + 0.7152*ng + 0.0722*nb
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

	face, err := getFont()
	if err != nil {
		log.Fatal("Failed to get font: ", err)
	}

	advance, _ := face.GlyphAdvance(' ')
	charWidth := advance.Ceil()
	charHeight := (face.Metrics().Ascent + face.Metrics().Descent).Ceil()
	imgWidth := width * charWidth
	imgHeight := height * charHeight

	outImg := image.NewRGBA(image.Rect(0, 0, imgWidth, imgHeight))

	draw.Draw(outImg, outImg.Bounds(), &image.Uniform{color.Black}, image.Point{}, draw.Src)

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

func getFont() (font.Face, error) {
	fontName := "photos/Ac437_IBM_EGA_8x8.ttf"

	fontData, err := os.ReadFile(fontName)
	if err != nil {
		return nil, err
	}

	ft, err := opentype.Parse(fontData)
	if err != nil {
		return nil, err
	}

	face, err := opentype.NewFace(ft, &opentype.FaceOptions{
		Size: 8,
		DPI:  72,
	})

	return face, err
}
