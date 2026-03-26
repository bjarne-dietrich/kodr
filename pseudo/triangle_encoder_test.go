package pseudo_test

import (
	"bytes"
	"errors"
	"math"
	mathrand "math/rand"
	"testing"

	"github.com/itzmeanjan/kodr"
	"github.com/itzmeanjan/kodr/kodr_internals"
	"github.com/itzmeanjan/kodr/kodr_internals/testutils"
	"github.com/itzmeanjan/kodr/pseudo"
)

func TestNewTrianglePseudoRLNC(t *testing.T) {
	t.Run("Encoder", func(t *testing.T) {
		var (
			pieceCount  uint = 1 << 8
			pieceLength uint = 8192
		)

		pieces := testutils.GeneratePieces(pieceCount, pieceLength)
		enc := pseudo.NewTrianglePRLNCEncoder(pieces)
		dec := pseudo.NewTrianglePRLNCDecoder(pieceCount)

		encoderFlow(t, enc, dec, pieceCount, pieces)
	})

	t.Run("EncoderWithPieceCount", func(t *testing.T) {
		size := uint(2<<10 + mathrand.Intn(2<<10))
		pieceCount := uint(2<<1 + mathrand.Intn(2<<8))
		data := testutils.RandomData(size)

		enc, err := pseudo.NewTrianglePRLNCEncoderWithPieceCount(data, pieceCount)
		if err != nil {
			t.Fatalf("Error: %s\n", err.Error())
		}

		pieces, _, err := kodr_internals.OriginalPiecesFromDataAndPieceCount(data, pieceCount)
		if err != nil {
			t.Fatal(err.Error())
		}

		dec := pseudo.NewTrianglePRLNCDecoder(pieceCount)
		encoderFlow(t, enc, dec, pieceCount, pieces)
	})

	t.Run("EncoderWithPieceSize", func(t *testing.T) {
		size := uint(2<<10 + mathrand.Intn(2<<10))
		pieceSize := uint(2<<5 + mathrand.Intn(2<<5))
		pieceCount := uint(math.Ceil(float64(size) / float64(pieceSize)))
		data := testutils.RandomData(size)

		enc, err := pseudo.NewTrianglePRLNCEncoderWithPieceSize(data, pieceSize)
		if err != nil {
			t.Fatalf("Error: %s\n", err.Error())
		}

		pieces, _, err := kodr_internals.OriginalPiecesFromDataAndPieceSize(data, pieceSize)
		if err != nil {
			t.Fatal(err.Error())
		}

		dec := pseudo.NewTrianglePRLNCDecoder(pieceCount)
		encoderFlow(t, enc, dec, pieceCount, pieces)
	})
}

func encoderFlow(t *testing.T, enc *pseudo.TrianglePRLNCEncoder, dec *pseudo.TrianglePRLNCDecoder, pieceCount uint, pieces []kodr_internals.Piece) {
	for {
		c_piece := enc.CodedPiece()

		if mathrand.Intn(2) == 0 {
			continue
		}

		if err := dec.AddPiece(c_piece); err != nil && errors.Is(err, kodr.ErrAllUsefulPiecesReceived) {
			break
		}
	}

	d_pieces, err := dec.GetPieces()
	if err != nil {
		t.Fatal(err.Error())
	}

	if len(pieces) != len(d_pieces) {
		t.Fatal("didn't decode all !")
	}

	for i := range pieceCount {
		if !bytes.Equal(pieces[i], d_pieces[i]) {
			t.Fatal("decoded data doesn't match !")
		}
	}
}

func TestTrianglePseudoRLNCEncoder_Padding(t *testing.T) {
	t.Run("WithPieceCount", func(t *testing.T) {
		for range 1 << 5 {
			size := uint(2<<10 + mathrand.Intn(2<<10))
			pieceCount := uint(2<<1 + mathrand.Intn(2<<8))
			data := testutils.RandomData(size)

			enc, err := pseudo.NewTrianglePRLNCEncoderWithPieceCount(data, pieceCount)
			if err != nil {
				t.Fatalf("Error: %s\n", err.Error())
			}

			extra := enc.Padding()
			pieceSize := (size + extra) / pieceCount
			c_piece := enc.CodedPiece()
			if uint(len(c_piece.Piece)) != pieceSize {
				t.Fatalf("expected pieceSize to be %dB, found to be %dB\n", pieceSize, len(c_piece.Piece))
			}
		}
	})

	t.Run("WithPieceSize", func(t *testing.T) {
		for range 1 << 5 {
			size := uint(2<<10 + mathrand.Intn(2<<10))
			pieceSize := uint(2<<5 + mathrand.Intn(2<<5))
			pieceCount := uint(math.Ceil(float64(size) / float64(pieceSize)))
			data := testutils.RandomData(size)

			enc, err := pseudo.NewTrianglePRLNCEncoderWithPieceSize(data, pieceSize)
			if err != nil {
				t.Fatalf("Error: %s\n", err.Error())
			}

			extra := enc.Padding()
			c_pieceSize := (size + extra) / pieceCount
			c_piece := enc.CodedPiece()
			if pieceSize != c_pieceSize || uint(len(c_piece.Piece)) != pieceSize {
				t.Fatalf("expected pieceSize to be %dB, found to be %dB\n", c_pieceSize, len(c_piece.Piece))
			}
		}
	})
}
