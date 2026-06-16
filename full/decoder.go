package full

import (
	"github.com/itzmeanjan/kodr"
	"github.com/itzmeanjan/kodr/kodr_internals"
	"github.com/itzmeanjan/kodr/kodr_internals/base"
)

type FullRLNCDecoder struct {
	base.BaseDecoder
}

// AddPieceBytes tries to parse a kodr_internals.CodedPiece from piceBytes
// and adds it to the decoder.
func (d *FullRLNCDecoder) AddPieceBytes(pieceBytes []byte) error {
	if d.IsDecoded() {
		return kodr.ErrAllUsefulPiecesReceived
	}

	expected := d.GetExpectedPieceCount()
	pieceLength := d.PieceLength()

	if pieceLength != 0 {
		if len(pieceBytes) != int(pieceLength)+int(expected) {
			return kodr.ErrCodedDataLengthMismatch
		}
	}

	codedPiece := &kodr_internals.CodedPiece{Vector: pieceBytes[:expected], Piece: pieceBytes[expected:]}

	return d.AddPiece(codedPiece)
}

// If minimum #-of linearly independent coded pieces required
// for decoding coded pieces --- is provided with,
// it returns a decoder, which keeps applying
// full RLNC decoding step on received coded pieces
//
// As soon as minimum #-of linearly independent pieces are obtained
// which is generally equal to original #-of pieces, decoded pieces
// can be read back
func NewFullRLNCDecoder(pieceCount uint) *FullRLNCDecoder {
	return &FullRLNCDecoder{*base.NewBaseDecoder(pieceCount)}
}
