package main

import (
	"aoc/util"
	"fmt"
	"io"
	"slices"
	"strconv"
	"strings"
)

func main() {
	util.RunMain(mainErr)
}

func mainErr(r io.Reader) error {
	var result int

	for str := range util.AllTrimmedLines(r) {
		var e equation

		p := strings.Split(str, ":")
		if len(p) != 2 {
			return fmt.Errorf("unexpected split on ':' for %q", str)
		}

		var err error
		e.test, err = strconv.Atoi(p[0])
		if err != nil {
			return fmt.Errorf("parsing result portion of %q: %s", str, err)
		}

		for _, argStr := range strings.Split(p[1], " ") {
			if argStr == "" {
				continue
			}
			arg, err := strconv.Atoi(argStr)
			if err != nil {
				return fmt.Errorf("parsing argument %q in %q: %s", argStr, str, err)
			}
			e.inputs = append(e.inputs, arg)
		}

		if !validEquation(e) {
			continue
		}
		result += e.test
	}

	fmt.Printf("result: %d\n", result)

	return nil
}

type equation struct {
	test   int
	inputs []int
}

// validEquation returns true if there is some combination of
// addition or multiplication operations which can be inserted
// between the 'inputs' numbers to create an equation which results
// in 'test'. Operations are strictly evaluated left-to-right.
func validEquation(e equation) bool {
	/*

		We solve via induction and a tree-search. The base case:
			test: 0
			inputs: []

		is trivially valid. Then, we have three inductive rules to build
		up valid equations.

		Addition:

			{x+n, n:xs} is valid if {x, xs} is valid

		Multiplication:

			{x*n, n:xs} is valid if {x, xs} is valid

		Concatenation:

			{x || n, n:xs} is valid if {x, xs} is valid

		We can work these rules backwards starting from a proposed
		solution:

			- {n, x:xs} is valid if {n-x, xs} is valid
			- {n, x:xs} is valid if {n/x, xs} is valid (and if n%x == 0)
			- (n, x:xs) is valid if {unconcat(n,x), xs} is valid
			- {0, []} is valid

		(in the implementation we reverse the operands to making popping
		values off slightly easier in Go)

	*/
	slices.Reverse(e.inputs)
	return validEquationParts(e.test, e.inputs)
}

func validEquationParts(test int, reversedInputs []int) bool {
	if len(reversedInputs) == 0 {
		return test == 0
	}
	if test < 1 {
		return false
	}

	nextArg := reversedInputs[0]
	nextInputs := reversedInputs[1:]

	// try addition
	if test != 0 && validEquationParts(test-nextArg, nextInputs) {
		return true
	}

	// try multiplication
	if test%nextArg == 0 && validEquationParts(test/nextArg, nextInputs) {
		return true
	}

	// try concatenation
	prefix, canUnconcat := unconcat(test, nextArg)
	if !canUnconcat {
		return false
	}
	return validEquationParts(prefix, nextInputs)
}

func unconcat(x int, n int) (int, bool) {
	xStr := strconv.Itoa(x)
	nStr := strconv.Itoa(n)

	if !strings.HasSuffix(xStr, nStr) {
		return 0, false
	}

	prefixStr := strings.TrimSuffix(xStr, nStr)
	prefix, err := strconv.Atoi(prefixStr)
	if err != nil {
		return 0, false
	}
	return prefix, true
}
