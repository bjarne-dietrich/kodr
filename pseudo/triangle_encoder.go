package pseudo

import (
	"encoding/binary"

	"github.com/itzmeanjan/kodr/kodr_internals"
	"github.com/itzmeanjan/kodr/kodr_internals/base"
	"github.com/itzmeanjan/kodr/kodr_internals/operations"
)

type TrianglePRLNCEncoder struct {
	base.BaseEncoder
}

// CodedPiece returns subsequent coded pieces
// For triangle pseudo coding, first N-piece are returned in uncoded form
// Instead of the coding Vector, CodedPiece().Vector contains the index of the piece as an Uvarint
func (p *TrianglePRLNCEncoder) CodedPiece() *kodr_internals.CodedPiece {

	pieceID := p.GetCurrentPieceIdAndIncrement()
	pieceCount := p.PieceCount()

	vector := TriangleCodingVector(pieceID, pieceCount)
	piece := make(kodr_internals.Piece, p.PieceSize())

	if pieceID < pieceCount {
		for i := range pieceID + 1 {
			operations.XorAssignSlice(piece, *p.GetPiece(i))
		}
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
	return codedPiece
}

// NewTrianglePRLNCEncoder can be used when you've already split original data chunk
// into pieces of same length ( in terms of bytes ), and returns a TrianglePRLNCEncoder,
// which delivers coded pieces on-the-fly
func NewTrianglePRLNCEncoder(pieces []kodr_internals.Piece) *TrianglePRLNCEncoder {
	return &TrianglePRLNCEncoder{BaseEncoder: *base.NewBaseEncoder(pieces)}
}

// NewTrianglePRLNCEncoderWithPieceCount returns a TrianglePRLNCEncoder
// and splits the data into pieceCount same sized pieces,
// appending zero-padding to data if needed.
func NewTrianglePRLNCEncoderWithPieceCount(data []byte, pieceCount uint) (*TrianglePRLNCEncoder, error) {
	encoder, err := base.NewBaseEncoderWithPieceCount(data, pieceCount)
	if err != nil {
		return nil, err
	}
	return &TrianglePRLNCEncoder{*encoder}, nil
}

// NewTrianglePRLNCEncoderWithPieceSize returns a TrianglePRLNCEncoder
// and splits the data into pieces with a size of pieceSize,
// appending zero-padding to data if needed.
func NewTrianglePRLNCEncoderWithPieceSize(data []byte, pieceSize uint) (*TrianglePRLNCEncoder, error) {
	encoder, err := base.NewBaseEncoderWithPieceSize(data, pieceSize)
	if err != nil {
		return nil, err
	}
	return &TrianglePRLNCEncoder{*encoder}, nil
}
