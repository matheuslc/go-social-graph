package service

//go:generate mockgen -source=./find_user_service.go -destination=../mock/service/find_user_service.go

import (
	"gosocialgraph/pkg/entity"
	"gosocialgraph/pkg/repository"

	"github.com/google/uuid"
)

// FindUserResponse defines the usecase response.
type FindUserResponse struct {
	entity.User `json:"user"`
}

// FindUserService defines the service struct and its dependencies
type FindUserService struct {
	UserRepository repository.UserReader
}

// FindUserRunner
type FindUserRunner interface {
	Run(userID uuid.UUID) (FindUserResponse, error)
}

// Run executes the use case
func (sv FindUserService) Run(userID uuid.UUID) (FindUserResponse, error) {
	user, err := sv.UserRepository.Find(userID)

	if err != nil {
		return FindUserResponse{}, err
	}

	return FindUserResponse{User: user}, nil
}
