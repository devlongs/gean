package types

// State represents the complete state of the Lean Ethereum blockchain
type State struct {
	// The current slot
	Slot uint64
	// The registry of validators
	Validators []Validator
	// The latest block header
	LatestBlockHeader BlockHeader
	// The block roots
	BlockRoots [][32]byte
	// TODO: Add additional state fields as needed
}

// Validator represents a validator in the system
type Validator struct {
	// The validator's public key
	Pubkey [48]byte
	// The validator's index
	Index uint64
	// The validator's balance
	Balance uint64
	// TODO: Add additional validator fields
}
