package types

type Slot uint64
type ValidatorIndex uint64
type Epoch uint64
type Root [32]byte

type Bytes4 [4]byte
type Bytes20 [20]byte
type Bytes32 = Root
type Bytes48 [48]byte
type Bytes52 [52]byte
type Bytes96 [96]byte

const SecondsPerSlot uint64 = 4

func (r Root) IsZero() bool {
	return r == Root{}
}

// IsJustifiableAfter checks if this slot is a valid candidate for justification
// after a given finalized slot according to the 3SF-mini specification.
//
// A slot is justifiable if its distance (delta) from the finalized slot is:
//  1. Less than or equal to 5 (first 5 slots always justifiable)
//  2. A perfect square (e.g., 9, 16, 25...)
//  3. A pronic number of the form n*(n+1) (e.g., 6, 12, 20, 30...)
func (s Slot) IsJustifiableAfter(finalizedSlot Slot) bool {
	if s < finalizedSlot {
		return false
	}

	delta := uint64(s - finalizedSlot)

	// Rule 1: The first 5 slots after finalization are always justifiable
	if delta <= 5 {
		return true
	}

	// Rule 2: Perfect square distances are justifiable
	// Check: isqrt(delta)^2 == delta
	if isPerfectSquare(delta) {
		return true
	}

	// Rule 3: Pronic number distances are justifiable
	// Pronic numbers have the form n*(n+1): 2, 6, 12, 20, 30, 42, 56, ...
	// Mathematical insight: For pronic delta = n*(n+1), we have:
	//   4*delta + 1 = 4*n*(n+1) + 1 = (2n+1)^2
	// So 4*delta+1 must be an odd perfect square
	check := 4*delta + 1
	if isPerfectSquare(check) && isqrt(check)%2 == 1 {
		return true
	}

	return false
}

// isqrt computes the integer square root (floor of sqrt).
func isqrt(n uint64) uint64 {
	if n == 0 {
		return 0
	}
	x := n
	y := (x + 1) / 2
	for y < x {
		x = y
		y = (x + n/x) / 2
	}
	return x
}

// isPerfectSquare checks if n is a perfect square.
func isPerfectSquare(n uint64) bool {
	root := isqrt(n)
	return root*root == n
}

func SlotToTime(slot Slot, genesisTime uint64) uint64 {
	return genesisTime + uint64(slot)*SecondsPerSlot
}

func TimeToSlot(time, genesisTime uint64) Slot {
	if time < genesisTime {
		return 0
	}
	return Slot((time - genesisTime) / SecondsPerSlot)
}
