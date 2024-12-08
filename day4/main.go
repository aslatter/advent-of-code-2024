package main

import (
	"aoc/util"
	"bufio"
	"fmt"
	"io"
	"iter"
	"slices"
)

func main() {
	util.RunMain(mainErr)
}

func mainErr(r io.Reader) error {
	// 'data' is a word-grid consisting of letters.
	var data [][]byte
	s := bufio.NewScanner(r)
	for s.Scan() {
		row := slices.Clone(s.Bytes())
		if len(row) == 0 {
			continue
		}
		data = append(data, slices.Clone(s.Bytes()))
	}
	input := newPuzzle(data)

	var result int
	// loop over all positions in the puzzle. (0,0) is the top-left of the input,
	// with 'x' increasing to the right and 'y' increasing as we go down.
	for p := range input.AllPoints() {
		// find all patterns like the word "MAS" in a diagonal cross.
		if p.c != 'A' {
			continue
		}

		/**

		Once we find an 'A' we extract the 'corners'
		arround the 'A', going clockwise.

		Valid inputs are:

		M M
		 A
		S S

		S M
		 A
		S M

		S S
		 A
		M M

		M S
		 A
		M S

		These correspond to the "corner strings":
		 MMSS
		 SMMS
		 SSMM
		 MSSM

		**/

		c := input.Corners(p)
		switch c {
		case "MMSS", "SMMS", "SSMM", "MSSM":
			result++
		}
	}

	fmt.Printf("result: %d\n", result)

	return nil
}

type puzzle struct {
	width  int
	height int
	raw    [][]byte
}

func newPuzzle(raw [][]byte) puzzle {
	// panic if we don't have data
	_ = raw[0][0]

	return puzzle{
		width:  len(raw[0]),
		height: len(raw),
		raw:    raw,
	}
}

type point struct {
	x int
	y int
	c byte
}

func (p *puzzle) AllPoints() iter.Seq[point] {
	return func(yield func(point) bool) {
		for i := range p.height {
			for j := range p.width {
				if !yield(point{
					x: j,
					y: i,
					c: p.raw[i][j],
				}) {
					return
				}
			}
		}
	}
}

// return the corners around the position, clockwise
// starting from the upper-left.
func (p *puzzle) Corners(pos point) string {
	if pos.x == 0 || pos.x == p.width-1 {
		return ""
	}
	if pos.y == 0 || pos.y == p.height-1 {
		return ""
	}

	return string([]byte{
		p.raw[pos.y-1][pos.x-1], // upper left
		p.raw[pos.y-1][pos.x+1], // upper right
		p.raw[pos.y+1][pos.x+1], // lower right
		p.raw[pos.y+1][pos.x-1], // lower left
	})
}
