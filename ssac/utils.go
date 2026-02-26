package ssac

import (
	"errors"
	"math"

	"github.com/itzmeanjan/kodr/kodr_internals"
)

const DefaultQ0 byte = 21
const DefaultQ1 byte = 43

func GetCodedPieceFromBytes(bytes []byte, q0, q1 byte, n uint, m uint) *kodr_internals.CodedPiece {
	vector, err := DecompressVector(bytes, q0, q1, n, m)
	if err != nil {
		return nil
	}

	posBits := uint(math.Ceil(math.Log2(float64(n))))
	vectorBits := (posBits + 1) * m
	vectorLen := int((vectorBits + 7) / 8) // ceil(bits/8)

	return &kodr_internals.CodedPiece{
		Vector: vector,
		Piece:  bytes[vectorLen:],
	}
}

func CompressVector(vector []byte) ([]byte, error) {
	if len(vector) == 0 {
		return nil, errors.New("empty vector")
	}

	positionLenBits := uint8(math.Ceil(math.Log2(float64(len(vector)))))
	var b BitWriter

	for i, v := range vector {
		if v == 0 {
			continue
		}

		var bit bool
		switch v {
		case DefaultQ0:
			bit = false
		case DefaultQ1:
			bit = true
		default:
			return nil, errors.New("vector contains non-supported coefficient")
		}

		b.WriteBit(bit)
		if err := b.WriteBits(uint32(i), positionLenBits); err != nil {
			return nil, err
		}
	}

	return b.Bytes(), nil
}

// DecompressVector decompresses a CompressedSparseRow vector
// it takes n as number of output vector elements and m as level of sparsity
func DecompressVector(compressed []byte, q0, q1 byte, n uint, m uint) ([]byte, error) {

	output := make([]byte, n)
	positionLenBits := uint8(math.Ceil(math.Log2(float64(n))))

	r := NewBitReader(compressed)

	for range m {
		bit, err := r.ReadBit()
		if err != nil {
			return nil, err
		}

		position, err := r.ReadBits(positionLenBits)
		if err != nil {
			return nil, err
		}

		if uint(position) >= n {
			return nil, errors.New("invalid data in vector")
		}

		if bit {
			output[int(position)] = q1
		} else {
			output[int(position)] = q0
		}
	}

	return output, nil

}
