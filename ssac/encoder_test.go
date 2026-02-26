package ssac_test

import (
	"bytes"
	"crypto/rand"
	"errors"
	"math"
	mathrand "math/rand"
	"testing"

	"github.com/itzmeanjan/kodr"
	"github.com/itzmeanjan/kodr/kodr_internals"
	"github.com/itzmeanjan/kodr/ssac"
)

// Generates `N`-bytes of random data from default
// randomization source
func generateData(n uint) []byte {
	data := make([]byte, n)
	_, _ = rand.Read(data)
	return data
}

// Generates N-many pieces each of M-bytes length, to be used
// for testing purposes
func generatePieces(pieceCount uint, pieceLength uint) []kodr_internals.Piece {
	pieces := make([]kodr_internals.Piece, 0, pieceCount)
	for range pieceCount {
		pieces = append(pieces, generateData(pieceLength))
	}
	return pieces
}

func TestNewSSACRLNC(t *testing.T) {
	t.Run("Encoder", func(t *testing.T) {
		var (
			pieceCount  uint = 16
			pieceLength uint = 8192
		)

		pieces := generatePieces(pieceCount, pieceLength)
		enc := ssac.NewSSACRLNCEncoder(pieces)
		dec := ssac.NewSSACRLNCDecoder(pieceCount)

		encoderFlow(t, enc, dec, pieceCount, pieces)
	})

	t.Run("EncoderWithPieceCount", func(t *testing.T) {
		size := uint(2<<10 + mathrand.Intn(2<<10))
		pieceCount := uint(2<<1 + mathrand.Intn(2<<8))
		data := generateData(size)

		enc, err := ssac.NewSSACRLNCEncoderWithPieceCount(data, pieceCount)
		if err != nil {
			t.Fatalf("Error: %s\n", err.Error())
		}

		pieces, _, err := kodr_internals.OriginalPiecesFromDataAndPieceCount(data, pieceCount)
		if err != nil {
			t.Fatal(err.Error())
		}

		dec := ssac.NewSSACRLNCDecoder(pieceCount)
		encoderFlow(t, enc, dec, pieceCount, pieces)
	})

	t.Run("EncoderWithPieceSize", func(t *testing.T) {
		size := uint(2<<10 + mathrand.Intn(2<<10))
		pieceSize := uint(2<<5 + mathrand.Intn(2<<5))
		pieceCount := uint(math.Ceil(float64(size) / float64(pieceSize)))
		data := generateData(size)

		enc, err := ssac.NewSSACRLNCEncoderWithPieceSize(data, pieceSize)
		if err != nil {
			t.Fatalf("Error: %s\n", err.Error())
		}

		pieces, _, err := kodr_internals.OriginalPiecesFromDataAndPieceSize(data, pieceSize)
		if err != nil {
			t.Fatal(err.Error())
		}

		dec := ssac.NewSSACRLNCDecoder(pieceCount)
		encoderFlow(t, enc, dec, pieceCount, pieces)
	})
}

func encoderFlow(t *testing.T, enc *ssac.SSACRLNCEncoder, dec *ssac.SSACRLNCDecoder, pieceCount uint, pieces []kodr_internals.Piece) {
	for {
		codedPiece := enc.CodedPiece()

		if mathrand.Intn(2) == 0 {
			continue
		}

		if err := dec.AddPiece(codedPiece); err != nil && errors.Is(err, kodr.ErrAllUsefulPiecesReceived) {
			break
		}
	}

	decodedPieces, err := dec.GetPieces()
	if err != nil {
		t.Fatal(err.Error())
	}

	if len(pieces) != len(decodedPieces) {
		t.Fatal("didn't decode all !")
	}

	for i := range pieceCount {
		if !bytes.Equal(pieces[i], decodedPieces[i]) {
			t.Fatal("decoded data doesn't match !")
		}
	}
}

func TestPseudoRLNCEncoder_Padding(t *testing.T) {
	t.Run("WithPieceCount", func(t *testing.T) {
		for range 1 << 5 {
			size := uint(2<<10 + mathrand.Intn(2<<10))
			pieceCount := uint(2<<1 + mathrand.Intn(2<<8))
			data := generateData(size)

			enc, err := ssac.NewSSACRLNCEncoderWithPieceCount(data, pieceCount)
			if err != nil {
				t.Fatalf("Error: %s\n", err.Error())
			}

			extra := enc.Padding()
			pieceSize := (size + extra) / pieceCount
			codedPiece := enc.CodedPiece()
			if uint(len(codedPiece.Piece)) != pieceSize {
				t.Fatalf("expected pieceSize to be %dB, found to be %dB\n", pieceSize, len(codedPiece.Piece))
			}
		}
	})

	t.Run("WithPieceSize", func(t *testing.T) {
		for range 1 << 5 {
			size := uint(2<<10 + mathrand.Intn(2<<10))
			pieceSize := uint(2<<5 + mathrand.Intn(2<<5))
			pieceCount := uint(math.Ceil(float64(size) / float64(pieceSize)))
			data := generateData(size)

			enc, err := ssac.NewSSACRLNCEncoderWithPieceSize(data, pieceSize)
			if err != nil {
				t.Fatalf("Error: %s\n", err.Error())
			}

			extra := enc.Padding()
			codedPieceSize := (size + extra) / pieceCount
			codedPiece := enc.CodedPiece()
			if pieceSize != codedPieceSize || uint(len(codedPiece.Piece)) != pieceSize {
				t.Fatalf("expected pieceSize to be %dB, found to be %dB\n", codedPieceSize, len(codedPiece.Piece))
			}
		}
	})
}
