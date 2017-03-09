package main

import (
	"testing"
)

func testName(t *testing.T, fullname, expected string) {
	name := ExtractName(fullname)
	if name != expected {
		t.Errorf("TestNames(): \nexpected %s\ngot      %s", expected, name)
	}
}

func TestNames(t *testing.T) {

	testName(t, "_ZN4Simd4Avx213Yuv444pToBgraEPKhmS2_mS2_mmmPhmh", "SimdAvx2Yuv444pToBgra")
	testName(t, "_ZN4Simd4Avx213Yuv420pToBgraEPKhmS2_mS2_mmmPhmh", "SimdAvx2Yuv420pToBgra")
	testName(t, "_ZN4Simd4Avx213Yuv422pToBgraEPKhmS2_mS2_mmmPhmh", "SimdAvx2Yuv422pToBgra")
	testName(t, "_ZN4Simd4Avx213ReduceGray2x2EPKhmmmPhmmm", "SimdAvx2ReduceGray2x2")
	testName(t, "_ZN4Simd4Avx216AbsDifferenceSumEPKhmS2_mmmPy", "SimdAvx2AbsDifferenceSum")
}
