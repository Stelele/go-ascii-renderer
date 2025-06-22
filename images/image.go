package main

import (
	"fmt"
	"image"
	"image/color"
	"image/draw"
	_ "image/jpeg"
	"image/png"
	"log"
	"math"
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

func getImage(loc string, w int, h int) ([]float64, []Grad, error) {
	img, err := getRawImage(loc, w, h)
	if err != nil {
		return nil, nil, fmt.Errorf("Could not get raw image", err)
	}

	grads, err := getEdges(loc, w, h)
	if err != nil {
		return nil, nil, fmt.Errorf("Could not get image edges", err)
	}

	output := getLuminence(img, w, h)
	return output, grads, nil
}

func getRawImage(loc string, w int, h int) (*image.NRGBA, error) {
	file, err := imaging.Open(loc)
	if err != nil {
		return nil, fmt.Errorf("Could not load image ", err)
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
	img, grads, err := getImage(imgLoc, width, height)
	if err != nil {
		log.Fatal("Failed to get image: ", err)
	}

	output := ""
	for y := range height {
		for x := range width {
			output += getBrightnessChar(img[y*width+x], grads[y*width+x])
		}
		output += "\n"
	}

	fmt.Println(output)
}

func saveImage(imgLoc string, outImgLoc string, width int, height int) {
	img, grads, err := getImage(imgLoc, width, height)
	if err != nil {
		log.Fatal("Failed to get image: ", err)
	}

	output := []string{}
	for y := range height {
		line := ""
		for x := range width {
			line += getBrightnessChar(img[y*width+x], grads[y*width+x])
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

func getEdges(imgLoc string, w int, h int) ([]Grad, error) {
	file, err := imaging.Open(imgLoc)
	if err != nil {
		return nil, fmt.Errorf("Could not load image ", err)
	}

	d1 := imaging.Blur(file, 3)
	d2 := imaging.Blur(file, 10)

	img := imaging.Resize(d1, w, h, imaging.Lanczos)
	img2 := imaging.Resize(d2, w, h, imaging.Lanczos)

	bounds := img.Bounds()
	imgOut := image.NewNRGBA(bounds)
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			c1 := img.NRGBAAt(x, y)
			c2 := img2.NRGBAAt(x, y)

			r := clamp(int(c1.R) - int(c2.R))
			g := clamp(int(c1.G) - int(c2.G))
			b := clamp(int(c1.B) - int(c2.B))

			imgOut.Set(x, y, color.NRGBA{r, g, b, 255})
		}
	}

	grads := make([]Grad, w*h)

	bounds = imgOut.Bounds()
	gw, gh := bounds.Dx(), bounds.Dy()

	// Sobel kernels
	gx := [3][3]float64{
		{-1, 0, 1},
		{-2, 0, 2},
		{-1, 0, 1},
	}
	gy := [3][3]float64{
		{-1, -2, -1},
		{0, 0, 0},
		{1, 2, 1},
	}

	for y := 1; y < gh-1; y++ {
		for x := 1; x < gw-1; x++ {
			var dx, dy float64

			for ky := -1; ky <= 1; ky++ {
				for kx := -1; kx <= 1; kx++ {
					px := getSingleLuminence(imgOut.At(x+kx, y+ky)) * 255
					dx += gx[ky+1][kx+1] * px
					dy += gy[ky+1][kx+1] * px
				}
			}

			grads[y*gw+x] = Grad{
				dx,
				dy,
				math.Sqrt(dx*dx + dy*dy),
				math.Atan2(dy, dx),
			}
		}
	}

	return grads, nil
}

func clamp(v int) uint8 {
	if v < 0 {
		return 0
	}
	if v > 255 {
		return 255
	}
	return uint8(v)
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
