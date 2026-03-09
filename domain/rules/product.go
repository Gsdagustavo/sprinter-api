package rules

import "github.com/Gsdagustavo/sprinter-api/domain/entities"

func ValidateProduct(product *entities.Product) bool {
	if len(product.Name) < 3 {
		return false
	}

	if len(product.Description) < 10 {
		return false
	}

	if product.Price < 0 {
		return false
	}

	if product.Quantity < 0 {
		return false
	}

	if product.Discount < 0 || product.Discount > 100 {
		return false
	}

	return true
}
