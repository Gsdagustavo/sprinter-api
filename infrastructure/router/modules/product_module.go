package modules

import (
	"log/slog"
	"net/http"

	"github.com/Gsdagustavo/sprinter-api/domain"
	"github.com/Gsdagustavo/sprinter-api/domain/logger"
	"github.com/Gsdagustavo/sprinter-api/infrastructure/router"
	"github.com/gorilla/mux"
)

type productModule struct {
	productUseCases domain.ProductUseCase
	name            string
	path            string
}

func NewProductModule(productUseCases domain.ProductUseCase) router.Module {
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

func (m productModule) Setup(r *mux.Router) ([]router.RouteDefinition, *mux.Router) {
	defs := []router.RouteDefinition{
		{
			Path:        "",
			Description: "List products",
			Handler:     m.listProducts,
			HttpMethods: []string{http.MethodGet},
		},
	}

	for _, d := range defs {
		r.HandleFunc(m.path+d.Path, d.Handler).Methods(d.HttpMethods...)
	}

	return defs, r
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
