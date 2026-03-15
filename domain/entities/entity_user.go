package entities

import "time"

// CarboCoin is a unitary value that represents the base currency for the application.
type CarboCoin int64

type User struct {
	// ID is the unique identifier of the user.
	ID int64 `json:"id,omitempty"`

	// Name is the name of the user.
	Name string `json:"name,omitempty"`

	// Email is the email of the user.
	Email string `json:"email,omitempty"`

	// Password is the password hash of the user.
	Password string `json:"password,omitempty"`

	// CarboCoins is the amount of carbo coins the user has.
	CarboCoins CarboCoin `json:"carbo_coins,omitempty"`

	// Carbon is the count of how many carbon the user has not emitted by doing activities.
	Carbon float64 `json:"carbon,omitempty"`

	// TraveledDistance is the distance done by the user with all activities.
	TraveledDistance float64 `json:"traveled_distance,omitempty"`

	// CreatedAt is the date of creation in database.
	CreatedAt time.Time `json:"created_at,omitempty"`

	// CreatedAt is the last date of update in database.
	ModifiedAt time.Time `json:"modified_at,omitempty"`
}

// UserCredentials is an entity that represents a user credentials .
type UserCredentials struct {
	// Name in the credential entity.
	Name string `json:"name,omitempty"`

	// Email in the credential entity.
	Email string `json:"email,omitempty"`

	// Password in the credential entity.
	Password string `json:"password,omitempty"`
}

// AuthenticationResponse is a response from an authentication attempt.
type AuthenticationResponse struct {
	// Token is the authentication token returned in the authentication response.
	Token string `json:"token,omitempty"`
}
