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

		// grab rules applicable to current update
		subRules := rules.subRules(pages)

		if updateFollowsRules(subRules, pages) {
			continue
		}

		sorted := subRules.topoSort(pages)

		middleOfBatch := sorted[len(sorted)/2]
		result += middleOfBatch
	}

	fmt.Printf("result: %d\n", result)

	return nil
}

func updateFollowsRules(r *rules, update []int) bool {
	// in order to track if a dependency is satisfied we need
	// to track which pages we've seen.
	seen := map[int]bool{}

	// now, loop through the update and make sure it follows the
	// rules
	for _, p := range update {
		for preReqPage := range r.deps[p] {
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

// subRules returns a set of rules which pertain to the subset
// of pages passed-in.
func (r *rules) subRules(pages []int) *rules {
	seen := map[int]bool{}
	for _, p := range pages {
		seen[p] = true
	}

	subRules := newRules()
	for p2, s := range r.deps {
		if !seen[p2] {
			continue
		}
		for p1 := range s {
			if !seen[p1] {
				continue
			}
			subRules.add(p1, p2)
		}
	}

	return subRules
}

func (r *rules) topoSort(pages []int) []int {
	var result []int
	processed := map[int]bool{}

	// pop pages off the front of 'pages'
	// add them to 'result' if deps are not met, or to the back of the list of they are

	for len(pages) > 0 {
		// we really shouldn't use a slice as a FIFO queue, but
		// this should be okay in practice.
		nextP := pages[0]
		pages = pages[1:]

		satisfied := true
		for dep := range r.deps[nextP] {
			if !processed[dep] {
				satisfied = false
				break
			}
		}

		if !satisfied {
			// try it again later
			if len(pages) == 0 {
				// uh oh - trivially unsatisfiable
				panic("oops")
			}
			pages = append(pages, nextP)
			continue
		}
		result = append(result, nextP)
		processed[nextP] = true
	}

	return result
}
