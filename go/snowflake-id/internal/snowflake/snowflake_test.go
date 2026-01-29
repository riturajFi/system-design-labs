package snowflake

import "testing"

func TestNew_InvalidNode(t *testing.T) {
	if _, err := New(maxNode + 1); err == nil {
		t.Fatalf("expected error for node id > maxNode")
	}
}

func TestNextID_Monotonic(t *testing.T) {
	g, err := New(1)
	if err != nil {
		t.Fatalf("New: %v", err)
	}

	var last uint64
	for i := 0; i < 10_000; i++ {
		id, err := g.NextID()
		if err != nil {
			t.Fatalf("NextID: %v", err)
		}
		if i > 0 && id <= last {
			t.Fatalf("expected strictly increasing ids, got %d then %d", last, id)
		}
		last = id
	}
}

