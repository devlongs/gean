package types

import (
	"testing"
)

func TestCheckpoint(t *testing.T) {
	c := Checkpoint{
		Root: Root{1, 2, 3},
		Slot: Slot(100),
	}

	if c.Slot != 100 {
		t.Errorf("expected slot 100, got %d", c.Slot)
	}
	if c.Root[0] != 1 {
		t.Errorf("expected root[0] = 1, got %d", c.Root[0])
	}
}

func TestValidator(t *testing.T) {
	var pubkey Bytes52
	pubkey[0] = 0xAB

	v := Validator{
		Pubkey: pubkey,
		Index:  ValidatorIndex(42),
	}

	if v.Index != 42 {
		t.Errorf("expected index 42, got %d", v.Index)
	}
	if v.Pubkey[0] != 0xAB {
		t.Errorf("expected pubkey[0] = 0xAB, got %x", v.Pubkey[0])
	}
}

func TestAttestationData(t *testing.T) {
	data := AttestationData{
		Slot: Slot(10),
		Head: Checkpoint{Root: Root{1}, Slot: Slot(10)},
		Target: Checkpoint{Root: Root{2}, Slot: Slot(8)},
		Source: Checkpoint{Root: Root{3}, Slot: Slot(4)},
	}

	if data.Slot != 10 {
		t.Errorf("expected slot 10, got %d", data.Slot)
	}
	if data.Head.Slot != 10 {
		t.Errorf("expected head slot 10, got %d", data.Head.Slot)
	}
	if data.Target.Slot != 8 {
		t.Errorf("expected target slot 8, got %d", data.Target.Slot)
	}
	if data.Source.Slot != 4 {
		t.Errorf("expected source slot 4, got %d", data.Source.Slot)
	}
}

func TestAttestation(t *testing.T) {
	att := Attestation{
		ValidatorID: ValidatorIndex(5),
		Data: AttestationData{
			Slot:   Slot(10),
			Head:   Checkpoint{Root: Root{1}, Slot: Slot(10)},
			Target: Checkpoint{Root: Root{2}, Slot: Slot(8)},
			Source: Checkpoint{Root: Root{3}, Slot: Slot(4)},
		},
	}

	if att.ValidatorID != 5 {
		t.Errorf("expected validator 5, got %d", att.ValidatorID)
	}
}

func TestBlockHeader(t *testing.T) {
	header := BlockHeader{
		Slot:          Slot(100),
		ProposerIndex: ValidatorIndex(7),
		ParentRoot:    Root{1, 2, 3},
		StateRoot:     Root{4, 5, 6},
		BodyRoot:      Root{7, 8, 9},
	}

	if header.Slot != 100 {
		t.Errorf("expected slot 100, got %d", header.Slot)
	}
	if header.ProposerIndex != 7 {
		t.Errorf("expected proposer 7, got %d", header.ProposerIndex)
	}
}

func TestBlock(t *testing.T) {
	block := Block{
		Slot:          Slot(100),
		ProposerIndex: ValidatorIndex(7),
		ParentRoot:    Root{1, 2, 3},
		StateRoot:     Root{4, 5, 6},
		Body: BlockBody{
			Attestations: []AggregatedAttestation{},
		},
	}

	if block.Slot != 100 {
		t.Errorf("expected slot 100, got %d", block.Slot)
	}
	if len(block.Body.Attestations) != 0 {
		t.Errorf("expected 0 attestations, got %d", len(block.Body.Attestations))
	}
}

func TestState(t *testing.T) {
	state := State{
		Config:      Config{GenesisTime: 1000},
		Slot:        Slot(50),
		Validators:  []Validator{},
		LatestJustified: Checkpoint{},
		LatestFinalized: Checkpoint{},
	}

	if state.Slot != 50 {
		t.Errorf("expected slot 50, got %d", state.Slot)
	}
	if state.Config.GenesisTime != 1000 {
		t.Errorf("expected genesis time 1000, got %d", state.Config.GenesisTime)
	}
}

func TestBlockWithAttestation(t *testing.T) {
	bwa := BlockWithAttestation{
		Block: Block{
			Slot:          Slot(100),
			ProposerIndex: ValidatorIndex(7),
			ParentRoot:    Root{1, 2, 3},
			StateRoot:     Root{4, 5, 6},
			Body:          BlockBody{Attestations: []AggregatedAttestation{}},
		},
		ProposerAttestation: Attestation{
			ValidatorID: ValidatorIndex(7),
			Data: AttestationData{
				Slot:   Slot(100),
				Head:   Checkpoint{Root: Root{1}, Slot: Slot(100)},
				Target: Checkpoint{Root: Root{2}, Slot: Slot(96)},
				Source: Checkpoint{Root: Root{3}, Slot: Slot(64)},
			},
		},
	}

	if bwa.Block.Slot != 100 {
		t.Errorf("expected block slot 100, got %d", bwa.Block.Slot)
	}
	if bwa.ProposerAttestation.ValidatorID != 7 {
		t.Errorf("expected proposer validator 7, got %d", bwa.ProposerAttestation.ValidatorID)
	}
}

func TestSignedBlockWithAttestation(t *testing.T) {
	sbwa := SignedBlockWithAttestation{
		Message: BlockWithAttestation{
			Block: Block{
				Slot:          Slot(100),
				ProposerIndex: ValidatorIndex(7),
				ParentRoot:    Root{1, 2, 3},
				StateRoot:     Root{4, 5, 6},
				Body:          BlockBody{Attestations: []AggregatedAttestation{}},
			},
			ProposerAttestation: Attestation{
				ValidatorID: ValidatorIndex(7),
				Data: AttestationData{
					Slot:   Slot(100),
					Head:   Checkpoint{Root: Root{1}, Slot: Slot(100)},
					Target: Checkpoint{Root: Root{2}, Slot: Slot(96)},
					Source: Checkpoint{Root: Root{3}, Slot: Slot(64)},
				},
			},
		},
		Signatures: []Bytes3116{},
	}

	if sbwa.Message.Block.Slot != 100 {
		t.Errorf("expected block slot 100, got %d", sbwa.Message.Block.Slot)
	}
	if len(sbwa.Signatures) != 0 {
		t.Errorf("expected 0 signatures, got %d", len(sbwa.Signatures))
	}
}
