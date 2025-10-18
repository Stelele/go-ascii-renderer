# go-ascii-renderer

A small collection of Go programs that render ASCII-art from shapes, images, and videos.

This repository contains three example programs demonstrating different ways to produce ASCII art:

- `basics/` — a terminal-based animated ASCII renderer that draws procedurally generated shapes (donut) in the terminal.
- `images/` — converts static images (JPEG/PNG) into ASCII art printed to the terminal and saves a processed image.
- `videos-none-real-time/` — converts a video file into a sequence of frames, processes each frame into ASCII-style images, then re-encodes them into a video using `ffmpeg`.

## Prerequisites

- Go (1.20+ recommended) installed and available in your PATH.
- ffmpeg installed and available in your PATH (required only for the `videos-none-real-time` program).

On Windows, you can download ffmpeg from https://ffmpeg.org/download.html and add it to your PATH.

## Quick build (project root)

Open a terminal in the repository root and run:

```powershell
go build ./...
```

This will build each module. Alternatively, you can run programs directly with `go run` from their directories.

## Module: basics

Location: `basics/`

What it does:
- Runs a simple animated ASCII renderer in the terminal. It reads your terminal size, hides the cursor, and prints frames continuously.

How to run:

```powershell
cd basics
go run .
```

Notes:
- The program uses ANSI escape sequences to clear the screen and hide the cursor. It expects a compatible terminal.
- The animation parameters (rotation speeds, frame timing) are implemented in `basics/main.go` and related files.

## Module: images

Location: `images/`

What it does:
- Loads a set of images from `images/photos/` named like `gophers-0.jpeg`, `gophers-1.jpeg`, etc., prints ASCII output to the terminal, and saves processed images to `images/photos/out/`.

How to run:

```powershell
cd images
go run .
```

What to provide:
- The program expects input images at `images/photos/gophers-<n>.jpeg`. You can replace or add images using those names.

Output:
- Terminal ASCII preview printed for each image.
- A processed image saved to `images/photos/out/gophers-<n>.png` for each input image.

## Module: videos-none-real-time

Location: `videos-none-real-time/`

What it does:
- Uses `ffmpeg` to extract frames from a video placed in `videos-none-real-time/input-videos/`.
- Processes each extracted frame into ASCII-style images (resized) in parallel.
- Uses `ffmpeg` to re-encode processed frames back into a video and writes it to `videos-none-real-time/output-videos/`.

How to run:

1. Put your input video in `videos-none-real-time/input-videos/` and name it (for example) `video3.mp4`. The example code uses `video3.mp4` by default.

2. From the module directory run:

```powershell
cd videos-none-real-time
go run .
```

What happens:
- Temporary directories `tmp-images/raw/` and `tmp-images/processed/` are created.
- `ffmpeg` is invoked to extract frames into `tmp-images/raw/`.
- Each frame is processed by the Go program and written to `tmp-images/processed/`.
- `ffmpeg` is invoked again to create a video from the processed frames and writes it to `output-videos/<original-video-name>`.
- Temporary folders are removed when done.

Notes and troubleshooting
- Ensure `ffmpeg` is installed and on your PATH. If `exec.Command("ffmpeg", ...)` fails, you'll see an error. Running `ffmpeg -version` in a terminal should print a version message.
- The default video filename is `video3.mp4` — modify `videos-none-real-time/main.go` if you want to change it.
- Processing may be CPU and disk intensive for large videos; the code processes frames in parallel using goroutines and sync.WaitGroup.

Development and testing
- You can run package-level builds and tests with `go build ./...` and `go test ./...` (if tests are present).
- To experiment, modify parameters such as frame size, framerate, or input file names directly in the module's `main.go`.

Enjoy exploring ASCII rendering with Go!
