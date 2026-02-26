package matrix

import (
	"github.com/itzmeanjan/kodr"
	"github.com/itzmeanjan/kodr/kodr_internals"
	"github.com/itzmeanjan/kodr/kodr_internals/operations"
)

type DecoderState struct {
	pieceCount uint
	coeffs     Matrix
	coded      Matrix
}

func (d *DecoderState) cleanForward() {
	rowCount := int(d.coeffs.Rows())
	columnCount := int(d.coeffs.Cols())
	pivotLimit := min(rowCount, columnCount)

	codedRowWidth := len(d.coded[0])

	for i := range pivotLimit {
		if d.coeffs[i][i] == 0 {
			nonZeroColumn := false
			pivot := i + 1
			for ; pivot < rowCount; pivot++ {
				if d.coeffs[pivot][i] != 0 {
					nonZeroColumn = true
					break
				}
			}

			if !nonZeroColumn {
				continue
			}

			// row switching
			d.coeffs[i], d.coeffs[pivot] = d.coeffs[pivot], d.coeffs[i]
			d.coded[i], d.coded[pivot] = d.coded[pivot], d.coded[i]
		}

		pivotRow := d.coeffs[i]
		pivotCodedRow := d.coded[i]

		// Safe because pivot is guaranteed non-zero
		pivotInverse, _ := operations.Inverse(pivotRow[i])

		for j := i + 1; j < rowCount; j++ {
			if d.coeffs[j][i] != 0 {
				mul := operations.Mul(d.coeffs[j][i], pivotInverse)
				operations.MulAddConst(d.coeffs[j][i:columnCount], pivotRow[i:columnCount], mul)
				operations.MulAddConst(d.coded[j][0:codedRowWidth], pivotCodedRow[0:codedRowWidth], mul)
			}
		}
	}
}

func (d *DecoderState) cleanBackward() {
	rowCount := int(d.coeffs.Rows())
	columnCount := int(d.coeffs.Cols())
	pivotLimit := min(rowCount, columnCount)
	codedRowWidth := len(d.coded[0])

	for i := pivotLimit - 1; i >= 0; i-- {
		pivotValue := d.coeffs[i][i]
		if pivotValue == 0 {
			continue
		}
		pivotInverse, _ := operations.Inverse(pivotValue)

		for j := 0; j < i; j++ {
			if d.coeffs[j][i] == 0 {
				continue
			}
			mul := operations.Mul(d.coeffs[j][i], pivotInverse)
			operations.MulAddConst(d.coeffs[j][i:columnCount], d.coeffs[i][i:columnCount], mul)
			operations.MulAddConst(d.coded[j][:codedRowWidth], d.coded[i][:codedRowWidth], mul)
		}

		if pivotValue == 1 {
			continue
		}
		d.coeffs[i][i] = 1

		operations.MulConst(d.coeffs[i][i+1:columnCount], d.coeffs[i][i+1:columnCount], pivotInverse)
		operations.MulConst(d.coded[i][:codedRowWidth], d.coded[i][:codedRowWidth], pivotInverse)
	}
}

func (d *DecoderState) removeZeroRows() {
	colCount := len(d.coeffs[0])
	nonZeroRows := 0

	for readIndex := 0; readIndex < len(d.coeffs); readIndex++ {
		allZero := true
		for j := range colCount {
			if d.coeffs[readIndex][j] != 0 {
				allZero = false
				break
			}
		}
		if allZero {
			continue
		}

		if readIndex != nonZeroRows {
			d.coeffs[nonZeroRows] = d.coeffs[readIndex]
			d.coded[nonZeroRows] = d.coded[readIndex]
		}
		nonZeroRows++
	}
	d.coeffs = d.coeffs[:nonZeroRows]
	d.coded = d.coded[:nonZeroRows]
}

// Calculates Reduced Row Echelon Form of coefficient
// matrix, while also modifying coded piece matrix
// First it forward, backward cleans up matrix
// i.e. cells other than pivots are zeroed,
// later it checks if some rows of coefficient matrix
// are linearly dependent or not, if yes it removes those,
// while respective rows of coded piece matrix is also
// removed --- considered to be `not useful piece`
//
// Note: All operations are in-place, no more memory
// allocations are performed
func (d *DecoderState) Rref() {
	d.cleanForward()
	d.cleanBackward()
	d.removeZeroRows()
}

// Expected to be invoked after RREF-ed, in other words
// it won't rref matrix first to calculate rank,
// rather that needs to first invoked
func (d *DecoderState) Rank() uint {
	return d.coeffs.Rows()
}

// Current state of coding coefficient matrix
func (d *DecoderState) CoefficientMatrix() Matrix {
	return d.coeffs
}

// Current state of coded piece matrix, which is updated
// along side coding coefficient matrix ( during rref )
func (d *DecoderState) CodedPieceMatrix() Matrix {
	return d.coded
}

// Adds a new coded piece to decoder state, which will hopefully
// help in decoding pieces, if linearly independent with other rows
// i.e. read pieces
func (d *DecoderState) AddPiece(codedPiece *kodr_internals.CodedPiece) {
	d.coeffs = append(d.coeffs, codedPiece.Vector)
	d.coded = append(d.coded, codedPiece.Piece)
}

// Request decoded piece by index ( 0 based, definitely )
//
// If piece not yet decoded/ requested index is >= #-of
// pieces coded together, returns error message indicating so
//
// # Otherwise piece is returned, without any error
//
// Note: This method will copy decoded piece into newly allocated memory
// when whole decoding hasn't yet happened, to prevent any chance
// that user mistakenly modifies slice returned ( read piece )
// & that affects next round of decoding ( when new piece is received )
func (d *DecoderState) GetPiece(idx uint) (kodr_internals.Piece, error) {
	if idx >= d.pieceCount {
		return nil, kodr.ErrPieceOutOfBound
	}
	if idx >= d.coeffs.Rows() {
		return nil, kodr.ErrPieceNotDecodedYet
	}

	if d.Rank() >= d.pieceCount {
		return d.coded[idx], nil
	}

	cols := int(d.coeffs.Cols())
	decoded := true

OUT:
	for i := range cols {
		switch i {
		case int(idx):
			if d.coeffs[idx][i] != 1 {
				decoded = false
				break OUT
			}

		default:
			if d.coeffs[idx][i] != 0 {
				decoded = false
				break OUT
			}

		}
	}

	if !decoded {
		return nil, kodr.ErrPieceNotDecodedYet
	}

	buf := make([]byte, d.coded.Cols())
	copy(buf, d.coded[idx])
	return buf, nil
}

func NewDecoderStateWithPieceCount(pieceCount uint) *DecoderState {
	coeffs := make([][]byte, 0, pieceCount)
	coded := make([][]byte, 0, pieceCount)
	return &DecoderState{pieceCount: pieceCount, coeffs: coeffs, coded: coded}
}

func NewDecoderState(coeffs, coded Matrix) *DecoderState {
	return &DecoderState{pieceCount: uint(len(coeffs)), coeffs: coeffs, coded: coded}
}
