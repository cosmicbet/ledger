package types

import (
	"github.com/cosmos/cosmos-sdk/types/query"
)

// NewTicketsRequest returns a new QueryTicketsRequest with the provided data
func NewTicketsRequest(pagination *query.PageRequest) *QueryTicketsRequest {
	return &QueryTicketsRequest{
		Pagination: pagination,
	}
}
