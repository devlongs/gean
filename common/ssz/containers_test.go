package ssz

import (
	"testing"

	"github.com/devlongs/gean/common/types"
)

func TestHashTreeRootCheckpoint(t *testing.T) {
	c := &types.Checkpoint{
		Root: types.Root{1, 2, 3},
		Slot: types.Slot(100),
	}

	root := HashTreeRootCheckpoint(c)

	// Verify it produces a non-zero root
	if root == ZeroHash {
		t.Error("expected non-zero root for checkpoint")
	}

	// Same checkpoint should produce same root
	root2 := HashTreeRootCheckpoint(c)
	if root != root2 {
		t.Error("same checkpoint should produce same root")
	}

	// Different checkpoint should produce different root
	c2 := &types.Checkpoint{
		Root: types.Root{4, 5, 6},
		Slot: types.Slot(200),
	}
	root3 := HashTreeRootCheckpoint(c2)
	if root == root3 {
		t.Error("different checkpoints should produce different roots")
	}
}

func TestHashTreeRootValidator(t *testing.T) {
	var pubkey types.Bytes52
	pubkey[0] = 0xAB

	v := &types.Validator{
		Pubkey: pubkey,
		Index:  types.ValidatorIndex(42),
	}

	root := HashTreeRootValidator(v)

	if root == ZeroHash {
		t.Error("expected non-zero root for validator")
	}
}

func TestHashTreeRootAttestationData(t *testing.T) {
	data := &types.AttestationData{
		Slot: types.Slot(10),
		Head: types.Checkpoint{Root: types.Root{1}, Slot: types.Slot(10)},
		Target: types.Checkpoint{Root: types.Root{2}, Slot: types.Slot(8)},
		Source: types.Checkpoint{Root: types.Root{3}, Slot: types.Slot(4)},
	}

	root := HashTreeRootAttestationData(data)

	if root == ZeroHash {
		t.Error("expected non-zero root for attestation data")
	}
}

func TestHashTreeRootBlockHeader(t *testing.T) {
	header := &types.BlockHeader{
		Slot:          types.Slot(100),
		ProposerIndex: types.ValidatorIndex(7),
		ParentRoot:    types.Root{1, 2, 3},
		StateRoot:     types.Root{4, 5, 6},
		BodyRoot:      types.Root{7, 8, 9},
	}

	root := HashTreeRootBlockHeader(header)

	if root == ZeroHash {
		t.Error("expected non-zero root for block header")
	}

	// Same header should produce same root
	root2 := HashTreeRootBlockHeader(header)
	if root != root2 {
		t.Error("same header should produce same root")
	}
}

func TestHashTreeRootBlock(t *testing.T) {
	block := &types.Block{
		Slot:          types.Slot(100),
		ProposerIndex: types.ValidatorIndex(7),
		ParentRoot:    types.Root{1, 2, 3},
		StateRoot:     types.Root{4, 5, 6},
		Body: types.BlockBody{
			Attestations: []types.AggregatedAttestation{},
		},
	}

	root := HashTreeRootBlock(block, 128)

	if root == ZeroHash {
		t.Error("expected non-zero root for block")
	}
}

func TestHashTreeRootBytes(t *testing.T) {
	// Test small data (fits in one chunk)
	small := []byte{1, 2, 3, 4}
	root := HashTreeRootBytes(small)

	var expected types.Root
	copy(expected[:], small)
	if root != expected {
		t.Error("small byte array should be zero-padded to 32 bytes")
	}

	// Test larger data (multiple chunks)
	large := make([]byte, 64)
	for i := range large {
		large[i] = byte(i)
	}
	root2 := HashTreeRootBytes(large)
	if root2 == ZeroHash {
		t.Error("expected non-zero root for large byte array")
	}
}

func TestHashTreeRootContainer(t *testing.T) {
	// Empty container
	empty := HashTreeRootContainer([]types.Root{})
	if empty != ZeroHash {
		t.Error("empty container should have zero root")
	}

	// Single field
	single := HashTreeRootContainer([]types.Root{{1, 2, 3}})
	expected := types.Root{1, 2, 3}
	if single != expected {
		t.Error("single field container should return that field")
	}

	// Two fields should be hashed together
	two := HashTreeRootContainer([]types.Root{{1}, {2}})
	if two == ZeroHash {
		t.Error("expected non-zero root for two-field container")
	}
}

func TestHashTreeRootList(t *testing.T) {
	// Empty list
	empty := HashTreeRootList([]types.Root{}, 128)
	// Empty list mixes in length 0
	if empty == ZeroHash {
		t.Error("empty list should have non-zero root due to length mixing")
	}

	// List with elements
	elements := []types.Root{{1}, {2}, {3}}
	root := HashTreeRootList(elements, 128)
	if root == ZeroHash {
		t.Error("expected non-zero root for list with elements")
	}
}

func TestHashTreeRootBlockWithAttestation(t *testing.T) {
	bwa := &types.BlockWithAttestation{
		Block: types.Block{
			Slot:          types.Slot(100),
			ProposerIndex: types.ValidatorIndex(7),
			ParentRoot:    types.Root{1, 2, 3},
			StateRoot:     types.Root{4, 5, 6},
			Body:          types.BlockBody{Attestations: []types.AggregatedAttestation{}},
		},
		ProposerAttestation: types.Attestation{
			ValidatorID: types.ValidatorIndex(7),
			Data: types.AttestationData{
				Slot:   types.Slot(100),
				Head:   types.Checkpoint{Root: types.Root{1}, Slot: types.Slot(100)},
				Target: types.Checkpoint{Root: types.Root{2}, Slot: types.Slot(96)},
				Source: types.Checkpoint{Root: types.Root{3}, Slot: types.Slot(64)},
			},
		},
	}

	root := HashTreeRootBlockWithAttestation(bwa, 128)
	if root == ZeroHash {
		t.Error("expected non-zero root for block with attestation")
	}

	// Same input should produce same root
	root2 := HashTreeRootBlockWithAttestation(bwa, 128)
	if root != root2 {
		t.Error("same block with attestation should produce same root")
	}
}

func TestHashTreeRootSignedBlockWithAttestation(t *testing.T) {
	sbwa := &types.SignedBlockWithAttestation{
		Message: types.BlockWithAttestation{
			Block: types.Block{
				Slot:          types.Slot(100),
				ProposerIndex: types.ValidatorIndex(7),
				ParentRoot:    types.Root{1, 2, 3},
				StateRoot:     types.Root{4, 5, 6},
				Body:          types.BlockBody{Attestations: []types.AggregatedAttestation{}},
			},
			ProposerAttestation: types.Attestation{
				ValidatorID: types.ValidatorIndex(7),
				Data: types.AttestationData{
					Slot:   types.Slot(100),
					Head:   types.Checkpoint{Root: types.Root{1}, Slot: types.Slot(100)},
					Target: types.Checkpoint{Root: types.Root{2}, Slot: types.Slot(96)},
					Source: types.Checkpoint{Root: types.Root{3}, Slot: types.Slot(64)},
				},
			},
		},
		Signatures: []types.Bytes3116{},
	}

	root := HashTreeRootSignedBlockWithAttestation(sbwa, 128, 4096)
	if root == ZeroHash {
		t.Error("expected non-zero root for signed block with attestation")
	}
}

func TestHashTreeRootConfig(t *testing.T) {
	config := &types.Config{GenesisTime: 1700000000}

	root := HashTreeRootConfig(config)
	if root == ZeroHash {
		t.Error("expected non-zero root for config")
	}

	// Different genesis time should produce different root
	config2 := &types.Config{GenesisTime: 1800000000}
	root2 := HashTreeRootConfig(config2)
	if root == root2 {
		t.Error("different configs should produce different roots")
	}
}

func TestHashTreeRootState(t *testing.T) {
	// Create minimal valid state
	justifiedSlots, _ := types.BitlistFromBits([]bool{true, false, true}, 262144)
	justificationVotes, _ := types.BitlistFromBits([]bool{}, 262144*4096)

	state := &types.State{
		Config:          types.Config{GenesisTime: 1700000000},
		Slot:            types.Slot(100),
		LatestBlockHeader: types.BlockHeader{
			Slot:          types.Slot(99),
			ProposerIndex: types.ValidatorIndex(5),
			ParentRoot:    types.Root{1},
			StateRoot:     types.Root{2},
			BodyRoot:      types.Root{3},
		},
		LatestJustified:    types.Checkpoint{Root: types.Root{10}, Slot: types.Slot(96)},
		LatestFinalized:    types.Checkpoint{Root: types.Root{20}, Slot: types.Slot(64)},
		HistoricalRoots:    []types.Root{{1}, {2}, {3}},
		JustifiedSlots:     justifiedSlots,
		Validators:         []types.Validator{},
		JustificationRoots: []types.Root{},
		JustificationVotes: justificationVotes,
	}

	root := HashTreeRootState(state, 262144, 4096)
	if root == ZeroHash {
		t.Error("expected non-zero root for state")
	}

	// Same state should produce same root
	root2 := HashTreeRootState(state, 262144, 4096)
	if root != root2 {
		t.Error("same state should produce same root")
	}
}

func TestHashTreeRootStateWithValidators(t *testing.T) {
	justifiedSlots, _ := types.BitlistFromBits([]bool{true}, 262144)
	justificationVotes, _ := types.BitlistFromBits([]bool{}, 262144*4096)

	var pubkey1, pubkey2 types.Bytes52
	pubkey1[0] = 0xAA
	pubkey2[0] = 0xBB

	state := &types.State{
		Config:          types.Config{GenesisTime: 1700000000},
		Slot:            types.Slot(100),
		LatestBlockHeader: types.BlockHeader{
			Slot:          types.Slot(99),
			ProposerIndex: types.ValidatorIndex(0),
			ParentRoot:    types.Root{},
			StateRoot:     types.Root{},
			BodyRoot:      types.Root{},
		},
		LatestJustified:    types.Checkpoint{},
		LatestFinalized:    types.Checkpoint{},
		HistoricalRoots:    []types.Root{},
		JustifiedSlots:     justifiedSlots,
		Validators: []types.Validator{
			{Pubkey: pubkey1, Index: types.ValidatorIndex(0)},
			{Pubkey: pubkey2, Index: types.ValidatorIndex(1)},
		},
		JustificationRoots: []types.Root{},
		JustificationVotes: justificationVotes,
	}

	root := HashTreeRootState(state, 262144, 4096)
	if root == ZeroHash {
		t.Error("expected non-zero root for state with validators")
	}
}
