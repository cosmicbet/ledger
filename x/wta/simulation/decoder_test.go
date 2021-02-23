package simulation_test

import (
	"fmt"
	"testing"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/cosmicbet/ledger/app"
	"github.com/cosmicbet/ledger/x/wta/simulation"
	"github.com/cosmicbet/ledger/x/wta/types"

	"github.com/cosmos/cosmos-sdk/types/kv"

	"github.com/stretchr/testify/require"
)

func TestDecodeStore(t *testing.T) {
	encodingCfg := app.MakeEncodingConfig()
	cdc := encodingCfg.Marshaler
	dec := simulation.NewDecodeStore(cdc)

	drawEndTime := time.Date(2020, 1, 1, 00, 00, 00, 000, time.UTC)

	ticket := types.NewTicket(
		"ticket-1",
		time.Date(2020, 1, 1, 00, 00, 00, 000, time.UTC),
		"owner-1",
	)

	historicalDraw := types.NewHistoricalDrawData(
		types.NewDraw(
			1,
			1,
			sdk.NewCoins(sdk.NewInt64Coin(sdk.DefaultBondDenom, 100)),
			time.Date(2020, 1, 1, 00, 00, 00, 000, time.UTC),
		),
		types.NewTicket(
			"ticket-n",
			time.Date(2020, 1, 5, 00, 00, 00, 000, time.UTC),
			"owner-n",
		),
	)

	kvPairs := kv.Pairs{Pairs: []kv.Pair{
		{
			Key:   types.CurrentDrawEndTimeStoreKey,
			Value: types.MustMarshalDrawEndTime(drawEndTime),
		},
		{
			Key:   types.TicketsStoreKey(ticket.Id),
			Value: cdc.MustMarshalBinaryBare(&ticket),
		},
		{
			Key:   types.HistoricalDataStoreKey(historicalDraw.Draw.EndTime),
			Value: cdc.MustMarshalBinaryBare(&historicalDraw),
		},
	}}

	tests := []struct {
		name        string
		expectedLog string
	}{
		{"Draw end time", fmt.Sprintf("CurrentDrawEndTimeA: %s\nCurrentDrawEndTimeB: %s\n",
			drawEndTime.Format(time.RFC3339), drawEndTime.Format(time.RFC3339))},
		{"Ticket", fmt.Sprintf("TicketA: %s\nTicketB: %s\n", &ticket, &ticket)},
		{"Historical draw", fmt.Sprintf("HistoricalDataA: %s\nHistoricalDataB: %s\n", &historicalDraw, &historicalDraw)},
		{"other", ""},
	}

	for i, tt := range tests {
		i, tt := i, tt
		t.Run(tt.name, func(t *testing.T) {
			switch i {
			case len(tests) - 1:
				require.Panics(t, func() { dec(kvPairs.Pairs[i], kvPairs.Pairs[i]) }, tt.name)
			default:
				require.Equal(t, tt.expectedLog, dec(kvPairs.Pairs[i], kvPairs.Pairs[i]), tt.name)
			}
		})
	}
}
