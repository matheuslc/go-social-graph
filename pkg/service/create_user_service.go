package service

//go:generate mockgen -source=./create_user_service.go -destination=../mock/service/create_user_service.go

import (
	"gosocialgraph/pkg/entity"
	"gosocialgraph/pkg/repository"
)

// CreateUserService executes the business logic for create a user
type CreateUserService struct {
	UserRepository repository.UserReaderWriter
}

// Run executes everything together
func (sv *CreateUserService) Run(username string) (entity.User, error) {
	_, err := sv.UserRepository.FindByUsername(username)
	if err != nil {
		return entity.User{}, err
	}

	newUser, err := sv.UserRepository.Create(username)
	if err != nil {
		return entity.User{}, err
	}

	return newUser, nil
}
