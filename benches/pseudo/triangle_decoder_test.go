package pseudo_test

import (
	"testing"
	"time"

	"github.com/itzmeanjan/kodr/kodr_internals"
	"github.com/itzmeanjan/kodr/kodr_internals/testutils"
	"github.com/itzmeanjan/kodr/pseudo"
)

func BenchmarkTriangleRLNCDecoder(t *testing.B) {

	t.Run("500k", func(b *testing.B) {
		b.Run("2Pieces", func(b *testing.B) { decodeTriangle(b, 2, 1<<19) })
		b.Run("3Pieces", func(b *testing.B) { decodeTriangle(b, 3, 1<<19) })
		b.Run("4Pieces", func(b *testing.B) { decodeTriangle(b, 4, 1<<19) })
		b.Run("15Pieces", func(b *testing.B) { decodeTriangle(b, 15, 1<<19) })
		b.Run("77Pieces", func(b *testing.B) { decodeTriangle(b, 77, 1<<19) })
	})

	t.Run("1M", func(b *testing.B) {
		b.Run("16Pieces", func(b *testing.B) { decodeTriangle(b, 1<<4, 1<<20) })
		b.Run("32Pieces", func(b *testing.B) { decodeTriangle(b, 1<<5, 1<<20) })
		b.Run("64Pieces", func(b *testing.B) { decodeTriangle(b, 1<<6, 1<<20) })
		b.Run("128Pieces", func(b *testing.B) { decodeTriangle(b, 1<<7, 1<<20) })
		b.Run("256Pieces", func(b *testing.B) { decodeTriangle(b, 1<<8, 1<<20) })
	})

	t.Run("2M", func(b *testing.B) {
		b.Run("16Pieces", func(b *testing.B) { decodeTriangle(b, 1<<4, 1<<21) })
		b.Run("32Pieces", func(b *testing.B) { decodeTriangle(b, 1<<5, 1<<21) })
		b.Run("64Pieces", func(b *testing.B) { decodeTriangle(b, 1<<6, 1<<21) })
		b.Run("128Pieces", func(b *testing.B) { decodeTriangle(b, 1<<7, 1<<21) })
		b.Run("256Pieces", func(b *testing.B) { decodeTriangle(b, 1<<8, 1<<21) })
	})

	t.Run("16M", func(b *testing.B) {
		b.Run("16 Pieces", func(b *testing.B) { decodeTriangle(b, 1<<4, 1<<24) })
		b.Run("32 Pieces", func(b *testing.B) { decodeTriangle(b, 1<<5, 1<<24) })
		b.Run("64 Pieces", func(b *testing.B) { decodeTriangle(b, 1<<6, 1<<24) })
		b.Run("128 Pieces", func(b *testing.B) { decodeTriangle(b, 1<<7, 1<<24) })
		b.Run("256 Pieces", func(b *testing.B) { decodeTriangle(b, 1<<8, 1<<24) })
	})

	t.Run("32M", func(b *testing.B) {
		b.Run("16 Pieces", func(b *testing.B) { decodeTriangle(b, 1<<4, 1<<25) })
		b.Run("32 Pieces", func(b *testing.B) { decodeTriangle(b, 1<<5, 1<<25) })
		b.Run("64 Pieces", func(b *testing.B) { decodeTriangle(b, 1<<6, 1<<25) })
		b.Run("128 Pieces", func(b *testing.B) { decodeTriangle(b, 1<<7, 1<<25) })
		b.Run("256 Pieces", func(b *testing.B) { decodeTriangle(b, 1<<8, 1<<25) })
	})
}

func decodeTriangle(t *testing.B, pieceCount uint, total uint) {
	data := testutils.RandomData(total)

	enc, err := pseudo.NewTrianglePRLNCEncoderWithPieceCount(data, pieceCount)
	if err != nil {
		t.Fatalf("Error: %s\n", err.Error())
	}

	pieces := make([]*kodr_internals.CodedPiece, 0, 2*pieceCount)
	for range 2 * pieceCount {
		pieces = append(pieces, enc.CodedPiece())
	}

	t.ResetTimer()

	totalDuration := 0 * time.Second
	for t.Loop() {
		totalDuration += decodeLoopInner(t, pieceCount, pieces, pseudo.NewTrianglePRLNCDecoder(pieceCount))
	}

	t.ReportMetric(0, "ns/op")
	t.ReportMetric(totalDuration.Seconds()/float64(t.N), "seconds/decode")
}
