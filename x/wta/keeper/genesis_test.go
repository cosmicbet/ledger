package keeper_test

import (
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/cosmicbet/ledger/x/wta/types"
)

func (suite *KeeperTestSuite) Test_ExportGenesis() {
	usecases := []struct {
		name            string
		drawEndDate     time.Time
		tickets         []types.Ticket
		historicalDraws []types.HistoricalDrawData
		params          types.Params
	}{
		{
			name:            "empty tickets and historical data",
			drawEndDate:     time.Date(2020, 1, 1, 00, 00, 00, 000, time.UTC),
			tickets:         nil,
			historicalDraws: nil,
			params: types.NewParams(
				sdk.NewInt(98),
				sdk.NewInt(1),
				sdk.NewInt(1),
				time.Minute*5,
				sdk.NewInt64Coin("stake", 10),
			),
		},
		{
			name:        "non empty tickets and historical data",
			drawEndDate: time.Date(2020, 1, 1, 00, 00, 00, 000, time.UTC),
			tickets: []types.Ticket{
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
			},
			historicalDraws: []types.HistoricalDrawData{
				types.NewHistoricalDrawData(
					types.NewDraw(
						1,
						1,
						sdk.NewCoins(sdk.NewInt64Coin("stake", 10)),
						time.Date(2020, 1, 1, 00, 00, 00, 000, time.UTC),
					),
					types.NewTicket(
						"old-ticket",
						time.Date(2019, 12, 31, 23, 59, 59, 000, time.UTC),
						"old-winner",
					),
				),
			},
			params: types.NewParams(
				sdk.NewInt(95),
				sdk.NewInt(3),
				sdk.NewInt(2),
				time.Minute*3,
				sdk.NewInt64Coin("stake", 10),
			),
		},
	}

	for _, uc := range usecases {
		suite.SetupTest()
		suite.Run(uc.name, func() {
			suite.keeper.SaveCurrentDrawEndTime(suite.ctx, uc.drawEndDate)
			suite.keeper.SaveTickets(suite.ctx, uc.tickets)
			for _, h := range uc.historicalDraws {
				suite.keeper.SaveHistoricalDraw(suite.ctx, h)
			}
			suite.keeper.SetParams(suite.ctx, uc.params)

			exported := suite.keeper.ExportGenesis(suite.ctx)
			suite.Require().Equal(uc.drawEndDate, exported.DrawEndTime)
			suite.Require().Equal(uc.tickets, exported.Tickets)
			suite.Require().Equal(uc.historicalDraws, exported.PastDraws)
			suite.Require().Equal(uc.params, exported.Params)
		})
	}
}

func (suite *KeeperTestSuite) Test_ImportGenesis() {
	usecases := []struct {
		name    string
		genesis *types.GenesisState
	}{
		{
			name: "empty tickets and historical data",
			genesis: types.NewGenesisState(
				time.Date(2020, 1, 1, 00, 00, 00, 000, time.UTC),
				nil,
				nil,
				types.NewParams(
					sdk.NewInt(98),
					sdk.NewInt(1),
					sdk.NewInt(1),
					time.Minute*5,
					sdk.NewInt64Coin("stake", 10),
				),
			),
		},
		{
			name: "non empty tickets and historical data",
			genesis: types.NewGenesisState(
				time.Date(2020, 1, 1, 00, 00, 00, 000, time.UTC),
				[]types.Ticket{
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
				},
				[]types.HistoricalDrawData{
					types.NewHistoricalDrawData(
						types.NewDraw(
							1,
							1,
							sdk.NewCoins(sdk.NewInt64Coin("stake", 10)),
							time.Date(2020, 1, 1, 00, 00, 00, 000, time.UTC),
						),
						types.NewTicket(
							"old-ticket",
							time.Date(2019, 12, 31, 23, 59, 59, 000, time.UTC),
							"old-winner",
						),
					),
				},
				types.NewParams(
					sdk.NewInt(95),
					sdk.NewInt(3),
					sdk.NewInt(2),
					time.Minute*3,
					sdk.NewInt64Coin("stake", 10),
				),
			),
		},
	}

	for _, uc := range usecases {
		suite.SetupTest()
		suite.Run(uc.name, func() {
			suite.keeper.InitGenesis(suite.ctx, *uc.genesis)

			draw := suite.keeper.GetCurrentDraw(suite.ctx)
			suite.Require().Equal(uc.genesis.DrawEndTime, draw.EndTime)

			tickets := suite.keeper.GetTickets(suite.ctx)
			suite.Require().Equal(uc.genesis.Tickets, tickets)

			pastDraws := suite.keeper.GetHistoricalDrawsData(suite.ctx)
			suite.Require().Equal(uc.genesis.PastDraws, pastDraws)

			params := suite.keeper.GetParams(suite.ctx)
			suite.Require().Equal(uc.genesis.Params, params)
		})
	}
}
