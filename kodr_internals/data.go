package kodr_internals

import (
	"crypto/rand"
	"math"

	"github.com/itzmeanjan/kodr"
	"github.com/itzmeanjan/kodr/kodr_internals/gf256"
)

// A Piece of data is nothing but a byte array
type Piece []byte

// Multiply performs the operation Piece += piece * by‚
// Two pieces are coded together by performing
// symbol by symbol finite field arithmetic, where
// a single byte is a symbol.
// The second piece is multiplied by coding coefficient 'c'
func (p *Piece) Multiply(piece Piece, c byte) {
	for i := range piece {
		res := gf256.New((*p)[i])

		l := gf256.New(piece[i])
		r := gf256.New(c)

		res.AddAssign(l.Mul(r))
		(*p)[i] = res.Get()
	}
}

// The CodingVector is one component of coded piece
// It holds information regarding how original pieces are
// combined
type CodingVector []byte

// CodedPiece holds a Piece along with its CodingVector
// to be used by recoder / decoder
type CodedPiece struct {
	Vector CodingVector
	Piece  Piece
}

// Len returns the total length of coded piece which is len(coding_vector) + len(piece)
func (c *CodedPiece) Len() uint {
	return uint(len(c.Vector) + len(c.Piece))
}

// Flatten copies CodingVector and CodedPiece into single byte
// slice { CodingVector..., CodedPiece... }
func (c *CodedPiece) Flatten() []byte {
	res := make([]byte, c.Len())
	copy(res[:len(c.Vector)], c.Vector)
	copy(res[len(c.Vector):], c.Piece)
	return res
}

// IsSystematic returns true if the CodingVector of this CodedPiece
// has only non-zero elements and exactly one element '1'.
// This effectively make the piece an uncoded chunk of original data.
func (c *CodedPiece) IsSystematic() bool {
	pos := -1
	for i := range len(c.Vector) {
		switch c.Vector[i] {
		case 0:
			continue

		case 1:
			if pos != -1 {
				return false
			}
			pos = i

		default:
			return false

		}
	}

	return pos >= 0 && pos < len(c.Vector)
}

// GenerateCodingVector returns a‚ random CodingVector of specified length.
// It uses crypto.rand.Rand and panics if crypto.rand.Read fails which
// seems impossible at time of writing.
func GenerateCodingVector(n uint) CodingVector {
	vector := make(CodingVector, n)
	_, err := rand.Read(vector)
	if err != nil {
		panic(err)
	}
	return vector
}

// OriginalPiecesFromDataAndPieceSize takes a data byte slice and a desired pieceSize in bytes.
// It'll split data into n pieceSize sized chunks and use zero-padding if needed.
// It returns a slice of n pieceSize chunks typed as Piece,
// the number of padding bytes appended to original data before splitting
// and an Error if pieceSize >= len(data), if pieceSize == 0
// or if appending padding to data fails.
func OriginalPiecesFromDataAndPieceSize(data []byte, pieceSize uint) ([]Piece, uint, error) {
	if pieceSize == 0 {
		return nil, 0, kodr.ErrZeroPieceSize
	}

	if int(pieceSize) >= len(data) {
		return nil, 0, kodr.ErrBadPieceCount
	}

	pieceCount := int(math.Ceil(float64(len(data)) / float64(pieceSize)))
	padding := uint(pieceCount*int(pieceSize) - len(data))

	var data_ []byte
	if padding > 0 {
		data_ = make([]byte, pieceCount*int(pieceSize))
		if n := copy(data_, data); n != len(data) {
			return nil, 0, kodr.ErrCopyFailedDuringPieceConstruction
		}
	} else {
		data_ = data
	}

	pieces := make([]Piece, pieceCount)
	for i := 0; i < pieceCount; i++ {
		piece := data_[int(pieceSize)*i : int(pieceSize)*(i+1)]
		pieces[i] = piece
	}

	return pieces, padding, nil
}

// OriginalPiecesFromDataAndPieceCount takes a data byte slice and a desired pieceCount.
// It'll split data into pieceCount equal sized chunks and use zero-padding if needed.
// It returns a slice of pieceCount chunks typed as Piece,
// the number of padding bytes appended to original data before splitting
// and an Error if pieceCount > len(data), if pieceCount < 2
// or if appending padding to data fails.
func OriginalPiecesFromDataAndPieceCount(data []byte, pieceCount uint) ([]Piece, uint, error) {
	if pieceCount < 2 {
		return nil, 0, kodr.ErrBadPieceCount
	}

	if int(pieceCount) > len(data) {
		return nil, 0, kodr.ErrPieceCountMoreThanTotalBytes
	}

	pieceSize := (uint(len(data)) + (pieceCount - 1)) / pieceCount
	padding := pieceCount*pieceSize - uint(len(data))

	var data_ []byte
	if padding > 0 {
		data_ = make([]byte, pieceSize*pieceCount)
		if n := copy(data_, data); n != len(data) {
			return nil, 0, kodr.ErrCopyFailedDuringPieceConstruction
		}
	} else {
		data_ = data
	}

	// padding field will always be 0, because it's already extended ( if required ) to be
	// properly divisible by `pieceSize`
	split, _, err := OriginalPiecesFromDataAndPieceSize(data_, pieceSize)
	return split, padding, err
}

// CodedPiecesFromBytes takes a byte slice containing
// pieceCount CodedPiece in CodedPiece.Flatten form.
// It returns pieceCount CodedPiece or an error
// if pieceCount or piecesCodedTogether have a mismatch to len(data)
func CodedPiecesFromBytes(data []byte, pieceCount uint, piecesCodedTogether uint) ([]*CodedPiece, error) {
	codedPieceLength := len(data) / int(pieceCount)
	if codedPieceLength*int(pieceCount) != len(data) {
		return nil, kodr.ErrCodedDataLengthMismatch
	}

	if !(piecesCodedTogether < uint(codedPieceLength)) {
		return nil, kodr.ErrCodingVectorLengthMismatch
	}

	codedPieces := make([]*CodedPiece, pieceCount)
	for i := 0; i < int(pieceCount); i++ {
		codedPiece := data[codedPieceLength*i : codedPieceLength*(i+1)]
		codedPieces[i] = &CodedPiece{
			Vector: codedPiece[:piecesCodedTogether],
			Piece:  codedPiece[piecesCodedTogether:],
		}
	}

	return codedPieces, nil
}
