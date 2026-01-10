package types

// Attestation represents a validator's vote on a block
type Attestation struct {
	// The slot in which the attestation was created
	Slot uint64
	// The index of the validator making the attestation
	ValidatorIndex uint64
	// The root of the block being attested to
	BlockRoot [32]byte
	// The validator's signature
	Signature Signature
}

// Attestations represents a collection of attestations
type Attestations []Attestation
