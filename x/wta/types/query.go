package types

import (
	"github.com/cosmos/cosmos-sdk/types/query"
)

// NewTicketsRequest returns a new QueryTicketsRequest with the provided pagination data
func NewTicketsRequest(pagination *query.PageRequest) *QueryTicketsRequest {
	return &QueryTicketsRequest{
		Pagination: pagination,
	}
}

// NewPastDrawsRequest returns a new QueryPastDrawsRequest with the provided pagination data
func NewPastDrawsRequest(pagination *query.PageRequest) *QueryPastDrawsRequest {
	return &QueryPastDrawsRequest{
		Pagination: pagination,
	}
}
