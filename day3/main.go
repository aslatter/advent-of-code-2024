package main

import (
	"aoc/util"
	"fmt"
	"io"
	"regexp"
	"strconv"
)

func main() {
	util.RunMain(mainErr)
}

func mainErr(r io.Reader) error {
	rx := regexp.MustCompile(`mul\((\d+),(\d+)\)`)

	data, err := io.ReadAll(r)
	if err != nil {
		return fmt.Errorf("reading input data: %s", err)
	}

	var result int

	for _, match := range rx.FindAllStringSubmatch(string(data), -1) {
		if len(match) != 3 {
			return fmt.Errorf("unexpected number of matches in %v", match)
		}

		fmt.Printf("match info: %+v\n", match)

		lhs, err := strconv.Atoi(match[1])
		if err != nil {
			return fmt.Errorf("parsing lhs of %q: %s", match[0], err)
		}
		rhs, err := strconv.Atoi(match[2])
		if err != nil {
			return fmt.Errorf("parsing rhs of %q: %s", match[0], err)
		}

		product := lhs * rhs
		result += product
	}

	fmt.Printf("result: %d\n", result)

	return nil
}
