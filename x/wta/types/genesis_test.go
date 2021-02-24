package types_test

import (
	"github.com/cosmicbet/ledger/x/wta/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestValidateGenesis(t *testing.T) {
	usecases := []struct {
		name      string
		genesis   *types.GenesisState
		shouldErr bool
	}{
		{
			name: "zero draw end time should error",
			genesis: types.NewGenesisState(
				time.Time{},
				nil,
				nil,
				types.DefaultParams(),
			),
			shouldErr: true,
		},
		{
			name: "past draw end time should error",
			genesis: types.NewGenesisState(
				time.Now().Add(-time.Hour*1),
				nil,
				nil,
				types.DefaultParams(),
			),
			shouldErr: true,
		},
		{
			name: "invalid ticket data",
			genesis: types.NewGenesisState(
				time.Now().Add(-time.Hour*1),
				[]types.Ticket{
					types.NewTicket(
						"ticket-id",
						time.Time{},
						"invalid-owner",
					),
				},
				nil,
				types.DefaultParams(),
			),
			shouldErr: true,
		},
		{
			name: "ticket creation time after draw end time",
			genesis: types.NewGenesisState(
				time.Now().Add(-time.Hour*2),
				[]types.Ticket{
					types.NewTicket(
						"ticket-id",
						time.Now().Add(-time.Hour*2+time.Minute),
						"cosmos14zfwkjm35j05ydm3s3qu4he39yjxe9575echwl",
					),
				},
				nil,
				types.DefaultParams(),
			),
			shouldErr: true,
		},
		{
			name: "ticket creation in the future",
			genesis: types.NewGenesisState(
				time.Now().Add(time.Hour*48),
				[]types.Ticket{
					types.NewTicket(
						"ticket-id",
						time.Now().Add(time.Hour*24),
						"cosmos14zfwkjm35j05ydm3s3qu4he39yjxe9575echwl",
					),
				},
				nil,
				types.DefaultParams(),
			),
			shouldErr: true,
		},
		{
			name: "duplicated ticket ids",
			genesis: types.NewGenesisState(
				time.Now().Add(time.Hour),
				[]types.Ticket{
					types.NewTicket(
						"ticket-id",
						time.Now(),
						"cosmos14zfwkjm35j05ydm3s3qu4he39yjxe9575echwl",
					),
					types.NewTicket(
						"ticket-id",
						time.Now().Add(-time.Minute),
						"cosmos14zfwkjm35j05ydm3s3qu4he39yjxe9575echwl",
					),
				},
				nil,
				types.DefaultParams(),
			),
			shouldErr: true,
		},
		{
			name: "invalid historical data",
			genesis: types.NewGenesisState(
				time.Now().Add(time.Hour),
				nil,
				[]types.HistoricalDrawData{
					types.NewHistoricalDrawData(
						types.NewDraw(
							1,
							1,
							sdk.NewCoins(sdk.NewInt64Coin(sdk.DefaultBondDenom, 10)),
							time.Time{},
						),
						types.NewTicket(
							"ticket-id",
							time.Time{},
							"winner",
						),
					),
				},
				types.DefaultParams(),
			),
			shouldErr: true,
		},
		{
			name: "invalid params",
			genesis: types.NewGenesisState(
				time.Now().Add(time.Hour),
				nil,
				nil,
				types.NewParams(
					sdk.NewInt(98),
					sdk.NewInt(2),
					sdk.NewInt(2),
					time.Minute,
					sdk.NewInt64Coin(sdk.DefaultBondDenom, 10),
				),
			),
			shouldErr: true,
		},
		{
			name: "valid genesis",
			genesis: types.NewGenesisState(
				time.Now().Add(time.Hour),
				[]types.Ticket{
					types.NewTicket(
						"ticket-id",
						time.Now().Add(-5*time.Minute),
						"cosmos14zfwkjm35j05ydm3s3qu4he39yjxe9575echwl",
					),
				},
				[]types.HistoricalDrawData{
					types.NewHistoricalDrawData(
						types.NewDraw(
							1,
							1,
							sdk.NewCoins(sdk.NewInt64Coin(sdk.DefaultBondDenom, 1000)),
							time.Now().Add(-7*24*time.Hour),
						),
						types.NewTicket(
							"winning-ticet",
							time.Now().Add(-9*25*time.Hour),
							"cosmos14zfwkjm35j05ydm3s3qu4he39yjxe9575echwl",
						),
					),
				},
				types.NewParams(
					sdk.NewInt(92),
					sdk.NewInt(7),
					sdk.NewInt(1),
					time.Hour*12,
					sdk.NewInt64Coin(sdk.DefaultBondDenom, 100),
				),
			),
			shouldErr: false,
		},
		{
			name:      "default genesis",
			genesis:   types.DefaultGenesisState(),
			shouldErr: false,
		},
	}

	for _, uc := range usecases {
		uc := uc
		t.Run(uc.name, func(t *testing.T) {
			err := types.ValidateGenesis(uc.genesis)

			if uc.shouldErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}
