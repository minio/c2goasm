//+build !noasm
//+build !appengine

package c2goasmtest

import (
	"unsafe"
)

//go:noescape
func _MaddArgs10(vec1, vec2, vec3, vec4, vec5, vec6, vec7, vec8, vec9, result unsafe.Pointer)

func MaddArgs10(f1, f2, f3, f4, f5, f6, f7, f8, f9 [8]float32) [8]float32 {

	_f10 := [8]float32{}

	_MaddArgs10(unsafe.Pointer(&f1), unsafe.Pointer(&f2), unsafe.Pointer(&f3), unsafe.Pointer(&f4), unsafe.Pointer(&f5), unsafe.Pointer(&f6), unsafe.Pointer(&f7), unsafe.Pointer(&f8), unsafe.Pointer(&f9), unsafe.Pointer(&_f10))

	return _f10
}
