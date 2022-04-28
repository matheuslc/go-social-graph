package service

import (
	"gosocialgraph/pkg/entity"
	"gosocialgraph/pkg/repository"

	"github.com/google/uuid"
)

// FindUserIntent defines what you need to execute the usecase
type FindUserIntent struct {
	UserID uuid.UUID `json:"user_id"`
}

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
	Run(intent FindUserIntent) (FindUserResponse, error)
}

// Run executes the use case
func (sv FindUserService) Run(intent FindUserIntent) (FindUserResponse, error) {
	user, err := sv.UserRepository.Find(intent.UserID.String())

	if err != nil {
		return FindUserResponse{}, err
	}

	return FindUserResponse{User: user}, nil
}
