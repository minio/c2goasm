//+build !noasm
//+build !appengine

package c2goasmtest

import (
	"unsafe"
)

//go:noescape
func _MaddMemcpy(vec1, vec2, vec3 unsafe.Pointer, size1, size2 uint64, result unsafe.Pointer)

func MaddMemcpy(f1, f2, f3 *[8]float32, size1, size2 uint64) [8]float32 {

	_f4 := [8]float32{}

	_MaddMemcpy(unsafe.Pointer(f1), unsafe.Pointer(f2), unsafe.Pointer(f3), size1, size2, unsafe.Pointer(&_f4))

	return _f4
}
