package pseudo_test

import (
	"testing"

	"github.com/itzmeanjan/kodr/kodr_internals/testutils"
	"github.com/itzmeanjan/kodr/pseudo"
)

func TestDiagonalPseudoRLNCDecoderFlow(t *testing.T) {
	pieceCount := uint(128)
	pieceLength := uint(8192)
	pieces := testutils.GeneratePieces(pieceCount, pieceLength)

	testutils.BaseTestFlow(t, pseudo.NewDiagonalPRLNCEncoder(pieces), pseudo.NewDiagonalPRLNCDecoder(pieceCount), pieceCount, pieces)
}
