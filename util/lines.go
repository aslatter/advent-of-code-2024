package util

import (
	"bufio"
	"fmt"
	"io"
	"iter"
	"strings"
)

// AllTrimmedLines iterates over all lines in the [io.Reader] r, trimming
// leading and trailing whitespace.
func AllTrimmedLines(r io.Reader) iter.Seq[string] {
	return func(yield func(string) bool) {
		s := bufio.NewScanner(r)
		for s.Scan() {
			line := s.Text()
			line = strings.TrimSpace(line)
			if line == "" {
				continue
			}
			if !yield(line) {
				break
			}
		}
		if err := s.Err(); err != nil {
			panic(fmt.Errorf("reading input: %s", err))
		}
	}
}
