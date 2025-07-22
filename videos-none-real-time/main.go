package main

import (
	"fmt"
	"io/fs"
	"log"
	"os"
	"os/exec"
)

func main() {

	vid := "video1.mp4"

	tmpRawPath := "tmp-images/raw/"
	tmpProcessedPath := "tmp-images/processed/"
	outputVideosPath := "output-videos/"

	checkIfFolderExistsAndCreate(tmpRawPath)
	checkIfFolderExistsAndCreate(tmpProcessedPath)
	checkIfFolderExistsAndCreate(outputVideosPath)

	cmd := exec.Command("ffmpeg", "-i", "input-videos/"+vid, tmpRawPath+vid+"-%04d.png")

	_, err := cmd.CombinedOutput()
	if err != nil {
		panic(err)
	}

	println("Processing images...")
	i := 1
	for {
		imgName := vid + "-" + fmt.Sprintf("%04d", i) + ".png"
		imgPath := tmpRawPath + imgName
		outImgPath := tmpProcessedPath + imgName

		if _, err := os.Stat(imgPath); os.IsNotExist(err) {
			break
		}

		saveImage(imgPath, outImgPath, 80, 45)

		i += 1
	}
	println("Images processed.")

	println("Creating video from processed images...")
	cmd = exec.Command("ffmpeg", "-framerate", "30", "-i", tmpProcessedPath+vid+"-%04d.png", "-c:v", "libx264", "-pix_fmt", "yuv420p", outputVideosPath+vid)
	_, err = cmd.CombinedOutput()
	if err != nil {
		panic(err)
	}
	println("Video created successfully at", outputVideosPath+vid)
}

func checkIfFolderExistsAndCreate(path string) {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		err := os.MkdirAll(path, fs.ModeDir)
		if err != nil {
			log.Fatal("Failed to create directory", err)
		}
	}
}
