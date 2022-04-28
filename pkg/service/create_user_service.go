package service

import (
	"errors"
	"gosocialgraph/pkg/entity"
	"gosocialgraph/pkg/repository"
)

// CreateUserIntent defines the fields necessary to create a new user
type CreateUserIntent struct {
	Username string `json:"username"`
}

// CreateUserService executes the business logic for create a user
type CreateUserService struct {
	UserRepository repository.UserReaderWriter
}

// Run executes everything together
func (sv *CreateUserService) Run(intent CreateUserIntent) (entity.User, error) {
	persistedUser, err := sv.UserRepository.FindByUsername(intent.Username)

	if err != nil {
		return entity.User{}, err
	}

	if persistedUser {
		return entity.User{}, errors.New("User already exists")
	}

	newUser, err := sv.UserRepository.Create(intent.Username)
	if err != nil {
		return entity.User{}, err
	}

	return newUser, nil
}
