package ssac_test

import (
	"slices"
	"testing"

	ssac "github.com/itzmeanjan/kodr/ssac"
)

func TestBitWriter(t *testing.T) {
	bytes := []byte{0b00110100, 0b11011101}
	var b ssac.BitWriter

	b.WriteBit(false)
	_ = b.WriteBits(0b01101, 5)
	_ = b.WriteBits(0b001101, 6)
	b.WriteBit(true)
	_ = b.WriteBits(0b1, 1)
	_ = b.WriteBits(0b01, 2)

	if !slices.Equal(b.Bytes(), bytes) {
		t.Logf("Have: %b", b.Bytes())
		t.Logf("Expected: %b", bytes)
		t.Errorf("Bytes mismatch")
	}
}

func TestBitWriter_BigN(t *testing.T) {
	var b ssac.BitWriter

	if err := b.WriteBits(0, 255); err == nil {
		t.Fatal("expected error")
	}
}

func TestBitReader(t *testing.T) {
	bytes := []byte{0b00110100, 0b11011101}
	r := ssac.NewBitReader(bytes)

	if bit, err := r.ReadBit(); err != nil {
		t.Fatal(err)
	} else if bit {
		t.Fatal("expected false")
	}

	if bits, err := r.ReadBits(5); err != nil {
		t.Fatal(err)
	} else if bits != 0b01101 {
		t.Fatal("wrong bits returned")
	}

	if bits, err := r.ReadBits(6); err != nil {
		t.Fatal(err)
	} else if bits != 0b001101 {
		t.Fatal("wrong bits returned")
	}

	if bit, err := r.ReadBit(); err != nil {
		t.Fatal(err)
	} else if !bit {
		t.Fatal("expected true")
	}

	if bits, err := r.ReadBits(1); err != nil {
		t.Fatal(err)
	} else if bits != 1 {
		t.Fatal("wrong bits returned")
	}

	if bits, err := r.ReadBits(2); err != nil {
		t.Fatal(err)
	} else if bits != 1 {
		t.Fatal("wrong bits returned")
	}

}

func TestBitReader_BigN(t *testing.T) {
	bytes := []byte{0b00110100, 0b11011101, 0b00110100, 0b11011101, 0b00110100, 0b11011101, 0b00110100, 0b11011101}
	r := ssac.NewBitReader(bytes)

	if _, err := r.ReadBits(35); err == nil {
		t.Fatal("expected error")
	}
}

func TestBitReader_ReadTooMuch(t *testing.T) {
	bytes := []byte{0b00110100, 0b11011101}
	r := ssac.NewBitReader(bytes)

	if _, err := r.ReadBits(17); err == nil {
		t.Fatal("expected error")
	}
}

func TestBitReader_ReadTooMuch2(t *testing.T) {
	bytes := []byte{0b00110100, 0b11011101}
	r := ssac.NewBitReader(bytes)

	if _, err := r.ReadBits(16); err != nil {
		t.Fatal(err)
	}
	if _, err := r.ReadBit(); err == nil {
		t.Fatal("expected error")
	}
}

func TestBitReader_ReadTooMuch3(t *testing.T) {
	bytes := []byte{0b00110100, 0b11011101}
	r := ssac.NewBitReader(bytes)

	if _, err := r.ReadBits(16); err != nil {
		t.Fatal(err)
	}
	if _, err := r.ReadBits(2); err == nil {
		t.Fatal("expected error")
	}
}
