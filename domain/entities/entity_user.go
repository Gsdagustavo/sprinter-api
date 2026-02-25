package entities

import "time"

type User struct {
	// ID is the unique identifier of the user
	ID int64 `json:"id,omitempty"`

	// UUID is another identifier for the user
	UUID string `json:"uuid,omitempty"`

	// Name is the name of the user like "John Doe"
	Name string `json:"name,omitempty"`

	// Email is the email of the user
	Email string `json:"email,omitempty"`

	// Password is the password of the user (always saved in hash)
	Password string `json:"password,omitempty"`

	// CarboCoins is the official trade coin for the app (trade for products)
	CarboCoins int64 `json:"carbo_coins,omitempty"`

	// Carbon is the count of how many carbon the user has not emitted by doing activities
	Carbon float64 `json:"carbon,omitempty"`

	// TraveledDistance is the distance done by the user with all activities
	TraveledDistance float64 `json:"traveled_distance,omitempty"`

	// CreatedAt is the time that the field was created at the database
	CreatedAt time.Time `json:"created_at,omitempty"`

	// ModifiedAt is the time that the field has been modified in the database
	ModifiedAt time.Time `json:"modified_at,omitempty"`
}

type UserCredentials struct {
	Name     string `json:"name,omitempty"`
	Email    string `json:"email,omitempty"`
	Password string `json:"password,omitempty"`
}

type AuthenticationResponse struct {
	Token string `json:"token,omitempty"`
}
