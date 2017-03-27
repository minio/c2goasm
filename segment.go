package main

import (
	"fmt"
	"strings"
	"regexp"
)

var regexpRet = regexp.MustCompile(`^\s*ret`)

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

// Segment the soure into multiple routines
func SegmentSource(src []string) []Segment {

	globals := SplitOnGlobals(src)

	if len(globals) == 0 {
		return []Segment{}
	}

	segments := []Segment{}


	for splitBegin, g := globals[0], 0; g < len(globals); g++ {
		splitEnd := len(src)
		if g < len(globals) - 1 {
			splitEnd = globals[g + 1]
		}

		scrambledName := strings.TrimSpace(strings.Split(src[splitBegin], ".globl")[1])
		entryName := ExtractName(scrambledName)

		// Search for ret statement from the back
		for index := splitEnd-1; index >= splitBegin; index-- {
			if match := regexpRet.FindStringSubmatch(src[index]); len(match) > 0 {

				// Found closing ret statement, start searching back to first non closing statement
				i := 1
				for ; index-i >= 0; i++ {
					if !IsEpilogueInstruction(src[index-i]) {
						break
					}
				}

				startLine := splitBegin + findLabel(src[splitBegin:splitEnd], scrambledName) + 1
				segments = append(segments, Segment{Name: entryName, Start: startLine, End: index - i + 1})
			}
		}

		splitBegin = splitEnd
	}

	return segments
}

func SegmentEatPrologue(lines []string, stack *Stack) int {

	index, line := 0, ""

	for index, line = range lines {

		// Remove ## comments
		if parts := strings.SplitN(line, `##`, 2); len(parts) > 1 {
			if strings.TrimSpace(parts[0]) == "" {
				continue
			}
			line = parts[0]
		}

		if !stack.IsPrologueInstruction(line) {
			break
		}
	}

	return index
}

func findLabel(lines []string, label string) int {

	labelDef := label + ":"

	for index, line := range lines {
		if strings.HasPrefix(line, labelDef) {
			return index
		}
	}

	panic(fmt.Sprintf("Failed to find label: %s", labelDef))
}
