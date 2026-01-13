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
type Bytes3116 [3116]byte // XMSS signature size

const SecondsPerSlot uint64 = 4

func (r Root) IsZero() bool {
	return r == Root{}
}

// IsJustifiableAfter implements 3SF-mini finality rules.
// A slot is justifiable if delta (distance from finalized) is:
//   - <= 5, OR a perfect square, OR a pronic number n*(n+1)
func (s Slot) IsJustifiableAfter(finalizedSlot Slot) bool {
	if s < finalizedSlot {
		return false
	}

	delta := uint64(s - finalizedSlot)

	if delta <= 5 {
		return true
	}

	if isPerfectSquare(delta) {
		return true
	}

	// Pronic check: 4*delta+1 must be an odd perfect square
	check := 4*delta + 1
	if isPerfectSquare(check) && isqrt(check)%2 == 1 {
		return true
	}

	return false
}

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
