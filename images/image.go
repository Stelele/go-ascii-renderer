package main

import (
	"fmt"
	"image"
	_ "image/jpeg"
	_ "image/png"
	"os"
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

	bounds := img.Bounds()
	width, height := bounds.Max.X, bounds.Max.Y

	cellX := width / w
	cellY := height / h

	idx := 0
	for y := 0; y < height; y += cellY {
		for x := 0; x < width; x += cellX {
			cW, cH := x+cellX, y+cellY
			if cW >= width {
				cW = width
			}
			if cH >= height {
				cH = height
			}

			sum := 0.
			for iy := y; iy < cH; iy++ {
				for ix := x; ix < cW; ix++ {
					c := img.At(ix, iy)
					r, g, b, _ := c.RGBA()
					ra := float64(r) / float64(0xffff)
					rg := float64(g) / float64(0xffff)
					rb := float64(b) / float64(0xffff)

					luminance := 0.2126*ra + 0.7152*rg + 0.0722*rb
					sum += luminance
				}
			}

			if idx < w*h {
				output[idx] = sum / float64(cellX*cellY)
				idx += 1
			}
		}
	}

	return output, nil
}
