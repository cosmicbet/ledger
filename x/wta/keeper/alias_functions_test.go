package keeper_test

import (
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/cosmicbet/ledger/x/wta/types"
)

func (suite *KeeperTestSuite) Test_IterateTickets() {
	tickets := []types.Ticket{
		types.NewTicket(
			"1",
			time.Date(2020, 1, 1, 00, 00, 00, 000, time.UTC),
			"owner-1",
		),
		types.NewTicket(
			"2",
			time.Date(2020, 1, 2, 00, 00, 00, 000, time.UTC),
			"owner-2",
		),
		types.NewTicket(
			"3",
			time.Date(2020, 1, 3, 00, 00, 00, 000, time.UTC),
			"owner-3",
		),
	}
	suite.keeper.SaveTickets(suite.ctx, tickets)

	var received []types.Ticket
	suite.keeper.IterateTickets(suite.ctx, func(index int64, ticket types.Ticket) (stop bool) {
		received = append(received, ticket)
		return index == 1
	})
	suite.Require().Equal(tickets[:2], received) // The last ticket should not be added
}

func (suite *KeeperTestSuite) Test_GetTickets() {
	tickets := []types.Ticket{
		types.NewTicket(
			"1",
			time.Date(2020, 1, 1, 00, 00, 00, 000, time.UTC),
			"owner-1",
		),
		types.NewTicket(
			"2",
			time.Date(2020, 1, 2, 00, 00, 00, 000, time.UTC),
			"owner-2",
		),
		types.NewTicket(
			"3",
			time.Date(2020, 1, 3, 00, 00, 00, 000, time.UTC),
			"owner-3",
		),
	}
	suite.keeper.SaveTickets(suite.ctx, tickets)

	stored := suite.keeper.GetTickets(suite.ctx)
	suite.Require().Equal(tickets, stored)
}

func (suite *KeeperTestSuite) Test_IterateHistoricalDrawsData() {
	data := []types.HistoricalDrawData{
		types.NewHistoricalDrawData(
			types.NewDraw(
				1,
				1,
				sdk.NewCoins(sdk.NewInt64Coin("stake", 10)),
				time.Date(2020, 1, 1, 00, 00, 00, 000, time.UTC),
			),
			types.NewTicket(
				"ticket-1",
				time.Date(2019, 12, 31, 23, 59, 59, 999, time.UTC),
				"winner-1",
			),
		),
		types.NewHistoricalDrawData(
			types.NewDraw(
				10,
				100,
				sdk.NewCoins(sdk.NewInt64Coin("stake", 1000)),
				time.Date(2020, 1, 2, 00, 00, 00, 000, time.UTC),
			),
			types.NewTicket(
				"ticket-2",
				time.Date(2020, 1, 1, 23, 59, 59, 999, time.UTC),
				"winner-3",
			),
		),
	}
	for _, data := range data {
		suite.keeper.SaveHistoricalDraw(suite.ctx, data)
	}

	var pastDraws []types.HistoricalDrawData
	suite.keeper.IterateHistoricalDrawsData(suite.ctx, func(index int64, data types.HistoricalDrawData) (stop bool) {
		pastDraws = append(pastDraws, data)
		return index == 0
	})
	suite.Require().Equal(data[:1], pastDraws) // Last data should not be added as we return false before that
}

func (suite *KeeperTestSuite) Test_GetHistoricalDrawsData() {
	data := []types.HistoricalDrawData{
		types.NewHistoricalDrawData(
			types.NewDraw(
				1,
				1,
				sdk.NewCoins(sdk.NewInt64Coin("stake", 10)),
				time.Date(2020, 1, 1, 00, 00, 00, 000, time.UTC),
			),
			types.NewTicket(
				"ticket-1",
				time.Date(2019, 12, 31, 23, 59, 59, 999, time.UTC),
				"winner-1",
			),
		),
		types.NewHistoricalDrawData(
			types.NewDraw(
				10,
				100,
				sdk.NewCoins(sdk.NewInt64Coin("stake", 1000)),
				time.Date(2020, 1, 2, 00, 00, 00, 000, time.UTC),
			),
			types.NewTicket(
				"ticket-2",
				time.Date(2020, 1, 1, 23, 59, 59, 999, time.UTC),
				"winner-3",
			),
		),
	}
	for _, data := range data {
		suite.keeper.SaveHistoricalDraw(suite.ctx, data)
	}

	stored := suite.keeper.GetHistoricalDrawsData(suite.ctx)
	suite.Require().Equal(data, stored)
}
