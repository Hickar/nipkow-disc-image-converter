package main

import (
	"bytes"
	"fmt"
	"io/fs"
	"log"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strconv"
)

var (
	inputFile        string  = "./input/mickey.mp4"
	outputDir        string  = "./output"
	outputFileFormat string  = "out%d.bmp"
	tmpFile          string  = "./tmp/min.mp4"
	scaleWidth       int     = 64
	scaleHeight      int     = 32
	scaleRatio       float32 = 1 / 2
	fps              int     = 15
)

func execute(name string, arg ...string) error {
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	cmd := exec.Command(name, arg...)
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	if err := cmd.Start(); err != nil {
		log.Printf("stdErr: %s", stderr.String())
		return err
	}
	if err := cmd.Wait(); err != nil {
		log.Printf("stdErr: %s", stderr.String())
		return err
	}

	fmt.Printf("Stdout: %s", stdout.String())
	return nil
}

func minify(input, output, ratio string) error {
	return execute("ffmpeg", "-i", input, "-vf", ratio, output)
}

func framify(input, output string, fps int) error {
	return execute("ffmpeg", "-i", input, "-vf", fmt.Sprintf("fps=%d", fps), output)
}

func grayscale(input, output string) error {
	return execute("ffmpeg", "-i", input, "-vf", "format=gray", output)
}

func main() {
	ratio := fmt.Sprintf("scale=%d:%d,setdar=ratio=%.5f", scaleWidth, scaleHeight, scaleRatio)
	err := minify(inputFile, tmpFile, ratio)
	if err != nil {
		log.Fatal(err)
	}

	err = framify(tmpFile, path.Join(outputDir, outputFileFormat), fps)
	if err != nil {
		log.Fatal(err)
	}

	root, err := os.Getwd()
	fileSystem := os.DirFS(root)
	if err != nil {
		log.Fatal(err)
	}

	files, err := fs.ReadDir(fileSystem, outputDir)
	for i, file := range files {
		ext := filepath.Ext(file.Name())
		inputFullPath := filepath.Join(outputDir, file.Name())
		outputFullPath := filepath.Join(outputDir, strconv.Itoa(i)+ext)

		matched, _ := filepath.Match("*out*.bmp", inputFullPath)
		if matched {
			err := grayscale(inputFullPath, outputFullPath)
			if err != nil {
				log.Fatalf("Can't apply grayscale filter to %s:\n%s", inputFullPath, err)
			}
		}
	}
}
