package main

import (
	"fmt"
	"regexp"
	"strings"
)

var regexpRet = regexp.MustCompile(`^\s*ret`)

type Segment struct {
	Name       string
	Start, End int
	stack      Stack
}

type Global struct {
	dotGlobalLine   int
	globalName      string
	globalLabelLine int
}

func SplitOnGlobals(lines []string) []Global {

	var result []Global

	for index, line := range lines {
		if strings.Contains(line, ".globl") {

			scrambled := strings.TrimSpace(strings.Split(line, ".globl")[1])
			name := ExtractName(scrambled)

			labelLine := findLabel(lines, scrambled)

			result = append(result, Global{dotGlobalLine: index, globalName: name, globalLabelLine: labelLine})
		}
	}

	return result
}

// Segment the source into multiple routines
func SegmentSource(src []string) []Segment {

	globals := SplitOnGlobals(src)

	if len(globals) == 0 {
		return []Segment{}
	}

	segments := []Segment{}

	splitBegin := globals[0].dotGlobalLine
	for iglobal, global := range globals {
		splitEnd := len(src)
		if iglobal < len(globals)-1 {
			splitEnd = globals[iglobal+1].dotGlobalLine
		}

		// Search for `ret` statement starting from the back
		for index := splitEnd - 1; index >= splitBegin; index-- {
			if match := regexpRet.FindStringSubmatch(src[index]); len(match) > 0 {

				// Found closing ret statement, start searching back to first non closing statement
				i := 1
				for ; index-i >= 0; i++ {
					if !IsEpilogueInstruction(src[index-i]) {
						break
					}
				}

				segments = append(segments, Segment{Name: global.globalName, Start: global.globalLabelLine + 1, End: index - i + 1})
			}
		}

		splitBegin = splitEnd
	}

	return segments
}

func SegmentEatPrologue(lines []string, stack *Stack) int {

	index, line := 0, ""

	for index, line = range lines {

		var skip bool
		line, skip = stripComments(line) // Remove ## comments
		if skip {
			continue
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
