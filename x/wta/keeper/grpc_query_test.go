package keeper_test

import (
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/query"

	"github.com/cosmicbet/ledger/x/wta/keeper"
	"github.com/cosmicbet/ledger/x/wta/types"
)

func (suite *KeeperTestSuite) Test_Querier_Tickets() {
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

	usecases := []struct {
		name       string
		req        *types.QueryTicketsRequest
		shouldErr  bool
		expTickets []types.Ticket
	}{
		{
			name:      "empty request",
			req:       nil,
			shouldErr: true,
		},
		{
			name: "small pagination",
			req: types.NewTicketsRequest(&query.PageRequest{
				Offset: 1,
				Limit:  1,
			}),
			shouldErr:  false,
			expTickets: []types.Ticket{tickets[1]},
		},
		{
			name: "large pagination",
			req: types.NewTicketsRequest(&query.PageRequest{
				Offset: 0,
				Limit:  1000,
			}),
			shouldErr:  false,
			expTickets: tickets,
		},
	}

	for _, uc := range usecases {
		suite.SetupTest()
		suite.Run(uc.name, func() {
			suite.keeper.SaveTickets(suite.ctx, tickets)

			querier := keeper.NewQuerierImpl(suite.keeper)
			res, err := querier.Tickets(sdk.WrapSDKContext(suite.ctx), uc.req)

			if uc.shouldErr {
				suite.Require().Error(err)
			} else {
				suite.Require().NoError(err)
				suite.Require().Equal(uc.expTickets, res.Tickets)
			}
		})
	}
}

func (suite *KeeperTestSuite) Test_Querier_NextDraw() {
	usecases := []struct {
		name        string
		drawEndTime time.Time
		prize       sdk.Coins
		tickets     []types.Ticket
		req         *types.QueryNextDrawRequest
		shouldErr   bool
		expDraw     types.Draw
	}{
		{
			name:      "invalid request",
			req:       nil,
			shouldErr: true,
		},
		{
			name:        "valid request",
			drawEndTime: time.Date(2020, 1, 2, 00, 00, 00, 000, time.UTC),
			prize:       sdk.NewCoins(sdk.NewInt64Coin("stake", 1000)),
			tickets: []types.Ticket{
				types.NewTicket(
					"ticket-1",
					time.Date(2020, 1, 1, 00, 00, 00, 000, time.UTC),
					"owner-1",
				),
				types.NewTicket(
					"ticket-2",
					time.Date(2020, 1, 1, 00, 00, 00, 000, time.UTC),
					"owner-1",
				),
				types.NewTicket(
					"ticket-3",
					time.Date(2020, 1, 1, 00, 00, 00, 000, time.UTC),
					"owner-2",
				),
			},
			req:       &types.QueryNextDrawRequest{},
			shouldErr: false,
			expDraw: types.NewDraw(
				2,
				3,
				sdk.NewCoins(sdk.NewInt64Coin("stake", 1000)),
				time.Date(2020, 1, 2, 00, 00, 00, 000, time.UTC),
			),
		},
	}

	for _, uc := range usecases {
		suite.SetupTest()
		suite.Run(uc.name, func() {
			suite.SaveDrawData(suite.ctx, uc.drawEndTime, uc.prize)
			suite.keeper.SaveTickets(suite.ctx, uc.tickets)

			querier := keeper.NewQuerierImpl(suite.keeper)
			res, err := querier.NextDraw(sdk.WrapSDKContext(suite.ctx), uc.req)

			if uc.shouldErr {
				suite.Require().Error(err)
			} else {
				suite.Require().NoError(err)
				suite.Require().Equal(uc.expDraw, res.Draw)
			}
		})
	}
}

func (suite *KeeperTestSuite) Test_Querier_PastDraws() {
	draws := []types.HistoricalDrawData{
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

	usecases := []struct {
		name      string
		req       *types.QueryPastDrawsRequest
		shouldErr bool
		expDraws  []types.HistoricalDrawData
	}{
		{
			name:      "invalid request",
			req:       nil,
			shouldErr: true,
		},
		{
			name: "small pagination",
			req: types.NewPastDrawsRequest(&query.PageRequest{
				Offset: 1,
				Limit:  1,
			}),
			expDraws: []types.HistoricalDrawData{draws[1]},
		},
		{
			name: "large pagination",
			req: types.NewPastDrawsRequest(&query.PageRequest{
				Offset: 0,
				Limit:  100,
			}),
			expDraws: draws,
		},
	}

	for _, uc := range usecases {
		suite.SetupTest()
		suite.Run(uc.name, func() {
			for _, d := range draws {
				suite.keeper.SaveHistoricalDraw(suite.ctx, d)
			}

			querier := keeper.NewQuerierImpl(suite.keeper)
			res, err := querier.PastDraws(sdk.WrapSDKContext(suite.ctx), uc.req)

			if uc.shouldErr {
				suite.Require().Error(err)
			} else {
				suite.Require().NoError(err)
				suite.Require().Equal(uc.expDraws, res.Draws)
			}
		})
	}
}

func (suite *KeeperTestSuite) Test_Querier_Params() {
	distributionParams := types.NewDistributionParams(
		sdk.NewDecWithPrec(95, 2),
		sdk.NewDecWithPrec(3, 2),
		sdk.NewDecWithPrec(2, 2),
	)
	drawParams := types.NewDrawParams(time.Minute * 3)
	ticketParams := types.NewTicketParams(sdk.NewInt64Coin("stake", 10))

	usecases := []struct {
		name      string
		req       *types.QueryParamsRequest
		shouldErr bool
	}{
		{
			name:      "invalid request",
			req:       nil,
			shouldErr: true,
		},
		{
			name:      "valid request",
			req:       &types.QueryParamsRequest{},
			shouldErr: false,
		},
	}

	for _, uc := range usecases {
		suite.SetupTest()
		suite.Run(uc.name, func() {
			suite.keeper.SetDistributionParams(suite.ctx, distributionParams)
			suite.keeper.SetDrawParams(suite.ctx, drawParams)
			suite.keeper.SetTicketParams(suite.ctx, ticketParams)

			querier := keeper.NewQuerierImpl(suite.keeper)
			res, err := querier.Params(sdk.WrapSDKContext(suite.ctx), uc.req)

			if uc.shouldErr {
				suite.Require().Error(err)
			} else {
				suite.Require().NoError(err)
				suite.Require().Equal(distributionParams, res.DistributionParams)
				suite.Require().Equal(drawParams, res.DrawParams)
				suite.Require().Equal(ticketParams, res.TicketParams)
			}
		})
	}
}
