package types_test

import (
	"github.com/cosmicbet/ledger/x/wta/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestTicket_Validate(t *testing.T) {
	usecases := []struct {
		name      string
		ticket    types.Ticket
		shouldErr bool
	}{
		{
			name:      "invalid id",
			ticket:    types.NewTicket("", time.Now(), "cosmos14zfwkjm35j05ydm3s3qu4he39yjxe9575echwl"),
			shouldErr: true,
		},
		{
			name:      "invalid time",
			ticket:    types.NewTicket("ticket-id", time.Time{}, "cosmos14zfwkjm35j05ydm3s3qu4he39yjxe9575echwl"),
			shouldErr: true,
		},
		{
			name:      "invalid owner",
			ticket:    types.NewTicket("ticket-id", time.Now(), ""),
			shouldErr: true,
		},
		{
			name:      "valid ticket",
			ticket:    types.NewTicket("ticket-id", time.Now(), "cosmos14zfwkjm35j05ydm3s3qu4he39yjxe9575echwl"),
			shouldErr: false,
		},
	}

	for _, uc := range usecases {
		uc := uc
		t.Run(uc.name, func(t *testing.T) {
			err := uc.ticket.Validate()

			if uc.shouldErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestIsTicketIDDuplicated(t *testing.T) {
	usecases := []struct {
		name          string
		id            string
		tickets       []types.Ticket
		expDuplicated bool
	}{
		{
			name: "duplicated id",
			id:   "ticket-id",
			tickets: []types.Ticket{
				types.NewTicket("ticket-id", time.Now(), "owner-1"),
				types.NewTicket("ticket-id", time.Now(), "owner-1"),
			},
			expDuplicated: true,
		},
		{
			name: "non duplicated id",
			id:   "ticket-id-1",
			tickets: []types.Ticket{
				types.NewTicket("ticket-id-1", time.Now(), "owner-1"),
				types.NewTicket("ticket-id-2", time.Now(), "owner-1"),
			},
			expDuplicated: false,
		},
	}

	for _, uc := range usecases {
		uc := uc
		t.Run(uc.name, func(t *testing.T) {
			require.Equal(t, uc.expDuplicated, types.IsTicketIDDuplicated(uc.id, uc.tickets))
		})
	}
}

func TestDraw_Validate(t *testing.T) {
	usecases := []struct {
		name      string
		draw      types.Draw
		shouldErr bool
	}{
		{
			name: "invalid number of tickets and participants",
			draw: types.NewDraw(
				10,
				5,
				sdk.NewCoins(sdk.NewInt64Coin(sdk.DefaultBondDenom, 10)),
				time.Now(),
			),
			shouldErr: true,
		},
		{
			name: "invalid prize",
			draw: types.NewDraw(
				1,
				1,
				sdk.Coins{sdk.Coin{Denom: "./+", Amount: sdk.NewInt(10)}},
				time.Now(),
			),
			shouldErr: true,
		},
		{
			name: "invalid time",
			draw: types.NewDraw(
				1,
				1,
				sdk.NewCoins(sdk.NewInt64Coin(sdk.DefaultBondDenom, 10)),
				time.Time{},
			),
			shouldErr: true,
		},
		{
			name: "valid ticket",
			draw: types.NewDraw(
				1,
				1,
				sdk.NewCoins(sdk.NewInt64Coin(sdk.DefaultBondDenom, 10)),
				time.Now(),
			),
			shouldErr: false,
		},
	}

	for _, uc := range usecases {
		uc := uc
		t.Run(uc.name, func(t *testing.T) {
			err := uc.draw.Validate()

			if uc.shouldErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})

	}
}
