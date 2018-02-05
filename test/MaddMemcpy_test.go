package c2goasmtest

import (
	"testing"
)

func TestMaddMemcpy(t *testing.T) {

	f1 := [8]float32{}
	f2 := [8]float32{}
	f3 := [8]float32{}

	for i := 0; i < 8; i++ {
		f1[i] = float32(i)
		f2[i] = float32(i * 2)
		f3[i] = float32(i * 3)
	}

	f4 := MaddMemcpy(&f1, &f2, &f3, 6*4, 7*4)

	for i := 0; i < 8; i++ {
		expected := float32(0)
		if i < 6 {
			expected = f1[i]*f1[i] + f1[i]
		} else if i < 7 {
			expected = f1[i]*f2[i] + f1[i]
		} else {
			expected = f1[i]*f2[i] + f3[i]
		}
		if f4[i] != expected {
			t.Errorf("TestMaddMemcpy(): \nexpected %f\ngot      %f", expected, f4[i])
		}
	}
}
