package entities

// PaginatedList is a generic struct that defines basic attributes for pagination.
type PaginatedList[T any] struct {
	// The list of items in.
	Items []T `json:"items"`

	// RequestedItems is the number of items that the client requested for.
	RequestedItems int64 `json:"requested_items"`

	// TotalCount is the total number of items in the database that match the given criteria
	TotalCount int64 `json:"total_count"`

	// Pages is the total number of pages considering the limit given by the client.
	Pages int64 `json:"pages"`
}
