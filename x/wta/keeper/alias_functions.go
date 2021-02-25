package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/cosmicbet/ledger/x/wta/types"
)

// IterateTickets iterates through the tickets and performs the provided function
func (k Keeper) IterateTickets(ctx sdk.Context, fn func(index int64, ticket types.Ticket) (stop bool)) {
	store := ctx.KVStore(k.storeKey)

	iterator := sdk.KVStorePrefixIterator(store, types.TicketsStorePrefix)
	defer iterator.Close()

	i := int64(0)
	for ; iterator.Valid(); iterator.Next() {
		ticket := types.MustUnmarshalTicket(k.cdc, iterator.Value())

		stop := fn(i, ticket)
		if stop {
			break
		}
		i++
	}
}

// GetTickets returns the list of tickets currently stored
func (k Keeper) GetTickets(ctx sdk.Context) []types.Ticket {
	var tickets []types.Ticket
	k.IterateTickets(ctx, func(_ int64, ticket types.Ticket) (stop bool) {
		tickets = append(tickets, ticket)
		return false
	})
	return tickets
}

// GetDrawParticipantsAndTickets returns the list of participants that have entered the draw,
// and the list of all tickets sold for such draw
func (k Keeper) GetDrawParticipantsAndTickets(ctx sdk.Context) (participants []string, ticketsSold []types.Ticket) {
	participantsAddresses := map[string]bool{}
	k.IterateTickets(ctx, func(index int64, ticket types.Ticket) (stop bool) {
		if !participantsAddresses[ticket.Owner] {
			participants = append(participants, ticket.Owner)
			participantsAddresses[ticket.Owner] = true
		}

		ticketsSold = append(ticketsSold, ticket)
		return false
	})

	return participants, ticketsSold
}

// IterateHistoricalDrawsData iterates through the historical data and performs the provided function
func (k Keeper) IterateHistoricalDrawsData(ctx sdk.Context, fn func(index int64, data types.HistoricalDrawData) (stop bool)) {
	store := ctx.KVStore(k.storeKey)

	iterator := sdk.KVStorePrefixIterator(store, types.HistoricalDrawStorePrefix)
	defer iterator.Close()

	i := int64(0)
	for ; iterator.Valid(); iterator.Next() {
		historicalDraw := types.MustUnmarshalHistoricalDrawData(k.cdc, iterator.Value())

		stop := fn(i, historicalDraw)
		if stop {
			break
		}
		i++
	}
}

// GetHistoricalDrawsData returns the list of historical draws data currently stored
func (k Keeper) GetHistoricalDrawsData(ctx sdk.Context) []types.HistoricalDrawData {
	var historicalDraws []types.HistoricalDrawData
	k.IterateHistoricalDrawsData(ctx, func(_ int64, data types.HistoricalDrawData) (stop bool) {
		historicalDraws = append(historicalDraws, data)
		return false
	})
	return historicalDraws
}
