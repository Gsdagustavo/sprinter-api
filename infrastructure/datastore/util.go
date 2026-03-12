package datastore

import (
	"fmt"
	"math"
	"strings"

	"github.com/Gsdagustavo/sprinter-api/domain/entities"
)

// GetQueryCount returns the count version of the given query
func GetQueryCount(query string) string {
	idx := strings.Index(query, "ORDER BY")
	if idx < 0 {
		idx = len(query)
	}
	base := "SELECT COUNT(*) FROM ("
	return fmt.Sprintf(base+"%s\n)", query[:idx])
}

// GetPaginated returns the given query with pagination based on the given filter
func GetPaginated(query string, filter entities.GeneralFilter) string {
	return fmt.Sprintf("%s LIMIT %d OFFSET %d", query, filter.Limit, filter.Limit*(filter.Page-1))
}

// GetTotalPages returns the total number of pages in the database based on the given filter
func GetTotalPages(totalCount int64, filter entities.GeneralFilter) int64 {
	return int64(math.Ceil(float64(totalCount) / float64(filter.Limit)))
}
