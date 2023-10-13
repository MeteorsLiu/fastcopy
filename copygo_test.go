package fastcopy

import (
	"fmt"
	"testing"
)

func TestMemmove(t *testing.T) {
	t.Parallel()
	size := 256
	if testing.Short() {
		size = 128 + 16
	}
	src := make([]byte, size)
	dst := make([]byte, size)
	for i := 0; i < size; i++ {
		src[i] = byte(128 + (i & 127))
	}
	for i := 0; i < size; i++ {
		dst[i] = byte(i & 127)
	}
	for n := 0; n <= size; n++ {
		for x := 0; x <= size-n; x++ { // offset in src
			for y := 0; y <= size-n; y++ { // offset in dst
				Copy(dst[y:y+n], src[x:x+n])
				for i := 0; i < y; i++ {
					if dst[i] != byte(i&127) {
						t.Fatalf("prefix dst[%d] = %d", i, dst[i])
					}
				}
				for i := y; i < y+n; i++ {
					if dst[i] != byte(128+((i-y+x)&127)) {
						t.Fatalf("copied dst[%d] = %d", i, dst[i])
					}
					dst[i] = byte(i & 127) // reset dst
				}
				for i := y + n; i < size; i++ {
					if dst[i] != byte(i&127) {
						t.Fatalf("suffix dst[%d] = %d", i, dst[i])
					}
				}
			}
		}
	}
}

func TestMemmoveAlias(t *testing.T) {
	t.Parallel()
	size := 256
	if testing.Short() {
		size = 128 + 16
	}
	buf := make([]byte, size)
	for i := 0; i < size; i++ {
		buf[i] = byte(i)
	}
	for n := 0; n <= size; n++ {
		for x := 0; x <= size-n; x++ { // src offset
			for y := 0; y <= size-n; y++ { // dst offset
				Copy(buf[y:y+n], buf[x:x+n])
				for i := 0; i < y; i++ {
					if buf[i] != byte(i) {
						t.Fatalf("prefix buf[%d] = %d", i, buf[i])
					}
				}
				for i := y; i < y+n; i++ {
					if buf[i] != byte(i-y+x) {
						t.Fatalf("copied buf[%d] = %d", i, buf[i])
					}
					buf[i] = byte(i) // reset buf
				}
				for i := y + n; i < size; i++ {
					if buf[i] != byte(i) {
						t.Fatalf("suffix buf[%d] = %d", i, buf[i])
					}
				}
			}
		}
	}
}

func benchmarkSizes(b *testing.B, sizes []int, fn func(b *testing.B, n int)) {
	for _, n := range sizes {
		b.Run(fmt.Sprint(n), func(b *testing.B) {
			b.SetBytes(int64(n))
			fn(b, n)
		})
	}
}

var bufSizes = []int{
	0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16,
	32, 64, 128, 256, 512, 1024, 2048, 4096,
}

var bufSizesOverlap = []int{
	32, 64, 128, 256, 512, 1024, 2048, 4096,
}

//var bufSizesOverlap = largeBufSizes

func BenchmarkMemmove(b *testing.B) {
	benchmarkSizes(b, bufSizes, func(b *testing.B, n int) {
		x := make([]byte, n)
		y := make([]byte, n)
		for i := 0; i < b.N; i++ {
			Copy(x, y)
		}
	})
}

func BenchmarkMemmoveOverlap(b *testing.B) {
	benchmarkSizes(b, bufSizesOverlap, func(b *testing.B, n int) {
		x := make([]byte, n+16)
		for i := 0; i < b.N; i++ {
			Copy(x[16:n+16], x[:n])
		}
	})
}

func BenchmarkMemmoveUnalignedDst(b *testing.B) {
	benchmarkSizes(b, bufSizes, func(b *testing.B, n int) {
		x := make([]byte, n+1)
		y := make([]byte, n)
		for i := 0; i < b.N; i++ {
			Copy(x[1:], y)
		}
	})
}

func BenchmarkMemmoveUnalignedDstOverlap(b *testing.B) {
	benchmarkSizes(b, bufSizesOverlap, func(b *testing.B, n int) {
		x := make([]byte, n+16)
		for i := 0; i < b.N; i++ {
			Copy(x[16:n+16], x[1:n+1])
		}
	})
}

func BenchmarkMemmoveUnalignedSrc(b *testing.B) {
	benchmarkSizes(b, bufSizes, func(b *testing.B, n int) {
		x := make([]byte, n)
		y := make([]byte, n+1)
		for i := 0; i < b.N; i++ {
			Copy(x, y[1:])
		}
	})
}

func BenchmarkMemmoveUnalignedSrcDst(b *testing.B) {
	for _, n := range []int{16, 64, 256, 4096, 65536} {
		buf := make([]byte, (n+8)*2)
		x := buf[:len(buf)/2]
		y := buf[len(buf)/2:]
		for _, off := range []int{0, 1, 4, 7} {
			b.Run(fmt.Sprint("f_", n, off), func(b *testing.B) {
				b.SetBytes(int64(n))
				for i := 0; i < b.N; i++ {
					Copy(x[off:n+off], y[off:n+off])
				}
			})

			b.Run(fmt.Sprint("b_", n, off), func(b *testing.B) {
				b.SetBytes(int64(n))
				for i := 0; i < b.N; i++ {
					Copy(y[off:n+off], x[off:n+off])
				}
			})
		}
	}
}

func BenchmarkMemmoveUnalignedSrcOverlap(b *testing.B) {
	benchmarkSizes(b, bufSizesOverlap, func(b *testing.B, n int) {
		x := make([]byte, n+1)
		for i := 0; i < b.N; i++ {
			Copy(x[1:n+1], x[:n])
		}
	})
}

func BenchmarkGoMemmove(b *testing.B) {
	benchmarkSizes(b, bufSizes, func(b *testing.B, n int) {
		x := make([]byte, n)
		y := make([]byte, n)
		for i := 0; i < b.N; i++ {
			copy(x, y)
		}
	})
}

func BenchmarkGoMemmoveOverlap(b *testing.B) {
	benchmarkSizes(b, bufSizesOverlap, func(b *testing.B, n int) {
		x := make([]byte, n+16)
		for i := 0; i < b.N; i++ {
			copy(x[16:n+16], x[:n])
		}
	})
}

func BenchmarkGoMemmoveUnalignedDst(b *testing.B) {
	benchmarkSizes(b, bufSizes, func(b *testing.B, n int) {
		x := make([]byte, n+1)
		y := make([]byte, n)
		for i := 0; i < b.N; i++ {
			copy(x[1:], y)
		}
	})
}

func BenchmarkGoMemmoveUnalignedDstOverlap(b *testing.B) {
	benchmarkSizes(b, bufSizesOverlap, func(b *testing.B, n int) {
		x := make([]byte, n+16)
		for i := 0; i < b.N; i++ {
			copy(x[16:n+16], x[1:n+1])
		}
	})
}

func BenchmarkGoMemmoveUnalignedSrc(b *testing.B) {
	benchmarkSizes(b, bufSizes, func(b *testing.B, n int) {
		x := make([]byte, n)
		y := make([]byte, n+1)
		for i := 0; i < b.N; i++ {
			copy(x, y[1:])
		}
	})
}

func BenchmarkGoMemmoveUnalignedSrcDst(b *testing.B) {
	for _, n := range []int{16, 64, 256, 4096, 65536} {
		buf := make([]byte, (n+8)*2)
		x := buf[:len(buf)/2]
		y := buf[len(buf)/2:]
		for _, off := range []int{0, 1, 4, 7} {
			b.Run(fmt.Sprint("f_", n, off), func(b *testing.B) {
				b.SetBytes(int64(n))
				for i := 0; i < b.N; i++ {
					copy(x[off:n+off], y[off:n+off])
				}
			})

			b.Run(fmt.Sprint("b_", n, off), func(b *testing.B) {
				b.SetBytes(int64(n))
				for i := 0; i < b.N; i++ {
					copy(y[off:n+off], x[off:n+off])
				}
			})
		}
	}
}

func BenchmarkGoMemmoveUnalignedSrcOverlap(b *testing.B) {
	benchmarkSizes(b, bufSizesOverlap, func(b *testing.B, n int) {
		x := make([]byte, n+1)
		for i := 0; i < b.N; i++ {
			copy(x[1:n+1], x[:n])
		}
	})
}
