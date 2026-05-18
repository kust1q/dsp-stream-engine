package dsp

/*
#cgo CFLAGS: -I../../cpp
#cgo LDFLAGS: -L../../cpp/build -ldsp -lstdc++
#include "dsp.h"
*/
import "C"
import (
	"unsafe"
)

// BiquadFilter represents a stateful biquad filter.
type BiquadFilter struct {
	handle C.BiquadHandle
}

// NewBiquadFilter creates a new biquad filter with the given coefficients.
func NewBiquadFilter(b0, b1, b2, a1, a2 float32) *BiquadFilter {
	handle := C.Biquad_Create(C.float(b0), C.float(b1), C.float(b2), C.float(a1), C.float(a2))
	return &BiquadFilter{handle: handle}
}

// Process applies the filter to the provided buffer in-place.
func (f *BiquadFilter) Process(buffer []float32) {
	if len(buffer) == 0 {
		return
	}
	C.Biquad_Process(f.handle, (*C.float)(unsafe.Pointer(&buffer[0])), C.int(len(buffer)))
}

// Close frees the underlying C++ filter instance.
func (f *BiquadFilter) Close() {
	if f.handle != nil {
		C.Biquad_Destroy(f.handle)
		f.handle = nil
	}
}
