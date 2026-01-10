package types

// Block represents a complete block in the Lean Ethereum blockchain
type Block struct {
	// The slot in which the block was proposed
	Slot uint64
	// The index of the validator that proposed the block
	ProposerIndex uint64
	// The root of the parent block
	ParentRoot [32]byte
	// The root of the state after applying transactions in this block
	StateRoot [32]byte
	// The block's payload
	Body BlockBody
}

// BlockBody contains the payload data of a block
type BlockBody struct {
	// Validator attestations carried in the block body
	Attestations []Attestation
}

// BlockHeader contains metadata about a block
type BlockHeader struct {
	// The slot in which the block was proposed
	Slot uint64
	// The index of the validator that proposed the block
	ProposerIndex uint64
	// The root of the parent block
	ParentRoot [32]byte
	// The root of the state after applying transactions in this block
	StateRoot [32]byte
	// The root of the block body
	BodyRoot [32]byte
}
