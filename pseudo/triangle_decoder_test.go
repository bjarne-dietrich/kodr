package pseudo_test

import (
	"testing"

	"github.com/itzmeanjan/kodr/kodr_internals/testutils"
	"github.com/itzmeanjan/kodr/pseudo"
)

func TestTrianglePseudoRLNCDecoderFlow(t *testing.T) {
	pieceCount := uint(128)
	pieceLength := uint(8192)
	pieces := testutils.GeneratePieces(pieceCount, pieceLength)

	testutils.BaseTestFlow(t, pseudo.NewTrianglePRLNCEncoder(pieces), pseudo.NewTrianglePRLNCDecoder(pieceCount), pieceCount, pieces)
}

func TestTrianglePseudoRLNCDecoderFlowSmall(t *testing.T) {
	var pieceCount uint = 4
	var pieceLength uint = 15
	pieces := testutils.GeneratePieces(pieceCount, pieceLength)

	testutils.BaseTestFlow(t, pseudo.NewTrianglePRLNCEncoder(pieces), pseudo.NewTrianglePRLNCDecoder(pieceCount), pieceCount, pieces)
}
