package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

const (
	TypeMsgBuyTickets = "buy_tickets"
)

var _ sdk.Msg = &MsgBuyTickets{}

// NewMsgBuyTickets allows to build a new MsgBuyTickets instance
func NewMsgBuyTickets(quantity int32, user string) *MsgBuyTickets {
	return &MsgBuyTickets{
		Quantity: quantity,
		Buyer:    user,
	}
}

// Route implements sdk.Msg
func (msg MsgBuyTickets) Route() string {
	return RouterKey
}

// Type implements sdk.Msg
func (msg MsgBuyTickets) Type() string {
	return TypeMsgBuyTickets
}

// ValidateBasic implements sdk.Msg
func (msg MsgBuyTickets) ValidateBasic() error {
	if msg.Quantity <= 0 {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidRequest, "invalid tickets quantity: %d", msg.Quantity)
	}

	if _, err := sdk.AccAddressFromBech32(msg.Buyer); err != nil {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, "invalid buyer address")
	}

	return nil
}

// GetSignBytes implements sdk.Msg
func (msg MsgBuyTickets) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(&msg)
	return sdk.MustSortJSON(bz)
}

// GetSigners implements sdk.Msg
func (msg MsgBuyTickets) GetSigners() []sdk.AccAddress {
	buyerAddr, err := sdk.AccAddressFromBech32(msg.Buyer)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{buyerAddr}
}
