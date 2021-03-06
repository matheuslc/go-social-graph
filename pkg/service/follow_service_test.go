package service_test

import (
	"errors"
	mock "gosocialgraph/pkg/mock/repository"
	"gosocialgraph/pkg/service"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
)

func TestFollowIntent(t *testing.T) {
	to := uuid.New()
	from := uuid.New()

	intent, err := service.NewFollowIntent(to, from)

	if err != nil && intent.To == to && intent.From == from {
		t.Errorf("Could not create a valid follow intent")
	}
}

func TestSelfFollowIntent(t *testing.T) {
	to := uuid.New()
	from := to

	_, err := service.NewFollowIntent(to, from)
	if err == nil {
		t.Errorf("Expect a error when trying to create a follow intent with the same id for both")
	}
}

func TestFollowRun(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()

	to := uuid.New()
	from := uuid.New()

	repo := mock.NewMockFollower(controller)

	repo.EXPECT().Follow(to.String(), from.String()).Return(true, nil)
	sv := service.FollowService{
		UserRepository: repo,
	}

	intent := service.FollowIntent{To: to, From: from}

	result, err := sv.Run(intent)
	if err != nil && result != true {
		t.Errorf("Got an error when trying to follow a user. Error: %s", err)
	}
}

func TestFollowFailRun(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()

	to := uuid.New()
	from := uuid.New()

	repo := mock.NewMockFollower(controller)

	repo.EXPECT().Follow(to.String(), from.String()).Return(false, errors.New("Fake error"))
	sv := service.FollowService{
		UserRepository: repo,
	}

	intent := service.FollowIntent{To: to, From: from}

	result, err := sv.Run(intent)
	if err == nil && result == true {
		t.Errorf("Expected an error when trying to follow a user. Error: %s", err)
	}
}
