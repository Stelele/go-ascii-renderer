package main

import (
	"fmt"
	"log"
	"time"
)

func main() {
	clearScreen()

	for i := range 5 {
		width, height, err := getTerminalSize()
		if err != nil {
			log.Fatal("Failed to fetch terminal dimensions", err)
			return
		}

		hideCursor()
		resetCursor()

		imgLoc := fmt.Sprintf("photos/gophers-%d.jpeg", i)
		outImgloc := fmt.Sprintf("photos/out/gophers-%d.png", i)

		printImage(imgLoc, width, height)
		saveImage(imgLoc, outImgloc, 240, 135)

		time.Sleep(time.Second)
	}
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
