package types

import "testing"

func TestSlotToTime(t *testing.T) {
	genesis := uint64(1700000000)
	if SlotToTime(0, genesis) != 1700000000 {
		t.Error("slot 0")
	}
	if SlotToTime(1, genesis) != 1700000004 {
		t.Error("slot 1")
	}
	if SlotToTime(100, genesis) != 1700000400 {
		t.Error("slot 100")
	}
}

func TestTimeToSlot(t *testing.T) {
	genesis := uint64(1700000000)
	if TimeToSlot(1700000000, genesis) != 0 {
		t.Error("time at genesis")
	}
	if TimeToSlot(1700000004, genesis) != 1 {
		t.Error("time +4s")
	}
	if TimeToSlot(1699999999, genesis) != 0 {
		t.Error("time before genesis")
	}
}

func TestRootIsZero(t *testing.T) {
	var zero Root
	if !zero.IsZero() {
		t.Error("zero root")
	}
	if (Root{1}).IsZero() {
		t.Error("non-zero root")
	}
}

func TestIsJustifiableAfter(t *testing.T) {
	finalized := Slot(10)

	tests := []struct {
		slot     Slot
		expected bool
		reason   string
	}{
		// Rule 1: delta <= 5 always justifiable
		{10, true, "delta=0"},
		{11, true, "delta=1"},
		{12, true, "delta=2"},
		{13, true, "delta=3"},
		{14, true, "delta=4"},
		{15, true, "delta=5"},

		// delta=6 is pronic (2*3), justifiable
		{16, true, "delta=6 (pronic 2*3)"},

		// delta=7 is neither square nor pronic
		{17, false, "delta=7 (not justifiable)"},

		// delta=8 is neither square nor pronic
		{18, false, "delta=8 (not justifiable)"},

		// delta=9 is perfect square (3^2)
		{19, true, "delta=9 (perfect square 3^2)"},

		// delta=10 is neither
		{20, false, "delta=10 (not justifiable)"},

		// delta=12 is pronic (3*4)
		{22, true, "delta=12 (pronic 3*4)"},

		// delta=16 is perfect square (4^2)
		{26, true, "delta=16 (perfect square 4^2)"},

		// delta=20 is pronic (4*5)
		{30, true, "delta=20 (pronic 4*5)"},

		// delta=25 is perfect square (5^2)
		{35, true, "delta=25 (perfect square 5^2)"},

		// delta=30 is pronic (5*6)
		{40, true, "delta=30 (pronic 5*6)"},

		// delta=36 is perfect square (6^2)
		{46, true, "delta=36 (perfect square 6^2)"},

		// delta=42 is pronic (6*7)
		{52, true, "delta=42 (pronic 6*7)"},
	}

	for _, tt := range tests {
		got := tt.slot.IsJustifiableAfter(finalized)
		if got != tt.expected {
			t.Errorf("Slot(%d).IsJustifiableAfter(%d) = %v, want %v (%s)",
				tt.slot, finalized, got, tt.expected, tt.reason)
		}
	}
}

func TestIsJustifiableAfter_BeforeFinalized(t *testing.T) {
	// Slot before finalized should return false
	finalized := Slot(100)
	candidate := Slot(50)
	if candidate.IsJustifiableAfter(finalized) {
		t.Error("slot before finalized should not be justifiable")
	}
}

func TestIsqrt(t *testing.T) {
	tests := []struct {
		n        uint64
		expected uint64
	}{
		{0, 0},
		{1, 1},
		{2, 1},
		{3, 1},
		{4, 2},
		{8, 2},
		{9, 3},
		{15, 3},
		{16, 4},
		{24, 4},
		{25, 5},
		{100, 10},
		{1000000, 1000},
	}

	for _, tt := range tests {
		got := isqrt(tt.n)
		if got != tt.expected {
			t.Errorf("isqrt(%d) = %d, want %d", tt.n, got, tt.expected)
		}
	}
}

func TestIsPerfectSquare(t *testing.T) {
	squares := []uint64{0, 1, 4, 9, 16, 25, 36, 49, 64, 81, 100, 10000}
	nonSquares := []uint64{2, 3, 5, 6, 7, 8, 10, 11, 12, 99, 101}

	for _, n := range squares {
		if !isPerfectSquare(n) {
			t.Errorf("isPerfectSquare(%d) = false, want true", n)
		}
	}

	for _, n := range nonSquares {
		if isPerfectSquare(n) {
			t.Errorf("isPerfectSquare(%d) = true, want false", n)
		}
	}
}
