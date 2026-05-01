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

	// Routes returns all the module's routes
	Routes() []RouteDefinition

	// Middlewares returns all the module's middlewares
	Middlewares() []mux.MiddlewareFunc
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

	// Public defines whether the route is public
	Public bool
}
