package entity

import (
	"time"

	"github.com/google/uuid"
)

type Post struct {
	ID        uuid.UUID `json:"id"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`
	Parent    *UserPost `json:"parent"`
}

type UserPost struct {
	User User `json:"user"`
	Post Post `json:"post"`
}
