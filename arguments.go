package main

import (
	"strings"
)

func InsertArguments() []string {

	str := ` 	// Load golang arguments into respective registers for standard 64-bit function call interface
 	MOVQ arg1+0(FP), DI  // rdi = arg1
 	MOVQ arg2+8(FP), SI  // rsi = arg2
 	MOVQ arg3+16(FP), DX // rdx = arg3
 	MOVQ arg4+24(FP), CX // rcx = arg4
 	MOVQ arg5+32(FP), R8 // r8 = arg5
 	MOVQ arg6+40(FP), R9 // r9 = arg6

 	// Setup base pointer for loading constants
 	LEAQ LCDATA1<>(SB), BP`

	return strings.Split(str, "\n")
}

func Aap() {

}
