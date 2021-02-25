package keeper_test

import (
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/cosmicbet/ledger/x/wta/types"
)

func (suite *KeeperTestSuite) Test_ExportGenesis() {
	usecases := []struct {
		name               string
		drawEndDate        time.Time
		tickets            []types.Ticket
		historicalDraws    []types.HistoricalDrawData
		distributionParams types.DistributionParams
		drawParams         types.DrawParams
		ticketParams       types.TicketParams
	}{
		{
			name:            "empty tickets and historical data",
			drawEndDate:     time.Date(2020, 1, 1, 00, 00, 00, 000, time.UTC),
			tickets:         nil,
			historicalDraws: nil,
			distributionParams: types.NewDistributionParams(
				sdk.NewDecWithPrec(98, 2),
				sdk.NewDecWithPrec(1, 2),
				sdk.NewDecWithPrec(1, 2),
			),
			drawParams:   types.NewDrawParams(time.Minute * 5),
			ticketParams: types.NewTicketParams(sdk.NewInt64Coin("stake", 10)),
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
			distributionParams: types.NewDistributionParams(
				sdk.NewDecWithPrec(95, 2),
				sdk.NewDecWithPrec(3, 2),
				sdk.NewDecWithPrec(2, 2),
			),
			drawParams:   types.NewDrawParams(time.Minute * 3),
			ticketParams: types.NewTicketParams(sdk.NewInt64Coin("stake", 10)),
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
			suite.keeper.SetDistributionParams(suite.ctx, uc.distributionParams)
			suite.keeper.SetDrawParams(suite.ctx, uc.drawParams)
			suite.keeper.SetTicketParams(suite.ctx, uc.ticketParams)

			exported := suite.keeper.ExportGenesis(suite.ctx)
			suite.Require().Equal(uc.drawEndDate, exported.DrawEndTime)
			suite.Require().Equal(uc.tickets, exported.Tickets)
			suite.Require().Equal(uc.historicalDraws, exported.PastDraws)
			suite.Require().Equal(uc.distributionParams, exported.DistributionParams)
			suite.Require().Equal(uc.drawParams, exported.DrawParams)
			suite.Require().Equal(uc.ticketParams, exported.TicketParams)
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
				types.NewDistributionParams(
					sdk.NewDecWithPrec(98, 2),
					sdk.NewDecWithPrec(1, 2),
					sdk.NewDecWithPrec(1, 2),
				),
				types.NewDrawParams(time.Minute*5),
				types.NewTicketParams(sdk.NewInt64Coin("stake", 10)),
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
				types.NewDistributionParams(
					sdk.NewDecWithPrec(95, 2),
					sdk.NewDecWithPrec(3, 2),
					sdk.NewDecWithPrec(2, 2),
				),
				types.NewDrawParams(time.Minute*3),
				types.NewTicketParams(sdk.NewInt64Coin("stake", 10)),
			),
		},
	}

	for _, uc := range usecases {
		suite.SetupTest()
		suite.Run(uc.name, func() {
			suite.keeper.InitGenesis(suite.ctx, *uc.genesis)

			draw := suite.keeper.GetCurrentDraw(suite.ctx)
			suite.Require().Equal(uc.genesis.DrawEndTime, draw.EndTime)

			suite.Require().Equal(uc.genesis.Tickets, suite.keeper.GetTickets(suite.ctx))
			suite.Require().Equal(uc.genesis.PastDraws, suite.keeper.GetHistoricalDrawsData(suite.ctx))

			suite.Require().Equal(uc.genesis.DistributionParams, suite.keeper.GetDistributionParams(suite.ctx))
			suite.Require().Equal(uc.genesis.DrawParams, suite.keeper.GetDrawParams(suite.ctx))
			suite.Require().Equal(uc.genesis.TicketParams, suite.keeper.GetTicketParams(suite.ctx))
		})
	}
}
