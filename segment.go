package main

import (
	"fmt"
	"strings"
)

type Segment struct {
	Name       string
	Start, End int
	stack      Stack
}

type Exit struct {
	Name  string
	End   int
	stack Stack
}

func segmentEqual(a, b []Segment) bool {

	if a == nil && b == nil {
		return true
	}

	if a == nil || b == nil {
		return false
	}

	if len(a) != len(b) {
		return false
	}

	for i := range a {
		if !(a[i].Name == b[i].Name && a[i].Start == b[i].Start && a[i].End == b[i].End) {
			return false
		}
	}

	return true
}

func SegmentSource(src []string) []Segment {

	exits := []Exit{}

	exitGatherUntilRet := -1

	// Collect exit of subroutines
	for index, line := range src {

		if strings.Contains(line, ".exit") {
			exitName := ExtractName(strings.Split(line, "## %")[1])

			for _, e := range exits {
				if e.Name == exitName {
					panic(fmt.Sprintf("Exit name %s already found", exitName))
				}
			}

			exits = append(exits, Exit{Name: exitName, End: index + 1})

			// Gather stack information
			exitGatherUntilRet = len(exits) - 1
		}

		if exitGatherUntilRet != -1 && strings.Contains(line, "ret") {

			// Lines of postamble
			stackLines := src[exits[exitGatherUntilRet].End : index+1]

			exits[exitGatherUntilRet].stack = ExtractStackInfo(stackLines)

			exitGatherUntilRet = -1
		}
	}

	segments := []Segment{}

	entryGatherUntilRet := -1

	// Find start of subroutines
	for index, line := range src {

		if strings.Contains(line, "## @") {
			entryName := ExtractName(strings.Split(line, "## @")[1])

			for _, s := range segments {
				if s.Name == entryName {
					panic(fmt.Sprintf("Entry name %s already found", entryName))
				}
			}

			var stack Stack
			end := -1
			for _, e := range exits {
				if e.Name == entryName {
					end = e.End
					stack = e.stack
					break
				}
			}

			segments = append(segments, Segment{Name: entryName, Start: index + 1, End: end, stack: stack})

			if end == -1 {
				// No exit label found, start searching for ret statement
				entryGatherUntilRet = len(segments) - 1
			}
		}

		if entryGatherUntilRet != -1 && strings.Contains(line, "ret") {

			// Found closing ret statement, start searching back to first non closing statement
			i := 1
			for ; index-i >= 0; i++ {
				if !IsStdCallEpilogue(src[index-i]) {
					break
				}
			}
			segments[entryGatherUntilRet].End = index - i + 1

			stackLines := src[index-i+1 : index+1]

			fmt.Println(stackLines)
			segments[entryGatherUntilRet].stack = ExtractStackInfo(stackLines)

			entryGatherUntilRet = -1
		}
	}

	return segments
}
