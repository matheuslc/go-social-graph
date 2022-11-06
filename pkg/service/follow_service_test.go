package service_test

import (
	"errors"
	mock "gosocialgraph/pkg/mock/repository"
	"gosocialgraph/pkg/service"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
)

func TestFollowRun(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()

	to := uuid.New()
	from := uuid.New()

	repo := mock.NewMockFollower(controller)

	repo.EXPECT().Follow(to.String(), from.String()).Return(nil)
	sv := service.FollowService{
		UserRepository: repo,
	}

	err := sv.Run(to, from)
	if err != nil {
		t.Errorf("Got an error when trying to follow a user. Error: %s", err)
	}
}

func TestFollowFailRun(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()

	to := uuid.New()
	from := uuid.New()

	repo := mock.NewMockFollower(controller)

	repo.EXPECT().Follow(to.String(), from.String()).Return(errors.New("Fake error"))
	sv := service.FollowService{
		UserRepository: repo,
	}

	err := sv.Run(to, from)
	if err == nil {
		t.Errorf("Expected an error when trying to follow a user. Error: %s", err)
	}
}
