package main

import (
	"aoc/util"
	"fmt"
	"io"
	"maps"
	"slices"
)

func main() {
	util.RunMain(mainErr)
}

func mainErr(r io.Reader) error {
	// parse the input
	var level [][]bool
	var startX, startY int

	for lineStr := range util.AllTrimmedLines(r) {
		var row []bool
		for _, c := range lineStr {
			var isObstructed bool
			if c == '#' {
				isObstructed = true
			}
			if c == '^' {
				startX = len(row)
				startY = len(level)
			}
			row = append(row, isObstructed)

		}
		level = append(level, row)
	}

	s := state{
		level: level,
	}
	s.setPosition(startX, startY, orientationNorth)

	// find positions to evaluate
	trials := findBasePositions(&s)

	// evaluate if adding a barrier would trap the guard
	var result int
	for _, p := range trials {
		if evaluateIntervention(&s, p) {
			result++
		}
	}

	fmt.Printf("result: %d\n", result)

	return nil
}

type pos struct {
	x int
	y int
}

type posAndDirection struct {
	x         int
	y         int
	direction orientation
}

type state struct {
	level     [][]bool
	visited   map[posAndDirection]bool
	x         int
	y         int
	direction orientation
	result    disposition
}

func (s *state) clone() *state {
	var newS state

	newS.level = make([][]bool, 0, len(s.level))
	for _, r := range s.level {
		newS.level = append(newS.level, slices.Clone(r))
	}

	newS.visited = maps.Clone(s.visited)

	newS.x = s.x
	newS.y = s.y
	newS.direction = s.direction

	return &newS
}

func (s *state) setPosition(x int, y int, direction orientation) bool {
	if s.visited == nil {
		s.visited = map[posAndDirection]bool{}
	}

	l := posAndDirection{
		x:         x,
		y:         y,
		direction: direction,
	}

	if s.visited[l] {
		return false
	}

	s.visited[l] = true

	s.x = x
	s.y = y
	s.direction = direction

	return true
}

func (s *state) next() bool {
	nextX := s.x
	nextY := s.y
	nextDirection := s.direction

	switch s.direction {
	case orientationNorth:
		nextY--
	case orientationEast:
		nextX++
	case orientationSouth:
		nextY++
	case orientationWest:
		nextX--
	}

	// did we move off the map?
	if nextX < 0 || nextY < 0 || nextX == len(s.level[0]) || nextY == len(s.level) {
		s.result = dispositionExited
		return false
	}

	// did we attempt to move into an obstruction?
	if s.level[nextY][nextX] {
		// reset position, turn to the right
		nextX = s.x
		nextY = s.y
		nextDirection++
		nextDirection %= 4
	}

	// did we hit a loop?
	if !s.setPosition(nextX, nextY, nextDirection) {
		s.result = dispositionLooped
		return false
	}

	// keep going
	return true
}

type orientation int

const (
	orientationNorth orientation = iota
	orientationEast
	orientationSouth
	orientationWest
)

type disposition int

const (
	dispositionInProgress disposition = iota
	dispositionExited
	dispositionLooped
)

// findBasePositions finds positions visited by the guard if
// we don't intervene. We don't need to evaluate any other positions
// because those barriers will not be reachable by the guard.
func findBasePositions(base *state) []pos {
	s := base.clone()
	for s.next() {
	}
	var result []pos
	for pp := range s.visited {
		result = append(result, pos{x: pp.x, y: pp.y})
	}

	// remove duplicates
	slices.SortFunc(result, func(a, b pos) int {
		d := a.x - b.x
		if d == 0 {
			return a.y - b.y
		}
		return d
	})
	result = slices.Compact(result)

	return result
}

// evaluateIntervention returns 'true' if placing a barrier
// at position 'pos' would result in trapping the guard in
// a loop.
func evaluateIntervention(base *state, p pos) bool {
	s := base.clone()

	// add barrier
	s.level[p.y][p.x] = true

	// run simulation
	for s.next() {
	}

	// we either stopped because the guard left or
	// because the guard is trapped.
	return s.result == dispositionLooped
}
