package service

//go:generate mockgen -source=./create_user_service.go -destination=../mock/service/create_user_service.go

import (
	"fmt"
	"gosocialgraph/pkg/entity"
	"gosocialgraph/pkg/repository"

	"github.com/google/uuid"
)

// CreateUserService executes the business logic for create a user
type CreateUserService struct {
	UserRepository repository.UserReaderWriter
}

// Run executes everything together
func (sv *CreateUserService) Run(username, password string) (entity.User, error) {
	found, err := sv.UserRepository.FindByUsername(username)
	if err != nil {
		return entity.User{}, err
	}

	fmt.Println("found", found)

	if found.ID != uuid.Nil {
		return entity.User{}, fmt.Errorf("User already exist")
	}

	newUser, err := sv.UserRepository.Create(username, password)
	if err != nil {
		return entity.User{}, err
	}

	return newUser, nil
}
