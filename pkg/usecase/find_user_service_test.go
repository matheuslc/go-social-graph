package usecase

import (
	"errors"
	mock "gosocialgraph/pkg/mock/user"
	"gosocialgraph/pkg/user"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
)

func TestFindUserRun(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()

	userID := uuid.New()
	expectedUserFound := user.User{
		ID: userID,
	}

	repo := mock.NewMockReader(controller)
	repo.EXPECT().Find(userID.String()).Return(expectedUserFound, nil)

	intent := FindUserIntent{
		UserID: userID,
	}
	service := FindUserService{
		UserRepository: repo,
	}

	_, err := service.Run(intent)
	if err != nil {
		t.Errorf("User not found. Error %s", err)
	}
}

func TestFindUserFailRun(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()

	userID := uuid.New()
	empty := user.User{}

	repo := mock.NewMockReader(controller)
	repo.EXPECT().Find(userID.String()).Return(empty, errors.New("Fake error"))

	intent := FindUserIntent{
		UserID: userID,
	}
	service := FindUserService{
		UserRepository: repo,
	}

	_, err := service.Run(intent)
	if err == nil {
		t.Errorf("Expected an error when trying to create an user. Error %s", err)
	}
}
