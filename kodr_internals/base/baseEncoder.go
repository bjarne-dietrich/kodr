package base

import (
	"github.com/itzmeanjan/kodr/kodr_internals"
)

type BaseEncoder struct {
	currentPieceId uint
	pieces         []kodr_internals.Piece
	extra          uint
}

// PieceCount returns total #-of pieces being coded together --- denoting
// these many linearly independent pieces are required
// successfully decoding back to original pieces
func (e *BaseEncoder) PieceCount() uint {
	return uint(len(e.pieces))
}

// PieceSize returns size of one piece
// Total data being coded = pieceSize * pieceCount + padding
func (e *BaseEncoder) PieceSize() uint {
	return uint(len(e.pieces[0]))
}

// Padding returns the number of extra padding bytes added at end of original
// data slice for making all pieces of same size
func (e *BaseEncoder) Padding() uint {
	return e.extra
}

func (e *BaseEncoder) GetCurrentPieceId() uint {
	return e.currentPieceId
}

func (e *BaseEncoder) IncrementCurrentPieceId() {
	e.currentPieceId++
}

func (e *BaseEncoder) GetPiece(pieceId uint) *kodr_internals.Piece {
	if pieceId < uint(len(e.pieces)) {
		return &e.pieces[pieceId]
	}

	return nil
}

func NewBaseEncoder(pieces []kodr_internals.Piece) *BaseEncoder {
	return &BaseEncoder{currentPieceId: 0, pieces: pieces}
}

func NewBaseEncoderWithPieceCount(data []byte, pieceCount uint) (*BaseEncoder, error) {
	pieces, padding, err := kodr_internals.OriginalPiecesFromDataAndPieceCount(data, pieceCount)
	if err != nil {
		return nil, err
	}

	enc := NewBaseEncoder(pieces)
	enc.extra = padding
	return enc, nil
}

func NewBaseEncoderWithPieceSize(data []byte, pieceSize uint) (*BaseEncoder, error) {
	pieces, padding, err := kodr_internals.OriginalPiecesFromDataAndPieceSize(data, pieceSize)
	if err != nil {
		return nil, err
	}

	enc := NewBaseEncoder(pieces)
	enc.extra = padding
	return enc, nil
}
