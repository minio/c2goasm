package main

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

type StackArgs struct {
	Number        int
	OffsetToFirst int
}

func ArgumentsOnStack(lines []string) StackArgs {

	offsets := make(map[uint]bool)

	for _, l := range lines {
		l, _ = stripComments(l)
		if match := regexpRbpLoadHigher.FindStringSubmatch(l); len(match) > 1 {
			offset, _ := strconv.Atoi(match[1])
			if _, found := offsets[uint(offset)]; !found {
				offsets[uint(offset)] = true
			}
		}
	}

	offset := ^uint(0)
	for o, _ := range offsets {
		if o < offset {
			offset = o
		}
	}
	return StackArgs{Number: len(offsets), OffsetToFirst: int(offset)}
}

func GetGolangArgs(goCompanion, protoName string) int {

	gocode, err := readLines(goCompanion)
	if err != nil {
		panic(fmt.Sprintf("Failed to read companion go code: %v", err))
	}

	regexpFuncAndArgs := regexp.MustCompile(`func\s+(.*)\((.*)\)`)

	for _, g := range gocode {

		// Search for name of function and arguments
		if match := regexpFuncAndArgs.FindStringSubmatch(g); len(match) > 2 {
			if match[1] == "_" + protoName {
				return len(strings.Split(match[2], ","))
			}
		}
	}

	panic(fmt.Sprintf("Failed to find function prototype for %s", protoName))
}
