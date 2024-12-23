package main

import (
	"aoc/util"
	"fmt"
	"io"
	"iter"
	"math/big"
)

func main() {
	util.RunMain(mainErr)
}

func mainErr(r io.Reader) error {

	// parse the locations of all antennae, grouped
	// by label.
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

	// log all positions which are anti-nodes.
	antinodes := map[pos]bool{}

	for _, allPos := range antennae {
		for pair := range pairs(allPos) {
			for a := range nodes(maxX, maxY, pair[0], pair[1]) {
				antinodes[a] = true
			}
		}
	}

	result := len(antinodes)
	fmt.Printf("result: %d\n", result)

	return nil
}

type pos struct{ x, y int }

// pairs returns all pair-wise matchings of the points
// in 'x'.
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

// nodes returns all points on a straight line between p1 and
// p2, while also in-bounds on the map.
func nodes(maxX int, maxY int, p1 pos, p2 pos) iter.Seq[pos] {
	return func(yield func(pos) bool) {
		dx := p2.x - p1.x
		dy := p2.y - p1.y
		// too lazy to do gcd
		r := big.NewRat(int64(dx), int64(dy))
		dx = int(r.Num().Int64())
		dy = int(r.Denom().Int64())

		// p1 is always in-line with p1 and p2
		if !yield(p1) {
			return
		}

		// p2 will be hit by moving from p1

		n := 0
		for {
			n++
			p := pos{
				x: p1.x + n*dx,
				y: p1.y + n*dy,
			}
			if p.x < 0 || p.y < 0 || p.x > maxX || p.y > maxY {
				break
			}
			if !yield(p) {
				return
			}
		}

		n = 0
		for {
			n--
			p := pos{
				x: p1.x + n*dx,
				y: p1.y + n*dy,
			}
			if p.x < 0 || p.y < 0 || p.x > maxX || p.y > maxY {
				break
			}
			if !yield(p) {
				return
			}
		}
	}
}
