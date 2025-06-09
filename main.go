package main

import (
	"fmt"
	"time"
)

var (
	f      int     = 0
	t      float64 = 0.
	width  int
	height int
	dt     float64 = 1 / 60.
)

func main() {
	var err error
	width, height, err = getTerminalSize()
	if err != nil {
		fmt.Println("Couldn't get terminal dimensions")
		return
	}

	clearScreen()
	for {
		hideCursor()
		resetCursor()
		preProcessFrame(A, B)
		update(donut)
		time.Sleep(16 * time.Millisecond)
		f += 1
		t += dt
		A += 0.03
		B += 0.01
	}
}

func update(up func(int, int) float64) {
	output := ""

	for y := range height {
		for x := range width {
			output += getBrightnessChar(up(x, y))
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
