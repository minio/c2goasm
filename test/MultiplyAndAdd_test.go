package c2goasmtest

import (
	"testing"
)

func TestMultiplyAndAdd(t *testing.T) {

	f1 := [8]float32{}
	f2 := [8]float32{}
	f3 := [8]float32{}

	for i := 0; i < 8; i++ {
		f1[i] = float32(i)
		f2[i] = float32(i * 2)
		f3[i] = float32(i * 3)
	}

	f4 := MultiplyAndAdd(&f1, &f2, &f3)

	for i := 0; i < 8; i++ {
		expected := f1[i]*f2[i] + f3[i]
		if f4[i] != expected {
			t.Errorf("TestMultiplyAndAdd(): \nexpected %f\ngot      %f", expected, f4[i])
		}
	}
}
