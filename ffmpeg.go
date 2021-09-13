package main

import "fmt"

func minify(input, output, ratio string) error {
	return execute("ffmpeg", "-i", input, "-vf", ratio, output)
}

func framify(input, output string, fps int) error {
	return execute("ffmpeg", "-i", input, "-vf", fmt.Sprintf("fps=%d", fps), output)
}

func grayscale(input, output string) error {
	return execute("ffmpeg", "-i", input, "-vf", "format=gray", output)
}