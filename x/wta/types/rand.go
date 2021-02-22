package types

import (
	"crypto/sha256"
	"encoding/binary"
	"math/rand"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// NewRandFromSeed returns a new deterministic rand.Rand using the given byte slice as seed
func NewRandFromSeed(seed []byte) *rand.Rand {
	hash := sha256.Sum256(seed)
	return rand.New(rand.NewSource(int64(binary.BigEndian.Uint64(hash[8:]))))
}

// newSeedFromCtx returns a new seed based on the given context
func newSeedFromCtx(ctx sdk.Context) []byte {
	return append(ctx.BlockHeader().LastCommitHash, ctx.TxBytes()...)
}

// NewRandFromCtx returns a new rand.Rand based on the given context
func NewRandFromCtx(ctx sdk.Context) *rand.Rand {
	return NewRandFromSeed(newSeedFromCtx(ctx))
}

// NewRandFromCtxAndIndex returns a new rand.Rand based on the given context and index
func NewRandFromCtxAndIndex(ctx sdk.Context, i int) *rand.Rand {
	var index = make([]byte, 8)
	binary.BigEndian.PutUint64(index, uint64(i))

	seed := append(newSeedFromCtx(ctx), index...)
	return NewRandFromSeed(seed)
}
