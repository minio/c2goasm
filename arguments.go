package main

import (
	"fmt"
)

// Write the prologue for the subroutine
func WriteGoasmPrologue(number int) []string {

	registers := []string{"DI", "SI", "DX", "CX", "R8", "R9"}

	var result []string
	for i, reg := range registers {

		result = append(result, fmt.Sprintf("    MOVQ arg%d+%d(FP), %s", i+1, i*8, reg))
		if i +1 == number {
			break
		}
	}

	result = append(result, "", fmt.Sprintf("    LEAQ LCDATA1<>(SB), BP  // Setup base pointer for loading constants"), "")

	return result
}
