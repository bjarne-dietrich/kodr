package pseudo

import (
	"math"
	"slices"

	"github.com/itzmeanjan/kodr/kodr_internals"
)

// NonPrimitiveCodingVector returns a pseudo coding and is not to be used by hand.
// You would want to use diagonal.CodingVector or triangle.CodingVector
func NonPrimitiveCodingVector(idx uint, pieceCount uint) kodr_internals.CodingVector {

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
