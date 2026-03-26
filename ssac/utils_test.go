package ssac_test

import (
	"bytes"
	"math/rand/v2"
	"testing"

	"github.com/itzmeanjan/kodr/kodr_internals"
	ssac "github.com/itzmeanjan/kodr/ssac"
)

func Test_CompressVector(t *testing.T) {
	vector := []byte{0, ssac.DefaultQ0, ssac.DefaultQ0, 0, ssac.DefaultQ1, 0, 0, 0}
	expected := []byte{0b00010010, 0b11000000}

	csr, err := ssac.CompressVector(vector)
	if err != nil {
		t.Fatal(err)
	}

	if !bytes.Equal(csr, expected) {
		t.Fatalf("compressed vector does not match expected output")
	}

}

func Test_DecompressVector(t *testing.T) {
	expected := []byte{0, ssac.DefaultQ0, ssac.DefaultQ0, 0, ssac.DefaultQ1, 0, 0, 0}
	compressed := []byte{0b00010010, 0b11000000}

	vector, err := ssac.DecompressVector(compressed, ssac.DefaultQ0, ssac.DefaultQ1, 8, 3)
	if err != nil {
		t.Fatal(err)
	}

	if !bytes.Equal(vector, expected) {
		t.Log(vector)
		t.Fatalf("compressed vector does not match expected output")
	}

}

func TestVectorRoundtrip(t *testing.T) {
	pieceCount := uint(64)
	sparsity := uint(3)

	for n := 0; n < 10000; n++ {
		// build random sparse vector like encoder does
		v := make([]byte, pieceCount)
		q := []byte{ssac.DefaultQ0, ssac.DefaultQ1}
		filled := 0
		for filled < int(sparsity) {
			i := rand.IntN(int(pieceCount))
			if v[i] == 0 {
				v[i] = q[rand.IntN(2)]
				filled++
			}
		}

		cv, err := ssac.CompressVector(v)
		if err != nil {
			t.Fatalf("compress: %v", err)
		}

		// fake piece and go through full byte format
		p := make([]byte, 128)
		cp := &kodr_internals.CodedPiece{Vector: cv, Piece: p}
		got := ssac.GetCodedPieceFromBytes(cp.Flatten(), ssac.DefaultQ0, ssac.DefaultQ1, pieceCount, sparsity)
		if got == nil {
			t.Fatalf("got nil, want non-nil")
		}
		if !bytes.Equal(got.Vector, v) {
			t.Fatalf("vector mismatch")
		}
	}
}
