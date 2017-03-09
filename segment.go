package main

import (
	"fmt"
	"strings"
)

type Segment struct {
	Name string
	Start, End int
}

func segmentEqual(a, b []Segment) bool {

	if a == nil && b == nil {
		return true;
	}

	if a == nil || b == nil {
		return false;
	}

	if len(a) != len(b) {
		return false
	}

	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}

	return true
}

func SegmentSource(src []string) []Segment {

	segments := []Segment{}

	for index, line := range src {
		if strings.Contains(line, "## @") {
			entryName := ExtractName(strings.Split(line, "## @")[1])

			for _, s := range segments {
				if s.Name == entryName {
					panic(fmt.Sprintf("Entry name %s already found", entryName))
				}
			}

			segments = append(segments, Segment{Name: entryName, Start: index})
		}

		if strings.Contains(line, ".exit") {
			exitName := ExtractName(strings.Split(line, "## %")[1])

			isegment := -1
			var s Segment
			for isegment, s = range segments {
				if s.Name == exitName {
					break
				}
			}
			if isegment == -1 || isegment == len(segments) {
				panic(fmt.Sprintf("No entry name found for exit %s", exitName))
			}

			segments[isegment].End = index
		}
	}

	return segments
}
