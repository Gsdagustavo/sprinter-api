package router

import (
	"context"
	"errors"
	"net/http"

	"github.com/VitorFranciscoDev/sprinter-api/domain/entities"
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

	return contextUser.(*entities.User), nil
}
