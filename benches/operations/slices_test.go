package operations_test

import (
	"math/rand"
	"testing"

	"github.com/itzmeanjan/kodr/kodr_internals"
	"github.com/itzmeanjan/kodr/kodr_internals/operations"
)

func randBytes(r *rand.Rand, n int) []byte {
	b := make([]byte, n)
	_, _ = r.Read(b)
	return b
}

func BenchmarkPiece(b *testing.B) {
	b.Run("PieceSize 2^4", func(b *testing.B) {
		pieceSize := 1 << 4
		piecesAll(b, pieceSize)
	})
	b.Run("PieceSize 2^8", func(b *testing.B) {
		pieceSize := 1 << 8
		piecesAll(b, pieceSize)
	})
	b.Run("PieceSize 2^12", func(b *testing.B) {
		pieceSize := 1 << 12
		piecesAll(b, pieceSize)
	})
	b.Run("PieceSize 2^16", func(b *testing.B) {
		pieceSize := 1 << 16
		piecesAll(b, pieceSize)
	})
	b.Run("PieceSize 2^20", func(b *testing.B) {
		pieceSize := 1 << 20
		piecesAll(b, pieceSize)
	})
	b.Run("PieceSize 2^24", func(b *testing.B) {
		pieceSize := 1 << 24
		piecesAll(b, pieceSize)
	})
}

func piecesAll(b *testing.B, pieceSize int) {
	b.Run("GF256 Multiply", func(b *testing.B) { pieceMultiply(b, pieceSize) })
	b.Run("Mul Add Const", func(b *testing.B) { executionWrapper(b, pieceSize, operations.MulAddConst) })
	b.Run("Mul Add Const Table", func(b *testing.B) {
		executionWrapper(b, pieceSize, func(dst, src []byte, c byte) {
			table := operations.BuildMulTable(c)
			operations.MulAddConstTable(dst, src, &table)
		})
	})
	b.Run("Mul Add Const Table Cached", func(b *testing.B) { mulAddConstTableCached(b, pieceSize) })
	b.Run("Nibble Splitting", func(b *testing.B) { executionWrapper(b, pieceSize, operations.MulAddConstNibble) })

}

func pieceMultiply(b *testing.B, pieceSize int) {

	dst := make(kodr_internals.Piece, pieceSize)
	r := rand.New(rand.NewSource(42))
	dst0 := randBytes(r, pieceSize)
	c := byte(7)

	data := randBytes(r, pieceSize)
	copy(dst, dst0)

	b.ReportAllocs()
	b.SetBytes(int64(pieceSize * 2))
	b.ResetTimer()

	for b.Loop() {
		dst.Multiply(data, c)
	}
}

func mulAddConstTableCached(b *testing.B, pieceSize int) {
	r := rand.New(rand.NewSource(42))
	dst0 := randBytes(r, pieceSize)
	c := byte(7)

	data := randBytes(r, pieceSize)
	table := operations.BuildMulTable(c)

	b.ReportAllocs()
	b.SetBytes(int64(pieceSize * 2))
	b.ResetTimer()

	for b.Loop() {
		operations.MulAddConstTable(dst0, data, &table)
	}
}

func executionWrapper(b *testing.B, pieceSize int, loop func(dst, src []byte, c byte)) {
	r := rand.New(rand.NewSource(42))
	dst0 := randBytes(r, pieceSize)
	c := byte(7)

	data := randBytes(r, pieceSize)

	b.ReportAllocs()
	b.SetBytes(int64(pieceSize * 2))

	for b.Loop() {
		loop(dst0, data, c)
	}
}
