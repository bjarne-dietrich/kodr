package pseudo

import (
	"encoding/binary"
	"math"
	"slices"

	"github.com/itzmeanjan/kodr/kodr_internals"
)

// DiagonalCodingVector returns a diagonal coding vector for the given index
// it has a size of pieceCount. If idx < pieceCount the vector will
// only have one non-zero element which is at vector[idx].
// Otherwise, is will return what pseudo.NonPrimitiveCodingVector returns.
func DiagonalCodingVector(idx uint, pieceCount uint) kodr_internals.CodingVector {

	// C_idx,j = {1 if j == idx, 0 otherwise}
	// C0 (1 0 0 0 ...)
	// C1 (0 1 0 0 ...)
	if idx < pieceCount {
		vector := make(kodr_internals.CodingVector, pieceCount)
		vector[idx] = 1
		return vector
	} else {
		return nonPrimitiveCodingVector(idx, pieceCount)
	}

}

// GetDiagonalCodedPieceFromBytes expects a diagonal coded piece as bytes
// It parses the index, calculates the coding vector and returns the
// kodr_internals.CodedPiece with the remaining bytes as piece
func GetDiagonalCodedPieceFromBytes(bytes []byte, pieceCount uint) *kodr_internals.CodedPiece {
	index, n := binary.Uvarint(bytes)
	if n <= 0 {
		return nil
	}

	return &kodr_internals.CodedPiece{
		Vector: DiagonalCodingVector(uint(index), pieceCount),
		Piece:  bytes[n:],
	}
}

// TriangleCodingVector returns a triangle coding vector for the given index
// it has a size of pieceCount. If idx < pieceCount the vector will
// only have idx non-zero elements from vector[0] to vector[idx].
// Otherwise, is will return what pseudo.NonPrimitiveCodingVector returns.
func TriangleCodingVector(idx uint, pieceCount uint) kodr_internals.CodingVector {

	// C[idx,j] = {1 if j <= idx, 0 otherwise}
	// C0 (1 0 0 0 ...)
	// C1 (1 1 0 0 ...)
	if idx < pieceCount {
		vector := make(kodr_internals.CodingVector, pieceCount)
		for i := 0; i <= int(idx); i++ {
			vector[i] = 1
		}
		return vector
	} else {
		return nonPrimitiveCodingVector(idx, pieceCount)
	}
}

// GetTriangleCodedPieceFromBytes expects a triangle coded piece as bytes
// It parses the index, calculates the coding vector and returns the
// kodr_internals.CodedPiece with the remaining bytes as piece
func GetTriangleCodedPieceFromBytes(bytes []byte, pieceCount uint) *kodr_internals.CodedPiece {
	index, n := binary.Uvarint(bytes)
	if n <= 0 {
		return nil
	}

	return &kodr_internals.CodedPiece{
		Vector: TriangleCodingVector(uint(index), pieceCount),
		Piece:  bytes[n:],
	}
}

// nonPrimitiveCodingVector returns a pseudo coding vector.
func nonPrimitiveCodingVector(idx uint, pieceCount uint) kodr_internals.CodingVector {

	// with n being p.PieceCount:
	// Cn,j = {1}
	// Cn (1 1 1 1 ...)
	if idx == pieceCount {
		return slices.Repeat(kodr_internals.CodingVector{1}, int(pieceCount))
	}

	// a ≥ 2, 0 ≤ k ≤ a-2 such that i = N+2+((a-1)*(a-2))/2 + k
	// but dont use (a,b) where b = a-1 because of linear dependence

	N := pieceCount

	a := uint(math.Floor(1.5 + math.Sqrt(0.25+2*float64(idx-N-1))))
	k := idx - N - 1 - (a-1)*(a-2)/2

	vector := make(kodr_internals.CodingVector, pieceCount)
	for j := range vector {
		if j%int(a) == int(k) {
			vector[j] = 1
		}
	}

	return vector
}
