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
func NewMsgBuyTickets(quantity uint32, user string) *MsgBuyTickets {
	return &MsgBuyTickets{
		Quantity: quantity,
		Buyer:    user,
	}
}

// Route implements sdk.Msg
func (m *MsgBuyTickets) Route() string {
	return RouterKey
}

// Type implements sdk.Msg
func (m *MsgBuyTickets) Type() string {
	return TypeMsgBuyTickets
}

// ValidateBasic implements sdk.Msg
func (m *MsgBuyTickets) ValidateBasic() error {
	if m.Quantity <= 0 {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidRequest, "invalid tickets quantity: %d", m.Quantity)
	}

	if _, err := sdk.AccAddressFromBech32(m.Buyer); err != nil {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, "invalid buyer address")
	}

	return nil
}

// GetSignBytes implements sdk.Msg
func (m *MsgBuyTickets) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(m)
	return sdk.MustSortJSON(bz)
}

// GetSigners implements sdk.Msg
func (m *MsgBuyTickets) GetSigners() []sdk.AccAddress {
	buyerAddr, err := sdk.AccAddressFromBech32(m.Buyer)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{buyerAddr}
}
