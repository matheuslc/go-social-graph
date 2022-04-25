package user

import (
	"time"

	"github.com/google/uuid"
)

const USERNAME_MAX_LENGTH = 14

type User struct {
	Id        uuid.UUID `json:"id"`
	Username  string    `json:"username"`
	CreatedAt time.Time `json:"created_at"`
}

type UserStats struct {
	Followers  int `json:"followers"`
	Following  int `json:"following"`
	PostsCount int `json:"posts_count"`
}

type Follower interface {
	Follow(from User, to User) (bool, error)
}
