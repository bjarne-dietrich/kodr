package pseudo

import (
	"github.com/itzmeanjan/kodr"
	"github.com/itzmeanjan/kodr/kodr_internals"
	"github.com/itzmeanjan/kodr/kodr_internals/base"
)

type TrianglePRLNCDecoder struct {
	base.BaseDecoder
}

// AddPiece adds a kodr_internals.CodedPiece to the decoder.
func (p *TrianglePRLNCDecoder) AddPiece(piece *kodr_internals.CodedPiece) error {
	return p.AddPieceBytes(piece.Flatten())
}

// AddPieceBytes tries to parse a kodr_internals.CodedPiece from piceBytes
// and adds it to the decoder.
func (p *TrianglePRLNCDecoder) AddPieceBytes(pieceBytes []byte) error {
	if p.IsDecoded() {
		return kodr.ErrAllUsefulPiecesReceived
	}
	codedPiece := GetTriangleCodedPieceFromBytes(pieceBytes, p.GetExpectedPieceCount())
	return p.BaseDecoder.AddPiece(codedPiece)
}

func NewTrianglePRLNCDecoder(pieceCount uint) *TrianglePRLNCDecoder {
	return &TrianglePRLNCDecoder{*base.NewBaseDecoder(pieceCount)}
}
