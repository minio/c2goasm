package c2goasmtest

import (
	"testing"
)

func TestMaddConstant(t *testing.T) {

	f1 := [8]float32{}
	f2 := [8]float32{}

	for i := 0; i < 8; i++ {
		f1[i] = float32(i)
		f2[i] = float32(i * 2)
	}

	f3 := MaddConstant(f1, f2)

	for i := 0; i < 8; i++ {
		expected := f1[i]*f2[i] + float32(i+1)
		if f3[i] != expected {
			t.Errorf("TestMaddConstant(): \nexpected %f\ngot      %f", expected, f3[i])
		}
	}
}
