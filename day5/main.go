package main

import (
	"bufio"
	"fmt"
	"io"
	"strconv"
	"strings"

	"aoc/util"
)

func main() {
	util.RunMain(mainErr)
}

func mainErr(r io.Reader) error {
	rules := newRules()
	s := bufio.NewScanner(r)

	// phase one - parse rules in for "page|page"
	for s.Scan() {
		str := s.Text()
		if str == "" {
			break
		}

		var lhs int
		var rhs int
		_, err := fmt.Sscanf(str, "%d|%d", &lhs, &rhs)
		if err != nil {
			return fmt.Errorf("parsing string %q: %s", str, err)
		}
		rules.add(lhs, rhs)
	}

	var result int

	// phase two - parse updates and evaluate if they follow rules
	for s.Scan() {
		str := s.Text()
		if str == "" {
			break
		}

		var pages []int
		for _, pStr := range strings.Split(str, ",") {
			p, err := strconv.Atoi(pStr)
			if err != nil {
				return fmt.Errorf("parsing page %q in line %q: %s", pStr, str, err)
			}
			pages = append(pages, p)
		}

		if !updateFollowsRules(rules, pages) {
			continue
		}

		middleOfBatch := pages[len(pages)/2]
		result += middleOfBatch
	}

	fmt.Printf("result: %d\n", result)

	return nil
}

func updateFollowsRules(r *rules, update []int) bool {
	// we ignore rules which concern pages not in the update.
	// so we need to know which pages are in the update.
	inUpdate := map[int]bool{}
	for _, p := range update {
		inUpdate[p] = true
	}

	// in order to track if a dependency is satisfied we need
	// to track which pages we've seen.
	seen := map[int]bool{}

	// now, loop through the update and make sure it follows the
	// rules
	for _, p := range update {
		for preReqPage := range r.deps[p] {
			if !inUpdate[preReqPage] {
				continue
			}
			if !seen[preReqPage] {
				return false
			}
		}
		seen[p] = true
	}

	return true
}

type rules struct {
	deps map[int]map[int]bool
}

func newRules() *rules {
	return &rules{
		deps: map[int]map[int]bool{},
	}
}

// add a dependency from p1 to p2. That is, p1 must appear before
// p2 in any update.
func (r *rules) add(p1, p2 int) {
	nestedMap, ok := r.deps[p2]
	if !ok {
		nestedMap = map[int]bool{}
		r.deps[p2] = nestedMap
	}
	nestedMap[p1] = true
}
