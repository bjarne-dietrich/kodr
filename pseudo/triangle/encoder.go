package triangle

import (
	"encoding/binary"

	"github.com/itzmeanjan/kodr/kodr_internals"
	"github.com/itzmeanjan/kodr/kodr_internals/operations"
)

type TrianglePRLNCEncoder struct {
	currentPieceId uint
	pieces         []kodr_internals.Piece
	extra          uint
}

// PieceCount returns total #-of pieces being coded together --- denoting
// these many linearly independent pieces are required
// successfully decoding back to original pieces
func (p *TrianglePRLNCEncoder) PieceCount() uint {
	return uint(len(p.pieces))
}

// PieceSize returns size of one piece
// Total data being coded = pieceSize * pieceCount + padding
func (p *TrianglePRLNCEncoder) PieceSize() uint {
	return uint(len(p.pieces[0]))
}

// Padding returns the number of extra padding bytes added at end of original
// data slice for making all pieces of same size
func (p *TrianglePRLNCEncoder) Padding() uint {
	return p.extra
}

// CodedPiece returns subsequent coded pieces
// For triangle pseudo coding, first N-piece are returned in uncoded form
// Instead of the coding Vector, CodedPiece().Vector contains the index of the piece as an Uvarint
func (p *TrianglePRLNCEncoder) CodedPiece() *kodr_internals.CodedPiece {

	vector := CodingVector(p.currentPieceId, p.PieceCount())
	piece := make(kodr_internals.Piece, p.PieceSize())

	if p.currentPieceId < p.PieceCount() {
		for i := range p.currentPieceId + 1 {
			operations.XorAssignSlice(piece, p.pieces[i])
		}
	} else {
		for i := range p.pieces {
			if vector[i] == 1 {
				operations.XorAssignSlice(piece, p.pieces[i])
			}
		}
	}

	codedPiece := &kodr_internals.CodedPiece{
		Vector: binary.AppendUvarint(nil, uint64(p.currentPieceId)),
		Piece:  piece,
	}
	p.currentPieceId++
	return codedPiece
}

// NewTrianglePRLNCEncoder can be used when you've already split original data chunk
// into pieces of same length ( in terms of bytes ), and returns a TrianglePRLNCEncoder,
// which delivers coded pieces on-the-fly
func NewTrianglePRLNCEncoder(pieces []kodr_internals.Piece) *TrianglePRLNCEncoder {
	return &TrianglePRLNCEncoder{currentPieceId: 0, pieces: pieces}
}

// NewTrianglePRLNCEncoderWithPieceCount returns a TrianglePRLNCEncoder
// and splits the data into pieceCount same sized pieces,
// appending zero-padding to data if needed.
func NewTrianglePRLNCEncoderWithPieceCount(data []byte, pieceCount uint) (*TrianglePRLNCEncoder, error) {
	pieces, padding, err := kodr_internals.OriginalPiecesFromDataAndPieceCount(data, pieceCount)
	if err != nil {
		return nil, err
	}

	enc := NewTrianglePRLNCEncoder(pieces)
	enc.extra = padding
	return enc, nil
}

// NewTrianglePRLNCEncoderWithPieceSize returns a TrianglePRLNCEncoder
// and splits the data into pieces with a size of pieceSize,
// appending zero-padding to data if needed.
func NewTrianglePRLNCEncoderWithPieceSize(data []byte, pieceSize uint) (*TrianglePRLNCEncoder, error) {
	pieces, padding, err := kodr_internals.OriginalPiecesFromDataAndPieceSize(data, pieceSize)
	if err != nil {
		return nil, err
	}

	enc := NewTrianglePRLNCEncoder(pieces)
	enc.extra = padding
	return enc, nil
}
