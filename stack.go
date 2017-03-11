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

func (s *Stack) IsStdCallHeader(line string) (bool, string) {

	if strings.Contains(line, "push") {
		parts := strings.SplitN(line, "push", 2)
		fmt.Println("push:", parts[1])
		return true, ""
	} else if strings.Contains(line, "mov") {
		parts := strings.SplitN(line, "mov", 2)
		argument := parts[1]
		args := strings.SplitN(argument, ",", 2)
		if strings.TrimSpace(args[0]) == "rbp" && strings.TrimSpace(args[1]) == "rsp" {
			if s.SetRbp {
				return true, ""
			} else {
				panic(fmt.Sprintf("mov found but not expected to be set: %s", line))
			}
		} else {
			return false, ""
		}
	} else if strings.Contains(line, "sub") {
		parts := strings.SplitN(line, "sub", 2)
		argument := parts[1]
		args := strings.SplitN(argument, ",", 2)
		if strings.TrimSpace(args[0]) == "rsp" {
			space, _ := strconv.Atoi(strings.TrimSpace(args[1]))
			if s.FixedSize != 0 && s.FixedSize == space {
				return true, fmt.Sprintf("    SUB $%d, SP", s.FixedSize)
			} else {
				panic(fmt.Sprintf("'sub rsp' found but not for fixed stack size: %s", line))
			}
		} else {
			return false, ""
		}
	} else {
		panic(fmt.Sprintf("Unknown line for IsStdCallHeader: %s", line))
	}

	return false, ""
}

func IsStdCallFooter(line string) bool {

	if strings.Contains(line, "vzeroupper") {
		return true
	} else if strings.Contains(line, "pop") {
		return true
	}

	return false
}

func (s *Stack) Return() []string {

	var result []string

	if s.FixedSize != 0 {
		result = append(result, fmt.Sprintf("    ADD $%d, SP", s.FixedSize))
	}
	if s.VZeroUpper {
		result = append(result, "    VZEROUPPER")
	}
	result = append(result, "    RET")

	return result
}
