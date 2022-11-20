package service_test

import (
	"context"
	"errors"
	mock_repository "gosocialgraph/pkg/mock/repository"
	"gosocialgraph/pkg/service"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
)

func TestStatsRun(t *testing.T) {
	controller := gomock.NewController(t)
	repo := mock_repository.NewMockStats(controller)

	userID := uuid.New()
	expectedFollowers := 10
	expectedPosts := 20
	expectedFollowing := 100

	ctx := context.Background()

	repo.EXPECT().CountFollowers(ctx, userID).Return(int64(expectedFollowers), nil)
	repo.EXPECT().CountPosts(ctx, userID).Return(int64(expectedPosts), nil)
	repo.EXPECT().CountFollowing(ctx, userID).Return(int64(expectedFollowing), nil)

	sv := service.StatsService{repo}

	result, err := sv.Run(ctx, userID)
	if err != nil {
		t.Errorf("Expected stats from a user. Got an error %s", err)
	}

	if result.Followers != expectedFollowers {
		t.Errorf("Wrong number of followers. Expected %d. Got %d", result.Followers, expectedFollowers)
	}

	if result.Following != expectedFollowing {
		t.Errorf("Wrong number of following. Expected %d. Got %d", result.Following, expectedFollowing)
	}

	if result.PostsCount != expectedPosts {
		t.Errorf("Wrong number of posts. Expected %d. Got %d", result.PostsCount, expectedPosts)
	}
}

func TestStatsFailRun(t *testing.T) {
	controller := gomock.NewController(t)
	repo := mock_repository.NewMockStats(controller)

	userID := uuid.New()

	firstError := errors.New("first error")
	secondError := errors.New("second error")
	thirdError := errors.New("third error")
	ctx := context.Background()

	repo.EXPECT().CountFollowers(ctx, userID).Return(int64(0), firstError)
	repo.EXPECT().CountPosts(ctx, userID).Return(int64(0), secondError)
	repo.EXPECT().CountFollowing(ctx, userID).Return(int64(0), thirdError)

	sv := service.StatsService{repo}

	_, err := sv.Run(ctx, userID)
	if err == nil {
		t.Errorf("Expected at least an error. Service executed successfuly")
	}
}
