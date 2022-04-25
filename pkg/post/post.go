package post

import (
	"gosocialgraph/pkg/user"
	"time"

	"github.com/google/uuid"
)

type Post struct {
	Id        uuid.UUID `json:"id"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`
}

type Repost struct {
	Id        uuid.UUID `json:"id"`
	Parent    UserPost  `json:"parent"`
	Quote     string    `json:"quote"`
	CreatedAt time.Time `json:"created_at"`
}

type UserPost struct {
	User user.User   `json:"user"`
	Post interface{} `json:"post"`
}
