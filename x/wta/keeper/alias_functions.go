package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/cosmicbet/ledger/x/wta/types"
)

// Iterate through the tickets and perform the provided function
func (k Keeper) IterateTickets(ctx sdk.Context, fn func(index int64, ticket types.Ticket) (stop bool)) {
	store := ctx.KVStore(k.storeKey)

	iterator := sdk.KVStoreReversePrefixIterator(store, types.TicketsStorePrefix)
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
	k.IterateTickets(ctx, func(index int64, ticket types.Ticket) (stop bool) {
		tickets = append(tickets, ticket)
		return false
	})
	return tickets
}
