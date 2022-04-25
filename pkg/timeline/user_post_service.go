package timeline

import (
	"gosocialgraph/pkg/post"

	"github.com/google/uuid"
)

type UserPostService struct {
	Repository Reader
}

type UserPostsIntent struct {
	UserId uuid.UUID
}

type UserPostResponse struct {
	Posts []post.UserPost `json:"posts"`
}

func (sv UserPostService) Run(intent UserPostsIntent) (UserPostResponse, error) {
	response, err := sv.Repository.UserPosts(intent.UserId)

	if err != nil {
		return UserPostResponse{}, err
	}

	return UserPostResponse{
		Posts: response,
	}, nil
}
