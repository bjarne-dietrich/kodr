package ssac

import (
	"math/rand/v2"

	"github.com/itzmeanjan/kodr/kodr_internals"
	"github.com/itzmeanjan/kodr/kodr_internals/base"
	"github.com/itzmeanjan/kodr/kodr_internals/operations"
)

type SSACRLNCEncoder struct {
	base.BaseEncoder
	sparsityLevel uint
}

func (p *SSACRLNCEncoder) randomCodingVector() kodr_internals.CodingVector {
	pieceCount := p.PieceCount()

	vector := make(kodr_internals.CodingVector, pieceCount)
	n := 0

	q := []byte{DefaultQ0, DefaultQ1}
	for n < int(p.sparsityLevel) {
		// Random Index
		ri := rand.IntN(int(pieceCount))
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
	pieceCount := p.PieceCount()
	piece := make(kodr_internals.Piece, p.PieceSize())

	for i := range pieceCount {
		if vector[i] != 0 {
			operations.MulAddConst(piece, *p.GetPiece(i), vector[i])
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
	return &SSACRLNCEncoder{BaseEncoder: *base.NewBaseEncoder(pieces), sparsityLevel: 3}
}

// NewSSACRLNCEncoderWithPieceCount returns a SSACRLNCEncoder
// and splits the data into pieceCount same sized pieces,
// appending zero-padding to data if needed.
func NewSSACRLNCEncoderWithPieceCount(data []byte, pieceCount uint) (*SSACRLNCEncoder, error) {

	encoder, err := base.NewBaseEncoderWithPieceCount(data, pieceCount)
	if err != nil {
		return nil, err
	}
	return &SSACRLNCEncoder{BaseEncoder: *encoder, sparsityLevel: 3}, nil

}

// NewSSACRLNCEncoderWithPieceSize returns a SSACRLNCEncoder
// and splits the data into pieces with a size of pieceSize,
// appending zero-padding to data if needed.
func NewSSACRLNCEncoderWithPieceSize(data []byte, pieceSize uint) (*SSACRLNCEncoder, error) {
	encoder, err := base.NewBaseEncoderWithPieceSize(data, pieceSize)
	if err != nil {
		return nil, err
	}
	return &SSACRLNCEncoder{BaseEncoder: *encoder, sparsityLevel: 3}, nil
}
