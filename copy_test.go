package fastcopy

import (
	"testing"
)

const (
	SIZE = 32768
)

func checkDstByte(t *testing.T, dst []byte) {
	for n, i := range dst {
		if i != byte(n) {
			t.Error("move fail")
			return
		}
	}
}

func checkDstInt(t *testing.T, dst []int) {
	for n, i := range dst {
		if i != n {
			t.Error("move fail")
			return
		}
	}
}

func checkDstFloat32(t *testing.T, dst []float32) {
	for n, i := range dst {
		if i != float32(n) {
			t.Error("move fail")
			return
		}
	}
}

// 1 byte
func TestCopy(t *testing.T) {
	src := make([]byte, SIZE)
	dst := make([]byte, SIZE)

	for i := 0; i < len(src); i++ {
		src[i] = byte(i)
	}

	t.Log(hasERMS, isX64, Copy(dst, src))

	checkDstByte(t, dst)
}

func TestCopySlice(t *testing.T) {
	dst := make([]int, 2000)
	src := make([]int, 100)

	for i := 0; i < len(src); i++ {
		src[i] = i
	}
	var n int
	for i := 0; i < 20; i++ {
		n += Copy(dst[n:], src)
	}
	t.Log(hasERMS, isX64, dst)

}

// 8 byte
func TestCopyInt(t *testing.T) {
	src := make([]int, SIZE)
	dst := make([]int, SIZE)

	for i := 0; i < len(src); i++ {
		src[i] = i
	}

	t.Log(hasERMS, isX64, Copy(dst, src))

	checkDstInt(t, dst)
}

// 4 byte
func TestCopyFloat32(t *testing.T) {
	src := make([]float32, SIZE)
	dst := make([]float32, SIZE)

	for i := 0; i < len(src); i++ {
		src[i] = float32(i)
	}

	t.Log(hasERMS, isX64, Copy(dst, src))

	checkDstFloat32(t, dst)
}

func BenchmarkCopy(b *testing.B) {
	src := make([]int, 32768)
	dst := make([]int, 32768)

	for i := 0; i < len(src); i++ {
		src[i] = i
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		Copy(dst, src)
	}
}

func BenchmarkGoCopy(b *testing.B) {
	src := make([]int, 32768)
	dst := make([]int, 32768)

	for i := 0; i < len(src); i++ {
		src[i] = i
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		copy(dst, src)
	}
}

var largeBufSizes = []int{
	8192, 16384, 32768,
}

func BenchmarkCopyOutput(b *testing.B) {
	benchmarkSizes(b, largeBufSizes, func(b *testing.B, n int) {
		x := make([]byte, n)
		y := make([]byte, n)
		for i := 0; i < b.N; i++ {
			Copy(x, y)
		}
	})
}

func BenchmarkCopyGoOutput(b *testing.B) {
	benchmarkSizes(b, largeBufSizes, func(b *testing.B, n int) {
		x := make([]byte, n)
		y := make([]byte, n)
		for i := 0; i < b.N; i++ {
			copy(x, y)
		}
	})
}
