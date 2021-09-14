package main

import (
	"flag"
	"fmt"
	"io/fs"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path"
	"strconv"
)

var (
	inputFile        = flag.String("input", "", "Path to the input file for converting")
	outputFile       = flag.String("output", "output.bin", "Path to the output file. Default is <<current_dir>>/output.bin")
	outputFileFormat = "out%d.bmp"
)

var (
	scaleWidth  = flag.Int("width", 64, "Width of the output file. Default: 64")
	scaleHeight = flag.Int("height", 32, "Height of the output file. Default: 32")
	fps         = flag.Int("fps", 16, "FPS of the output file. FPS count must be less than 16. Default: 16")
)

func main() {
	flag.Parse()
	scaleRatio := float64(*scaleWidth) / float64(*scaleHeight)

	if *inputFile == "" {
		log.Fatalf("no input file was provided\n")
	}

	_, err := exec.LookPath("ffmpeg")
	if err != nil {
		log.Fatalf("\"FFMpeg\" is not installed, proceed to link to download: https://www.ffmpeg.org/\n")
	}

	tmpDir, err := ioutil.TempDir("", "nipkow")
	if err != nil {
		log.Fatalf("unable to create temporary directory: %s\n", err)
	}

	defer func(path string) {
		err := os.RemoveAll(path)
		if err != nil {
			log.Fatalf("unable to delete temprorary directory: %s\n", err)
		}
	}(tmpDir)

	ratio := fmt.Sprintf("scale=%d:%d,setdar=ratio=%f", *scaleWidth, *scaleHeight, scaleRatio)
	err = minify(*inputFile, path.Join(tmpDir, "min.mp4"), ratio)
	if err != nil {
		fmt.Printf("error: %s\n", err)
		return
	}

	err = framify(path.Join(tmpDir, "min.mp4"), path.Join(tmpDir, outputFileFormat), *fps)
	if err != nil {
		fmt.Printf("error: %s\n", err)
		return
	}

	fileSystem := os.DirFS(tmpDir)
	if err != nil {
		log.Fatal(err)
	}

	output, err := os.OpenFile(*outputFile, os.O_CREATE|os.O_WRONLY, 0777)
	files, err := fs.ReadDir(fileSystem, ".")

	for i, file := range files {
		matched, _ := path.Match("out*.bmp", file.Name())
		if matched {
			ext := path.Ext(file.Name())
			inputFullPath := path.Join(tmpDir, file.Name())
			outputFullPath := path.Join(tmpDir, strconv.Itoa(i)+ext)

			err := grayscale(inputFullPath, outputFullPath)
			if err != nil {
				fmt.Printf("can't apply grayscale filter to %s:\n%s\n", inputFullPath, err)
				return
			}

			img, err := loadBMP(outputFullPath)
			if err != nil {
				fmt.Printf("unable to load BMP file: %s\n", err)
			}

			bytesOffset := int64(i * (*scaleWidth) * (*scaleHeight))
			at, err := output.WriteAt(img.Pix, bytesOffset)
			if err != nil {
				fmt.Printf("error during writing to output file at %d: %s\n", at, err)
			}
		}
	}
}
