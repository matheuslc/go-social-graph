package service_test

import (
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

	intent, err := service.NewUnfollowIntent(to, from)
	if err != nil {
		t.Errorf("Could not create the intent to execute the usecase")
	}

	repo := mock.NewMockUnfollower(controller)
	repo.EXPECT().Unfollow(to.String(), from.String()).Return(true, nil)

	service := service.UnfollowService{
		Repository: repo,
	}

	_, err = service.Run(intent)
	if err != nil {
		t.Errorf("Could not unfollow an user. Error %s", err.Error())
	}
}

func TestUnfollowFailRun(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()

	to := uuid.New()
	from := uuid.New()

	repo := mock.NewMockUnfollower(controller)
	repo.EXPECT().Unfollow(to.String(), from.String()).Return(false, errors.New("Could not execute repository unfollow"))

	intent, _ := service.NewUnfollowIntent(to, from)
	service := service.UnfollowService{
		Repository: repo,
	}

	_, err := service.Run(intent)
	if err == nil {
		t.Errorf("Expect an error when trying to unfollow a user. The repository should error.")
	}
}

func TestUnfollowFailIntent(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()

	to := uuid.New()

	_, err := service.NewUnfollowIntent(to, to)
	if err == nil {
		t.Errorf("Expect an error when trying to unfollow yourself")
	}
}
