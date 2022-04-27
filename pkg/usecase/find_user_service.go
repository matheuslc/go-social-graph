package usecase

import (
	"gosocialgraph/pkg/user"

	"github.com/google/uuid"
)

// FindUserIntent defines what you need to execute the usecase
type FindUserIntent struct {
	UserID uuid.UUID `json:"user_id"`
}

// FindUserResponse defines the usecase response.
type FindUserResponse struct {
	user.User `json:"user"`
}

// FindUserService defines the service struct and its dependencies
type FindUserService struct {
	UserRepository user.Reader
}

// Run executes the use case
func (sv FindUserService) Run(intent FindUserIntent) (FindUserResponse, error) {
	user, err := sv.UserRepository.Find(intent.UserID.String())

	if err != nil {
		return FindUserResponse{}, err
	}

	return FindUserResponse{
		User: user,
	}, nil
}
