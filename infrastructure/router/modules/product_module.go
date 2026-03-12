package modules

import (
	"github.com/Gsdagustavo/sprinter-api/domain"
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

func (p productModule) Name() string {
	return p.name
}

func (p productModule) Path() string {
	return p.path
}

func (p productModule) Setup(r *mux.Router) ([]router.RouteDefinition, *mux.Router) {
	return nil, r
}
