package main

import (
	"fmt"
	"io/fs"
	"io/ioutil"
	"log"
	"os"
	"path"
	"path/filepath"
	"strconv"
)

var (
	inputFile        string = "./input/mickey.mp4"
	outputFileFormat string = "out%d.bmp"
)

var (
	scaleWidth  int     = 64
	scaleHeight int     = 32
	scaleRatio  float32 = 1 / 2
	fps         int     = 1
)

func main() {
	ratio := fmt.Sprintf("scale=%d:%d,setdar=ratio=%f", scaleWidth, scaleHeight, scaleRatio)

	tmpDir, err := ioutil.TempDir("", "nipkow")
	if err != nil {
		log.Fatalf("unable to create temporary directory: %s", err)
	}
	defer func(path string) {
		err := os.RemoveAll(path)
		if err != nil {
			log.Fatalf("unable to delete temprorary directory: %s", err)
		}
	}(tmpDir)

	err = minify(inputFile, filepath.Join(tmpDir, "min.mp4"), ratio)
	if err != nil {
		log.Fatal(err)
	}

	err = framify(filepath.Join(tmpDir, "min.mp4"), path.Join(tmpDir, outputFileFormat), fps)
	if err != nil {
		log.Fatal(err)
	}

	fileSystem := os.DirFS(tmpDir)
	if err != nil {
		log.Fatal(err)
	}

	output, err := os.OpenFile("output.bin", os.O_CREATE|os.O_WRONLY, 0777)
	files, err := fs.ReadDir(fileSystem, ".")

	for i, file := range files {
		matched, _ := filepath.Match("out*.bmp", file.Name())
		if matched {
			ext := filepath.Ext(file.Name())
			inputFullPath := filepath.Join(tmpDir, file.Name())
			outputFullPath := filepath.Join(tmpDir, strconv.Itoa(i)+ext)

			err := grayscale(inputFullPath, outputFullPath)
			if err != nil {
				log.Fatalf("can't apply grayscale filter to %s:\n%s", inputFullPath, err)
			}

			img, err := loadBMP(outputFullPath)
			if err != nil {
				log.Fatalf("unable to load BMP file: %s", err)
			}

			at, err := output.WriteAt(img.Pix, int64(i*2048))
			if err != nil {
				log.Fatalf("error during writing to output file at %d: %s", at, err)
			}

			fmt.Printf("output at: %d\n", at)
			fmt.Printf("written overall:%d\n\n", int64(i*2048))
		}
	}
}
