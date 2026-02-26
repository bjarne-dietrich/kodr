package ssac

import (
	"math/rand/v2"

	"github.com/itzmeanjan/kodr/kodr_internals"
	"github.com/itzmeanjan/kodr/kodr_internals/operations"
)

type SSACRLNCEncoder struct {
	pieces        []kodr_internals.Piece
	extra         uint
	sparsityLevel uint
}

// PieceCount returns total #-of pieces being coded together --- denoting
// these many linearly independent pieces are required
// successfully decoding back to original pieces
func (p *SSACRLNCEncoder) PieceCount() uint {
	return uint(len(p.pieces))
}

// PieceSize returns size of one piece
// Total data being coded = pieceSize * pieceCount + padding
func (p *SSACRLNCEncoder) PieceSize() uint {
	return uint(len(p.pieces[0]))
}

// Padding returns the number of extra padding bytes added at end of original
// data slice for making all pieces of same size
func (p *SSACRLNCEncoder) Padding() uint {
	return p.extra
}

func (p *SSACRLNCEncoder) randomCodingVector() kodr_internals.CodingVector {
	vector := make(kodr_internals.CodingVector, len(p.pieces))
	n := 0

	q := []byte{DefaultQ0, DefaultQ1}
	for n < int(p.sparsityLevel) {
		// Random Index
		ri := rand.IntN(len(p.pieces))
		if vector[ri] == 0 {
			// Random Value
			vector[ri] = q[rand.IntN(2)]
			n++
		}
	}

	return vector
}

// CodedPiece returns subsequent coded pieces
func (p *SSACRLNCEncoder) CodedPiece() *kodr_internals.CodedPiece {

	vector := p.randomCodingVector()
	piece := make(kodr_internals.Piece, p.PieceSize())

	for i := range p.pieces {
		if vector[i] != 0 {
			operations.MulAddConst(piece, p.pieces[i], vector[i])
		}
	}

	compressedVector, err := CompressVector(vector)
	if err != nil {
		panic(err)
	}

	codedPiece := &kodr_internals.CodedPiece{
		Vector: compressedVector,
		Piece:  piece,
	}
	return codedPiece
}

// NewSSACRLNCEncoder can be used when you've already split original data chunk
// into pieces of same length ( in terms of bytes ), and returns a SSACRLNCEncoder,
// which delivers coded pieces on-the-fly
func NewSSACRLNCEncoder(pieces []kodr_internals.Piece) *SSACRLNCEncoder {
	return &SSACRLNCEncoder{pieces: pieces, sparsityLevel: 3}
}

// NewSSACRLNCEncoderWithPieceCount returns a SSACRLNCEncoder
// and splits the data into pieceCount same sized pieces,
// appending zero-padding to data if needed.
func NewSSACRLNCEncoderWithPieceCount(data []byte, pieceCount uint) (*SSACRLNCEncoder, error) {
	pieces, padding, err := kodr_internals.OriginalPiecesFromDataAndPieceCount(data, pieceCount)
	if err != nil {
		return nil, err
	}

	enc := NewSSACRLNCEncoder(pieces)
	enc.extra = padding
	return enc, nil
}

// NewSSACRLNCEncoderWithPieceSize returns a SSACRLNCEncoder
// and splits the data into pieces with a size of pieceSize,
// appending zero-padding to data if needed.
func NewSSACRLNCEncoderWithPieceSize(data []byte, pieceSize uint) (*SSACRLNCEncoder, error) {
	pieces, padding, err := kodr_internals.OriginalPiecesFromDataAndPieceSize(data, pieceSize)
	if err != nil {
		return nil, err
	}

	enc := NewSSACRLNCEncoder(pieces)
	enc.extra = padding
	return enc, nil
}
