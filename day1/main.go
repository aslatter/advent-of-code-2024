package main

import (
	"fmt"
	"io"

	"aoc/util"
)

func main() {
	util.RunMain(mainErr)
}

func mainErr(r io.Reader) error {
	var lefts []int
	rights := map[int]int{}

	for line := range util.AllTrimmedLines(r) {
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
		rights[right]++
	}

	var score int

	for i := range len(lefts) {
		leftNum := lefts[i]
		leftOccurences := rights[leftNum]
		score += (leftNum * leftOccurences)
	}

	fmt.Println("result:", score)

	return nil
}
