package service_test

import (
	"errors"
	"gosocialgraph/pkg/entity"
	mock "gosocialgraph/pkg/mock/repository"
	"gosocialgraph/pkg/service"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
)

func TestFindUserRun(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()

	userID := uuid.New()
	expectedUserFound := entity.User{ID: userID}

	repo := mock.NewMockUserReader(controller)
	repo.EXPECT().Find(userID.String()).Return(expectedUserFound, nil)

	intent := service.FindUserIntent{UserID: userID}
	service := service.FindUserService{UserRepository: repo}

	_, err := service.Run(intent)
	if err != nil {
		t.Errorf("User not found. Error %s", err)
	}
}

func TestFindUserFailRun(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()

	userID := uuid.New()
	empty := entity.User{}

	repo := mock.NewMockUserReader(controller)
	repo.EXPECT().Find(userID.String()).Return(empty, errors.New("Fake error"))

	intent := service.FindUserIntent{
		UserID: userID,
	}
	service := service.FindUserService{
		UserRepository: repo,
	}

	_, err := service.Run(intent)
	if err == nil {
		t.Errorf("Expected an error when trying to create an user. Error %s", err)
	}
}
