package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/cosmicbet/ledger/x/wta/types"
)

// GetDistributionParams returns the current DistributionParams from the global param store
func (k Keeper) GetDistributionParams(ctx sdk.Context) types.DistributionParams {
	var p types.DistributionParams
	k.paramSubspace.Get(ctx, types.ParamStoreDistributionParamsKey, &p)
	return p
}

// GetDrawParams returns the current DrawParams from the global param store
func (k Keeper) GetDrawParams(ctx sdk.Context) types.DrawParams {
	var p types.DrawParams
	k.paramSubspace.Get(ctx, types.ParamStoreDrawParamsKey, &p)
	return p
}

// GetTicketParams returns the current TicketParams from the global param store
func (k Keeper) GetTicketParams(ctx sdk.Context) types.TicketParams {
	var p types.TicketParams
	k.paramSubspace.Get(ctx, types.ParamStoreTicketParamsKey, &p)
	return p
}

// SetDistributionParams sets DistributionParams to the global param store
func (k Keeper) SetDistributionParams(ctx sdk.Context, params types.DistributionParams) {
	k.paramSubspace.Set(ctx, types.ParamStoreDistributionParamsKey, &params)
}

// SetDrawParams sets DrawParams to the global param store
func (k Keeper) SetDrawParams(ctx sdk.Context, params types.DrawParams) {
	k.paramSubspace.Set(ctx, types.ParamStoreDrawParamsKey, &params)
}

// SetTicketParams sets TicketParams to the global param store
func (k Keeper) SetTicketParams(ctx sdk.Context, params types.TicketParams) {
	k.paramSubspace.Set(ctx, types.ParamStoreTicketParamsKey, &params)
}
