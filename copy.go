package fastcopy

import (
	"unsafe"
	_ "unsafe"

	"github.com/klauspost/cpuid/v2"
	"golang.org/x/exp/constraints"
)

type (
	CanMove interface {
		constraints.Integer | constraints.Float | constraints.Complex | ~string | uintptr
	}
)

var (
	// for debug: hasERMS = false
	hasERMS = cpuid.CPU.Has(cpuid.ERMS)
	isX64   = cpuid.CPU.X64Level() != 0
)

//go:linkname memmove runtime.memmove
func memmove(to, from unsafe.Pointer, n uintptr)

func copy_movsb(to, from unsafe.Pointer, n uintptr)
func copy_movsq(to, from unsafe.Pointer, n uintptr) (left, copied uintptr)

func Copy[T CanMove](dst, src []T) (n int) {
	n = len(src)
	if n > len(dst) {
		n = len(dst)
	}
	elem := unsafe.Sizeof(src[0])
	size := elem * uintptr(n)

	if size > 15500 && (hasERMS || isX64) {
		if hasERMS {
			copy_movsb(
				unsafe.Pointer(&dst[0]),
				unsafe.Pointer(&src[0]),
				size,
			)
		} else {
			left, copied := copy_movsq(
				unsafe.Pointer(&dst[0]),
				unsafe.Pointer(&src[0]),
				size,
			)
			if left > 0 {
				memmove(
					unsafe.Pointer(&dst[copied]),
					unsafe.Pointer(&src[copied]),
					left*elem,
				)
			}
		}
	} else {
		memmove(
			unsafe.Pointer(&dst[0]),
			unsafe.Pointer(&src[0]),
			size,
		)
	}
	return
}
