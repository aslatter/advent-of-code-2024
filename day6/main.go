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

	for s.next() {
	}

	result := len(s.visited)

	fmt.Printf("result: %d\n", result)

	return nil
}

type pos struct {
	x int
	y int
}

type state struct {
	level     [][]bool
	visited   map[pos]bool
	x         int
	y         int
	direction orientation
}

func (s *state) setPosition(x int, y int, direction orientation) {
	if s.visited == nil {
		s.visited = map[pos]bool{}
	}

	s.visited[pos{x: x, y: y}] = true
	s.x = x
	s.y = y
	s.direction = direction
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
		return false
	}

	// did we attempt to move into an obstruction?
	if s.level[nextY][nextX] {
		nextX = s.x
		nextY = s.y
		nextDirection++
		nextDirection %= 4
	}

	s.setPosition(nextX, nextY, nextDirection)

	// keep going
	return true
}

type orientation int

const (
	orientationNorth orientation = iota
	orientationEast  orientation = iota
	orientationSouth orientation = iota
	orientationWest  orientation = iota
)
