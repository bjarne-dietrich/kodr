package systematic

import (
	"github.com/itzmeanjan/kodr"
	"github.com/itzmeanjan/kodr/kodr_internals"
	"github.com/itzmeanjan/kodr/kodr_internals/base"
)

type SystematicRLNCDecoder struct {
	base.BaseDecoder
}

// AddPieceBytes tries to parse a kodr_internals.CodedPiece from piceBytes
// and adds it to the decoder.
func (d *SystematicRLNCDecoder) AddPieceBytes(pieceBytes []byte) error {
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

// Pieces coded by systematic mean, along with randomly coded pieces,
// are decoded with this decoder
//
// @note Actually FullRLNCDecoder could have been used for same purpose
// making this one redundant
//
// I'll consider improving decoding by exploiting
// systematic coded pieces ( vectors )/ removing this
// in some future date
func NewSystematicRLNCDecoder(pieceCount uint) *SystematicRLNCDecoder {
	return &SystematicRLNCDecoder{*base.NewBaseDecoder(pieceCount)}
}
