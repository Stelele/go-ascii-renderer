package main

import (
	"fmt"
	"math"
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
		update(test)
		time.Sleep(16 * time.Millisecond)
		f += 1
		t += dt
	}
}

func test(x int, y int) string {
	u0 := 2*(float64(x)/float64(width)) - 1
	v0 := 2*(float64(y)/float64(height)) - 1

	u := 2*(u0-math.Floor(u0)) - 1
	v := 2*(v0-math.Floor(v0)) - 1

	d := math.Sqrt(u*u + v*v)

	d = math.Sin(d*8+t*0.8) / 8
	d = math.Abs(d)

	if d < 0.01 {
		return " "
	}
	if d < 0.05 {
		return "."
	}
	if d < 0.1 {
		return ","
	}

	if d < 0.4 {
		return "o"
	}

	return "#"
}

func update(up func(int, int) string) {
	output := ""

	for y := range height {
		for x := range width {
			output += up(x, y)
		}
		output += "\n"
	}

	fmt.Println(output)
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
