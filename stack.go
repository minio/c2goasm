package main

import (
	"fmt"
	"strings"
	"strconv"
)

type Stack struct {
	Pushes []string
	SetRbp bool
	FixedSize int
	DynamicSize bool
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
		} else {
			panic(fmt.Sprintf("Unknown line for postamble: %s", line))
		}
	}

	return stack
}

