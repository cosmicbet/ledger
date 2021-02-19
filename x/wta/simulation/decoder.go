package simulation

import (
	"bytes"
	"fmt"

	"github.com/cosmicbet/ledger/x/wta/types"

	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/types/kv"
)

// NewDecodeStore returns a new decoder that unmarshals the KVPair's Value
// to the corresponding relationships type
func NewDecodeStore(cdc codec.Marshaler) func(kvA, kvB kv.Pair) string {
	return func(kvA, kvB kv.Pair) string {
		switch {
		case bytes.HasPrefix(kvA.Key, types.TicketsStorePrefix):
			var ticketA, ticketB types.Ticket
			cdc.MustUnmarshalBinaryBare(kvA.Value, &ticketA)
			cdc.MustUnmarshalBinaryBare(kvB.Value, &ticketB)
			return fmt.Sprintf("TicketA: %s\nTicketB: %s\n", &ticketA, &ticketB)

		case bytes.HasPrefix(kvA.Key, types.HistoricalDrawStorePrefix):
			var dataA, dataB types.HistoricalDrawData
			cdc.MustUnmarshalBinaryBare(kvA.Value, &dataA)
			cdc.MustUnmarshalBinaryBare(kvB.Value, &dataB)
			return fmt.Sprintf("HistoricalDataA: %s\nHistoricalDataB: %s\n", &dataA, &dataB)

		case bytes.Equal(kvA.Key, types.CurrentDrawStoreKey):
			var drawA, drawB types.Draw
			cdc.MustUnmarshalBinaryBare(kvA.Value, &drawA)
			cdc.MustUnmarshalBinaryBare(kvB.Value, &drawB)
			return fmt.Sprintf("CurrentDrawA: %s\nCurrentDrawB: %s\n", &drawA, &drawB)

		default:
			panic(fmt.Sprintf("unexpected %s key %X (%s)", types.ModuleName, kvA.Key, kvA.Key))
		}
	}
}
