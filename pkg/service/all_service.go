package service

//go:generate mockgen -source=./all_service.go -destination=../mock/service/all_service.go

import (
	"context"
	"gosocialgraph/pkg/entity"
	"gosocialgraph/pkg/repository"
)

// AllService defines the service struct holder and its dependencies
type AllService struct {
	Repository repository.TimelineReader
}

// AllPostResponse defiens the usecase response
type AllPostResponse struct {
	Posts []entity.UserPost `json:"posts"`
}

// AllPostRunner defines the contract for the service runner
type AllPostRunner interface {
	Run(ctx context.Context) (AllPostResponse, error)
}

// Run defines how the usecase can be executed
func (sv AllService) Run(ctx context.Context) (AllPostResponse, error) {
	response, err := sv.Repository.All(ctx)

	if err != nil {
		return AllPostResponse{}, err
	}

	return AllPostResponse{Posts: response}, nil
}
