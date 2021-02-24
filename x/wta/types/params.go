package types

import (
	"fmt"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	paramstypes "github.com/cosmos/cosmos-sdk/x/params/types"
)

const (
	// Default params space for the params keeper
	DefaultParamSpace = ModuleName
)

// Parameters store keys
var (
	PrizePercentageParamKey         = []byte("PrizePercentage")
	CommunityPoolPercentageParamKey = []byte("CommunityPoolPercentage")
	BurnPercentageParamKey          = []byte("BurnPercentage")
	DrawDurationParamKey            = []byte("DrawDuration")
	TicketPriceParamKey             = []byte("TicketPrice")
)

// ParamKeyTable Key declaration for parameters
func ParamKeyTable() paramstypes.KeyTable {
	return paramstypes.NewKeyTable().RegisterParamSet(&Params{})
}

// NewParams creates a new Params object
func NewParams(
	prizePercentage, communityPoolPercentage, burnPercentage sdk.Int,
	drawDuration time.Duration, ticketPrice sdk.Coin,
) Params {
	return Params{
		PrizePercentage:         prizePercentage,
		CommunityPoolPercentage: communityPoolPercentage,
		BurnPercentage:          burnPercentage,
		DrawDuration:            drawDuration,
		TicketPrice:             ticketPrice,
	}
}

// DefaultParams return default params object
func DefaultParams() Params {
	return Params{
		PrizePercentage:         sdk.NewInt(98),
		CommunityPoolPercentage: sdk.NewInt(1),
		BurnPercentage:          sdk.NewInt(1),
		DrawDuration:            time.Hour * 24,
		TicketPrice:             sdk.NewCoin(sdk.DefaultBondDenom, sdk.NewInt(10000000)),
	}
}

// ParamSetPairs implements the ParamSet interface and returns the key/value pairs
// of posts module's parameters.
func (params *Params) ParamSetPairs() paramstypes.ParamSetPairs {
	return paramstypes.ParamSetPairs{
		paramstypes.NewParamSetPair(PrizePercentageParamKey, &params.PrizePercentage, ValidatePercentageValue),
		paramstypes.NewParamSetPair(CommunityPoolPercentageParamKey, &params.CommunityPoolPercentage, ValidatePercentageValue),
		paramstypes.NewParamSetPair(BurnPercentageParamKey, &params.BurnPercentage, ValidatePercentageValue),
		paramstypes.NewParamSetPair(DrawDurationParamKey, &params.DrawDuration, ValidateDurationValue),
		paramstypes.NewParamSetPair(TicketPriceParamKey, &params.TicketPrice, ValidateTicketPriceValue),
	}
}

// Validate perform basic checks on all parameters to ensure they are correct
func (params Params) Validate() error {
	err := ValidatePercentageValue(params.PrizePercentage)
	if err != nil {
		return err
	}

	err = ValidatePercentageValue(params.CommunityPoolPercentage)
	if err != nil {
		return err
	}

	err = ValidatePercentageValue(params.BurnPercentage)
	if err != nil {
		return err
	}

	if !params.PrizePercentage.Add(params.CommunityPoolPercentage).Add(params.BurnPercentage).Equal(sdk.NewInt(100)) {
		return fmt.Errorf("percentages does not sum to 100")
	}

	err = ValidateDurationValue(params.DrawDuration)
	if err != nil {
		return err
	}

	err = ValidateTicketPriceValue(params.TicketPrice)
	if err != nil {
		return err
	}

	return nil
}

// ValidatePercentageValue validates a percentage value making sure it's not negative or exceeding 100
func ValidatePercentageValue(i interface{}) error {
	params, isCorrectParam := i.(sdk.Int)

	if !isCorrectParam {
		return fmt.Errorf("invalid parameters type: %s", i)
	}

	if params.IsZero() || params.IsNegative() || params.GT(sdk.NewInt(100)) {
		return fmt.Errorf("invalid percentage value: %s", params)
	}

	return nil
}

// ValidateDurationValue validates a duration value making sure it's not zero
func ValidateDurationValue(i interface{}) error {
	duration, isCorrectParam := i.(time.Duration)

	if !isCorrectParam {
		return fmt.Errorf("invalid parameters type: %s", i)
	}

	if duration == 0 || duration < time.Minute {
		return fmt.Errorf("invalid draw duration param: %s", duration)
	}

	return nil
}

// ValidateTicketPriceValue validates a ticket price value
func ValidateTicketPriceValue(i interface{}) error {
	price, isCorrectParam := i.(sdk.Coin)

	if !isCorrectParam {
		return fmt.Errorf("invalid parameters type: %s", i)
	}

	if err := price.Validate(); err != nil {
		return fmt.Errorf("invalid ticket price param: %s", err.Error())
	}

	if price.IsZero() {
		return fmt.Errorf("ticket price cannot be zero")
	}

	return nil
}
