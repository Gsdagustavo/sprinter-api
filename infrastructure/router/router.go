package router

import (
	"net/http"

	"github.com/gorilla/mux"
)

type Module interface {
	// Name returns the module name
	Name() string

	// Path returns the module base path
	Path() string

	// Setup sets all the route handlers
	//
	// Returns a list of all the routes defined by the module and optionally a new router base
	Setup(r *mux.Router) ([]RouteDefinition, *mux.Router)
}

type RouteDefinition struct {
	// Path is the path for the route
	Path string

	// Description is a small text describing the route
	Description string

	// Handler is the function handler for the route
	Handler http.HandlerFunc

	// HttpMethods is a list of HTTP methods accepted by the route
	HttpMethods []string
}
