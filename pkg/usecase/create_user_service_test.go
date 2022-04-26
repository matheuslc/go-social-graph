package usecase

import (
	mock "gosocialgraph/pkg/mock/user"
	"gosocialgraph/pkg/user"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
)

func TestCreateUserRun(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()

	username := "devtest"

	repo := mock.NewMockReaderWriter(controller)

	repo.EXPECT().FindByUsername(username).Return(false, nil)
	repo.EXPECT().Create(username).Return(user.User{
		ID:        uuid.New(),
		Username:  username,
		CreatedAt: time.Now(),
	}, nil)
	service := CreateUserService{
		UserRepository: repo,
	}

	intent := CreateUserIntent{Username: username}

	result, err := service.Run(intent)
	if err != nil || len(result.Username) == 0 || result.Username != username {
		t.Errorf("Got an error when trying to create a new user. Error: %s", err)
	}
}

func TestCreateUserAlreadyExistRun(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()

	username := "devtest"

	repo := mock.NewMockReaderWriter(controller)

	repo.EXPECT().FindByUsername(username).Return(true, nil)
	service := CreateUserService{
		UserRepository: repo,
	}

	intent := CreateUserIntent{Username: username}

	_, err := service.Run(intent)
	if err == nil {
		t.Errorf("Expect an error when trying to create a new user that already exist. Error: %s", err)
	}
}
