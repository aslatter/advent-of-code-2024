package main

import (
	"aoc/util"
	"fmt"
	"io"
)

func main() {
	util.RunMain(mainErr)
}

func mainErr(r io.Reader) error {
	// read all input as bytes
	diskMap, err := io.ReadAll(r)
	if err != nil {
		return fmt.Errorf("reading input bytes: %s", err)
	}

	// entries with odd indices are free-space, even indices are data

	// we iterate through the map updating our checksum as we go.
	// when we hit a 'zero' spot we iterate from the back, and when
	// the zero is exhausted we resume iterating from the front.

	type state struct {
		atFront       bool
		blockPosition int

		frontLocation  int
		frontValue     int
		frontRemaining int

		backLocation         int
		backValue            int
		backRemaining        int
		zerosToFillRemaining int
	}
	var s state
	s.backLocation = len(diskMap)
	s.frontLocation = -1

	// if we've initialized the 'back location' to a spot with zeros,
	// bump it out one more
	if s.backLocation%2 == 1 {
		s.backLocation++
	}

	var result int

	for {
		if !s.atFront {
			if s.zerosToFillRemaining == 0 {
				s.atFront = true
				s.frontLocation++
				s.frontValue = s.frontLocation / 2

				if s.frontLocation > s.backLocation {
					break
				}

				s.frontRemaining = int(diskMap[s.frontLocation]) - '0'
				if s.frontLocation == s.backLocation {
					// this block may have been partially processed.
					// steal remaining work from the "back" state.
					s.frontRemaining = s.backRemaining
					s.backRemaining = 0
				}

				continue
			}
			if s.backRemaining == 0 {
				s.backLocation -= 2
				s.backRemaining = int(diskMap[s.backLocation]) - '0'
				s.backValue = s.backLocation / 2

				if s.backLocation < s.frontLocation {
					break
				}
				continue
			}

			result += s.blockPosition * s.backValue
			s.zerosToFillRemaining--
			s.backRemaining--
			s.blockPosition++
			continue

		}
		if s.frontRemaining == 0 {
			s.atFront = false
			s.frontLocation++

			if s.frontLocation == len(diskMap) {
				// ???
				break
			}

			s.zerosToFillRemaining = int(diskMap[s.frontLocation]) - '0'
			continue
		}

		result += s.blockPosition * s.frontValue
		s.frontRemaining--
		s.blockPosition++
		continue
	}

	fmt.Printf("result: %d\n", result)

	return nil
}
