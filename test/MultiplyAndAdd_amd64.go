//+build !noasm
//+build !appengine

package c2goasmtest

import (
	"unsafe"
)

//go:noescape
func _MultiplyAndAdd(vec1, vec2, vec3, result unsafe.Pointer)

func MultiplyAndAdd(f1, f2, f3 float64) float64 {

	_f1 := f1
	_f2 := f2
	_f3 := f3
	_f4 := float64(0.0)
	_MultiplyAndAdd(unsafe.Pointer(&_f1), unsafe.Pointer(&_f2), unsafe.Pointer(&_f3), unsafe.Pointer(&_f4))

	return _f4
}