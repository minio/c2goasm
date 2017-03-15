package main

import (
	"fmt"
	"strconv"
	"strings"
	"regexp"
)

type Stack struct {
	Pushes       []string
	SetRbpIns    bool
	StackSize    int
	AlignedStack bool
	VZeroUpper   bool
}

var regexpAddRsp = regexp.MustCompile(`^\s*add\s*rsp, ([0-9]+)$`)
var regexpLeaRsp = regexp.MustCompile(`^\s*lea\s*rsp, `)
var regexpPop = regexp.MustCompile(`^\s*pop\s*([a-z0-9]+)$`)

func ExtractStackInfo(epilogue []string) Stack {

	stack := Stack{}

	// Iterate over epilogue, starting from last instruction
	for ipost := len(epilogue) - 1; ipost >= 0; ipost-- {
		line := epilogue[ipost]
		if match := regexpPop.FindStringSubmatch(line); len(match) > 1 {
			register := match[1]
			stack.Pushes = append(stack.Pushes, register)
			if register == "rbp" {
				stack.SetRbpIns = true
			}
		} else if match := regexpAddRsp.FindStringSubmatch(line); len(match) > 1 {
			stack.StackSize, _ = strconv.Atoi(match[1])
		} else if match := regexpLeaRsp.FindStringSubmatch(line); len(match) > 0 {
			stack.AlignedStack = true
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

func (s *Stack) IsStdCallPrologue(line string) bool {

	if strings.Contains(line, "push") {
		parts := strings.SplitN(line, "push", 2)
		fmt.Println("push:", parts[1])
		return true
	} else if strings.Contains(line, "mov") {
		parts := strings.SplitN(line, "mov", 2)
		argument := parts[1]
		args := strings.SplitN(argument, ",", 2)
		if strings.TrimSpace(args[0]) == "rbp" && strings.TrimSpace(args[1]) == "rsp" {
			if s.SetRbpIns {
				return true
			} else {
				panic(fmt.Sprintf("mov found but not expected to be set: %s", line))
			}
		} else {
			return false
		}
	} else if strings.Contains(line, "sub") {
		parts := strings.SplitN(line, "sub", 2)
		argument := parts[1]
		args := strings.SplitN(argument, ",", 2)
		if strings.TrimSpace(args[0]) == "rsp" {
			space, _ := strconv.Atoi(strings.TrimSpace(args[1]))
			if !s.AlignedStack && s.StackSize == space {
				return true
			} else {
				panic(fmt.Sprintf("'sub rsp' found but not for fixed stack size: %s", line))
			}
		} else {
			return false
		}
	} else {
		panic(fmt.Sprintf("Unknown line for IsStdCallHeader: %s", line))
	}

	return false
}

func IsStdCallEpilogue(line string) bool {

	if strings.Contains(line, "vzeroupper") {
		return true
	} else if strings.Contains(line, "pop") {
		return true
	} else if match := regexpAddRsp.FindStringSubmatch(line); len(match) > 0 {
		return true
	}

	return false
}
