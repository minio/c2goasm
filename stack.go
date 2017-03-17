package main

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

type Stack struct {
	Pops         []string
	SetRbpIns    bool
	StackSize    int
	AlignedStack bool
	VZeroUpper   bool
}

var regexpAddRsp = regexp.MustCompile(`^\s*add\s*rsp, ([0-9]+)$`)
var regexpAndRsp = regexp.MustCompile(`^\s*and\s*rsp, ([\-0-9]+)$`)
var regexpSubRsp = regexp.MustCompile(`^\s*sub\s*rsp, ([0-9]+)$`)
var regexpLeaRsp = regexp.MustCompile(`^\s*lea\s*rsp, `)
var regexpPop = regexp.MustCompile(`^\s*pop\s*([a-z0-9]+)$`)
var regexpPush = regexp.MustCompile(`^\s*push\s*([a-z0-9]+)$`)
var regexpMov = regexp.MustCompile(`^\s*mov\s*([a-z0-9]+), ([a-z0-9]+)$`)

func ExtractStackInfo(epilogue []string) Stack {

	stack := Stack{}

	// Iterate over epilogue, starting from last instruction
	for ipost := len(epilogue) - 1; ipost >= 0; ipost-- {
		line := epilogue[ipost]

		if !stack.ExtractEpilogue(line) {
			panic(fmt.Sprintf("Unknown line for epilogue: %s", line))
		}
	}

	return stack
}

func (stack *Stack) ExtractEpilogue(line string) bool {

	if match := regexpPop.FindStringSubmatch(line); len(match) > 1 {
		register := match[1]

		stack.Pops = append(stack.Pops, register)
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
		return false
	}

	return true
}

func IsEpilogueInstruction(line string) bool {

	return (&Stack{}).ExtractEpilogue(line)
}

func (s *Stack) IsPrologueInstruction(line string) bool {

	if match := regexpPush.FindStringSubmatch(line); len(match) > 1 {
		hasCorrespondingPop := listContains(match[1], s.Pops)
		if hasCorrespondingPop {
			return true
		} else if !hasCorrespondingPop && s.StackSize >= 8 {
			// Could not find a corresponding `pop` but rsp is modified directly (see test-case pro/epilogue6)
			s.StackSize -= 8
			return true
		} else {
			return false
		}
	} else if match := regexpMov.FindStringSubmatch(line); len(match) > 2 && match[1] == "rbp" && match[2] == "rsp" {
		if s.SetRbpIns {
			return true
		} else {
			panic(fmt.Sprintf("mov found but not expected to be set: %s", line))
		}
	} else if match := regexpAndRsp.FindStringSubmatch(line); len(match) > 1 {
		return true
	} else if match := regexpSubRsp.FindStringSubmatch(line); len(match) > 1 {
		space, _ := strconv.Atoi(match[1])
		if !s.AlignedStack && s.StackSize == space {
			return true
		} else if s.AlignedStack && s.StackSize == 0 {
			return true
		} else {
			panic(fmt.Sprintf("'sub rsp' found but in unexpected scenario: %s", line))
		}
	}

	return false
}

func listContains(value string, list []string) bool {
	for _, v := range list {
		if v == value {
			return true
		}
	}
	return false
}