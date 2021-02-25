package types_test

import (
	"testing"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"

	"github.com/cosmicbet/ledger/x/wta/types"
)

func TestValidateDistributionParams(t *testing.T) {
	usecases := []struct {
		name      string
		params    types.DistributionParams
		shouldErr bool
	}{
		{
			name: "invalid prize percentage",
			params: types.NewDistributionParams(
				sdk.NewDecWithPrec(0, 2),
				sdk.NewDecWithPrec(99, 2),
				sdk.NewDecWithPrec(1, 2),
			),
			shouldErr: true,
		},
		{
			name: "invalid fee percentage",
			params: types.NewDistributionParams(
				sdk.NewDecWithPrec(90, 2),
				sdk.NewDecWithPrec(101, 2),
				sdk.NewDecWithPrec(1, 2),
			),
			shouldErr: true,
		},
		{
			name: "invalid burn percentage",
			params: types.NewDistributionParams(
				sdk.NewDecWithPrec(1, 2),
				sdk.NewDecWithPrec(99, 2),
				sdk.NewDecWithPrec(-1, 2),
			),
			shouldErr: true,
		},
		{
			name: "invalid percentages sum (> 1)",
			params: types.NewDistributionParams(
				sdk.NewDecWithPrec(90, 2),
				sdk.NewDecWithPrec(5, 2),
				sdk.NewDecWithPrec(6, 2),
			),
			shouldErr: true,
		},
		{
			name: "invalid percentages sum (< 1)",
			params: types.NewDistributionParams(
				sdk.NewDecWithPrec(90, 2),
				sdk.NewDecWithPrec(5, 2),
				sdk.NewDecWithPrec(4, 2),
			),
			shouldErr: true,
		},
	}

	for _, uc := range usecases {
		uc := uc
		t.Run(uc.name, func(t *testing.T) {
			err := types.ValidateDistributionParams(uc.params)
			if uc.shouldErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestValidateDrawParams(t *testing.T) {
	usecases := []struct {
		name      string
		params    types.DrawParams
		shouldErr bool
	}{

		{
			name:      "zero duration",
			params:    types.NewDrawParams(time.Minute * 0),
			shouldErr: true,
		},
		{
			name:      "invalid duration",
			params:    types.NewDrawParams(time.Second * 30),
			shouldErr: true,
		},
	}

	for _, uc := range usecases {
		uc := uc
		t.Run(uc.name, func(t *testing.T) {
			err := types.ValidateDrawParams(uc.params)
			if uc.shouldErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestValidateTicketParams(t *testing.T) {
	usecases := []struct {
		name      string
		params    types.TicketParams
		shouldErr bool
	}{
		{
			name: "invalid ticket price",
			params: types.NewTicketParams(
				sdk.Coin{Denom: "./", Amount: sdk.NewInt(100)},
			),
			shouldErr: true,
		},
		{
			name: "valid params",
			params: types.NewTicketParams(
				sdk.NewInt64Coin("stake", 100),
			),
			shouldErr: false,
		},
	}

	for _, uc := range usecases {
		uc := uc
		t.Run(uc.name, func(t *testing.T) {
			err := types.ValidateTicketParams(uc.params)
			if uc.shouldErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}
