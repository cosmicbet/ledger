package types_test

import (
	"github.com/cosmicbet/ledger/x/wta/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestParams_Validate(t *testing.T) {
	usecases := []struct {
		name      string
		params    types.Params
		shouldErr bool
	}{
		{
			name: "invalid prize percentage",
			params: types.NewParams(
				sdk.NewInt(0),
				sdk.NewInt(99),
				sdk.NewInt(1),
				time.Minute*5,
				sdk.NewInt64Coin("stake", 100),
			),
			shouldErr: true,
		},
		{
			name: "invalid community pool percentage",
			params: types.NewParams(
				sdk.NewInt(90),
				sdk.NewInt(101),
				sdk.NewInt(1),
				time.Minute*5,
				sdk.NewInt64Coin("stake", 100),
			),
			shouldErr: true,
		},
		{
			name: "invalid burn percentage",
			params: types.NewParams(
				sdk.NewInt(1),
				sdk.NewInt(99),
				sdk.NewInt(-1),
				time.Minute*5,
				sdk.NewInt64Coin("stake", 100),
			),
			shouldErr: true,
		},
		{
			name: "invalid percentages sum (> 100)",
			params: types.NewParams(
				sdk.NewInt(90),
				sdk.NewInt(5),
				sdk.NewInt(6),
				time.Minute*5,
				sdk.NewInt64Coin("stake", 100),
			),
			shouldErr: true,
		},
		{
			name: "invalid percentages sum (<100)",
			params: types.NewParams(
				sdk.NewInt(90),
				sdk.NewInt(5),
				sdk.NewInt(4),
				time.Minute*5,
				sdk.NewInt64Coin("stake", 100),
			),
			shouldErr: true,
		},
		{
			name: "invalid duration",
			params: types.NewParams(
				sdk.NewInt(90),
				sdk.NewInt(5),
				sdk.NewInt(5),
				time.Minute*0,
				sdk.NewInt64Coin("stake", 100),
			),
			shouldErr: true,
		},
		{
			name: "invalid ticket price",
			params: types.NewParams(
				sdk.NewInt(90),
				sdk.NewInt(5),
				sdk.NewInt(5),
				time.Minute*5,
				sdk.Coin{Denom: "./", Amount: sdk.NewInt(100)},
			),
			shouldErr: true,
		},
		{
			name: "valid params",
			params: types.NewParams(
				sdk.NewInt(93),
				sdk.NewInt(4),
				sdk.NewInt(3),
				time.Minute*5,
				sdk.NewInt64Coin("stake", 100),
			),
			shouldErr: false,
		},
	}

	for _, uc := range usecases {
		uc := uc
		t.Run(uc.name, func(t *testing.T) {
			err := uc.params.Validate()
			if uc.shouldErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}
