package entities

import "time"

type Point struct {
	// ID is the unique identifier of the point
	ID int64 `json:"id,omitempty"`

	// UUID is another identifier for the point
	UUID string `json:"uuid,omitempty"`

	// ActivityID is the relation between the point and the activity
	ActivityID int64 `json:"activity_id,omitempty"`

	// Latitude is the latitude of the actual point
	Latitude float64 `json:"latitude,omitempty"`

	// Longitude is the longitude of the actual point
	Longitude float64 `json:"longitude,omitempty"`

	// CreatedAt is the time that the field was created at the database
	CreatedAt time.Time `json:"created_at,omitempty"`

	// ModifiedAt is the time that the field has been modified in the database
	ModifiedAt time.Time `json:"modified_at,omitempty"`
}
