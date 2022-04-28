package service

import (
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

// Run defines how the usecase can be executed
func (sv AllService) Run() (AllPostResponse, error) {
	response, err := sv.Repository.All()

	if err != nil {
		return AllPostResponse{}, err
	}

	return AllPostResponse{Posts: response}, nil
}
