package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/cosmicbet/ledger/x/wta/types"
)

// ExportGenesis exports the current state of the chain
func (k Keeper) ExportGenesis(ctx sdk.Context) *types.GenesisState {
	return types.NewGenesisState(
		k.GetCurrentDraw(ctx),
		k.GetTickets(ctx),
		k.GetParams(ctx),
	)
}

// InitGenesis initializes the given state
func (k Keeper) InitGenesis(ctx sdk.Context, state types.GenesisState) {
	err := k.SaveCurrentDraw(ctx, state.Draw)
	if err != nil {
		panic(err)
	}

	err = k.SaveTickets(ctx, state.Tickets)
	if err != nil {
		panic(err)
	}

	k.SetParams(ctx, state.Params)
}