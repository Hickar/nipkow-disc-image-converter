package main

import (
	"fmt"
	"image"
	"log"
	"os"

	"golang.org/x/image/bmp"
)

func loadBMP(filepath string) (*image.Paletted, error) {
	rawBMP, err := os.Open(filepath)
	if err != nil {
		return nil, fmt.Errorf("unexpected error during opening BMP file: %s", err)
	}

	bmpBytes, err := bmp.Decode(rawBMP)
	if err != nil {
		return nil, fmt.Errorf("unexpected error during BMP decoding: %s", err)
	}

	img, ok := bmpBytes.(*image.Paletted)
	if !ok {
		log.Fatalf("error on image type assertion: %s", filepath)
	}

	return img, nil
}