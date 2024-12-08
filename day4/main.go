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
		// the target word is 'XMAS', but we can find it forwards or backwards.
		if p.c != 'X' && p.c != 'S' {
			continue
		}
		var trialWords []string
		trialWords = append(trialWords, input.East(p))
		trialWords = append(trialWords, input.South(p))
		trialWords = append(trialWords, input.SouthEast(p))
		trialWords = append(trialWords, input.SouthWest(p))
		for _, w := range trialWords {
			if w == "XMAS" || w == "SAMX" {
				result++
			}
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

func (p *puzzle) South(pos point) string {
	var r []byte
	r = append(r, pos.c)
	if pos.y == p.height-1 {
		return string(r)
	}
	maxY := pos.y + 3
	if maxY > p.height-1 {
		maxY = p.height - 1
	}
	for i := pos.y + 1; i <= maxY; i++ {
		r = append(r, p.raw[i][pos.x])
	}
	return string(r)
}

func (p *puzzle) East(pos point) string {
	var r []byte
	r = append(r, pos.c)
	if pos.x == p.width-1 {
		return string(r)
	}
	maxX := pos.x + 3
	if maxX > p.width-1 {
		maxX = p.width - 1
	}
	for j := pos.x + 1; j <= maxX; j++ {
		r = append(r, p.raw[pos.y][j])
	}
	return string(r)
}

func (p *puzzle) SouthEast(pos point) string {
	var r []byte
	r = append(r, pos.c)
	maxDelta := 3
	if pos.x+maxDelta > p.width-1 {
		maxDelta = p.width - 1 - pos.x
	}
	if pos.y+maxDelta > p.height-1 {
		maxDelta = p.height - 1 - pos.y
	}
	for d := 1; d <= maxDelta; d++ {
		r = append(r, p.raw[pos.y+d][pos.x+d])
	}
	return string(r)
}

func (p *puzzle) SouthWest(pos point) string {
	var r []byte
	r = append(r, pos.c)
	maxDelta := 3
	if pos.x-maxDelta < 0 {
		maxDelta = pos.x
	}
	if pos.y+maxDelta > p.height-1 {
		maxDelta = p.height - 1 - pos.y
	}
	for d := 1; d <= maxDelta; d++ {
		r = append(r, p.raw[pos.y+d][pos.x-d])
	}
	return string(r)
}
