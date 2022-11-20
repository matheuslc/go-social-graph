package service_test

import (
	"context"
	"errors"
	mock "gosocialgraph/pkg/mock/repository"
	"gosocialgraph/pkg/service"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
)

func TestUnfollowRun(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()

	to := uuid.New()
	from := uuid.New()

	ctx := context.Background()
	repo := mock.NewMockUnfollower(controller)
	repo.EXPECT().Unfollow(ctx, to.String(), from.String()).Return(nil)

	service := service.UnfollowService{
		Repository: repo,
	}

	err := service.Run(ctx, to, from)
	if err != nil {
		t.Errorf("Could not unfollow an user. Error %s", err.Error())
	}
}

func TestUnfollowFailRun(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()

	to := uuid.New()
	from := uuid.New()
	ctx := context.Background()
	repo := mock.NewMockUnfollower(controller)
	repo.EXPECT().Unfollow(ctx, to.String(), from.String()).Return(errors.New("Could not execute repository unfollow"))

	service := service.UnfollowService{Repository: repo}

	err := service.Run(ctx, to, from)
	if err == nil {
		t.Errorf("Expect an error when trying to unfollow a user. The repository should error.")
	}
}
