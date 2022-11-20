package service_test

import (
	"context"
	"gosocialgraph/pkg/entity"
	u "gosocialgraph/pkg/entity"
	mock_service "gosocialgraph/pkg/mock/service"
	"gosocialgraph/pkg/service"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
)

func TestProfileRun(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()

	userID := uuid.New()

	// Mocking find user service
	user := u.User{ID: userID}
	userFound := service.FindUserResponse{User: user}

	ctx := context.Background()
	findUserService := mock_service.NewMockFindUserRunner(controller)
	findUserService.EXPECT().Run(ctx, userID).Return(userFound, nil)

	// Mocking stats service
	statsResponse := service.StatsResponse{
		ProfileStats: service.ProfileStats{Followers: 10, Following: 5, PostsCount: 10},
	}
	statsService := mock_service.NewMockStatsRunner(controller)
	statsService.EXPECT().Run(ctx, userID).Return(statsResponse, nil)

	// Mocking user post stats service
	firstPost := entity.Post{ID: uuid.New(), Content: "first fake post", CreatedAt: time.Now()}
	secondPost := entity.Post{ID: uuid.New(), Content: "second fake post", CreatedAt: time.Now()}

	posts := []entity.UserPost{}
	posts = append(posts, entity.UserPost{User: user, Post: firstPost}, entity.UserPost{User: user, Post: secondPost})

	userPostService := mock_service.NewMockUserPostRunner(controller)
	userPostService.EXPECT().Run(ctx, userID).Return(posts, nil)

	sv := service.ProfileService{
		FindUserService: findUserService,
		StatsService:    statsService,
		UserPostService: userPostService,
	}

	result, err := sv.Run(ctx, userID)
	if err != nil {
		t.Errorf("Could not get user profile properly")
	}

	if result.User.ID != user.ID {
		t.Errorf("User different than expected. Expect: %s. Got: %s", user.ID, result.User.ID)
	}

	if result.Stats != statsResponse.ProfileStats {
		t.Errorf("Stats different than expected")
	}

	lengthExpected := 2
	actualLength := len(result.Posts)
	if actualLength != lengthExpected {
		t.Errorf("User lenght different than the expected. Expect: %d. Got: %d", lengthExpected, actualLength)
	}
}
