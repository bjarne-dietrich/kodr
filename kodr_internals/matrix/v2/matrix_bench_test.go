package matrix_test

import (
	"crypto/rand"
	"testing"

	"github.com/itzmeanjan/kodr/kodr_internals/matrix/v2"
)

// Note: If fillWithZero is set, it's not really a random matrix
func randomMatrix(rows, cols int, fillWithZero bool) [][]byte {
	mat := make([][]byte, 0, rows)

	for range rows {
		row := make([]byte, cols)
		// already filled with zero
		if !fillWithZero {
			_, _ = rand.Read(row)
		}
		mat = append(mat, row)
	}
	return mat
}

func BenchmarkMatrixReduceAndPrune(b *testing.B) {
	b.Run("16x16", func(b *testing.B) { reduceAndPrune(b, 1<<4) })
	b.Run("32x32", func(b *testing.B) { reduceAndPrune(b, 1<<5) })
	b.Run("64x64", func(b *testing.B) { reduceAndPrune(b, 1<<6) })
	b.Run("128x128", func(b *testing.B) { reduceAndPrune(b, 1<<7) })
	b.Run("256x256", func(b *testing.B) { reduceAndPrune(b, 1<<8) })
	b.Run("512x512", func(b *testing.B) { reduceAndPrune(b, 1<<9) })
	b.Run("1024x1024", func(b *testing.B) { reduceAndPrune(b, 1<<10) })
	b.Run("2048x2048", func(b *testing.B) { reduceAndPrune(b, 1<<11) })
	b.Run("4096x4096", func(b *testing.B) { reduceAndPrune(b, 1<<12) })
}

func reduceAndPrune(b *testing.B, dim int) {
	b.SetBytes(int64(dim*dim) << 1)
	b.ReportAllocs()

	for b.Loop() {
		b.StopTimer()
		coefficients := randomMatrix(dim, dim, false)
		coded := randomMatrix(dim, dim, true)
		decoderState := matrix.NewDecoderState(coefficients, coded)
		b.StartTimer()

		decoderState.ReduceAndPrune()
	}
}

func BenchmarkMatrixMultiply(b *testing.B) {
	benchMultiplyDim := func(dim int) {

		b.Run("alloc", func(b *testing.B) {
			b.SetBytes(int64(2 * dim * dim))
			b.ReportAllocs()
			for b.Loop() {
				b.StopTimer()
				a := randomMatrix(dim, dim, false)
				x := randomMatrix(dim, dim, false)
				A := matrix.Matrix(a)
				b.StartTimer()

				_, err := A.Multiply(x)
				if err != nil {
					b.Fatal(err)
				}
			}
		})

		b.Run("reuse", func(b *testing.B) {
			b.SetBytes(int64(2 * dim * dim))
			b.ReportAllocs()
			// Inputs created once; still measures Multiply's internal allocations.
			a := randomMatrix(dim, dim, false)
			x := randomMatrix(dim, dim, false)
			A := matrix.Matrix(a)

			b.ResetTimer()
			for b.Loop() {
				_, err := A.Multiply(x)
				if err != nil {
					b.Fatal(err)
				}
			}
		})
	}

	b.Run("16x16", func(b *testing.B) { benchMultiplyDim(1 << 4) })
	b.Run("32x32", func(b *testing.B) { benchMultiplyDim(1 << 5) })
	b.Run("64x64", func(b *testing.B) { benchMultiplyDim(1 << 6) })
	b.Run("128x128", func(b *testing.B) { benchMultiplyDim(1 << 7) })
	b.Run("256x256", func(b *testing.B) { benchMultiplyDim(1 << 8) })
	b.Run("512x512", func(b *testing.B) { benchMultiplyDim(1 << 9) })
	b.Run("1024x1024", func(b *testing.B) { benchMultiplyDim(1 << 10) })
	b.Run("2048x2048", func(b *testing.B) { benchMultiplyDim(1 << 11) })
	b.Run("4096x4096", func(b *testing.B) { benchMultiplyDim(1 << 12) })
}
