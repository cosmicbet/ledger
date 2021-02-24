package keeper_test

import (
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"

	"github.com/cosmicbet/ledger/x/wta/keeper"
	"github.com/cosmicbet/ledger/x/wta/types"
)

func (suite *KeeperTestSuite) Test_MsgServer_BuyTickets() {
	addr, err := sdk.AccAddressFromBech32("cosmos14zfwkjm35j05ydm3s3qu4he39yjxe9575echwl")
	suite.Require().NoError(err)

	params := types.NewParams(
		sdk.NewInt(98),
		sdk.NewInt(1),
		sdk.NewInt(1),
		time.Minute*1,
		sdk.NewInt64Coin(sdk.DefaultBondDenom, 10),
	)

	usecases := []struct {
		name            string
		stored          []types.Ticket
		accBalance      sdk.Coins
		msg             *types.MsgBuyTickets
		shouldErr       bool
		expParticipants uint32
		expTicketsSold  uint32
	}{
		{
			name:      "invalid address",
			msg:       types.NewMsgBuyTickets(10, "address"),
			shouldErr: true,
		},
		{
			name:      "insufficient balance",
			msg:       types.NewMsgBuyTickets(1, addr.String()),
			shouldErr: true,
		},
		{
			name:            "buying without any stored ticket",
			stored:          nil,
			accBalance:      sdk.NewCoins(sdk.NewInt64Coin(sdk.DefaultBondDenom, 10000)),
			msg:             types.NewMsgBuyTickets(10, addr.String()),
			shouldErr:       false,
			expParticipants: 1,
			expTicketsSold:  10,
		},
		{
			name:       "buying more tickets",
			accBalance: sdk.NewCoins(sdk.NewInt64Coin(sdk.DefaultBondDenom, 10000)),
			stored: []types.Ticket{
				types.NewTicket(
					"ticket-1",
					time.Date(2020, 1, 1, 00, 00, 00, 000, time.UTC),
					addr.String(),
				),
				types.NewTicket(
					"ticket-2",
					time.Date(2020, 1, 1, 00, 00, 00, 000, time.UTC),
					addr.String(),
				),
			},
			msg:             types.NewMsgBuyTickets(5, addr.String()),
			shouldErr:       false,
			expParticipants: 1,
			expTicketsSold:  7,
		},
		{
			name:       "buying tickets as second participant",
			accBalance: sdk.NewCoins(sdk.NewInt64Coin(sdk.DefaultBondDenom, 10000)),
			stored: []types.Ticket{
				types.NewTicket(
					"ticket-1",
					time.Date(2020, 1, 1, 00, 00, 00, 000, time.UTC),
					"user-2",
				),
				types.NewTicket(
					"ticket-2",
					time.Date(2020, 1, 1, 00, 00, 00, 000, time.UTC),
					"user-2",
				),
			},
			msg:             types.NewMsgBuyTickets(5, addr.String()),
			shouldErr:       false,
			expParticipants: 2,
			expTicketsSold:  7,
		},
	}

	for _, uc := range usecases {
		suite.SetupTest()
		suite.Run(uc.name, func() {
			suite.keeper.SaveTickets(suite.ctx, uc.stored)
			suite.keeper.SetParams(suite.ctx, params)
			suite.bk.SetSupply(suite.ctx, banktypes.NewSupply(uc.accBalance))
			suite.Require().NoError(suite.bk.SetBalances(suite.ctx, addr, uc.accBalance))

			server := keeper.NewMsgServerImpl(suite.keeper)
			_, err := server.BuyTickets(sdk.WrapSDKContext(suite.ctx), uc.msg)

			if uc.shouldErr {
				suite.Require().Error(err)
			} else {
				suite.Require().NoError(err)

				participants, ticketsSold := suite.keeper.GetDrawParticipantsAndTicketsSold(suite.ctx)
				suite.Require().Equal(uc.expParticipants, participants)
				suite.Require().Equal(uc.expTicketsSold, ticketsSold)
			}
		})
	}
}
