package matrix

import (
	"testing"

	"github.com/itzmeanjan/kodr"
	"github.com/stretchr/testify/assert"
)

func TestMatrixReduceAndPrune(t *testing.T) {
	{
		matrix := Matrix{{70, 137, 2, 152}, {223, 92, 234, 98}, {217, 141, 33, 44}, {145, 135, 71, 45}}
		reducedMatrix := Matrix{{1, 0, 0, 105}, {0, 1, 0, 181}, {0, 0, 1, 42}}
		coded := Matrix{{0, 0, 0, 0}, {0, 0, 0, 0}, {0, 0, 0, 0}, {0, 0, 0, 0}}

		dec := NewDecoderState(matrix, coded)
		dec.ReduceAndPrune()
		assert.Equal(t, reducedMatrix, dec.coefficientMatrix, "Reduced matrix doesn't match!")
	}

	{
		matrix := Matrix{{68, 54, 6, 230}, {16, 56, 215, 78}, {159, 186, 146, 163}, {122, 41, 205, 133}}
		reducedMatrix := Matrix{{1, 0, 0, 0}, {0, 1, 0, 0}, {0, 0, 1, 0}, {0, 0, 0, 1}}
		coded := Matrix{{0, 0, 0, 0}, {0, 0, 0, 0}, {0, 0, 0, 0}, {0, 0, 0, 0}}

		dec := NewDecoderState(matrix, coded)
		dec.ReduceAndPrune()
		assert.Equal(t, reducedMatrix, dec.coefficientMatrix, "Reduced matrix doesn't match!")
	}

	{
		matrix := Matrix{{100, 31, 76, 199, 119}, {207, 34, 207, 208, 18}, {62, 20, 54, 6, 187}, {66, 8, 52, 73, 54}, {122, 138, 247, 211, 165}}
		reducedMatrix := Matrix{{1, 0, 0, 0, 0}, {0, 1, 0, 0, 0}, {0, 0, 1, 0, 0}, {0, 0, 0, 1, 0}, {0, 0, 0, 0, 1}}
		coded := Matrix{{0, 0, 0, 0}, {0, 0, 0, 0}, {0, 0, 0, 0}, {0, 0, 0, 0}, {0, 0, 0, 0}}

		dec := NewDecoderState(matrix, coded)
		dec.ReduceAndPrune()
		assert.Equal(t, reducedMatrix, dec.coefficientMatrix, "Reduced matrix doesn't match!")
	}
}

func TestMatrixRank(t *testing.T) {
	{
		matrix := Matrix{{70, 137, 2, 152}, {223, 92, 234, 98}, {217, 141, 33, 44}, {145, 135, 71, 45}}
		coded := Matrix{{0, 0, 0, 0}, {0, 0, 0, 0}, {0, 0, 0, 0}, {0, 0, 0, 0}}
		assert.Equal(t, uint(3), NewDecoderState(matrix, coded).CalculateRank(), "Wrong Rank returned")
	}

	{
		matrix := Matrix{{68, 54, 6, 230}, {16, 56, 215, 78}, {159, 186, 146, 163}, {122, 41, 205, 133}}
		coded := Matrix{{0, 0, 0, 0}, {0, 0, 0, 0}, {0, 0, 0, 0}, {0, 0, 0, 0}}
		assert.Equal(t, uint(4), NewDecoderState(matrix, coded).CalculateRank(), "Wrong Rank returned")
	}

	{
		matrix := Matrix{{100, 31, 76, 199, 119}, {207, 34, 207, 208, 18}, {62, 20, 54, 6, 187}, {66, 8, 52, 73, 54}, {122, 138, 247, 211, 165}}
		coded := Matrix{{0, 0, 0, 0}, {0, 0, 0, 0}, {0, 0, 0, 0}, {0, 0, 0, 0}, {0, 0, 0, 0}}
		assert.Equal(t, uint(5), NewDecoderState(matrix, coded).CalculateRank(), "Wrong Rank returned")
	}
}

func TestMatrixMultiplication(t *testing.T) {
	m1 := Matrix{{102, 82, 165, 0}}
	m2 := Matrix{{157, 233, 247}, {160, 28, 233}, {149, 234, 117}, {200, 181, 55}}
	m3 := Matrix{{1, 2, 3}}
	expected := Matrix{{186, 23, 11}}

	_, err := m3.Multiply(m2)
	assert.ErrorIs(t, err, kodr.ErrMatrixDimensionMismatch, "expected failed matrix multiplication error indication")

	mul, err := m1.Multiply(m2)
	assert.Nil(t, err)

	assert.Equal(t, expected, mul)

}
