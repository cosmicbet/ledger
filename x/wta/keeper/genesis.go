package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/cosmicbet/ledger/x/wta/types"
)

// ExportGenesis exports the current state of the chain
func (k Keeper) ExportGenesis(ctx sdk.Context) *types.GenesisState {
	return types.NewGenesisState(
		k.GetCurrentDraw(ctx).EndTime,
		k.GetTickets(ctx),
		k.GetHistoricalDrawsData(ctx),
		k.GetDistributionParams(ctx),
		k.GetDrawParams(ctx),
		k.GetTicketParams(ctx),
	)
}

// InitGenesis initializes the given state
func (k Keeper) InitGenesis(ctx sdk.Context, state types.GenesisState) {
	k.SaveCurrentDrawEndTime(ctx, state.DrawEndTime)
	k.SaveTickets(ctx, state.Tickets)

	for _, data := range state.PastDraws {
		k.SaveHistoricalDraw(ctx, data)
	}

	k.SetDistributionParams(ctx, state.DistributionParams)
	k.SetDrawParams(ctx, state.DrawParams)
	k.SetTicketParams(ctx, state.TicketParams)
}
