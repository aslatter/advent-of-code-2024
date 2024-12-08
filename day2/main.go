package main

import (
	"fmt"
	"io"
	"slices"
	"strconv"
	"strings"

	"aoc/util"
)

func main() {
	util.RunMain(mainErr)
}

func mainErr(r io.Reader) error {
	var result int

	for line := range util.AllTrimmedLines(r) {
		var r report
		pieces := strings.Split(line, " ")
		for _, p := range pieces {
			n, err := strconv.Atoi(p)
			if err != nil {
				return fmt.Errorf("parsing string %q for line %q: %s", p, line, err)
			}
			r = append(r, n)
		}

		isSafe := isSafeReport(r)
		if isSafe {
			result++
		}

		fmt.Printf("line %q is safe? %v\n", line, isSafe)
	}

	fmt.Printf("result: %d\n", result)

	return nil
}

type report []int

func isSafeReport(r report) bool {
	if isReportStrictlySafe(r) {
		return true
	}
	for i := range r {
		rr := slices.Delete(slices.Clone(r), i, i+1)
		if isReportStrictlySafe(rr) {
			return true
		}
	}
	return false
}

func isReportStrictlySafe(r report) bool {
	var isIncreasing bool
	if r[1] > r[0] {
		isIncreasing = true
	}

	// check numbers against each-other
	prev := r[0]

	if isIncreasing {
		for _, next := range r[1:] {
			if next <= prev {
				return false
			}
			if next-prev > 3 {
				return false
			}
			prev = next
		}
		return true
	}
	for _, next := range r[1:] {
		if next >= prev {
			return false
		}
		if prev-next > 3 {
			return false
		}
		prev = next
	}

	return true
}
