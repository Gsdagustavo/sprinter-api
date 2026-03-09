package entities

import "time"

type ActivityType int64

const (

	// Cycling represents the ActivityType value for cycling activities.
	Cycling ActivityType = 1

	// Running represents the ActivityType value for running activities.
	Running ActivityType = 2

	Walking ActivityType = 3
)

type Activity struct {
	// ID is the unique identifier of the activity
	ID int64 `json:"id,omitempty"`

	// UUID is another identifier for the activity
	UUID string `json:"uuid,omitempty"`

	// UserID is the relation between the activity and the user
	UserID int64 `json:"user_id,omitempty"`

	// Type is the activity that the user has done
	Type ActivityType `json:"type,omitempty"`

	// Route is the track of the activity
	Route []Point `json:"route,omitempty"`

	// StartTime is the time that the activity has started
	StartTime time.Time `json:"start_time,omitempty"`

	// EndTime is the time that the activity has ended
	EndTime time.Time `json:"end_time,omitempty"`

	// CreatedAt is the time that the field was created at the database
	CreatedAt time.Time `json:"created_at,omitempty"`

	// ModifiedAt is the time that the field has been modified in the database
	ModifiedAt time.Time `json:"modified_at,omitempty"`
}

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
