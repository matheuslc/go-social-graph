package service

//go:generate mockgen -source=./user_post_service.go -destination=../mock/service/user_post_service.go

import (
	"gosocialgraph/pkg/entity"
	"gosocialgraph/pkg/repository"

	"github.com/google/uuid"
)

type UserPostService struct {
	Repository repository.TimelineReader
}

type UserPostResponse struct {
	Posts []entity.UserPost `json:"posts"`
}

// UserPostRunner
type UserPostRunner interface {
	Run(userID uuid.UUID) (UserPostResponse, error)
}

func (sv UserPostService) Run(userID uuid.UUID) (UserPostResponse, error) {
	response, err := sv.Repository.UserPosts(userID)

	if err != nil {
		return UserPostResponse{}, err
	}

	return UserPostResponse{
		Posts: response,
	}, nil
}
