package wta

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/cosmicbet/ledger/x/wta/keeper"
	"github.com/cosmicbet/ledger/x/wta/types"
)

// BeginBlocker will check if there is a current draw for which a winner should be drawn.
// If there is, randomly gets the winner and rewards it with the prize of the draw itself.
// Then, creates a new draw.
func BeginBlocker(ctx sdk.Context, k keeper.Keeper) {
	draw := k.GetCurrentDraw(ctx)

	// Check to make sure it's fine to draw the winner
	if ctx.BlockTime().Before(draw.EndTime) {
		return
	}

	tickets := k.GetTickets(ctx)
	if len(tickets) > 0 {

		// Get a random winning ticket
		r := types.NewRandFromCtx(ctx)
		ticket := tickets[r.Intn(len(tickets))]

		winner, err := sdk.AccAddressFromBech32(ticket.Owner)
		if err != nil {
			panic(err)
		}

		// Send the prize to the winner
		err = k.TransferDrawPrize(ctx, draw.Prize, winner)
		if err != nil {
			panic(err)
		}
	}

	// Remove all the tickets
	k.WipeCurrentTickets(ctx)

	// Create a new draw
	endTime := ctx.BlockTime().Add(k.GetParams(ctx).DrawDuration)
	err := k.SaveCurrentDraw(ctx, types.NewDraw(sdk.NewCoins(), endTime))
	if err != nil {
		panic(err)
	}
}
