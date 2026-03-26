package base

import "github.com/itzmeanjan/kodr/kodr_internals"

type Decoder interface {
	AddPiece(piece *kodr_internals.CodedPiece) error
	GetPieces() ([]kodr_internals.Piece, error)
	IsDecoded() bool
}
