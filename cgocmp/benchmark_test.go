package cgocmp

import (
	"testing"
)

func BenchmarkMultiplyAndAdd(b *testing.B) {

	f1 := [8]float32{}
	f2 := [8]float32{}
	f3 := [8]float32{}

	for i := 0; i < 8; i++ {
		f1[i] = float32(i)
		f2[i] = float32(i * 2)
		f3[i] = float32(i * 3)
	}

	for i := 0; i < b.N; i++ {
		CgoMultiplyAndAdd(&f1, &f2, &f3)
	}
}
