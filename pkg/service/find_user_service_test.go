package service_test

import (
	"context"
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

	ctx := context.Background()
	repo := mock.NewMockUserReader(controller)
	repo.EXPECT().Find(ctx, userID).Return(expectedUserFound, nil)

	service := service.FindUserService{UserRepository: repo}

	_, err := service.Run(ctx, userID)
	if err != nil {
		t.Errorf("User not found. Error %s", err)
	}
}

func TestFindUserFailRun(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()

	userID := uuid.New()
	empty := entity.User{}

	ctx := context.Background()
	repo := mock.NewMockUserReader(controller)
	repo.EXPECT().Find(ctx, userID).Return(empty, errors.New("Fake error"))

	service := service.FindUserService{UserRepository: repo}

	_, err := service.Run(ctx, userID)
	if err == nil {
		t.Errorf("Expected an error when trying to create an user. Error %s", err)
	}
}
