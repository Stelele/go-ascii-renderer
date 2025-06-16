package main

import (
	"fmt"
	"image"
	_ "image/jpeg"
	_ "image/png"
	"os"

	"github.com/nfnt/resize"
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
