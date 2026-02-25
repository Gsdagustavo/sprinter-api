package usecases

import (
	"context"

	"github.com/VitorFranciscoDev/sprinter-api/domain/entities"
	"github.com/VitorFranciscoDev/sprinter-api/infrastructure/datastore"
)

type ProductUseCase struct {
	r datastore.ProductRepository
}

// NewProductUseCase instantiates a new Product Use Case
func NewProductUseCase(r datastore.ProductRepository) *ProductUseCase {
	return &ProductUseCase{r: r}
}

func (p *ProductUseCase) Add(ctx context.Context, product *entities.Product) error {
	return p.r.Add(ctx, product)
}

func (p *ProductUseCase) Delete(ctx context.Context, id int64) error {
	return p.r.Delete(ctx, id)
}

func (p *ProductUseCase) Update(ctx context.Context, product *entities.Product) error {
	return p.r.Update(ctx, product)
}

func (p *ProductUseCase) Get(ctx context.Context, id int64) (*entities.Product, error) {
	return p.r.Get(ctx, id)
}

func (p *ProductUseCase) GetAll(ctx context.Context) ([]*entities.Product, error) {
	return p.r.GetAll(ctx)
}
