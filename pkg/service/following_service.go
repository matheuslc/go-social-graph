package service

import (
	"gosocialgraph/pkg/entity"
	"gosocialgraph/pkg/repository"

	"github.com/google/uuid"
)

type FollowingService struct {
	Repository repository.TimelineReader
}

type FollowingIntent struct {
	UserId uuid.UUID
}

type FollowingResponse struct {
	Posts []entity.UserPost `json:"posts"`
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
