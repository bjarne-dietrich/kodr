package base

import (
	"sync"

	"github.com/itzmeanjan/kodr"
	"github.com/itzmeanjan/kodr/kodr_internals"
	"github.com/itzmeanjan/kodr/kodr_internals/matrix/v2"
)

type BaseDecoder struct {
	expected, useful, received uint
	state                      *matrix.DecoderState
	mutex                      sync.RWMutex
}

// Each piece of N-many bytes
//
// Note: If no pieces are yet added to decoder state, then
// returns 0, denoting **unknown**
func (d *BaseDecoder) PieceLength() uint {
	d.mutex.Lock()
	defer d.mutex.Unlock()
	return d.state.GetPieceLength()
}

func (d *BaseDecoder) GetExpectedPieceCount() uint {
	d.mutex.RLock()
	defer d.mutex.RUnlock()
	return d.expected
}

// Already decoded back to original pieces, with collected pieces ?
// If yes, no more pieces need to be collected
func (d *BaseDecoder) IsDecoded() bool {
	d.mutex.RLock()
	defer d.mutex.RUnlock()
	return d.useful >= d.expected
}

// How many more pieces are required to be collected so that
// whole data can be decoded successfully ?
//
// After collecting these many pieces, original data can be decoded
func (d *BaseDecoder) Required() uint {
	d.mutex.RLock()
	defer d.mutex.RUnlock()
	return d.expected - d.useful
}

// Add one more collected coded piece, which will be used for decoding
// back to original pieces
//
// If all required pieces are already collected i.e. successful decoding
// has happened --- new pieces to be discarded, with an error denoting same
func (d *BaseDecoder) AddPiece(piece *kodr_internals.CodedPiece) error {
	if d.IsDecoded() {
		return kodr.ErrAllUsefulPiecesReceived
	}
	d.mutex.Lock()
	defer d.mutex.Unlock()

	d.state.AddPiece(piece)
	d.received++
	if d.received < 2 {
		d.useful++
		return nil
	}

	d.useful = d.state.CalculateRank()
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
func (d *BaseDecoder) GetPiece(i uint) (kodr_internals.Piece, error) {
	d.mutex.RLock()
	defer d.mutex.RUnlock()
	return d.state.GetPiece(i)
}

// All original pieces in order --- only when full decoding has happened
func (d *BaseDecoder) GetPieces() ([]kodr_internals.Piece, error) {
	if !d.IsDecoded() {
		return nil, kodr.ErrMoreUsefulPiecesRequired
	}

	d.mutex.RLock()
	defer d.mutex.RUnlock()

	pieces := make([]kodr_internals.Piece, 0, d.useful)
	for i := range d.useful {
		piece, err := d.GetPiece(i)
		if err != nil {
			return nil, err
		}

		pieces = append(pieces, piece)
	}

	return pieces, nil
}

func NewBaseDecoder(pieceCount uint) *BaseDecoder {
	state := matrix.NewDecoderStateWithPieceCount(pieceCount)
	return &BaseDecoder{expected: pieceCount, state: state}
}
