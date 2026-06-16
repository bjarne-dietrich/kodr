package ssac

import (
	"github.com/itzmeanjan/kodr"
	"github.com/itzmeanjan/kodr/kodr_internals"
	"github.com/itzmeanjan/kodr/kodr_internals/base"
)

type SSACRLNCDecoder struct {
	base.BaseDecoder
	sparsityLevel uint
}

// Each piece of N-many byte
// Add one more collected coded piece, which will be used for decoding
// back to original pieces
//
// If all required pieces are already collected i.e. successful decoding
// has happened --- new pieces to be discarded, with an error denoting same
func (p *SSACRLNCDecoder) AddPiece(piece *kodr_internals.CodedPiece) error {
	return p.AddPieceBytes(piece.Flatten())
}

func (p *SSACRLNCDecoder) AddPieceBytes(pieceBytes []byte) error {
	if p.IsDecoded() {
		return kodr.ErrAllUsefulPiecesReceived
	}

	codedPiece := GetCodedPieceFromBytes(pieceBytes, DefaultQ0, DefaultQ1, p.GetExpectedPieceCount(), p.sparsityLevel)

	return p.BaseDecoder.AddPiece(codedPiece)
}

func NewSSACRLNCDecoder(pieceCount uint, sparsityLevel uint) *SSACRLNCDecoder {
	return &SSACRLNCDecoder{BaseDecoder: *base.NewBaseDecoder(pieceCount), sparsityLevel: sparsityLevel}
}
