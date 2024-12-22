package main

import (
	"aoc/util"
	"fmt"
	"io"
	"iter"
)

func main() {
	util.RunMain(mainErr)
}

func mainErr(r io.Reader) error {
	antennae := map[byte][]pos{}

	var y int
	var maxX int
	for str := range util.AllTrimmedLines(r) {
		maxX = len(str) - 1
		for x, b := range []byte(str) {
			switch {
			case b >= 'a' && b <= 'z', b >= 'A' && b <= 'Z', b >= '0' && b <= '9':
				antennae[b] = append(antennae[b], pos{x, y})
			}
		}
		y++
	}
	maxY := y - 1

	antinodes := map[pos]bool{}

	for _, allPos := range antennae {
		for pair := range pairs(allPos) {
			a := antinode(pair[0], pair[1])
			if a.x >= 0 && a.x <= maxX && a.y >= 0 && a.y <= maxY {
				antinodes[a] = true
			}
			a = antinode(pair[1], pair[0])
			if a.x >= 0 && a.x <= maxX && a.y >= 0 && a.y <= maxY {
				antinodes[a] = true
			}
		}
	}

	result := len(antinodes)
	fmt.Printf("result: %d\n", result)

	return nil
}

type pos struct{ x, y int }

func pairs(x []pos) iter.Seq[[2]pos] {
	return func(yield func([2]pos) bool) {
		for i := range len(x) {
			for j := i + 1; j < len(x); j++ {
				if !yield([2]pos{x[i], x[j]}) {
					return
				}
			}
		}
	}
}

func antinode(p1 pos, p2 pos) pos {
	dx := p1.x - p2.x
	dy := p1.y - p2.y
	return pos{
		x: p2.x - dx,
		y: p2.y - dy,
	}
}
