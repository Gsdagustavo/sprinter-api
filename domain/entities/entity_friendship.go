package entities

import "time"

// FriendshipRequestStatus represents the current state of a friendship request
type FriendshipRequestStatus string

const (
	// Accepted represents a friendship request that was accepted by the recipient
	Accepted FriendshipRequestStatus = "ACCEPTED"

	// Pending represents a friendship request that is waiting for the recipient's answer
	Pending FriendshipRequestStatus = "PENDING"

	// Denied represents a friendship request that was denied by the recipient
	Denied FriendshipRequestStatus = "DENIED"
)

// Friendship represents an established friendship between two users
type Friendship struct {
	// ID is the unique identifier of the friendship
	ID int64 `json:"id,omitempty"`

	// UUID is another identifier for the friendship
	UUID string `json:"uuid,omitempty"`

	// FirstUserID is the ID for one of the users of the friendship
	FirstUserID int64 `json:"first_user_id,omitempty"`

	// FirstUserID is the other ID for the users of the friendship
	SecondUserID int64 `json:"second_user_id,omitempty"`

	// CreatedAt is the time that the field was created at the database
	CreatedAt time.Time `json:"created_at,omitempty"`

	// ModifiedAt is the time that the field has been modified in the database
	ModifiedAt time.Time `json:"modified_at,omitempty"`
}

// FriendshipRequest represents a friendship request sent from one user to another
type FriendshipRequest struct {
	// SenderID is the ID of the user that send the friendship request
	SenderID int64 `json:"sender_id,omitempty"`

	// RecipientID is the ID of the user that received the friendship request
	RecipientID int64 `json:"recipient_id,omitempty"`

	// Status is the actual situation of the request [pending, accepted, denied]
	Status FriendshipRequestStatus `json:"status,omitempty"`

	// CreatedAt is the time that the field was created at the database
	CreatedAt time.Time `json:"created_at,omitempty"`

	// ModifiedAt is the time that the field has been modified in the database
	ModifiedAt time.Time `json:"modified_at,omitempty"`
}
