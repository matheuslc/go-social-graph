package timeline

import (
	"gosocialgraph/pkg/post"

	"github.com/google/uuid"
)

type FollowingService struct {
	Repository Reader
}

type FollowingIntent struct {
	UserId uuid.UUID
}

type FollowingResponse struct {
	Posts []post.UserPost `json:"posts"`
}

func (sv FollowingService) Run(intent FollowingIntent) (FollowingResponse, error) {
	response, err := sv.Repository.TimelineFor(intent.UserId)

	if err != nil {
		return FollowingResponse{}, err
	}

	return FollowingResponse{
		Posts: response,
	}, nil
}
