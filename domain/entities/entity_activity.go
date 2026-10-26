package entities

import "time"

// ActivityType represents the kind of physical activity done by the user
type ActivityType string

const (
	// Cycling represents an activity done by bicycle
	Cycling ActivityType = "CYCLING"

	// Running represents an activity done by running
	Running ActivityType = "RUNNING"

	// Walking represents an activity done by walking
	Walking ActivityType = "WALKING"
)

// Activity represents a physical activity done by a user
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
