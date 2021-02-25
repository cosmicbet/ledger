package types

import (
	"fmt"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	paramstypes "github.com/cosmos/cosmos-sdk/x/params/types"
)

const (
	// Default draw duration
	DefaultDrawDuration = time.Hour

	// Min draw duration
	MinDrawDuration = time.Minute
)

// Default wta params
var (
	DefaultPrizePercentage = sdk.NewDecWithPrec(98, 2)                                           // 98%
	DefaultBurnPercentage  = sdk.NewDecWithPrec(1, 2)                                            // 1%
	DefaultFeePercentage   = sdk.NewDecWithPrec(1, 2)                                            // 1%
	DefaultTicketPrice     = sdk.NewCoin(sdk.DefaultBondDenom, sdk.TokensFromConsensusPower(10)) // 10 Tokens
)

// Parameters store keys
var (
	ParamStoreDistributionParamsKey = []byte("DistributionParams")
	ParamStoreDrawParamsKey         = []byte("DrawParams")
	ParamStoreTicketParamsKey       = []byte("TicketParams")
)

// ParamKeyTable Key declaration for parameters
func ParamKeyTable() paramstypes.KeyTable {
	return paramstypes.NewKeyTable(
		paramstypes.NewParamSetPair(ParamStoreDistributionParamsKey, &DistributionParams{}, ValidateDistributionParams),
		paramstypes.NewParamSetPair(ParamStoreDrawParamsKey, &DrawParams{}, ValidateDrawParams),
		paramstypes.NewParamSetPair(ParamStoreTicketParamsKey, &TicketParams{}, ValidateTicketParams),
	)
}

// -------------------------------------------------------------------------------------------------------------------

func NewDistributionParams(prizePercentage, feePercentage, burnPercentage sdk.Dec) DistributionParams {
	return DistributionParams{
		PrizePercentage: prizePercentage,
		FeePercentage:   feePercentage,
		BurnPercentage:  burnPercentage,
	}
}

func DefaultDistributionParams() DistributionParams {
	return NewDistributionParams(
		DefaultPrizePercentage,
		DefaultBurnPercentage,
		DefaultFeePercentage,
	)
}

func ValidateDistributionParams(i interface{}) error {
	params, ok := i.(DistributionParams)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	err := validatePercentageValue(params.PrizePercentage)
	if err != nil {
		return err
	}

	err = validatePercentageValue(params.FeePercentage)
	if err != nil {
		return err
	}

	err = validatePercentageValue(params.BurnPercentage)
	if err != nil {
		return err
	}

	if !params.PrizePercentage.Add(params.FeePercentage).Add(params.BurnPercentage).Equal(sdk.NewDecWithPrec(100, 2)) {
		return fmt.Errorf("percentages does not sum to 1.00")
	}

	return nil
}

// validatePercentageValue validates a percentage value making sure it's not negative or exceeding 100
func validatePercentageValue(i interface{}) error {
	params, isCorrectParam := i.(sdk.Dec)

	if !isCorrectParam {
		return fmt.Errorf("invalid parameters type: %s", i)
	}

	if params.IsZero() || params.IsNegative() || params.GT(sdk.NewDecWithPrec(100, 2)) {
		return fmt.Errorf("invalid percentage value: %s", params)
	}

	return nil
}

// -------------------------------------------------------------------------------------------------------------------

func NewDrawParams(duration time.Duration) DrawParams {
	return DrawParams{
		Duration: duration,
	}
}

func DefaultDrawParams() DrawParams {
	return NewDrawParams(DefaultDrawDuration)
}

func ValidateDrawParams(i interface{}) error {
	params, ok := i.(DrawParams)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if params.Duration == 0 || params.Duration < MinDrawDuration {
		return fmt.Errorf("invalid draw duration param: %s", params.Duration)
	}

	return nil
}

// -------------------------------------------------------------------------------------------------------------------

func NewTicketParams(price sdk.Coin) TicketParams {
	return TicketParams{
		Price: price,
	}
}

func DefaultTicketParams() TicketParams {
	return NewTicketParams(DefaultTicketPrice)
}

func ValidateTicketParams(i interface{}) error {
	params, ok := i.(TicketParams)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if err := params.Price.Validate(); err != nil {
		return fmt.Errorf("invalid ticket price param: %s", err.Error())
	}

	if params.Price.IsZero() {
		return fmt.Errorf("ticket price cannot be zero")
	}

	return nil
}
