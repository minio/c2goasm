package c2goasmtest

import (
	"testing"
)

func TestMultiplyAndAddConstant(t *testing.T) {

	f1 := [8]float32{}
	f2 := [8]float32{}

	for i := 0; i < 8; i++ {
		f1[i] = float32(i)
		f2[i] = float32(i * 2)
	}

	f3 := MultiplyAndAddConstant(f1, f2)

	for i := 0; i < 8; i++ {
		expected := f1[i]*f2[i]+float32(i)
		if f3[i] != expected {
			t.Errorf("TestMultiplyAndAddConstant(): \nexpected %s\ngot      %s", expected, f3[i])
		}
	}
}
