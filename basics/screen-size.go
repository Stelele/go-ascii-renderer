package main

import (
	"fmt"
	"os"

	"golang.org/x/term"
)

func getTerminalSize() (width int, height int, err error) {
	width, height, err = term.GetSize(int(os.Stdout.Fd()))
	if err != nil {
		fmt.Println("Error getting terminal size:", err)
		return
	}

	return
}
