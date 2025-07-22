package main

import (
	"fmt"
	"io/fs"
	"log"
	"os"
	"time"
)

func main() {
	checkIfFolderExistsAndCreate("photos/out")
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

func checkIfFolderExistsAndCreate(path string) {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		err := os.MkdirAll(path, fs.ModeDir)
		if err != nil {
			log.Fatal("Failed to create directory", err)
		}
	}
}
