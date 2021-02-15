package keeper

import (
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	bankkeeper "github.com/cosmos/cosmos-sdk/x/bank/keeper"
	distrkeeper "github.com/cosmos/cosmos-sdk/x/distribution/keeper"
	paramstypes "github.com/cosmos/cosmos-sdk/x/params/types"

	"github.com/cosmicbet/ledger/x/wta/types"
)

// Keeper maintains the link to data storage and exposes getter/setter methods for the various parts of the state machine
type Keeper struct {
	storeKey      sdk.StoreKey
	cdc           codec.BinaryMarshaler
	paramSubspace paramstypes.Subspace

	bk bankkeeper.Keeper
	dk distrkeeper.Keeper
}

// NewKeeper creates new instances of the wta Keeper
func NewKeeper(
	cdc codec.BinaryMarshaler, storeKey sdk.StoreKey, paramSpace paramstypes.Subspace,
	bk bankkeeper.Keeper, dk distrkeeper.Keeper,
) Keeper {
	if !paramSpace.HasKeyTable() {
		paramSpace = paramSpace.WithKeyTable(types.ParamKeyTable())
	}

	return Keeper{
		storeKey:      storeKey,
		cdc:           cdc,
		paramSubspace: paramSpace,

		bk: bk,
		dk: dk,
	}
}

// ------------------------------------------------------------------------------------------------------------------

// WithdrawTicketsCost allows the provided buyer to buy the given quantity of tickets.
func (k Keeper) WithdrawTicketsCost(ctx sdk.Context, quantity int32, buyer sdk.AccAddress) error {
	// Check tickets quantity
	if quantity <= 0 {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidRequest, "invalid amount of tickets: %d", quantity)
	}

	ticketPrice := k.GetParams(ctx).TicketPrice
	ticketsTotal := sdk.NewCoin(ticketPrice.Denom, ticketPrice.Amount.MulRaw(int64(quantity)))

	// Check the user balance
	balance := k.bk.GetBalance(ctx, buyer, ticketsTotal.Denom)
	if balance.IsLT(ticketsTotal) {
		return sdkerrors.Wrapf(sdkerrors.ErrInsufficientFunds, "cannot purchase %d tickets", quantity)
	}

	params := k.GetParams(ctx)

	// Set the prize pool
	prizeAmount := sdk.NewCoin(ticketsTotal.Denom, ticketsTotal.Amount.Mul(params.PrizePercentage).QuoRaw(100))

	err := k.bk.SendCoinsFromAccountToModule(ctx, buyer, types.PrizeCollectorName, sdk.NewCoins(prizeAmount))
	if err != nil {
		return err
	}

	err = k.IncrementDrawPrize(ctx, prizeAmount)
	if err != nil {
		return err
	}

	// Send the pool amount to the pool
	poolAmount := sdk.NewCoin(ticketsTotal.Denom, ticketsTotal.Amount.Mul(params.CommunityPoolPercentage).QuoRaw(100))
	err = k.dk.FundCommunityPool(ctx, sdk.NewCoins(poolAmount), buyer)
	if err != nil {
		return err
	}

	// Send the tokens to burn to the proper module
	burnAmount := sdk.NewCoin(ticketsTotal.Denom, ticketsTotal.Amount.Mul(params.BurnPercentage).QuoRaw(100))
	err = k.bk.SendCoinsFromAccountToModule(ctx, buyer, types.PrizeBurnerName, sdk.NewCoins(burnAmount))
	if err != nil {
		return err
	}

	// Burn the tokens
	return k.bk.BurnCoins(ctx, types.PrizeBurnerName, sdk.NewCoins(burnAmount))
}

// SaveTickets sets the given tickets for the given user
func (k Keeper) SaveTickets(ctx sdk.Context, tickets []types.Ticket) error {
	store := ctx.KVStore(k.storeKey)
	for _, t := range tickets {
		bz, err := types.MarshalTicket(k.cdc, t)
		if err != nil {
			return err
		}

		store.Set(types.TicketsStoreKey(t.Id), bz)
	}
	return nil
}

// WipeCurrentTickets removes all the currently stored tickets
func (k Keeper) WipeCurrentTickets(ctx sdk.Context) {
	store := ctx.KVStore(k.storeKey)

	tickets := k.GetTickets(ctx)
	for _, t := range tickets {
		store.Delete(types.TicketsStoreKey(t.Id))
	}
}

// ------------------------------------------------------------------------------------------------------------------

// IncrementDrawPrize increments the current draw prize to the provided amount
func (k Keeper) IncrementDrawPrize(ctx sdk.Context, amount sdk.Coin) error {
	draw := k.GetCurrentDraw(ctx)
	draw.Prize = draw.Prize.Add(amount)

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypePrizeIncrease,
			sdk.NewAttribute(types.AttributeKeyPrizeAmount, draw.Prize.String()),
		),
	)

	return k.SaveCurrentDraw(ctx, draw)
}

// TransferDrawPrize transfers the provided prize to the specified winner account
func (k Keeper) TransferDrawPrize(ctx sdk.Context, prize sdk.Coins, winner sdk.AccAddress) error {
	return k.bk.SendCoinsFromModuleToAccount(ctx, types.PrizeCollectorName, winner, prize)
}

// SaveCurrentDraw stores the given draw as the next draw
func (k Keeper) SaveCurrentDraw(ctx sdk.Context, draw types.Draw) error {
	store := ctx.KVStore(k.storeKey)
	bz, err := types.MarshalDraw(k.cdc, draw)
	if err != nil {
		return err
	}

	store.Set(types.CurrentDrawStoreKey, bz)
	return nil
}

// GetCurrentDraw returns the Draw for which the tickets can be currently bought
func (k Keeper) GetCurrentDraw(ctx sdk.Context) types.Draw {
	store := ctx.KVStore(k.storeKey)
	return types.MustUnmarshalDraw(k.cdc, store.Get(types.CurrentDrawStoreKey))
}

// ------------------------------------------------------------------------------------------------------------------

// SaveHistoricalDraw saves the given draw as an historical draw
func (k Keeper) SaveHistoricalDraw(ctx sdk.Context, draw types.HistoricalDrawData) {
	store := ctx.KVStore(k.storeKey)
	store.Set(types.HistoricalDataStoreKey(draw.Draw.EndTime), types.MustMarshalHistoricalDraw(k.cdc, draw))
}