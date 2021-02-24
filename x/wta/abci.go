package wta

import (
	"time"

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

	participants, tickets := k.GetDrawParticipantsAndTickets(ctx)

	// We need at least two participants to make it fair
	if len(participants) > 2 {

		// Get a random winning ticket
		r := types.NewRandFromCtx(ctx)
		winningTicket := tickets[r.Intn(len(tickets))]

		winner, err := sdk.AccAddressFromBech32(winningTicket.Owner)
		if err != nil {
			panic(err)
		}

		// Send the prize to the winner
		err = k.TransferDrawPrize(ctx, draw.Prize, winner)
		if err != nil {
			panic(err)
		}

		ctx.EventManager().EmitEvent(
			sdk.NewEvent(
				types.EventTypeWinnerDrawn,
				sdk.NewAttribute(types.AttributeKeyWinnerAddress, winningTicket.Owner),
				sdk.NewAttribute(types.AttributeKeyWonAmount, draw.Prize.String()),
			),
		)

		// Save the past draw
		k.SaveHistoricalDraw(ctx, types.NewHistoricalDrawData(draw, winningTicket))

		// Remove all the tickets
		k.WipeCurrentTickets(ctx)
	}

	// Create a new draw
	endTime := ctx.BlockTime().Add(k.GetParams(ctx).DrawDuration)
	k.SaveCurrentDrawEndTime(ctx, endTime)

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeNewDraw,
			sdk.NewAttribute(types.AttributeKeyDrawClosing, endTime.Format(time.RFC3339)),
		),
	)
}
