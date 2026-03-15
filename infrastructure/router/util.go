package router

import (
	"context"
	"errors"
	"net/http"
	"strconv"
	"strings"

	"github.com/Gsdagustavo/sprinter-api/domain/entities"
	"github.com/Gsdagustavo/sprinter-api/domain/entities/derr"
)

const (
	// defaultLimit default limit parameter for pagination
	defaultLimit = 30

	// defaultLimit default page parameter for pagination
	defaultPage = 1
)

type userKey string

const contextUserKey userKey = "user"

var ErrUserNotFoundInRequest = errors.New("user not found in request")

// WithUser adds the given user to the request's context'
func WithUser(ctx context.Context, user *entities.User) context.Context {
	return context.WithValue(ctx, contextUserKey, user)
}

// GetUser attempts to retrieve the user in the request's context
func GetUser(r *http.Request) (*entities.User, error) {
	contextUser := r.Context().Value(contextUserKey)
	if contextUser == nil {
		return nil, ErrUserNotFoundInRequest
	}

	user, ok := contextUser.(*entities.User)
	if user == nil || !ok {
		return nil, ErrUserNotFoundInRequest
	}

	return user, nil
}

// GetDefaultFilterFromParams returns a default filter entity from the given HTTP request.
func GetDefaultFilterFromParams(r *http.Request) (*entities.GeneralFilter, error) {
	var filter entities.GeneralFilter
	query := r.URL.Query()

	// Get "limit" parameter for pagination
	limitStr := query.Get("limit")
	if limitStr != "" {
		limit, err := strconv.ParseInt(limitStr, 0, 64)
		if err != nil {
			return nil, derr.NewBadRequestError("failed to get limit parameter")
		}

		filter.Limit = limit
	}

	// Get "page" parameter for pagination
	pageStr := query.Get("page")
	if pageStr != "" {
		page, err := strconv.ParseInt(pageStr, 0, 64)
		if err != nil {
			return nil, derr.NewBadRequestError("failed to get page parameter")
		}

		filter.Page = page
	}

	// Ensure that both limit and page parameters are valid
	if filter.Limit <= 0 {
		filter.Limit = defaultLimit
	}

	if filter.Page <= 0 {
		filter.Page = defaultPage
	}

	// Get search and ordering parameters
	filter.Search = strings.TrimSpace(query.Get("search"))
	filter.OrderBy = strings.TrimSpace(query.Get("orderBy"))
	ordination := strings.TrimSpace(strings.ToUpper(query.Get("ordination")))
	if ordination == "ASC" || ordination == "DESC" {
		filter.Ordination = ordination
	} else {
		filter.Ordination = "ASC"
	}

	return &filter, nil
}
