package base

import "github.com/itzmeanjan/kodr/kodr_internals"

type Encoder interface {
	CodedPiece() *kodr_internals.CodedPiece
}
