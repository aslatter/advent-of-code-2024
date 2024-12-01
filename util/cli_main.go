package util

import (
	"flag"
	"fmt"
	"io"
	"os"
)

func RunMain(f func(r io.Reader) error) {
	var inputFileName string
	flag.StringVar(&inputFileName, "input", "", "file to read inputs from")
	flag.Parse()

	input, err := os.Open(inputFileName)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error opening file: %s\n", err)
		os.Exit(1)
	}
	err = f(input)
	_ = input.Close()
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %s\n", err)
		os.Exit(1)
	}
}
