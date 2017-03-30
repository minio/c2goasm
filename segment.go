package main

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"unicode"
)

var regexpRet = regexp.MustCompile(`^\s*ret`)

type Subroutine struct {
	name               string
	bodyStart, bodyEnd int
	epilogue           Epilogue
}

type Global struct {
	dotGlobalLine   int
	globalName      string
	globalLabelLine int
}

func splitOnGlobals(lines []string) []Global {

	var result []Global

	for index, line := range lines {
		if strings.Contains(line, ".globl") {

			scrambled := strings.TrimSpace(strings.Split(line, ".globl")[1])
			name := extractName(scrambled)

			labelLine := findLabel(lines, scrambled)

			result = append(result, Global{dotGlobalLine: index, globalName: name, globalLabelLine: labelLine})
		}
	}

	return result
}

// Segment the source into multiple routines
func segmentSource(src []string) []Subroutine {

	globals := splitOnGlobals(src)

	if len(globals) == 0 {
		return []Subroutine{}
	}

	subroutines := []Subroutine{}

	splitBegin := globals[0].dotGlobalLine
	for iglobal, global := range globals {
		splitEnd := len(src)
		if iglobal < len(globals)-1 {
			splitEnd = globals[iglobal+1].dotGlobalLine
		}

		// Search for `ret` statement starting from the back
		for lineRet := splitEnd - 1; lineRet >= splitBegin; lineRet-- {
			if match := regexpRet.FindStringSubmatch(src[lineRet]); len(match) > 0 {

				// Found closing ret statement, start searching back to first non closing statement
				i := 1
				for ; lineRet-i >= 0; i++ {
					if !isEpilogueInstruction(src[lineRet-i]) {
						break
					}
				}

				epilogueLines := src[lineRet-i+1 : lineRet+1]

				epilogue := extractEpilogueInfo(epilogueLines)

				subroutines = append(subroutines, Subroutine{name: global.globalName, bodyStart: global.globalLabelLine + 1, bodyEnd: lineRet - i + 1, epilogue: epilogue})
			}
		}

		splitBegin = splitEnd
	}

	return subroutines
}

func eatPrologueLines(lines []string, epilogue *Epilogue) int {

	index, line := 0, ""

	for index, line = range lines {

		var skip bool
		line, skip = stripComments(line) // Remove ## comments
		if skip {
			continue
		}

		if !epilogue.IsPrologueInstruction(line) {
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

func extractNamePart(part string) (int, string) {

	digits := 0
	for _, d := range part {
		if unicode.IsDigit(d) {
			digits += 1
		} else {
			break
		}
	}
	length, _ := strconv.Atoi(part[:digits])
	return digits + length, part[digits:(digits + length)]
}

func extractName(name string) string {

	var parts []string

	for index, ch := range name {
		if unicode.IsDigit(ch) {

			for index < len(name) {
				size, part := extractNamePart(name[index:])
				if size == 0 {
					break
				}

				parts = append(parts, part)
				index += size
			}

			break
		}
	}

	return strings.Join(parts, "")
}
