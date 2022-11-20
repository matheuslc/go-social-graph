package service

//go:generate mockgen -source=./create_user_service.go -destination=../mock/service/create_user_service.go

import (
	"context"
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
func (sv *CreateUserService) Run(ctx context.Context, username, email, password string) (entity.User, error) {
	found, err := sv.UserRepository.FindByUsername(ctx, username)
	if err != nil {
		return entity.User{}, err
	}

	if found.ID != uuid.Nil {
		return entity.User{}, fmt.Errorf("User already exist")
	}

	newUser, err := sv.UserRepository.Create(ctx, username, email, password)
	if err != nil {
		return entity.User{}, err
	}

	return newUser, nil
}
