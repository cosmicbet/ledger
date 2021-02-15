package keeper

import (
	"fmt"
	"strings"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/cosmicbet/ledger/x/wta/types"
)

// RegisterInvariants registers all staking invariants
func RegisterInvariants(ir sdk.InvariantRegistry, k Keeper) {
	ir.RegisterRoute(types.ModuleName, "tickets",
		TicketsInvariants(k))
}

// TicketsInvariants checks that the stored tickets are all valid and there are no duplicated id
func TicketsInvariants(k Keeper) sdk.Invariant {
	return func(ctx sdk.Context) (string, bool) {
		tickets := k.GetTickets(ctx)

		var duplicatedIds []string
		for _, t := range tickets {
			if types.IsTicketIDDuplicated(t.Id, tickets) {
				duplicatedIds = append(duplicatedIds, t.Id)
			}
		}

		broken := len(duplicatedIds) > 0

		// There should be no duplicated ticket ids
		return sdk.FormatInvariant(types.ModuleName, "duplicated ticket ids", fmt.Sprintf(
			"\tDuplicated ticket ids: %v\n"+
				strings.Join(duplicatedIds, ", "))), broken
	}
}
