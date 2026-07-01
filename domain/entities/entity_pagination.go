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

// CursorPaginatedList is a generic struct that defines cursor pagination attributes.
type CursorPaginatedList[T any] struct {
	// Items is the list of items returned for the current cursor.
	Items []T `json:"items"`

	// RequestedItems is the number of items that the client requested for.
	RequestedItems int64 `json:"requested_items"`

	// HasNext defines whether there are more items after this page.
	HasNext bool `json:"has_next"`

	// NextCursor is the cursor that should be used to request the next page.
	NextCursor *int64 `json:"next_cursor,omitempty"`
}

// GeneralFilter is a generic struct that defines basic fields for pagination, sorting and filtering.
type GeneralFilter struct {
	// Limit is the maximum number of items that should be returned on a single page.
	Limit int64 `json:"limit"`

	// Page is the current page that should be returned.
	Page int64 `json:"page"`

	// OrderBy is the column that is being ordered in the request.
	OrderBy string `json:"order_by"`

	// Ordination is the type of ordination. Could be either ASC or DESC.
	Ordination string `json:"ordination"`

	// Search is the search used to filter items.
	Search string `json:"search"`
}

// CursorFilter is a generic struct that defines basic fields for cursor pagination.
type CursorFilter struct {
	// Limit is the maximum number of items that should be returned on a single page.
	Limit int64 `json:"limit"`

	// Cursor is the last seen item identifier.
	Cursor int64 `json:"cursor"`
}
