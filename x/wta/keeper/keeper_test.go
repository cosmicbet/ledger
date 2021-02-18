package keeper_test

import (
	wtatypes "github.com/cosmicbet/ledger/x/wta/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	"time"
)

func (suite *KeeperTestSuite) Test_WithdrawTicketsCost() {
	usecases := []struct {
		name            string
		ticketPrice     sdk.Coin
		prizePercentage sdk.Int
		poolPercentage  sdk.Int
		burnPercentage  sdk.Int
		accountAddress  string
		accountBalance  sdk.Coins
		quantity        int32

		shouldErr      bool
		expAccBalance  sdk.Coins
		expPrizePool   sdk.Coins
		expPoolBalance sdk.Coins
		expSupply      sdk.Coins
	}{
		{
			name:            "insufficient balance (0)",
			ticketPrice:     sdk.NewInt64Coin("stake", 10),
			prizePercentage: sdk.NewInt(98),
			poolPercentage:  sdk.NewInt(1),
			burnPercentage:  sdk.NewInt(1),
			accountAddress:  "cosmos14zfwkjm35j05ydm3s3qu4he39yjxe9575echwl",
			accountBalance:  sdk.NewCoins(sdk.NewInt64Coin("fiches", 100)),
			quantity:        1,
			shouldErr:       true,
		},
		{
			name:            "insufficient balance (> 0)",
			ticketPrice:     sdk.NewInt64Coin("stake", 10),
			prizePercentage: sdk.NewInt(98),
			poolPercentage:  sdk.NewInt(1),
			burnPercentage:  sdk.NewInt(1),
			accountAddress:  "cosmos14zfwkjm35j05ydm3s3qu4he39yjxe9575echwl",
			accountBalance:  sdk.NewCoins(sdk.NewInt64Coin("stake", 9)),
			quantity:        1,
			shouldErr:       true,
		},
		{
			name:            "single ticket",
			ticketPrice:     sdk.NewInt64Coin("stake", 100),
			prizePercentage: sdk.NewInt(98),
			poolPercentage:  sdk.NewInt(1),
			burnPercentage:  sdk.NewInt(1),
			accountAddress:  "cosmos14zfwkjm35j05ydm3s3qu4he39yjxe9575echwl",
			accountBalance:  sdk.NewCoins(sdk.NewInt64Coin("stake", 100)),
			quantity:        1,
			shouldErr:       false,
			expAccBalance:   sdk.NewCoins(),
			expPrizePool:    sdk.NewCoins(sdk.NewInt64Coin("stake", 98)),
			expPoolBalance:  sdk.NewCoins(sdk.NewInt64Coin("stake", 1)),
			expSupply:       sdk.NewCoins(sdk.NewInt64Coin("stake", 99)),
		},
		{
			name:            "multiple tickets",
			ticketPrice:     sdk.NewInt64Coin("stake", 100),
			prizePercentage: sdk.NewInt(97),
			poolPercentage:  sdk.NewInt(2),
			burnPercentage:  sdk.NewInt(3),
			accountAddress:  "cosmos14zfwkjm35j05ydm3s3qu4he39yjxe9575echwl",
			accountBalance:  sdk.NewCoins(sdk.NewInt64Coin("stake", 1000)),
			quantity:        5,
			expPrizePool:    sdk.NewCoins(sdk.NewInt64Coin("stake", 485)),
			expAccBalance:   sdk.NewCoins(sdk.NewInt64Coin("stake", 500)),
			expPoolBalance:  sdk.NewCoins(sdk.NewInt64Coin("stake", 10)),
			expSupply:       sdk.NewCoins(sdk.NewInt64Coin("stake", 985)),
		},
	}

	for _, uc := range usecases {
		suite.SetupTest()
		suite.Run(uc.name, func() {
			// Set the params
			params := wtatypes.NewParams(uc.prizePercentage, uc.poolPercentage, uc.burnPercentage, 1*time.Minute, uc.ticketPrice)
			suite.Require().NoError(params.Validate())
			suite.keeper.SetParams(suite.ctx, params)

			// Get the account
			addr, err := sdk.AccAddressFromBech32(uc.accountAddress)
			suite.Require().NoError(err)

			// Set the coins supply and account balance
			suite.bk.SetSupply(suite.ctx, banktypes.NewSupply(uc.accountBalance))
			err = suite.bk.SetBalances(suite.ctx, addr, uc.accountBalance)
			suite.Require().NoError(err)

			// Buy the ticket
			err = suite.keeper.WithdrawTicketsCost(suite.ctx, uc.quantity, addr)

			if uc.shouldErr {
				suite.Require().Error(err)
			} else {
				suite.Require().NoError(err)

				accBalance := suite.bk.GetAllBalances(suite.ctx, addr)
				suite.Require().True(accBalance.IsEqual(uc.expAccBalance))

				prizeAcc := suite.ak.GetModuleAccount(suite.ctx, wtatypes.PrizeCollectorName)
				prizePool := suite.bk.GetAllBalances(suite.ctx, prizeAcc.GetAddress())
				suite.Require().True(prizePool.IsEqual(uc.expPrizePool))

				poolBalance := suite.dk.GetFeePoolCommunityCoins(suite.ctx)
				suite.Require().True(poolBalance.IsEqual(sdk.NewDecCoinsFromCoins(uc.expPoolBalance...)))

				supply := suite.bk.GetSupply(suite.ctx)
				suite.Require().True(supply.GetTotal().IsEqual(uc.expSupply))
			}
		})
	}
}

func (suite *KeeperTestSuite) Test_SaveTickets() {
	usecases := []struct {
		name    string
		tickets []wtatypes.Ticket
	}{
		{
			name:    "empty tickets slice",
			tickets: nil,
		},
		{
			name: "non empty tickets slice",
			tickets: []wtatypes.Ticket{
				wtatypes.NewTicket("1", time.Now(), "owner-1"),
				wtatypes.NewTicket("2", time.Now(), "owner-2"),
				wtatypes.NewTicket("3", time.Now(), "owner-3"),
			},
		},
	}

	for _, uc := range usecases {
		suite.SetupTest()
		suite.Run(uc.name, func() {
			suite.keeper.SaveTickets(suite.ctx, uc.tickets)

			var tickets []wtatypes.Ticket
			suite.keeper.IterateTickets(suite.ctx, func(_ int64, ticket wtatypes.Ticket) (stop bool) {
				tickets = append(tickets, ticket)
				return false
			})

			suite.Require().Len(tickets, len(uc.tickets))
			for _, ticket := range tickets {
				suite.Require().Contains(tickets, ticket)
			}
		})
	}
}

func (suite *KeeperTestSuite) Test_WipeCurrentTickets() {
	usecases := []struct {
		name          string
		storedTickets []wtatypes.Ticket
	}{
		{
			name:          "empty storage",
			storedTickets: nil,
		},
		{
			name: "non empty storage",
			storedTickets: []wtatypes.Ticket{
				wtatypes.NewTicket("1", time.Now(), "owner-1"),
				wtatypes.NewTicket("2", time.Now(), "owner-2"),
				wtatypes.NewTicket("3", time.Now(), "owner-3"),
			},
		},
	}

	for _, uc := range usecases {
		suite.SetupTest()
		suite.Run(uc.name, func() {
			suite.keeper.SaveTickets(suite.ctx, uc.storedTickets)
			suite.Require().Len(suite.keeper.GetTickets(suite.ctx), len(uc.storedTickets))

			suite.keeper.WipeCurrentTickets(suite.ctx)

			suite.Require().Empty(suite.keeper.GetTickets(suite.ctx))
		})
	}
}

func (suite *KeeperTestSuite) Test_IncrementDrawPrize() {
	usecases := []struct {
		name         string
		existingDraw wtatypes.Draw
		increment    sdk.Coin
		expectedDraw wtatypes.Draw
	}{
		{
			name: "empty prize draw",
			existingDraw: wtatypes.NewDraw(
				sdk.NewCoins(),
				time.Date(2020, 1, 1, 00, 00, 00, 000, time.UTC),
			),
			increment: sdk.NewInt64Coin("stake", 100),
			expectedDraw: wtatypes.NewDraw(
				sdk.NewCoins(sdk.NewInt64Coin("stake", 100)),
				time.Date(2020, 1, 1, 00, 00, 00, 000, time.UTC),
			),
		},
		{
			name: "non empty prize draw",
			existingDraw: wtatypes.NewDraw(
				sdk.NewCoins(sdk.NewInt64Coin("stake", 100)),
				time.Date(2020, 1, 1, 00, 00, 00, 000, time.UTC),
			),
			increment: sdk.NewInt64Coin("stake", 100),
			expectedDraw: wtatypes.NewDraw(
				sdk.NewCoins(sdk.NewInt64Coin("stake", 200)),
				time.Date(2020, 1, 1, 00, 00, 00, 000, time.UTC),
			),
		},
	}

	for _, uc := range usecases {
		suite.SetupTest()
		suite.Run(uc.name, func() {
			suite.keeper.SaveCurrentDraw(suite.ctx, uc.existingDraw)

			suite.keeper.IncrementDrawPrize(suite.ctx, uc.increment)

			draw := suite.keeper.GetCurrentDraw(suite.ctx)
			suite.Require().True(draw.Equal(uc.expectedDraw))
		})
	}
}

func (suite *KeeperTestSuite) Test_TransferDrawPrize() {
	usecases := []struct {
		name           string
		address        string
		initialBalance sdk.Coins
		prize          sdk.Coins
		expBalance     sdk.Coins
	}{
		{
			name:           "initial empty balance",
			address:        "cosmos14zfwkjm35j05ydm3s3qu4he39yjxe9575echwl",
			initialBalance: sdk.NewCoins(),
			prize:          sdk.NewCoins(sdk.NewInt64Coin("stake", 100)),
			expBalance:     sdk.NewCoins(sdk.NewInt64Coin("stake", 100)),
		},
		{
			name:           "initial non empty balance",
			address:        "cosmos14zfwkjm35j05ydm3s3qu4he39yjxe9575echwl",
			initialBalance: sdk.NewCoins(sdk.NewInt64Coin("stake", 100)),
			prize:          sdk.NewCoins(sdk.NewInt64Coin("stake", 100)),
			expBalance:     sdk.NewCoins(sdk.NewInt64Coin("stake", 200)),
		},
	}

	for _, uc := range usecases {
		suite.SetupTest()
		suite.Run(uc.name, func() {
			// Setup accounts
			addr, err := sdk.AccAddressFromBech32(uc.address)
			suite.Require().NoError(err)

			err = suite.bk.SetBalances(suite.ctx, addr, uc.initialBalance)
			suite.Require().NoError(err)

			moduleAddr := suite.ak.GetModuleAccount(suite.ctx, wtatypes.PrizeCollectorName)
			err = suite.bk.SetBalances(suite.ctx, moduleAddr.GetAddress(), uc.prize)

			// Transfer prize
			err = suite.keeper.TransferDrawPrize(suite.ctx, uc.prize, addr)
			suite.Require().NoError(err)

			// Check balances
			accBalance := suite.bk.GetAllBalances(suite.ctx, addr)
			suite.Require().True(accBalance.IsEqual(uc.expBalance))

			moduleBalance := suite.bk.GetAllBalances(suite.ctx, moduleAddr.GetAddress())
			suite.Require().True(moduleBalance.IsEqual(sdk.NewCoins()))
		})
	}
}

func (suite *KeeperTestSuite) Test_SaveCurrentDraw() {
	usecases := []struct {
		name     string
		existing *wtatypes.Draw
		toSave   wtatypes.Draw
	}{
		{
			name:     "saving draw when non existing",
			existing: nil,
			toSave: wtatypes.NewDraw(
				sdk.NewCoins(sdk.NewInt64Coin("stake", 10)),
				time.Date(2020, 1, 1, 00, 00, 00, 000, time.UTC),
			),
		},
		{
			name: "saving draw with existing one",
			existing: &wtatypes.Draw{
				Prize:   sdk.NewCoins(sdk.NewInt64Coin("stake", 10)),
				EndTime: time.Date(2020, 1, 1, 00, 00, 00, 000, time.UTC),
			},
			toSave: wtatypes.NewDraw(
				sdk.NewCoins(sdk.NewInt64Coin("uatom", 10000)),
				time.Date(2020, 1, 2, 00, 00, 00, 000, time.UTC),
			),
		},
	}

	for _, uc := range usecases {
		suite.SetupTest()
		suite.Run(uc.name, func() {
			if uc.existing != nil {
				suite.keeper.SaveCurrentDraw(suite.ctx, *uc.existing)
				suite.Require().True(suite.keeper.GetCurrentDraw(suite.ctx).Equal(*uc.existing))
			}

			suite.keeper.SaveCurrentDraw(suite.ctx, uc.toSave)
			suite.Require().True(suite.keeper.GetCurrentDraw(suite.ctx).Equal(uc.toSave))
		})
	}
}
