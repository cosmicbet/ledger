package keeper_test

import (
	"testing"
	"time"

	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/store"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authkeeper "github.com/cosmos/cosmos-sdk/x/auth/keeper"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	bankkeeper "github.com/cosmos/cosmos-sdk/x/bank/keeper"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	distrkeeper "github.com/cosmos/cosmos-sdk/x/distribution/keeper"
	distrtypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	paramskeeper "github.com/cosmos/cosmos-sdk/x/params/keeper"
	paramstypes "github.com/cosmos/cosmos-sdk/x/params/types"
	stakingkeeper "github.com/cosmos/cosmos-sdk/x/staking/keeper"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/stretchr/testify/suite"
	"github.com/tendermint/tendermint/libs/log"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"
	db "github.com/tendermint/tm-db"

	"github.com/cosmicbet/ledger/app"
	wtakeeper "github.com/cosmicbet/ledger/x/wta/keeper"
	wtatypes "github.com/cosmicbet/ledger/x/wta/types"
)

func TestKeeperTestSuite(t *testing.T) {
	suite.Run(t, new(KeeperTestSuite))
}

type KeeperTestSuite struct {
	suite.Suite

	cdc      codec.BinaryMarshaler
	ctx      sdk.Context
	storeKey sdk.StoreKey
	keeper   wtakeeper.Keeper
	ak       authkeeper.AccountKeeper
	bk       bankkeeper.Keeper
	dk       distrkeeper.Keeper
	pk       paramskeeper.Keeper
}

func (suite *KeeperTestSuite) SetupTest() {

	// Store keys
	keys := sdk.NewKVStoreKeys(
		authtypes.StoreKey, banktypes.StoreKey, distrtypes.StoreKey, paramstypes.StoreKey, stakingtypes.StoreKey,
		wtatypes.StoreKey,
	)
	suite.storeKey = keys[wtatypes.StoreKey]

	// Transient keys
	tKeys := sdk.NewTransientStoreKeys(paramstypes.TStoreKey)

	// Create an in-memory db
	memDB := db.NewMemDB()
	ms := store.NewCommitMultiStore(memDB)

	// Mount keys
	for _, key := range keys {
		ms.MountStoreWithDB(key, sdk.StoreTypeIAVL, memDB)
	}

	// Mount transient keys
	for _, key := range tKeys {
		ms.MountStoreWithDB(key, sdk.StoreTypeTransient, memDB)
	}

	// Load the database
	err := ms.LoadLatestVersion()
	if err != nil {
		panic(err)
	}

	// Create a custom ctx with custom time
	blockTime, _ := time.Parse(time.RFC3339, "2021-01-01T00:00:00.000Z")
	suite.ctx = sdk.NewContext(
		ms,
		tmproto.Header{ChainID: "test-chain-id", Time: blockTime},
		false,
		log.NewNopLogger(),
	)

	encodingConfig := app.MakeEncodingConfig()
	suite.cdc = encodingConfig.Marshaler

	// Build the keepers
	suite.pk = paramskeeper.NewKeeper(suite.cdc, encodingConfig.Amino, keys[paramstypes.StoreKey], tKeys[paramstypes.TStoreKey])

	suite.ak = authkeeper.NewAccountKeeper(
		suite.cdc, keys[authtypes.StoreKey], suite.pk.Subspace(authtypes.ModuleName),
		authtypes.ProtoBaseAccount, app.GetMaccPerms(),
	)

	suite.bk = bankkeeper.NewBaseKeeper(
		suite.cdc, keys[banktypes.StoreKey], suite.ak, suite.pk.Subspace(banktypes.ModuleName), app.BlockedAddrs(),
	)

	sk := stakingkeeper.NewKeeper(
		suite.cdc, keys[stakingtypes.StoreKey], suite.ak, suite.bk, suite.pk.Subspace(stakingtypes.ModuleName),
	)

	suite.dk = distrkeeper.NewKeeper(
		suite.cdc, keys[distrtypes.StoreKey], suite.pk.Subspace(distrtypes.ModuleName),
		suite.ak, suite.bk, &sk,
		authtypes.FeeCollectorName, app.BlockedAddrs(),
	)

	// Default fees to avoid errors
	suite.dk.SetFeePool(suite.ctx, distrtypes.InitialFeePool())

	suite.keeper = wtakeeper.NewKeeper(
		suite.cdc, keys[wtatypes.StoreKey], suite.pk.Subspace(wtatypes.DefaultParamSpace), suite.bk, suite.dk,
	)
}
