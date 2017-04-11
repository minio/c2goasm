package main

import (
	"fmt"
	"errors"
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

	offset := uint(0)
	for o := range offsets {
		if o > offset {
			offset = o
		}
	}
	if offset >= 16 {
		return StackArgs{OffsetToFirst: 16, Number: 1+int((offset-16)/8)}
	}
	return StackArgs{OffsetToFirst: 0, Number: 0}
}

func parseCompanionFile(goCompanion, protoName string) ([]string, []string) {

	gocode, err := readLines(goCompanion)
	if err != nil {
		panic(fmt.Sprintf("Failed to read companion go code: %v", err))
	}

	for _, goline := range gocode {

		ok, args, rets, err := getGolangArgs(protoName, goline)
		if err != nil {
			panic(fmt.Sprintf("Error: %v", err))
		} else if ok {
			return args, rets
		}
	}

	panic(fmt.Sprintf("Failed to find function prototype for %s", protoName))
}

var regexpFuncAndArgs = regexp.MustCompile(`^\s*func\s+([^\(]*)\(([^\)]*)\)(.*)`)
var regexpReturnVals = regexp.MustCompile(`^\((.*)\)`)

func getGolangArgs(protoName, goline string) (isFunc bool, args, rets []string, err error) {

	// Search for name of function and arguments
	if match := regexpFuncAndArgs.FindStringSubmatch(goline); len(match) > 2 {
		if match[1] == "_"+protoName {

			args, rets = []string{}, []string{}
			for _, arg := range strings.Split(match[2], ",") {
				args = append(args, strings.Fields(arg)[0])
			}

			trailer := strings.TrimSpace(match[3])
			if len(trailer) > 0 {
				// Trailing string found, search for return values
				if rmatch := regexpReturnVals.FindStringSubmatch(trailer); len(rmatch) > 1 {
					for _, ret := range strings.Split(rmatch[1], ",") {
						rets = append(rets, strings.Fields(ret)[0])
					}
				} else {
					return false, args, rets, errors.New(fmt.Sprintf("Badly formatted return argument (please use parenthesis and proper arguments naming): %s", trailer))
				}

			}

			return true, args, rets, nil
		}
	}

	return false, []string{}, []string{}, nil
}

func getTotalSizeOfArguments(argStart, argEnd int) uint {
	// TODO: Test if correct for non 64-bit arguments
	return uint((argEnd - argStart + 1)*8)
}

func getTotalSizeOfArgumentsAndReturnValues(argStart, argEnd int, returnValues []string) uint {
	// TODO: Test if correct for non 64-bit return values
	return getTotalSizeOfArguments(argStart, argEnd) + uint(len(returnValues)*8)
}

