package yaml

import (
	"testing"
)

func TestNode(t *testing.T) {
	n := &Node{}
	if !n.IsNull() || n.IsScalar() || n.IsSequence() || n.IsMap() {
		t.Fatalf("Node should be null.")
	}

	scalarvalues := []interface{}{
		"Scalar test.",
		int(1),
		int8(2),
		int16(3),
		int32(4),
		int64(5),
		uint(6),
		uint8(7),
		uint16(8),
		uint32(9),
		uint64(10),
		float32(1.0),
		float64(2.0),
	}

	for _, scalar := range scalarvalues {
		n.Set(scalar)
		if n.IsNull() || !n.IsScalar() || n.IsSequence() || n.IsMap() {
			t.Fatalf("Node should be scalar.")
		}
	}
}