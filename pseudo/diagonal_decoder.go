package pseudo

import (
	"github.com/itzmeanjan/kodr"
	"github.com/itzmeanjan/kodr/kodr_internals"
	"github.com/itzmeanjan/kodr/kodr_internals/base"
)

type DiagonalPRLNCDecoder struct {
	base.BaseDecoder
}

// AddPiece adds a kodr_internals.CodedPiece to the decoder.
func (p *DiagonalPRLNCDecoder) AddPiece(piece *kodr_internals.CodedPiece) error {
	return p.AddPieceBytes(piece.Flatten())
}

// AddPieceBytes tries to parse a kodr_internals.CodedPiece from piceBytes
// and adds it to the decoder.
func (p *DiagonalPRLNCDecoder) AddPieceBytes(pieceBytes []byte) error {
	if p.IsDecoded() {
		return kodr.ErrAllUsefulPiecesReceived
	}
	codedPiece := GetDiagonalCodedPieceFromBytes(pieceBytes, p.GetExpectedPieceCount())
	return p.BaseDecoder.AddPiece(codedPiece)
}

func NewDiagonalPRLNCDecoder(pieceCount uint) *DiagonalPRLNCDecoder {
	return &DiagonalPRLNCDecoder{*base.NewBaseDecoder(pieceCount)}
}
