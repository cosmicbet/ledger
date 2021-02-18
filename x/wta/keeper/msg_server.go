package keeper

import (
	"context"
	"encoding/hex"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	"github.com/cosmicbet/ledger/x/wta/types"
)

type msgServer struct {
	Keeper
}

// NewMsgServerImpl returns an implementation of the wta MsgServer interface
// for the provided Keeper.
func NewMsgServerImpl(keeper Keeper) types.MsgServer {
	return &msgServer{keeper}
}

// generateTickets generates n random tickets for the given user
func (k msgServer) generateTickets(ctx sdk.Context, n int32, user sdk.AccAddress) []types.Ticket {
	tickets := make([]types.Ticket, n)
	for i := range tickets {
		r := types.NewRandFromCtxAndIndex(ctx, i)

		// Get a 16-bytes random id
		var id = make([]byte, 16)
		r.Read(id)

		tickets[i] = types.NewTicket(
			hex.EncodeToString(id),
			ctx.BlockTime(),
			user.String(),
		)
	}

	return tickets
}

// BuyTickets implements MsgServer
func (k msgServer) BuyTickets(ctx context.Context, msg *types.MsgBuyTickets) (*types.MsgBuyTicketsResponse, error) {
	sdkCtx := sdk.UnwrapSDKContext(ctx)

	// Get user address
	user, err := sdk.AccAddressFromBech32(msg.Buyer)
	if err != nil {
		return nil, sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid buyer address")
	}

	// Withdraw the fees
	err = k.WithdrawTicketsCost(sdkCtx, msg.Quantity, user)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrLogic, err.Error())
	}

	tickets := k.generateTickets(sdkCtx, msg.Quantity, user)
	k.SaveTickets(sdkCtx, tickets)

	for _, t := range tickets {
		sdkCtx.EventManager().EmitEvent(
			sdk.NewEvent(
				types.EventTypeBuyTicket,
				sdk.NewAttribute(types.AttributeKeyTicketID, t.Id),
				sdk.NewAttribute(types.AttributeKeyTicketTimestamp, t.Timestamp.Format(time.RFC3339)),
				sdk.NewAttribute(types.AttributeKeyTicketBuyer, t.Owner),
			),
		)
	}

	sdkCtx.EventManager().EmitEvent(
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.ModuleName),
			sdk.NewAttribute(sdk.AttributeKeyAction, types.TypeMsgBuyTickets),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.Buyer),
		),
	)

	return &types.MsgBuyTicketsResponse{}, nil
}
