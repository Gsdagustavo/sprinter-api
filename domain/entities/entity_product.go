package entities

import "time"

// Product represents a product entity with details such as ID, name, price, and more.
type Product struct {
	// ID is the unique identifier of the product
	ID int64 `json:"id,omitempty"`

	// UUID is another identifier for the product
	UUID string `json:"uuid,omitempty"`

	// Name is the name of the product like "Ticket Superbowl"
	Name string `json:"name,omitempty"`

	// Description provides additional details or an overview about the product.
	Description string `json:"description,omitempty"`

	// ImageURL is the URL or location of the product's image.
	ImageURL string `json:"image_url,omitempty"`

	// Price is the price of the product in CarboCoins
	Price int64 `json:"price,omitempty"`

	// Stock is the stock of the product
	Stock int64 `json:"stock,omitempty"`

	// Discount is the percentage of discount that the product is on
	Discount int64 `json:"discount,omitempty"`

	// CreatedAt is the time that the field was created at the database
	CreatedAt time.Time `json:"created_at,omitempty"`

	// ModifiedAt is the time that the field has been modified in the database
	ModifiedAt time.Time `json:"modified_at,omitempty"`
}

// ProductCart represents an item in the shopping cart with its details and metadata.
type ProductCart struct {
	// ID is the unique identifier of the product in cart
	ID int64 `json:"id,omitempty"`

	// UUID is another identifier for the product in cart
	UUID string `json:"uuid,omitempty"`

	// ProductID is the relation between the product in cart and the product
	ProductID int64 `json:"product_id,omitempty"`

	// Quantity is the amount of that product that it is in the cart
	Quantity int64 `json:"quantity,omitempty"`

	// CreatedAt is the time that the field was created at the database
	CreatedAt time.Time `json:"created_at,omitempty"`

	// ModifiedAt is the time that the field has been modified in the database
	ModifiedAt time.Time `json:"modified_at,omitempty"`
}
