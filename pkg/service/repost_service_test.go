package service_test

import (
	mock_repository "gosocialgraph/pkg/mock/repository"
	"gosocialgraph/pkg/service"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
)

func TestRepostRun(t *testing.T) {
	controller := gomock.NewController(t)

	userID := uuid.New()
	parentPostID := uuid.New()
	quote := "retweet quote"
	intent := service.RepostIntent{UserID: userID, ParentPostID: parentPostID, Quote: quote}

	repo := mock_repository.NewMockReposter(controller)
	repo.EXPECT().Repost(intent.UserID.String(), intent.ParentPostID.String(), intent.Quote).Return(true, nil)

	sv := service.RepostService{Repository: repo}

	result, err := sv.Run(intent)
	if err != nil || result.Status != true {
		t.Errorf("Expected a repost to be successful made. Got %s", err)
	}
}
