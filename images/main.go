package main

import (
	"fmt"
	"log"
)

func main() {
	width, height, err := getTerminalSize()
	if err != nil {
		log.Fatal("Failed to fetch terminal dimensions", err)
		return
	}

	clearScreen()
	hideCursor()

	imgLoc := "photos/gophers.jpeg"
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

func getBrightnessChar(val float64) string {
	chars := []string{" ", ".", ",", "-", "~", ":", ";", "=", "!", "*", "#", "$", "@"}
	div := 1. / float64(len(chars))

	index := int(val / div)

	if index < 1 {
		return chars[0]
	}
	if index >= len(chars)-1 {
		return chars[len(chars)-1]
	}

	return chars[index]
}

func clearScreen() {
	fmt.Print("\x1b[2J\x1b[H")
}

func resetCursor() {
	fmt.Print("\x1b[H")
}

func hideCursor() {
	fmt.Print("\033[?25l")
}
