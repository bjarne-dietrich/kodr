package pseudo

import (
	"slices"
	"testing"
)

func TestNonPrimitiveCodingVector(t *testing.T) {

	expectEqual(t, nonPrimitiveCodingVector(10, 10), []byte{1, 0, 1, 0, 1, 0, 1, 0, 1, 0})
	expectEqual(t, nonPrimitiveCodingVector(11, 10), []byte{1, 0, 0, 1, 0, 0, 1, 0, 0, 1})
	expectEqual(t, nonPrimitiveCodingVector(12, 10), []byte{0, 1, 0, 0, 1, 0, 0, 1, 0, 0})
	expectEqual(t, nonPrimitiveCodingVector(13, 10), []byte{1, 0, 0, 0, 1, 0, 0, 0, 1, 0})
	expectEqual(t, nonPrimitiveCodingVector(14, 10), []byte{0, 1, 0, 0, 0, 1, 0, 0, 0, 1})
	expectEqual(t, nonPrimitiveCodingVector(15, 10), []byte{0, 0, 1, 0, 0, 0, 1, 0, 0, 0})

}

func expectEqual(t *testing.T, expected, got []byte) {
	if !slices.Equal(expected, got) {
		t.Fatalf("Expected Vector %v, Got: %v", expected, got)
	}
}
