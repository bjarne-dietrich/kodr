package pseudo

import (
	"encoding/binary"

	"github.com/itzmeanjan/kodr/kodr_internals"
	"github.com/itzmeanjan/kodr/kodr_internals/base"
	"github.com/itzmeanjan/kodr/kodr_internals/operations"
)

type DiagonalPRLNCEncoder struct {
	base.BaseEncoder
}

// CodedPiece returns subsequent coded pieces
// For diagonal pseudo coding, first N-piece are returned in uncoded form
// Instead of the coding Vector, CodedPiece().Vector contains the index of the piece as an Uvarint
func (p *DiagonalPRLNCEncoder) CodedPiece() *kodr_internals.CodedPiece {

	pieceID := p.GetCurrentPieceId()
	pieceCount := p.PieceCount()

	vector := DiagonalCodingVector(pieceID, pieceCount)
	piece := make(kodr_internals.Piece, p.PieceSize())

	// This could be omitted, but it should be a bit faster for large PieceCounts
	// else path could always be taken
	if pieceID < pieceCount {
		copy(piece, *p.GetPiece(pieceID))
	} else {
		for i := range pieceCount {
			if vector[i] == 1 {
				operations.XorAssignSlice(piece, *p.GetPiece(i))
			}
		}
	}

	codedPiece := &kodr_internals.CodedPiece{
		Vector: binary.AppendUvarint(nil, uint64(pieceID)),
		Piece:  piece,
	}
	p.IncrementCurrentPieceId()
	return codedPiece
}

// NewDiagonalPRLNCEncoder can be used when you've already split original data chunk
// into pieces of same length ( in terms of bytes ), and returns a DiagonalPRLNCEncoder,
// which delivers coded pieces on-the-fly
func NewDiagonalPRLNCEncoder(pieces []kodr_internals.Piece) *DiagonalPRLNCEncoder {
	return &DiagonalPRLNCEncoder{BaseEncoder: *base.NewBaseEncoder(pieces)}
}

// NewDiagonalPRLNCEncoderWithPieceCount returns a DiagonalPRLNCEncoder
// and splits the data into pieceCount same sized pieces,
// appending zero-padding to data if needed.
func NewDiagonalPRLNCEncoderWithPieceCount(data []byte, pieceCount uint) (*DiagonalPRLNCEncoder, error) {
	encoder, err := base.NewBaseEncoderWithPieceCount(data, pieceCount)
	if err != nil {
		return nil, err
	}
	return &DiagonalPRLNCEncoder{*encoder}, nil
}

// NewDiagonalPRLNCEncoderWithPieceSize returns a DiagonalPRLNCEncoder
// and splits the data into pieces with a size of pieceSize,
// appending zero-padding to data if needed.
func NewDiagonalPRLNCEncoderWithPieceSize(data []byte, pieceSize uint) (*DiagonalPRLNCEncoder, error) {
	encoder, err := base.NewBaseEncoderWithPieceSize(data, pieceSize)
	if err != nil {
		return nil, err
	}
	return &DiagonalPRLNCEncoder{*encoder}, nil
}
