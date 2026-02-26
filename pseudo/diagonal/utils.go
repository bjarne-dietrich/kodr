package diagonal

import (
	"encoding/binary"

	"github.com/itzmeanjan/kodr/kodr_internals"
	"github.com/itzmeanjan/kodr/pseudo"
)

// CodingVector returns a diagonal coding vector for the given index
// it has a size of pieceCount. If idx < pieceCount the vector will
// only have one non-zero element which is at vector[idx].
// Otherwise, is will return what pseudo.NonPrimitiveCodingVector returns.
func CodingVector(idx uint, pieceCount uint) kodr_internals.CodingVector {

	// Cidx,j = {1 if j == idx, 0 otherwise}
	// C0 (1 0 0 0 ...)
	// C1 (0 1 0 0 ...)
	if idx < pieceCount {
		vector := make(kodr_internals.CodingVector, pieceCount)
		vector[idx] = 1
		return vector
	} else {
		return pseudo.NonPrimitiveCodingVector(idx, pieceCount)
	}

}

// GetCodedPieceFromBytes expects a diagonal coded piece as bytes
// It parses the index, calculates the coding vector and returns the
// kodr_internals.CodedPiece with the remaining bytes as piece
func GetCodedPieceFromBytes(bytes []byte, pieceCount uint) *kodr_internals.CodedPiece {
	index, n := binary.Uvarint(bytes)
	if n <= 0 {
		return nil
	}

	return &kodr_internals.CodedPiece{
		Vector: CodingVector(uint(index), pieceCount),
		Piece:  bytes[n:],
	}
}
