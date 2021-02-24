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
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}

	sdkCtx := sdk.UnwrapSDKContext(ctx)

	store := sdkCtx.KVStore(k.storeKey)
	ticketsStore := prefix.NewStore(store, types.TicketsStorePrefix)

	var tickets []types.Ticket
	pageRes, err := query.FilteredPaginate(ticketsStore, req.Pagination, func(_ []byte, value []byte, accumulate bool) (bool, error) {
		ticket, err := types.UnmarshalTicket(k.cdc, value)
		if err != nil {
			return false, err
		}

		if accumulate {
			tickets = append(tickets, ticket)
		}

		return true, nil
	})

	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &types.QueryTicketsResponse{Tickets: tickets, Pagination: pageRes}, err
}

// NextDraw queries the details of the next draw
func (k querier) NextDraw(ctx context.Context, req *types.QueryNextDrawRequest) (*types.QueryNextDrawResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}

	sdkCtx := sdk.UnwrapSDKContext(ctx)
	draw := k.GetCurrentDraw(sdkCtx)
	return &types.QueryNextDrawResponse{Draw: draw}, nil
}

// PastDraws queries the details of the past draws
func (k querier) PastDraws(ctx context.Context, req *types.QueryPastDrawsRequest) (*types.QueryPastDrawsResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}

	sdkCtx := sdk.UnwrapSDKContext(ctx)

	store := sdkCtx.KVStore(k.storeKey)
	drawsStore := prefix.NewStore(store, types.HistoricalDrawStorePrefix)

	var draws []types.HistoricalDrawData
	pageRes, err := query.FilteredPaginate(drawsStore, req.Pagination, func(_ []byte, value []byte, accumulate bool) (bool, error) {
		data, err := types.UnmarshalHistoricalDraw(k.cdc, value)
		if err != nil {
			return false, err
		}

		if accumulate {
			draws = append(draws, data)
		}

		return true, nil
	})

	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &types.QueryPastDrawsResponse{Draws: draws, Pagination: pageRes}, nil
}

// Params queries the currently stored parameters
func (k Keeper) Params(ctx context.Context, req *types.QueryParamsRequest) (*types.QueryParamsResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}

	sdkCtx := sdk.UnwrapSDKContext(ctx)
	params := k.GetParams(sdkCtx)
	return &types.QueryParamsResponse{Params: params}, nil
}
