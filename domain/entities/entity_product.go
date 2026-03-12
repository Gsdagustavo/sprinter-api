package entities

import "time"

// Product represents a product entity.
type Product struct {
	// ID is the unique identifier of the product.
	ID int64 `json:"id,omitempty"`

	// Name is the name of the product.
	Name string `json:"name,omitempty"`

	// Description provides additional details or an overview about the product.
	Description string `json:"description,omitempty"`

	// ImageURL is the URL or location of the product's image.
	ImageURL string `json:"image_url,omitempty"`

	// Price is the price of the product in CarboCoins.
	Price int64 `json:"price,omitempty"`

	// Stock is the available stock of the product.
	Stock int64 `json:"stock,omitempty"`

	// CreatedAt is the date of creation in database.
	CreatedAt time.Time `json:"created_at,omitempty"`

	// CreatedAt is the last date of update in database.
	ModifiedAt time.Time `json:"modified_at,omitempty"`
}

// CartItem represents an item in a shopping cart.
type CartItem struct {
	// ID is the unique identifier of the product in cart
	ID int64 `json:"id,omitempty"`

	// UserID is the identifier of the user that this cart item belongs to.
	UserID int64

	// Product is the product in the cart.
	Product Product `json:"product"`

	// Quantity is the amount of that product that it is in the cart
	Quantity int64 `json:"quantity,omitempty"`

	// CreatedAt is the date of creation in database.
	CreatedAt time.Time `json:"created_at,omitempty"`

	// CreatedAt is the last date of update in database.
	ModifiedAt time.Time `json:"modified_at,omitempty"`
}
