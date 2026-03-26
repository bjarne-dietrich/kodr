package matrix

import (
	"github.com/itzmeanjan/kodr"
	"github.com/itzmeanjan/kodr/kodr_internals"
	"github.com/itzmeanjan/kodr/kodr_internals/operations"
)

type DecoderState struct {
	pieceCount        uint
	coefficientMatrix Matrix
	codedPiecesMatrix Matrix
}

// forwardEliminate performs the forward (Gaussian) elimination phase on the
// coefficient matrix to create an upper-triangular / row-echelon form.
//
// For each pivot position i, it ensures a non-zero pivot by swapping with a
// lower row when necessary, then eliminates all entries below the pivot in
// column i. The same row operations are applied to the coded piece matrix so
// both matrices remain consistent.
//
// All operations are in-place and do not allocate.
func (d *DecoderState) forwardEliminate() {
	rowCount := int(d.coefficientMatrix.Rows())
	columnCount := int(d.coefficientMatrix.Cols())
	pivotLimit := min(rowCount, columnCount)

	codedRowWidth := len(d.codedPiecesMatrix[0])

	for i := range pivotLimit {
		if d.coefficientMatrix[i][i] == 0 {
			nonZeroColumn := false
			pivot := i + 1
			for ; pivot < rowCount; pivot++ {
				if d.coefficientMatrix[pivot][i] != 0 {
					nonZeroColumn = true
					break
				}
			}

			if !nonZeroColumn {
				continue
			}

			// row switching
			d.coefficientMatrix[i], d.coefficientMatrix[pivot] = d.coefficientMatrix[pivot], d.coefficientMatrix[i]
			d.codedPiecesMatrix[i], d.codedPiecesMatrix[pivot] = d.codedPiecesMatrix[pivot], d.codedPiecesMatrix[i]
		}

		pivotRow := d.coefficientMatrix[i]
		pivotCodedRow := d.codedPiecesMatrix[i]

		// Safe because pivot is guaranteed non-zero
		pivotInverse, _ := operations.Inverse(pivotRow[i])

		for j := i + 1; j < rowCount; j++ {
			if d.coefficientMatrix[j][i] != 0 {
				mul := operations.Mul(d.coefficientMatrix[j][i], pivotInverse)
				operations.MulAddConst(d.coefficientMatrix[j][i:columnCount], pivotRow[i:columnCount], mul)
				operations.MulAddConst(d.codedPiecesMatrix[j][0:codedRowWidth], pivotCodedRow[0:codedRowWidth], mul)
			}
		}
	}
}

// backSubstitute performs the backward (Gauss–Jordan) phase after forward
// elimination.
//
// Walking pivots from bottom to top, it eliminates all entries above each
// pivot, and normalizes the pivot row so the pivot becomes 1 (when non-zero).
// The same row operations are applied to the coded piece matrix to keep it in
// lockstep with the coefficient matrix.
//
// All operations are in-place and do not allocate.
func (d *DecoderState) backSubstitute() {
	rowCount := int(d.coefficientMatrix.Rows())
	columnCount := int(d.coefficientMatrix.Cols())
	pivotLimit := min(rowCount, columnCount)
	codedRowWidth := len(d.codedPiecesMatrix[0])

	for i := pivotLimit - 1; i >= 0; i-- {
		pivotValue := d.coefficientMatrix[i][i]
		if pivotValue == 0 {
			continue
		}
		pivotInverse, _ := operations.Inverse(pivotValue)

		for j := 0; j < i; j++ {
			if d.coefficientMatrix[j][i] == 0 {
				continue
			}
			mul := operations.Mul(d.coefficientMatrix[j][i], pivotInverse)
			operations.MulAddConst(d.coefficientMatrix[j][i:columnCount], d.coefficientMatrix[i][i:columnCount], mul)
			operations.MulAddConst(d.codedPiecesMatrix[j][:codedRowWidth], d.codedPiecesMatrix[i][:codedRowWidth], mul)
		}

		if pivotValue == 1 {
			continue
		}
		d.coefficientMatrix[i][i] = 1

		operations.MulConst(d.coefficientMatrix[i][i+1:columnCount], d.coefficientMatrix[i][i+1:columnCount], pivotInverse)
		operations.MulConst(d.codedPiecesMatrix[i][:codedRowWidth], d.codedPiecesMatrix[i][:codedRowWidth], pivotInverse)
	}
}

// dropZeroRows removes rows whose coefficient row is entirely zero, and drops
// the corresponding rows from the coded piece matrix as well.
//
// It compacts both matrices in-place (stable with respect to the remaining
// rows) by overwriting removed rows, then reslices to the new row count.
// No allocations are performed.
func (d *DecoderState) dropZeroRows() {
	colCount := len(d.coefficientMatrix[0])
	nonZeroRows := 0

	for readIndex := 0; readIndex < len(d.coefficientMatrix); readIndex++ {
		allZero := true
		for j := range colCount {
			if d.coefficientMatrix[readIndex][j] != 0 {
				allZero = false
				break
			}
		}
		if allZero {
			continue
		}

		if readIndex != nonZeroRows {
			d.coefficientMatrix[nonZeroRows] = d.coefficientMatrix[readIndex]
			d.codedPiecesMatrix[nonZeroRows] = d.codedPiecesMatrix[readIndex]
		}
		nonZeroRows++
	}
	d.coefficientMatrix = d.coefficientMatrix[:nonZeroRows]
	d.codedPiecesMatrix = d.codedPiecesMatrix[:nonZeroRows]
}

// ReduceAndPrune reduces the coefficient matrix to reduced row echelon form
// (Gauss–Jordan elimination) while applying the same row operations to the
// codedPiecesMatrix.
//
// After reduction, it removes any all-zero / dependent rows from the
// coefficient matrix and drops the corresponding rows from codedPiecesMatrix
// (discarding non-useful pieces).
//
// All operations are performed in-place and do not allocate.
func (d *DecoderState) ReduceAndPrune() {
	d.forwardEliminate()
	d.backSubstitute()
	d.dropZeroRows()
}

// GetNumberOfPieces returns the number of rows in coefficientMatrix.
func (d *DecoderState) GetNumberOfPieces() uint {
	return d.coefficientMatrix.Rows()
}

// CalculateRank invokes ReduceAndPrune and returns GetNumberOfPieces.
func (d *DecoderState) CalculateRank() uint {
	d.ReduceAndPrune()
	return d.GetNumberOfPieces()
}

// GetPieceLength returns the length of one received data piece.
// If no piece was received yet, 0 is returned.
// (All Pieces do have the same length)
func (d *DecoderState) GetPieceLength() uint {
	if len(d.codedPiecesMatrix) == 0 {
		return 0
	}
	return d.codedPiecesMatrix.Cols()
}

// Adds a new codedPiecesMatrix piece to decoder state, which will hopefully
// help in decoding pieces, if linearly independent with other rows
// i.e. read pieces
func (d *DecoderState) AddPiece(codedPiece *kodr_internals.CodedPiece) {
	d.coefficientMatrix = append(d.coefficientMatrix, codedPiece.Vector)
	d.codedPiecesMatrix = append(d.codedPiecesMatrix, codedPiece.Piece)
}

// Request decoded piece by index ( 0 based, definitely )
//
// If piece not yet decoded/ requested index is >= #-of
// pieces codedPiecesMatrix together, returns error message indicating so
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
	if idx >= d.coefficientMatrix.Rows() {
		return nil, kodr.ErrPieceNotDecodedYet
	}

	if d.GetNumberOfPieces() >= d.pieceCount {
		return d.codedPiecesMatrix[idx], nil
	}

	cols := int(d.coefficientMatrix.Cols())
	decoded := true

OUT:
	for i := range cols {
		switch i {
		case int(idx):
			if d.coefficientMatrix[idx][i] != 1 {
				decoded = false
				break OUT
			}

		default:
			if d.coefficientMatrix[idx][i] != 0 {
				decoded = false
				break OUT
			}

		}
	}

	if !decoded {
		return nil, kodr.ErrPieceNotDecodedYet
	}

	buf := make([]byte, d.codedPiecesMatrix.Cols())
	copy(buf, d.codedPiecesMatrix[idx])
	return buf, nil
}

func NewDecoderStateWithPieceCount(pieceCount uint) *DecoderState {
	coefficients := make([][]byte, 0, pieceCount)
	coded := make([][]byte, 0, pieceCount)
	return &DecoderState{pieceCount: pieceCount, coefficientMatrix: coefficients, codedPiecesMatrix: coded}
}

func NewDecoderState(coefficients, coded Matrix) *DecoderState {
	return &DecoderState{pieceCount: uint(len(coefficients)), coefficientMatrix: coefficients, codedPiecesMatrix: coded}
}
