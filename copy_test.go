package fastcopy

import (
	"testing"
	"time"
)

const (
	GB   = 1024 * 1024 * 1024
	SIZE = 32768
)

func getLog(t *testing.T, name string, written uint64, last time.Time) {
	b := written / uint64(time.Since(last).Seconds())
	t.Log(name+" Output: ", b/GB, "Gb/s ", b, "B/s")
}

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
	src := make([]byte, 16383)
	dst := make([]byte, 16383)

	for i := 0; i < len(src); i++ {
		src[i] = byte(i)
	}

	t.Log(hasERMS, isX64, Copy(dst, src))

	checkDstByte(t, dst)
}

// 8 byte
func TestCopyInt(t *testing.T) {
	src := make([]int, 16383)
	dst := make([]int, 16383)

	for i := 0; i < len(src); i++ {
		src[i] = i
	}

	t.Log(hasERMS, isX64, Copy(dst, src))

	checkDstInt(t, dst)
}

// 4 byte
func TestCopyFloat32(t *testing.T) {
	src := make([]float32, 16383)
	dst := make([]float32, 16383)

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

func TestOutput(t *testing.T) {
	src := make([]byte, SIZE)
	for i := 0; i < len(src); i++ {
		src[i] = 1
	}
	dst := make([]byte, SIZE)
	zero := make([]byte, SIZE)
	written := uint64(0)
	now := time.Now()

	for time.Since(now).Seconds() < 11 {
		Copy(dst, src)
		written += SIZE

		Copy(dst, zero)
		written += SIZE

	}

	getLog(t, "optimized copy", written, now)

	clear(dst)
	written = uint64(0)
	now = time.Now()

	for time.Since(now).Seconds() < 11 {
		written += uint64(copy(dst, src))

		written += uint64(copy(dst, zero))
	}

	getLog(t, "copy", written, now)

}