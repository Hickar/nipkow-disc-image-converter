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

	"golang.org/x/image/bmp"
)

var (
	inputFile        string  = "./input/mickey.mp4"
	outputDir        string  = "output"
	outputFileFormat string  = "out%d.bmp"
	tmpFile          string  = "./tmp/min.mp4"
)

var (
	scaleWidth       int     = 64
	scaleHeight      int     = 32
	scaleRatio       float32 = 1 / 2
	fps              int     = 1
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

func main() {
	ratio := fmt.Sprintf("scale=%d:%d,setdar=ratio=%f", scaleWidth, scaleHeight, scaleRatio)
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
		matched, _ := filepath.Match("out*.bmp", file.Name())
		if matched {
			ext := filepath.Ext(file.Name())
			inputFullPath := filepath.Join(outputDir, file.Name())
			outputFullPath := filepath.Join(outputDir, strconv.Itoa(i)+ext)

			err := grayscale(inputFullPath, outputFullPath)
			if err != nil {
				log.Fatalf("can't apply grayscale filter to %s:\n%s", inputFullPath, err)
			}

			rawBMP, err := os.Open(outputFullPath)
			if err != nil {
				log.Fatalf("unexpected error during opening BMP file: %s", err)
			}

			image, err := bmp.Decode(rawBMP)
			fmt.Print(image)
			if err != nil {
				log.Fatalf("unexpected error during BMP decoding: %s", err)
			}

		}
	}
}
