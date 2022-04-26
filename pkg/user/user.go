package user

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

// UsernameMaxLength defines the maximum characters an username may have
// It can be extracted to an environment variable
const UsernameMaxLength = 14

// User defines the fields that makes something looks like a user
type User struct {
	ID        uuid.UUID `json:"id"`
	Username  string    `json:"username"`
	CreatedAt time.Time `json:"created_at"`
}

// NewUser creates a new valid user
func NewUser(username string) (User, error) {
	if len(username) > UsernameMaxLength {
		return User{}, errors.New("Username is bigger than allowed. Maximum lenght: %")
	}

	return User{Username: username}, nil
}
