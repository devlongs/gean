package ssz

import (
	"crypto/sha256"
	"encoding/binary"

	"github.com/devlongs/gean/common/types"
)

const BytesPerChunk = 32

var ZeroHash = types.Root{}

func Hash(data []byte) types.Root {
	return types.Root(sha256.Sum256(data))
}

func HashNodes(a, b types.Root) types.Root {
	h := sha256.New()
	h.Write(a[:])
	h.Write(b[:])
	var result types.Root
	copy(result[:], h.Sum(nil))
	return result
}

func HashTreeRootUint64(value uint64) types.Root {
	var buf [32]byte
	binary.LittleEndian.PutUint64(buf[:8], value)
	return types.Root(buf)
}

func HashTreeRootBitvector(bv *types.Bitvector) types.Root {
	chunks := packBits(bv.Len(), func(i int) bool { return bv.Get(i) })
	limit := (bv.Len() + 255) / 256
	return Merkleize(chunks, limit)
}

func HashTreeRootBitlist(bl *types.Bitlist) types.Root {
	chunks := packBits(bl.Len(), func(i int) bool { return bl.Get(i) })
	limit := (bl.Limit() + 255) / 256
	root := Merkleize(chunks, limit)
	return MixInLength(root, uint64(bl.Len()))
}

func packBits(n int, get func(int) bool) []types.Root {
	if n == 0 {
		return nil
	}
	byteLen := (n + 7) / 8
	data := make([]byte, byteLen)
	for i := 0; i < n; i++ {
		if get(i) {
			data[i/8] |= 1 << (i % 8)
		}
	}
	padded := make([]byte, ((byteLen+31)/32)*32)
	copy(padded, data)
	chunks := make([]types.Root, len(padded)/32)
	for i := range chunks {
		copy(chunks[i][:], padded[i*32:(i+1)*32])
	}
	return chunks
}

func Merkleize(chunks []types.Root, limit int) types.Root {
	n := len(chunks)
	if n == 0 {
		if limit > 0 {
			return zeroTreeRoot(nextPowerOfTwo(limit))
		}
		return ZeroHash
	}

	width := nextPowerOfTwo(n)
	if limit > 0 && limit >= n {
		width = nextPowerOfTwo(limit)
	}

	if width == 1 {
		return chunks[0]
	}

	level := make([]types.Root, width)
	copy(level, chunks)

	for len(level) > 1 {
		next := make([]types.Root, len(level)/2)
		for i := range next {
			next[i] = HashNodes(level[i*2], level[i*2+1])
		}
		level = next
	}
	return level[0]
}

func MixInLength(root types.Root, length uint64) types.Root {
	var lenChunk types.Root
	binary.LittleEndian.PutUint64(lenChunk[:8], length)
	return HashNodes(root, lenChunk)
}

func nextPowerOfTwo(x int) int {
	if x <= 1 {
		return 1
	}
	n := x - 1
	n |= n >> 1
	n |= n >> 2
	n |= n >> 4
	n |= n >> 8
	n |= n >> 16
	return n + 1
}

func zeroTreeRoot(width int) types.Root {
	if width <= 1 {
		return ZeroHash
	}
	h := ZeroHash
	for width > 1 {
		h = HashNodes(h, h)
		width /= 2
	}
	return h
}

func HashTreeRootBytes(data []byte) types.Root {
	if len(data) <= 32 {
		var chunk types.Root
		copy(chunk[:], data)
		return chunk
	}
	numChunks := (len(data) + 31) / 32
	chunks := make([]types.Root, numChunks)
	for i := range chunks {
		start := i * 32
		end := start + 32
		if end > len(data) {
			end = len(data)
		}
		copy(chunks[i][:], data[start:end])
	}
	return Merkleize(chunks, 0)
}

func HashTreeRootContainer(fields []types.Root) types.Root {
	return Merkleize(fields, 0)
}

func HashTreeRootList(elements []types.Root, limit int) types.Root {
	root := Merkleize(elements, limit)
	return MixInLength(root, uint64(len(elements)))
}

func HashTreeRootCheckpoint(c *types.Checkpoint) types.Root {
	fields := []types.Root{
		c.Root,
		HashTreeRootUint64(uint64(c.Slot)),
	}
	return HashTreeRootContainer(fields)
}

func HashTreeRootValidator(v *types.Validator) types.Root {
	fields := []types.Root{
		HashTreeRootBytes(v.Pubkey[:]),
		HashTreeRootUint64(uint64(v.Index)),
	}
	return HashTreeRootContainer(fields)
}

func HashTreeRootAttestationData(a *types.AttestationData) types.Root {
	fields := []types.Root{
		HashTreeRootUint64(uint64(a.Slot)),
		HashTreeRootCheckpoint(&a.Head),
		HashTreeRootCheckpoint(&a.Target),
		HashTreeRootCheckpoint(&a.Source),
	}
	return HashTreeRootContainer(fields)
}

func HashTreeRootAttestation(a *types.Attestation) types.Root {
	fields := []types.Root{
		HashTreeRootUint64(uint64(a.ValidatorID)),
		HashTreeRootAttestationData(&a.Data),
	}
	return HashTreeRootContainer(fields)
}

func HashTreeRootAggregatedAttestation(a *types.AggregatedAttestation, limit int) types.Root {
	fields := []types.Root{
		HashTreeRootBitlist(a.AggregationBits),
		HashTreeRootAttestationData(&a.Data),
	}
	return HashTreeRootContainer(fields)
}

func HashTreeRootBlockBody(b *types.BlockBody, attestationLimit int) types.Root {
	attRoots := make([]types.Root, len(b.Attestations))
	for i := range b.Attestations {
		attRoots[i] = HashTreeRootAggregatedAttestation(&b.Attestations[i], attestationLimit)
	}
	fields := []types.Root{
		HashTreeRootList(attRoots, attestationLimit),
	}
	return HashTreeRootContainer(fields)
}

func HashTreeRootBlockHeader(h *types.BlockHeader) types.Root {
	fields := []types.Root{
		HashTreeRootUint64(uint64(h.Slot)),
		HashTreeRootUint64(uint64(h.ProposerIndex)),
		h.ParentRoot,
		h.StateRoot,
		h.BodyRoot,
	}
	return HashTreeRootContainer(fields)
}

func HashTreeRootBlock(b *types.Block, attestationLimit int) types.Root {
	fields := []types.Root{
		HashTreeRootUint64(uint64(b.Slot)),
		HashTreeRootUint64(uint64(b.ProposerIndex)),
		b.ParentRoot,
		b.StateRoot,
		HashTreeRootBlockBody(&b.Body, attestationLimit),
	}
	return HashTreeRootContainer(fields)
}

func HashTreeRootBlockWithAttestation(bwa *types.BlockWithAttestation, attestationLimit int) types.Root {
	fields := []types.Root{
		HashTreeRootBlock(&bwa.Block, attestationLimit),
		HashTreeRootAttestation(&bwa.ProposerAttestation),
	}
	return HashTreeRootContainer(fields)
}

func HashTreeRootSignedBlockWithAttestation(sbwa *types.SignedBlockWithAttestation, attestationLimit, validatorLimit int) types.Root {
	sigRoots := make([]types.Root, len(sbwa.Signatures))
	for i := range sbwa.Signatures {
		sigRoots[i] = HashTreeRootBytes(sbwa.Signatures[i][:])
	}
	fields := []types.Root{
		HashTreeRootBlockWithAttestation(&sbwa.Message, attestationLimit),
		HashTreeRootList(sigRoots, validatorLimit),
	}
	return HashTreeRootContainer(fields)
}

func HashTreeRootConfig(c *types.Config) types.Root {
	fields := []types.Root{
		HashTreeRootUint64(c.GenesisTime),
	}
	return HashTreeRootContainer(fields)
}

func HashTreeRootState(s *types.State, historicalRootsLimit, validatorLimit int) types.Root {
	validatorRoots := make([]types.Root, len(s.Validators))
	for i := range s.Validators {
		validatorRoots[i] = HashTreeRootValidator(&s.Validators[i])
	}

	justificationRootsList := make([]types.Root, len(s.JustificationRoots))
	copy(justificationRootsList, s.JustificationRoots)

	fields := []types.Root{
		HashTreeRootConfig(&s.Config),
		HashTreeRootUint64(uint64(s.Slot)),
		HashTreeRootBlockHeader(&s.LatestBlockHeader),
		HashTreeRootCheckpoint(&s.LatestJustified),
		HashTreeRootCheckpoint(&s.LatestFinalized),
		HashTreeRootList(s.HistoricalRoots, historicalRootsLimit),
		HashTreeRootBitlist(s.JustifiedSlots),
		HashTreeRootList(validatorRoots, validatorLimit),
		HashTreeRootList(justificationRootsList, historicalRootsLimit),
		HashTreeRootBitlist(s.JustificationVotes),
	}

	return HashTreeRootContainer(fields)
}
