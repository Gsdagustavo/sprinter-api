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

	// ApiMethods are the API methods allowed by the route
	//
	// Example:
	//   - /activity/list: only allows SELECT method
	//   - /activity/type/register: only allows INSERT method
	//   - /activity/type/update: only allows UPDATE method
	ApiMethods []RouteApiMethod

	// Module is the module which this route is defined in
	//
	// Note: This is needed because some routes can be defined in the same Module object (the ones that define the
	// route handlers), but used in a different panel module (the ones displayed on the panel sidebar).
	//
	// Example:
	//   - /activity/list: defined in /activity module, used in "Activities" module
	//   - /activity/type/list: defined in /activity module, used in "Activity types" module
	//   - /activity/locality/list: defined in /activity module, used in "Localities" module
	//
	// These three endpoints/routes are defined in the same module, but used in different panel modules.
	//
	// Note: This field can (and should) be omitted if the route does not require access level checking, such as:
	//   - /login: any user can log in
	//   - /user/getData: any user can get their data
	//   - /user/updatePassword: any user can update their password
	//   - /user/acceptPolicies: any user can accept the terms and conditions
	Module ApiModuleName
}

// RouteApiMethod is the API method used by a RouteDefinition
type RouteApiMethod int
type ApiModuleName string

const (
	// RouteApiMethodSelect is used on routes that return a list of items
	RouteApiMethodSelect RouteApiMethod = iota

	// RouteApiMethodInsert is used on routes that insert a new item
	RouteApiMethodInsert

	// RouteApiMethodUpdate is used on routes that update an existing item
	RouteApiMethodUpdate

	// RouteApiMethodDelete is used on routes that delete or activate/deactivate an existing item
	RouteApiMethodDelete
)

// NOTE: These MUST match exactly the module name stored in the "module" database table
const (
	ApiModuleUser             ApiModuleName = "user"
	ApiModulePromoter         ApiModuleName = "promoter"
	ApiModuleSupplier         ApiModuleName = "supplier"
	ApiModuleActivity         ApiModuleName = "activity"
	ApiModuleActivityReason   ApiModuleName = "activity_reason"
	ApiModuleVisit            ApiModuleName = "visit"
	ApiModuleActivityType     ApiModuleName = "activity_type"
	ApiModuleActivityLocality ApiModuleName = "activity_locality"
)
