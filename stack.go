package main

import (
	"fmt"
	"strconv"
	"strings"
)

type Stack struct {
	Pushes      []string
	SetRbp      bool
	FixedSize   int
	DynamicSize bool
	VZeroUpper  bool
}

func ExtractStackInfo(postamble []string) Stack {

	stack := Stack{}

	// Iterate over postable, starting from last instruction
	for ipost := len(postamble) - 1; ipost >= 0; ipost-- {
		line := postamble[ipost]
		if strings.Contains(line, "pop") {
			parts := strings.SplitN(line, "pop", 2)
			register := strings.TrimSpace(parts[1])

			stack.Pushes = append(stack.Pushes, register)
			if register == "rbp" {
				stack.SetRbp = true
			}
		} else if strings.Contains(line, "add") {
			parts := strings.SplitN(line, "add", 2)
			argument := parts[1]
			args := strings.SplitN(argument, ",", 2)

			if strings.TrimSpace(args[0]) == "rsp" {
				stack.FixedSize, _ = strconv.Atoi(strings.TrimSpace(args[1]))
			} else {
				panic(fmt.Sprintf("Unexpected add statement for postamble: %s", line))
			}
		} else if strings.Contains(line, "lea") {
			parts := strings.SplitN(line, "lea", 2)
			argument := parts[1]
			args := strings.SplitN(argument, ",", 2)

			if strings.TrimSpace(args[0]) == "rsp" {
				stack.DynamicSize = true
			} else {
				panic(fmt.Sprintf("Unexpected add statement for postamble: %s", line))
			}
		} else if strings.Contains(line, "vzeroupper") {
			stack.VZeroUpper = true
		} else if strings.Contains(line, "ret") {
			// no action to take
		} else {
			panic(fmt.Sprintf("Unknown line for postamble: %s", line))
		}
	}

	return stack
}

func (s *Stack) CreateGoPreable() []string {

	var result []string

	if s.FixedSize != 0 {
		result = append(result, fmt.Sprintf("    SUB $%s, SP", s.FixedSize))
	}

	return result
}

func (s *Stack) CreateGoPostable() []string {

	var result []string

	if s.FixedSize != 0 {
		result = append(result, fmt.Sprintf("    ADD $%s, SP", s.FixedSize))
	}
	if s.VZeroUpper {
		result = append(result, "    VZEROUPPER")
	}
	result = append(result, "    RET")

	return result
}
