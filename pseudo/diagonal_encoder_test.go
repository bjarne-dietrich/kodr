package pseudo_test

import (
	"math"
	mathrand "math/rand"
	"testing"

	"github.com/itzmeanjan/kodr/kodr_internals"
	"github.com/itzmeanjan/kodr/kodr_internals/testutils"
	"github.com/itzmeanjan/kodr/pseudo"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewPseudoRLNC(t *testing.T) {
	t.Run("Encoder", func(t *testing.T) {
		var (
			pieceCount  uint = 1 << 8
			pieceLength uint = 8192
		)

		pieces := testutils.GeneratePieces(pieceCount, pieceLength)

		testutils.BaseTestFlow(t, pseudo.NewDiagonalPRLNCEncoder(pieces), pseudo.NewDiagonalPRLNCDecoder(pieceCount), pieceCount, pieces)
	})

	t.Run("EncoderWithPieceCount", func(t *testing.T) {
		size := uint(2<<10 + mathrand.Intn(2<<10))
		pieceCount := uint(2<<1 + mathrand.Intn(2<<8))
		data := testutils.RandomData(size)

		encoder, err := pseudo.NewDiagonalPRLNCEncoderWithPieceCount(data, pieceCount)
		require.NoError(t, err)

		pieces, _, err := kodr_internals.OriginalPiecesFromDataAndPieceCount(data, pieceCount)
		require.NoError(t, err)

		decoder := pseudo.NewDiagonalPRLNCDecoder(pieceCount)
		testutils.BaseTestFlow(t, encoder, decoder, pieceCount, pieces)
	})

	t.Run("EncoderWithPieceSize", func(t *testing.T) {
		size := uint(2<<10 + mathrand.Intn(2<<10))
		pieceSize := uint(2<<5 + mathrand.Intn(2<<5))
		pieceCount := uint(math.Ceil(float64(size) / float64(pieceSize)))
		data := testutils.RandomData(size)

		encoder, err := pseudo.NewDiagonalPRLNCEncoderWithPieceSize(data, pieceSize)
		require.NoError(t, err)

		pieces, _, err := kodr_internals.OriginalPiecesFromDataAndPieceSize(data, pieceSize)
		require.NoError(t, err)

		decoder := pseudo.NewDiagonalPRLNCDecoder(pieceCount)
		testutils.BaseTestFlow(t, encoder, decoder, pieceCount, pieces)
	})
}

func TestPseudoRLNCEncoder_Padding(t *testing.T) {
	t.Run("WithPieceCount", func(t *testing.T) {
		for range 1 << 5 {
			size := uint(2<<10 + mathrand.Intn(2<<10))
			pieceCount := uint(2<<1 + mathrand.Intn(2<<8))
			data := testutils.RandomData(size)

			encoder, err := pseudo.NewDiagonalPRLNCEncoderWithPieceCount(data, pieceCount)
			require.NoError(t, err)

			extra := encoder.Padding()
			pieceSize := (size + extra) / pieceCount
			codedPiece := encoder.CodedPiece()
			assert.Equal(t, pieceSize, uint(len(codedPiece.Piece)))
		}
	})

	t.Run("WithPieceSize", func(t *testing.T) {
		for range 1 << 5 {
			size := uint(2<<10 + mathrand.Intn(2<<10))
			pieceSize := uint(2<<5 + mathrand.Intn(2<<5))
			pieceCount := uint(math.Ceil(float64(size) / float64(pieceSize)))
			data := testutils.RandomData(size)

			enc, err := pseudo.NewDiagonalPRLNCEncoderWithPieceSize(data, pieceSize)
			require.NoError(t, err)

			extra := enc.Padding()
			codedPieceSize := (size + extra) / pieceCount
			codedPiece := enc.CodedPiece()
			assert.Equal(t, pieceSize, codedPieceSize)
			assert.Equal(t, pieceSize, uint(len(codedPiece.Piece)))
		}
	})
}
