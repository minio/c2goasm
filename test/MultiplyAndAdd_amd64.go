//+build !noasm
//+build !appengine

package c2goasmtest

import (
	"unsafe"
	"fmt"
)

//go:noescape
func _MultiplyAndAdd(vec1, vec2, vec3, result unsafe.Pointer)

type drijvend struct {
	a,b,c,d,e,f,g,h float32
}

func MultiplyAndAdd(f1, f2, f3 float64) float64 {

	_f1 := drijvend{}
	_f2 := drijvend{}
	_f3 := drijvend{}
	_f4 := drijvend{}

	_f1.b = 1.0
	_f2.b = 2.0
	_f3.b = 3.0

	_MultiplyAndAdd(unsafe.Pointer(&_f1), unsafe.Pointer(&_f2), unsafe.Pointer(&_f3), unsafe.Pointer(&_f4))

	fmt.Println(_f4.b)

	return 0.0
}