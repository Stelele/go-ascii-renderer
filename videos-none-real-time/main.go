package main

import (
	"fmt"
	"io/fs"
	"log"
	"os"
	"os/exec"
	"sync"
)

func main() {

	vid := "video3.mp4"

	tmpRawPath := "tmp-images/raw/"
	tmpProcessedPath := "tmp-images/processed/"
	outputVideosPath := "output-videos/"

	println("Setting up directories...")
	checkIfFolderExistsAndCreate(tmpRawPath)
	checkIfFolderExistsAndCreate(tmpProcessedPath)
	checkIfFolderExistsAndCreate(outputVideosPath)

	if _, err := os.Stat(outputVideosPath + vid); err == nil {
		os.Remove(outputVideosPath + vid)
	}

	cmd := exec.Command("ffmpeg", "-i", "input-videos/"+vid, tmpRawPath+vid+"-%04d.png")

	_, err := cmd.CombinedOutput()
	if err != nil {
		panic(err)
	}

	println("Processing images...")
	i := 1
	r := 16
	w := 16 * r
	h := 9 * r

	var wg sync.WaitGroup
	for {
		imgName := vid + "-" + fmt.Sprintf("%04d", i) + ".png"
		imgPath := tmpRawPath + imgName
		outImgPath := tmpProcessedPath + imgName

		if _, err := os.Stat(imgPath); os.IsNotExist(err) {
			break
		}

		wg.Add(1)
		go func(imgPath, outImgPath string, w, h int) {
			defer wg.Done()
			saveImage(imgPath, outImgPath, w, h)

		}(imgPath, outImgPath, w, h)

		i += 1
	}

	wg.Wait()
	println("Images processed.")

	println("Creating video from processed images...")
	cmd = exec.Command("ffmpeg", "-framerate", "30", "-i", tmpProcessedPath+vid+"-%04d.png", "-c:v", "libx264", "-pix_fmt", "yuv420p", outputVideosPath+vid)
	_, err = cmd.CombinedOutput()
	if err != nil {
		panic(err)
	}
	println("Video created successfully at", outputVideosPath+vid)

	println("Cleaning up temporary files...")
	err = os.RemoveAll(tmpRawPath)
	if err != nil {
		log.Fatal("Failed to remove temporary raw images directory:", err)
	}
	err = os.RemoveAll(tmpProcessedPath)
	if err != nil {
		log.Fatal("Failed to remove temporary processed images directory:", err)
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
