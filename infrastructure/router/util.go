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

	limitStr := query.Get("limit")
	if limitStr != "" {
		limit, err := strconv.ParseInt(limitStr, 0, 64)
		if err != nil {
			return nil, derr.NewBadRequestError("failed to get limit parameter")
		}

		if limit < 0 {
			limit = defaultLimit
		}

		filter.Limit = limit
	}

	pageStr := query.Get("page")
	if pageStr != "" {
		page, err := strconv.ParseInt(pageStr, 0, 64)
		if err != nil {
			return nil, derr.NewBadRequestError("failed to get page parameter")
		}

		if page < 0 {
			page = defaultPage
		}

		filter.Page = page
	}

	filter.Search = strings.TrimSpace(query.Get("search"))

	return &filter, nil
}
