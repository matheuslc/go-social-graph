package usecase

import (
	"errors"
	"gosocialgraph/pkg/user"
)

// CreateUserIntent defines the fields necessary to create a new user
type CreateUserIntent struct {
	Username string `json:"username"`
}

// CreateUserService executes the business logic for create a user
type CreateUserService struct {
	UserRepository user.ReaderWriter
}

// Run executes everything together
func (sv *CreateUserService) Run(intent CreateUserIntent) (user.User, error) {
	persistedUser, err := sv.UserRepository.FindByUsername(intent.Username)

	if err != nil {
		return user.User{}, err
	}

	if persistedUser {
		return user.User{}, errors.New("User already exists")
	}

	newUser, err := sv.UserRepository.Create(intent.Username)
	if err != nil {
		return user.User{}, err
	}

	return newUser, nil
}
