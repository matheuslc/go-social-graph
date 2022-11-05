package service

//go:generate mockgen -source=./following_service.go -destination=../mock/service/following_service.go

import (
	"gosocialgraph/pkg/entity"
	"gosocialgraph/pkg/repository"

	"github.com/google/uuid"
)

// FollowingService defines the struct holder for the service and its dependencies
type FollowingService struct {
	Repository repository.TimelineReader
}

// FollowingIntent defines the contract for requisting a usecase execution
type FollowingIntent struct {
	UserID uuid.UUID
}

// FollowingResponse defines the response from the usecase
type FollowingResponse struct {
	Posts []entity.UserPost `json:"posts"`
}

// Run executes the use case
func (sv FollowingService) Run(intent FollowingIntent) (FollowingResponse, error) {
	response, err := sv.Repository.TimelineFor(intent.UserID)
	if err != nil {
		return FollowingResponse{}, err
	}

	return FollowingResponse{Posts: response}, nil
}
