package main

import (
	"testing"
)

func TestMaddArgs10(t *testing.T) {

	f1 := [8]float32{}
	f2 := [8]float32{}
	f3 := [8]float32{}
	f4 := [8]float32{}
	f5 := [8]float32{}
	f6 := [8]float32{}
	f7 := [8]float32{}
	f8 := [8]float32{}
	f9 := [8]float32{}

	for i := 0; i < 8; i++ {
		f1[i] = float32(i)
		f2[i] = float32(i * 2)
		f3[i] = float32(i * 2)
		f4[i] = float32(i * 2)
		f5[i] = float32(i * 2)
		f6[i] = float32(i * 2)
		f7[i] = float32(i * 2)
		f8[i] = float32(i * 2)
		f9[i] = float32(i * 2)
	}

	f10 := MaddArgs10(f1, f2, f3, f4, f5, f6, f7, f8, f9)

	for i := 0; i < 8; i++ {
		expected := f1[i]*f2[i]+float32(i+1)
		if f3[i] != expected {
			t.Errorf("TestMaddConstant(): \nexpected %f\ngot      %f", expected, f3[i])
		}
	}
}
