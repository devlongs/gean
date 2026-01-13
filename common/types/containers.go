package types

type Checkpoint struct {
	Root Root
	Slot Slot
}

type Validator struct {
	Pubkey Bytes52
	Index  ValidatorIndex
}

type AttestationData struct {
	Slot   Slot
	Head   Checkpoint
	Target Checkpoint
	Source Checkpoint
}

type Attestation struct {
	ValidatorID ValidatorIndex
	Data        AttestationData
}

type SignedAttestation struct {
	ValidatorID ValidatorIndex
	Message     AttestationData
	Signature   Bytes3116
}

type AggregatedAttestation struct {
	AggregationBits *Bitlist
	Data            AttestationData
}

type BlockWithAttestation struct {
	Block               Block
	ProposerAttestation Attestation
}

type SignedBlockWithAttestation struct {
	Message    BlockWithAttestation
	Signatures []Bytes3116
}

type BlockBody struct {
	Attestations []AggregatedAttestation
}

type BlockHeader struct {
	Slot          Slot
	ProposerIndex ValidatorIndex
	ParentRoot    Root
	StateRoot     Root
	BodyRoot      Root
}

type Block struct {
	Slot          Slot
	ProposerIndex ValidatorIndex
	ParentRoot    Root
	StateRoot     Root
	Body          BlockBody
}

type Config struct {
	GenesisTime uint64
}

type State struct {
	Config             Config
	Slot               Slot
	LatestBlockHeader  BlockHeader
	LatestJustified    Checkpoint
	LatestFinalized    Checkpoint
	HistoricalRoots    []Root
	JustifiedSlots     *Bitlist
	Validators         []Validator
	JustificationRoots []Root
	JustificationVotes *Bitlist
}
