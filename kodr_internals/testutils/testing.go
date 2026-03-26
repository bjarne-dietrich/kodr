package testutils

import (
	"crypto/rand"
	"errors"
	mathrand "math/rand"
	"testing"

	"github.com/itzmeanjan/kodr"
	"github.com/itzmeanjan/kodr/kodr_internals"
	"github.com/itzmeanjan/kodr/kodr_internals/base"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func BaseTestFlow(t *testing.T, encoder base.Encoder, decoder base.Decoder, pieceCount uint, referencePieces []kodr_internals.Piece) {
	for {
		codedPiece := encoder.CodedPiece()

		if mathrand.Intn(2) == 0 {
			continue
		}

		if err := decoder.AddPiece(codedPiece); err != nil && errors.Is(err, kodr.ErrAllUsefulPiecesReceived) {
			break
		}
	}

	decodedPieces, err := decoder.GetPieces()
	require.NoError(t, err)
	assert.Equal(t, len(referencePieces), len(decodedPieces), "didnt decode at all!")
	assert.Equal(t, referencePieces, decodedPieces, "didnt decode at all!")
}

// RandomData returns `N`-bytes of random data
func RandomData(n uint) []byte {
	data := make([]byte, n)
	// can safely ignore error
	_, err := rand.Read(data)
	if err != nil {
		panic(err)
	}
	return data
}

// GeneratePieces returns N-many random pieces each of M-bytes length, to be used
// for testing purposes
func GeneratePieces(pieceCount uint, pieceLength uint) []kodr_internals.Piece {
	pieces := make([]kodr_internals.Piece, 0, pieceCount)
	for range pieceCount {
		pieces = append(pieces, RandomData(pieceLength))
	}
	return pieces
}
