package keeper

import (
	"time"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	authkeeper "github.com/cosmos/cosmos-sdk/x/auth/keeper"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
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

	ak authkeeper.AccountKeeper
	bk bankkeeper.Keeper
	dk distrkeeper.Keeper

	feeCollectorName string // name of the FeeCollector ModuleAccount
}

// NewKeeper creates new instances of the wta Keeper
func NewKeeper(
	cdc codec.BinaryMarshaler, storeKey sdk.StoreKey, paramSpace paramstypes.Subspace,
	ak authkeeper.AccountKeeper, bk bankkeeper.Keeper, dk distrkeeper.Keeper,
	feeCollectorName string,
) Keeper {
	if !paramSpace.HasKeyTable() {
		paramSpace = paramSpace.WithKeyTable(types.ParamKeyTable())
	}

	return Keeper{
		storeKey:      storeKey,
		cdc:           cdc,
		paramSubspace: paramSpace,

		ak: ak,
		bk: bk,
		dk: dk,

		feeCollectorName: feeCollectorName,
	}
}

// ------------------------------------------------------------------------------------------------------------------

// WithdrawTicketsCost allows the provided buyer to buy the given quantity of tickets.
func (k Keeper) WithdrawTicketsCost(ctx sdk.Context, quantity uint32, buyer sdk.AccAddress) error {
	// Check tickets quantity
	if quantity <= 0 {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidRequest, "invalid amount of tickets: %d", quantity)
	}

	ticketPrice := k.GetTicketParams(ctx).Price
	ticketsTotal := sdk.NewCoin(ticketPrice.Denom, ticketPrice.Amount.MulRaw(int64(quantity)))

	// Check the user balance
	balance := k.bk.GetBalance(ctx, buyer, ticketsTotal.Denom)
	if balance.IsLT(ticketsTotal) {
		return sdkerrors.Wrapf(sdkerrors.ErrInsufficientFunds, "cannot purchase %d tickets", quantity)
	}

	params := k.GetDistributionParams(ctx)

	// Update the prize pool
	prizeAmount := ticketsTotal.Amount.ToDec().Mul(params.PrizePercentage).RoundInt()
	prizeCoin := sdk.NewCoin(ticketsTotal.Denom, prizeAmount)
	err := k.bk.SendCoinsFromAccountToModule(ctx, buyer, types.PrizeCollectorName, sdk.NewCoins(prizeCoin))
	if err != nil {
		return err
	}

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypePrizeIncrease,
			sdk.NewAttribute(types.AttributeKeyPrizeAmount, prizeCoin.String()),
		),
	)

	// Send the fee amount to the fee pool
	feeAmount := ticketsTotal.Amount.ToDec().Mul(params.FeePercentage).RoundInt()
	feeCoin := sdk.NewCoin(ticketsTotal.Denom, feeAmount)
	err = k.bk.SendCoinsFromAccountToModule(ctx, buyer, k.feeCollectorName, sdk.NewCoins(feeCoin))
	if err != nil {
		return err
	}

	// Burn the tokens
	burnCoin := ticketsTotal.Sub(prizeCoin).Sub(feeCoin)
	err = k.bk.SendCoinsFromAccountToModule(ctx, buyer, types.PrizeBurnerName, sdk.NewCoins(burnCoin))
	if err != nil {
		return err
	}

	return k.bk.BurnCoins(ctx, types.PrizeBurnerName, sdk.NewCoins(burnCoin))
}

// SaveTickets sets the given tickets for the given user
func (k Keeper) SaveTickets(ctx sdk.Context, tickets []types.Ticket) {
	store := ctx.KVStore(k.storeKey)
	for _, t := range tickets {
		store.Set(types.TicketsStoreKey(t.Id), types.MustMarshalTicket(k.cdc, t))
	}
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

// TransferDrawPrize transfers the provided prize to the specified winner account
func (k Keeper) TransferDrawPrize(ctx sdk.Context, prize sdk.Coins, winner sdk.AccAddress) error {
	return k.bk.SendCoinsFromModuleToAccount(ctx, types.PrizeCollectorName, winner, prize)
}

// SaveCurrentDraw stores the given draw as the next draw
func (k Keeper) SaveCurrentDrawEndTime(ctx sdk.Context, endTime time.Time) {
	store := ctx.KVStore(k.storeKey)
	store.Set(types.CurrentDrawEndTimeStoreKey, types.MustMarshalDrawEndTime(endTime))
}

// GetCurrentDraw returns the Draw for which the tickets can be currently bought
func (k Keeper) GetCurrentDraw(ctx sdk.Context) types.Draw {
	store := ctx.KVStore(k.storeKey)
	endTime := types.MustUnmarshalDrawEndTime(store.Get(types.CurrentDrawEndTimeStoreKey))

	acc := authtypes.NewModuleAddress(types.PrizeCollectorName)
	prize := k.bk.GetAllBalances(ctx, acc)

	participants, ticketsSold := k.GetDrawParticipantsAndTickets(ctx)
	return types.NewDraw(uint32(len(participants)), uint32(len(ticketsSold)), prize, endTime)
}

// ------------------------------------------------------------------------------------------------------------------

// SaveHistoricalDraw saves the given draw as an historical draw
func (k Keeper) SaveHistoricalDraw(ctx sdk.Context, draw types.HistoricalDrawData) {
	store := ctx.KVStore(k.storeKey)
	store.Set(types.HistoricalDataStoreKey(draw.Draw.EndTime), types.MustMarshalHistoricalDraw(k.cdc, draw))
}
