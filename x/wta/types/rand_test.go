package types_test

import (
	"github.com/cosmicbet/ledger/x/wta/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"
	"testing"
)

func TestNewRandFromSeed(t *testing.T) {
	// Same seed generates same rands
	r1 := types.NewRandFromSeed([]byte("seed"))
	r2 := types.NewRandFromSeed([]byte("seed"))
	for i := 0; i < 100000; i++ {
		require.Equal(t, r1.Intn(1000000), r2.Intn(1000000))
	}
}

func TestNewRandFromCtx(t *testing.T) {
	// Same ctx generates same rands
	ctx := sdk.NewContext(
		nil,
		tmproto.Header{LastCommitHash: []byte("last_commit_hash")},
		false,
		nil,
	).WithTxBytes([]byte("tx_bytes"))

	r1 := types.NewRandFromCtx(ctx)
	r2 := types.NewRandFromCtx(ctx)
	for i := 0; i < 100000; i++ {
		require.Equal(t, r1.Intn(1000000), r2.Intn(1000000))
	}
}

func TestNewRandFromCtxAndIndex(t *testing.T) {
	// Same ctx and index generates same rands
	ctx := sdk.NewContext(
		nil,
		tmproto.Header{LastCommitHash: []byte("last_commit_hash")},
		false,
		nil,
	).WithTxBytes([]byte("tx_bytes"))

	r1 := types.NewRandFromCtxAndIndex(ctx, 1000000)
	r2 := types.NewRandFromCtxAndIndex(ctx, 1000000)
	for i := 0; i < 100000; i++ {
		require.Equal(t, r1.Intn(1000000), r2.Intn(1000000))
	}
}
