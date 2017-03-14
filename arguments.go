package main

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

func ArgumentsOnStack(lines []string) int {

	regexpRbpLoadHigher := regexp.MustCompile(`\[rbp \+ ([0-9]+)\]$`)

	offsets := make(map[int]bool)

	for _, l := range lines {

		if match := regexpRbpLoadHigher.FindStringSubmatch(l); len(match) > 1 {
			offset, _ := strconv.Atoi(match[1])
			if _, found := offsets[offset]; !found {
				offsets[offset] = true
			}
		}
	}

	fmt.Println(offsets)

	return len(offsets)
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
