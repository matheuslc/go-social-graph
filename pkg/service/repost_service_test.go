package service_test

import (
	"context"
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

	ctx := context.Background()
	repo := mock_repository.NewMockReposter(controller)
	repo.EXPECT().Repost(ctx, userID.String(), parentPostID.String(), quote).Return(nil)

	sv := service.RepostService{Repository: repo}

	err := sv.Run(ctx, userID, parentPostID, quote)
	if err != nil {
		t.Errorf("Expected a repost to be successful made. Got %s", err)
	}
}
