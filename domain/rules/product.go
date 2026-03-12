package rules

import (
	"github.com/Gsdagustavo/sprinter-api/domain/entities"
	"github.com/Gsdagustavo/sprinter-api/domain/entities/derr"
)

func ValidateProduct(product *entities.Product) error {
	if len(product.Name) < 3 {
		return derr.InvalidProductName
	}

	if len(product.Description) < 10 {
		return derr.InvalidProductDescription
	}

	if product.Price < 0 {
		return derr.InvalidProductPrice
	}

	if product.Stock < 0 {
		return derr.InvalidProductStock
	}

	return nil
}
