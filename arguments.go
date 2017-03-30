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

func argumentsOnStack(lines []string) StackArgs {

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
	for o := range offsets {
		if o < offset {
			offset = o
		}
	}
	return StackArgs{Number: len(offsets), OffsetToFirst: int(offset)}
}

func parseCompanionFile(goCompanion, protoName string) int {

	gocode, err := readLines(goCompanion)
	if err != nil {
		panic(fmt.Sprintf("Failed to read companion go code: %v", err))
	}

	for _, goline := range gocode {

		ok, args := getGolangArgs(protoName, goline)
		if ok {
			return args
		}
	}

	panic(fmt.Sprintf("Failed to find function prototype for %s", protoName))
}

var regexpFuncAndArgs = regexp.MustCompile(`func\s+(.*)\((.*)\)`)

func getGolangArgs(protoName, goline string) (bool, int) {

	// Search for name of function and arguments
	if match := regexpFuncAndArgs.FindStringSubmatch(goline); len(match) > 2 {
		if match[1] == "_"+protoName {
			return true, len(strings.Split(match[2], ","))
		}
	}

	return false, 0
}

func getTotalSizeOfArguments(argStart, argEnd int) int {
	// TODO: Test if correct for non 64-bit arguments
	return (argEnd - argStart + 1)*8
}