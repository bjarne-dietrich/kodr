package ssac

import (
	"errors"
)

type BitWriter struct {
	buf     []byte
	curByte byte
	nbits   uint8 // number of bits currently filled in curByte
}

func (w *BitWriter) WriteBit(bit bool) {
	if bit {
		w.curByte |= 1 << (7 - w.nbits)
	}
	w.nbits++

	if w.nbits == 8 {
		w.buf = append(w.buf, w.curByte)
		w.curByte = 0
		w.nbits = 0
	}
}

func (w *BitWriter) WriteBits(value uint32, n uint8) error {
	if n > 32 {
		return errors.New("too many bits")
	}
	for i := n; i > 0; i-- {
		bit := (value>>(i-1))&1 == 1
		w.WriteBit(bit)
	}
	return nil
}

func (w *BitWriter) Bytes() []byte {
	if w.nbits > 0 {
		return append(w.buf, w.curByte)
	}
	return w.buf
}

type BitReader struct {
	buffer     []byte
	currentBit uint32
}

func NewBitReader(buf []byte) *BitReader {
	return &BitReader{
		buffer:     buf,
		currentBit: 0,
	}
}

func (r *BitReader) ReadBit() (bool, error) {
	if len(r.buffer) <= int(r.currentBit/8) {
		return false, errors.New("buffer too small")
	}

	b := r.buffer[r.currentBit/8]
	b >>= 7 - (r.currentBit % 8)
	r.currentBit++
	return (b & 1) == 1, nil

}

func (r *BitReader) ReadBits(n uint8) (uint32, error) {

	if len(r.buffer)*8 < int(r.currentBit)+int(n) {
		return 0, errors.New("buffer too small")
	}
	if n > 32 {
		return 0, errors.New("n must be <= 32")
	}

	var v uint32
	for i := uint8(0); i < n; i++ {
		byteIdx := r.currentBit / 8
		bitOffset := 7 - (r.currentBit % 8) // MSB first
		bit := (r.buffer[byteIdx] >> bitOffset) & 1

		v = (v << 1) | uint32(bit)
		r.currentBit++
	}

	return v, nil

}
