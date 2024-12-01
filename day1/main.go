package main

import (
	"bufio"
	"fmt"
	"io"
	"slices"
	"strings"

	"aoc/util"
)

func main() {
	util.RunMain(mainErr)
}

func mainErr(r io.Reader) error {
	var lefts []int
	var rights []int

	s := bufio.NewScanner(r)
	for s.Scan() {
		line := s.Text()
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		var left int
		var right int
		n, err := fmt.Sscanf(line, "%d %d", &left, &right)
		if err != nil {
			return fmt.Errorf("parsing input line %q: %s", line, err)
		}
		if n != 2 {
			return fmt.Errorf("unexpected number of number parsed from line %q: %d", line, n)
		}

		lefts = append(lefts, left)
		rights = append(rights, right)
	}

	slices.Sort(lefts)
	slices.Sort(rights)

	var score int

	for i := range len(lefts) {
		diff := lefts[i] - rights[i]
		if diff < 0 {
			diff = -diff
		}
		score += diff
	}

	fmt.Println("result:", score)

	return nil
}
