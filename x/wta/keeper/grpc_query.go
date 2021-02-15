package keeper

import (
	"context"

	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/query"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/cosmicbet/ledger/x/wta/types"
)

var _ types.QueryServer = querier{}

// querier is used as Keeper will have duplicate methods if used directly, and gRPC names take precedence over keeper
type querier struct {
	Keeper
}

// NewMsgServerImpl returns an implementation of the wta QueryServer interface
// for the provided Keeper.
func NewQuerierImpl(k Keeper) types.QueryServer {
	return querier{Keeper: k}
}

// Tickets queries all tickets for the next draw
func (k querier) Tickets(ctx context.Context, req *types.QueryTicketsRequest) (*types.QueryTicketsResponse, error) {
	sdkCtx := sdk.UnwrapSDKContext(ctx)

	store := sdkCtx.KVStore(k.storeKey)
	ticketsStore := prefix.NewStore(store, types.TicketsStorePrefix)

	var tickets types.Tickets
	pageRes, err := query.FilteredPaginate(ticketsStore, req.Pagination, func(_ []byte, value []byte, accumulate bool) (bool, error) {
		ticket, err := types.UnmarshalTicket(k.cdc, value)
		if err != nil {
			return false, err
		}

		if accumulate {
			tickets.Tickets = append(tickets.Tickets, ticket)
		}

		return true, nil
	})

	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &types.QueryTicketsResponse{Tickets: tickets, Pagination: pageRes}, err
}

// Draw queries the details of the next draw
func (k querier) Draw(ctx context.Context, _ *types.QueryDrawRequest) (*types.QueryDrawResponse, error) {
	sdkCtx := sdk.UnwrapSDKContext(ctx)

	draw := k.GetCurrentDraw(sdkCtx)

	return &types.QueryDrawResponse{Draw: draw}, nil
}

// Params queries the currently stored parameters
func (k Keeper) Params(goCtx context.Context, _ *types.QueryParamsRequest) (*types.QueryParamsResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	params := k.GetParams(ctx)
	return &types.QueryParamsResponse{Params: params}, nil
}
