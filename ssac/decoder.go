package ssac

import (
	"github.com/itzmeanjan/kodr"
	"github.com/itzmeanjan/kodr/kodr_internals"
	"github.com/itzmeanjan/kodr/kodr_internals/matrix/v2"
)

type SSACRLNCDecoder struct {
	expected, useful, received uint
	state                      *matrix.DecoderState
	sparsityLevel              uint
}

// Each piece of N-many bytes
//
// Note: If no pieces are yet added to decoder state, then
// returns 0, denoting **unknown**
func (p *SSACRLNCDecoder) PieceLength() uint {
	return p.state.GetPieceLength()
}

// Already decoded back to original pieces, with collected pieces ?
//
// If yes, no more pieces need to be collected
func (p *SSACRLNCDecoder) IsDecoded() bool {
	return p.state.GetNumberOfPieces() >= p.expected
}

// How many more pieces are required to be collected so that
// whole data can be decoded successfully ?
//
// After collecting these many pieces, original data can be decoded
func (p *SSACRLNCDecoder) Required() uint {
	return p.expected - p.useful
}

// Add one more collected coded piece, which will be used for decoding
// back to original pieces
//
// If all required pieces are already collected i.e. successful decoding
// has happened --- new pieces to be discarded, with an error denoting same
func (p *SSACRLNCDecoder) AddPiece(piece *kodr_internals.CodedPiece) error {
	return p.AddPieceBytes(piece.Flatten())
}

func (p *SSACRLNCDecoder) AddPieceBytes(pieceBytes []byte) error {
	if p.IsDecoded() {
		return kodr.ErrAllUsefulPiecesReceived
	}

	codedPiece := GetCodedPieceFromBytes(pieceBytes, DefaultQ0, DefaultQ1, p.expected, p.sparsityLevel)

	p.state.AddPiece(codedPiece)
	p.received++
	if p.received < 2 {
		p.useful++
		return nil
	}

	p.useful = p.state.CalculateRank()
	return nil
}

// GetPiece - Get a decoded piece by index, may ( not ) succeed !
//
// Note: It's not necessary that full decoding needs to happen
// for this method to return something useful
//
// If M-many pieces are received among N-many expected ( read M <= N )
// then pieces with index in [0...M] ( remember upper bound exclusive )
// can be attempted to be consumed, given algebraic structure has revealed
// requested piece at index `i`
func (p *SSACRLNCDecoder) GetPiece(i uint) (kodr_internals.Piece, error) {
	return p.state.GetPiece(i)
}

// All original pieces in order --- only when full decoding has happened
func (p *SSACRLNCDecoder) GetPieces() ([]kodr_internals.Piece, error) {
	if !p.IsDecoded() {
		return nil, kodr.ErrMoreUsefulPiecesRequired
	}

	pieces := make([]kodr_internals.Piece, 0, p.useful)
	for i := range p.useful {
		piece, err := p.GetPiece(i)
		if err != nil {
			return nil, err
		}

		pieces = append(pieces, piece)
	}

	return pieces, nil
}

func NewSSACRLNCDecoder(pieceCount uint) *SSACRLNCDecoder {
	state := matrix.NewDecoderStateWithPieceCount(pieceCount)
	return &SSACRLNCDecoder{expected: pieceCount, state: state, sparsityLevel: 3}
}
