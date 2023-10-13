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

// implement in memmove_movsb.s. prevent converting from unsafe.Pointer to uintptr .
func can_movs(to, from unsafe.Pointer, n uintptr) (ok bool)
func copy_movsb(to, from unsafe.Pointer, n uintptr)
func copy_movsq(to, from unsafe.Pointer, n uintptr) (left, copied uintptr)

func Copy[T CanMove](dst, src []T) (n int) {
	n = min(len(src), len(dst))
	if n == 0 {
		return
	}
	dstPtr := unsafe.Pointer(&dst[0])
	srcPtr := unsafe.Pointer(&src[0])
	if srcPtr == dstPtr {
		n = 0
		return
	}
	nptr := uintptr(n)
	elem := unsafe.Sizeof(src[0])
	size := elem * nptr
	if size > 15500 && (hasERMS || isX64) && can_movs(dstPtr, srcPtr, nptr) {
		if hasERMS {
			copy_movsb(dstPtr, srcPtr, size)
		} else {
			left, copied := copy_movsq(dstPtr, srcPtr, size)
			if left > 0 {
				memmove(
					unsafe.Pointer(&dst[copied]),
					unsafe.Pointer(&src[copied]),
					left*elem,
				)
			}
		}
	} else {
		memmove(dstPtr, srcPtr, size)
	}
	return
}

func CopyMOVSB[T CanMove](dst, src []T) (n int) {
	n = min(len(src), len(dst))
	if n == 0 {
		return
	}
	dstPtr := unsafe.Pointer(&dst[0])
	srcPtr := unsafe.Pointer(&src[0])
	if srcPtr == dstPtr {
		n = 0
		return
	}
	elem := unsafe.Sizeof(src[0])
	nptr := uintptr(n)
	size := elem * nptr
	if can_movs(dstPtr, srcPtr, nptr) {
		copy_movsb(dstPtr, srcPtr, size)
	} else {
		memmove(dstPtr, srcPtr, size)
	}
	return
}

func CopyMOVSQ[T CanMove](dst, src []T) (n int) {
	n = min(len(src), len(dst))
	if n == 0 {
		return
	}
	dstPtr := unsafe.Pointer(&dst[0])
	srcPtr := unsafe.Pointer(&src[0])
	if srcPtr == dstPtr {
		n = 0
		return
	}
	elem := unsafe.Sizeof(src[0])
	nptr := uintptr(n)
	size := elem * nptr
	if can_movs(dstPtr, srcPtr, nptr) {
		left, copied := copy_movsq(
			dstPtr,
			srcPtr,
			size,
		)
		if left > 0 {
			memmove(
				unsafe.Pointer(&dst[copied]),
				unsafe.Pointer(&src[copied]),
				left*elem,
			)
		}
	} else {
		memmove(dstPtr, srcPtr, size)
	}
	return
}
