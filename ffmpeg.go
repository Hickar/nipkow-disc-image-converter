package main

import (
	"bytes"
	"fmt"
	"os/exec"
)

func execute(name string, arg ...string) error {
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	cmd := exec.Command(name, arg...)
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	if err := cmd.Start(); err != nil {
		fmt.Printf("stdErr: %s", stderr.String())
		return err
	}

	if err := cmd.Wait(); err != nil {
		fmt.Printf("stdErr: %s\n", stderr.String())
		return err
	}

	return nil
}

func minify(input, output string, ratio string) error {
	return execute("ffmpeg", "-i", input, "-vf", ratio, output)
}

func framify(input, output string, fps int) error {
	return execute("ffmpeg", "-i", input, "-vf", fmt.Sprintf("fps=%d", fps), output)
}

func grayscale(input, output string) error {
	return execute("ffmpeg", "-i", input, "-vf", "format=gray", output)
}