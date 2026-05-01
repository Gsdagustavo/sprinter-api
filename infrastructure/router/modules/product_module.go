package modules

import (
	"log/slog"
	"net/http"

	"github.com/Gsdagustavo/sprinter-api/domain/usecases"
	"github.com/Gsdagustavo/sprinter-api/infrastructure/router"
	"github.com/Gsdagustavo/sprinter-api/infrastructure/router/logger"
	"github.com/gorilla/mux"
)

type productModule struct {
	productUseCases usecases.ProductUseCase
	name            string
	path            string
}

func NewProductModule(productUseCases usecases.ProductUseCase) router.Module {
	return productModule{
		productUseCases: productUseCases,
		name:            "Product",
		path:            "/product",
	}
}

func (m productModule) Name() string {
	return m.name
}

func (m productModule) Path() string {
	return m.path
}

func (m productModule) Routes() []router.RouteDefinition {
	return []router.RouteDefinition{
		{
			Path:        "",
			Description: "List products",
			Handler:     m.listProducts,
			HttpMethods: []string{http.MethodGet},
			Public: true,
		},
	}
}

func (m productModule) Middlewares() []mux.MiddlewareFunc {
	return []mux.MiddlewareFunc{}
}

func (m productModule) listProducts(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	filter, err := router.GetDefaultFilterFromParams(r)
	if err != nil {
		slog.ErrorContext(ctx, "failed to get default filter params", logger.Err(err))
		router.HandleError(w, err)
		return
	}

	products, err := m.productUseCases.GetProducts(ctx, *filter)
	if err != nil {
		slog.ErrorContext(ctx, "failed to get products", logger.Err(err))
		router.HandleError(w, err)
		return
	}

	err = router.Write(w, products)
	if err != nil {
		slog.ErrorContext(ctx, "failed to write response", logger.Err(err))
	}
}
