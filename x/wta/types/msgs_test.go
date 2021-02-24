package types_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/cosmicbet/ledger/x/wta/types"
)

func TestMsgBuyTickets_ValidateBasic(t *testing.T) {
	usecases := []struct {
		name      string
		msg       *types.MsgBuyTickets
		shouldErr bool
	}{
		{
			name:      "invalid quantity",
			msg:       types.NewMsgBuyTickets(0, "cosmos14zfwkjm35j05ydm3s3qu4he39yjxe9575echwl"),
			shouldErr: true,
		},
		{
			name:      "invalid buyer",
			msg:       types.NewMsgBuyTickets(1, "buyer"),
			shouldErr: true,
		},
		{
			name:      "valid message",
			msg:       types.NewMsgBuyTickets(1, "cosmos14zfwkjm35j05ydm3s3qu4he39yjxe9575echwl"),
			shouldErr: false,
		},
	}

	for _, uc := range usecases {
		uc := uc
		t.Run(uc.name, func(t *testing.T) {
			err := uc.msg.ValidateBasic()

			if uc.shouldErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}
