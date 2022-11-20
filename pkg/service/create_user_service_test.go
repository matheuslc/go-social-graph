package service_test

import (
	"gosocialgraph/pkg/entity"
	mock "gosocialgraph/pkg/mock/repository"
	"gosocialgraph/pkg/service"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
)

func TestCreateUserRun(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()

	username := "devtest"
	email := "foo@gmail.com"
	password := "foo"

	repo := mock.NewMockUserReaderWriter(controller)
	repo.EXPECT().FindByUsername(username).Return(entity.User{}, nil)
	repo.EXPECT().Create(username, email, password).Return(entity.User{
		ID:        uuid.New(),
		Username:  username,
		Email:     email,
		CreatedAt: time.Now(),
	}, nil)

	sv := service.CreateUserService{UserRepository: repo}

	result, err := sv.Run(username, email, password)
	if err != nil || len(result.Username) == 0 || result.Username != username {
		t.Errorf("Got an error when trying to create a new user. Error: %s", err)
	}
}

func TestCreateUserAlreadyExistRun(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()

	username := "devtest"
	email := "foo@gmail.com"
	password := "foo"
	repo := mock.NewMockUserReaderWriter(controller)
	repo.EXPECT().FindByUsername(username).Return(entity.User{
		ID:       uuid.New(),
		Username: username,
	}, nil)

	sv := service.CreateUserService{UserRepository: repo}

	_, err := sv.Run(username, email, password)
	if err == nil {
		t.Errorf("Expect an error when trying to create a new user that already exist. Error: %s", err)
	}
}
