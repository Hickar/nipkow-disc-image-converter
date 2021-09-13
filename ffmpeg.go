package main

import (
	"bytes"
	"fmt"
	"log"
	"os/exec"
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
		log.Printf("stdErr: %s\n", stderr.String())
		return err
	}

	fmt.Printf("stdOut: %s\n", stdout.String())
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