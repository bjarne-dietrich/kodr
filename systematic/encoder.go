package systematic

import (
	"github.com/itzmeanjan/kodr/kodr_internals"
	"github.com/itzmeanjan/kodr/kodr_internals/base"
	"github.com/itzmeanjan/kodr/kodr_internals/operations"
)

type SystematicRLNCEncoder struct {
	base.BaseEncoder
}

// How many bytes of data, constructed by concatenating
// coded pieces together, required at minimum for decoding
// back to original pieces ?
//
// As I'm coding N-many pieces together, I need at least N-many
// linearly independent pieces, which are concatenated together
// to form a byte slice & can be used for original data reconstruction.
//
// So it computes N * codedPieceLen
func (s *SystematicRLNCEncoder) DecodableLen() uint {
	return s.PieceCount() * s.CodedPieceLen()
}

// If N-many original pieces are coded together
// what could be length of one such coded piece
// obtained by invoking `CodedPiece` ?
//
// Here N = len(pieces), original pieces which are
// being coded together
func (s *SystematicRLNCEncoder) CodedPieceLen() uint {
	return s.PieceCount() + s.PieceSize()
}

// Generates a systematic coded piece's coding vector, which has
// only one non-zero element ( 1 )
func (s *SystematicRLNCEncoder) systematicCodingVector(idx uint) kodr_internals.CodingVector {
	if !(idx < s.PieceCount()) {
		return nil
	}

	vector := make(kodr_internals.CodingVector, s.PieceCount())
	vector[idx] = 1
	return vector
}

// For systematic coding, first N-piece are returned in uncoded form
// i.e. coding vectors are having only single non-zero element ( 1 )
// in respective index of piece.
//
// Piece index `i` ( returned from this method ), where i < N
// is going to have coding vector = [N]byte, where only i'th index
// of this vector will have 1, all other fields will have 0.
//
// Here N = #-of pieces being coded together
//
// Later pieces are coded as they're done in Full RLNC scheme
// `i` keeps incrementing by +1, until it reaches N
func (s *SystematicRLNCEncoder) CodedPiece() *kodr_internals.CodedPiece {

	pieceID := s.GetCurrentPieceIdAndIncrement()
	pieceCount := s.PieceCount()

	if pieceID < pieceCount {
		// `nil` coding vector can be returned, which is
		// not being checked at all, as in that case we'll
		// never get into `if` branch
		vector := s.systematicCodingVector(pieceID)
		piece := make(kodr_internals.Piece, s.PieceSize())
		copy(piece, *s.GetPiece(pieceID))

		return &kodr_internals.CodedPiece{
			Vector: vector,
			Piece:  piece,
		}
	}

	vector := kodr_internals.GenerateCodingVector(pieceCount)
	piece := make(kodr_internals.Piece, s.PieceSize())

	for i := range pieceCount {
		operations.MulAddConst(piece, *s.GetPiece(i), vector[i])
	}

	return &kodr_internals.CodedPiece{
		Vector: vector,
		Piece:  piece,
	}
}

// When you've already split original data chunk into pieces
// of same length ( in terms of bytes ), this function can be used
// for creating one systematic RLNC encoder, which delivers coded pieces
// on-the-fly
func NewSystematicRLNCEncoder(pieces []kodr_internals.Piece) *SystematicRLNCEncoder {
	return &SystematicRLNCEncoder{*base.NewBaseEncoder(pieces)}
}

// If you know #-of pieces you want to code together, invoking
// this function splits whole data chunk into N-pieces, with padding
// bytes appended at end of last piece, if required & prepares
// full RLNC encoder for obtaining coded pieces
func NewSystematicRLNCEncoderWithPieceCount(data []byte, pieceCount uint) (*SystematicRLNCEncoder, error) {
	encoder, err := base.NewBaseEncoderWithPieceCount(data, pieceCount)
	if err != nil {
		return nil, err
	}
	return &SystematicRLNCEncoder{*encoder}, nil
}

// If you want to have N-bytes piece size for each, this
// function generates M-many pieces each of N-bytes size, which are ready
// to be coded together with full RLNC
func NewSystematicRLNCEncoderWithPieceSize(data []byte, pieceSize uint) (*SystematicRLNCEncoder, error) {
	encoder, err := base.NewBaseEncoderWithPieceSize(data, pieceSize)
	if err != nil {
		return nil, err
	}
	return &SystematicRLNCEncoder{*encoder}, nil
}
