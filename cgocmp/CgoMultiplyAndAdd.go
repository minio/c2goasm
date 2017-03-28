package cgocmp

// #include "cpp/MultiplyAndAdd.h"
// #cgo CFLAGS: -mavx -mfma
import "C"
import (
	"unsafe"
)

func CgoMultiplyAndAdd(f1, f2, f3 *[8]float32) [8]float32 {

	_f4 := [8]float32{}

	C.MultiplyAndAdd((*C.float)(unsafe.Pointer(f1)), (*C.float)(unsafe.Pointer(f2)), (*C.float)(unsafe.Pointer(f3)), (*C.float)(unsafe.Pointer(&_f4)))

	return _f4
}
