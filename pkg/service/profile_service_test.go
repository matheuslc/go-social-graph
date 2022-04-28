package service_test

import (
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
	intent := service.FindUserIntent{UserID: userID}

	findUserService := mock_service.NewMockFindUserRunner(controller)
	findUserService.EXPECT().Run(intent).Return(userFound, nil)

	// Mocking stats service
	statsIntent := service.StatsIntent{UserID: userID}
	statsResponse := service.StatsResponse{
		ProfileStats: service.ProfileStats{Followers: 10, Following: 5, PostsCount: 10},
	}
	statsService := mock_service.NewMockStatsRunner(controller)
	statsService.EXPECT().Run(statsIntent).Return(statsResponse, nil)

	// Mocking user post stats service
	userPostIntent := service.UserPostsIntent{UserID: userID}

	firstPost := entity.Post{ID: uuid.New(), Content: "first fake post", CreatedAt: time.Now()}
	secondPost := entity.Post{ID: uuid.New(), Content: "second fake post", CreatedAt: time.Now()}

	posts := []entity.UserPost{}
	posts = append(posts, entity.UserPost{User: user, Post: firstPost}, entity.UserPost{User: user, Post: secondPost})

	userPostResponse := service.UserPostResponse{Posts: posts}
	userPostService := mock_service.NewMockUserPostRunner(controller)
	userPostService.EXPECT().Run(userPostIntent).Return(userPostResponse, nil)

	sv := service.ProfileService{
		FindUserService: findUserService,
		StatsService:    statsService,
		UserPostService: userPostService,
	}

	profileIntent := service.ProfileIntent{UserID: userID}

	result, err := sv.Run(profileIntent)
	if err != nil {
		t.Errorf("Could not get user profile properly")
	}

	if result.User != user {
		t.Errorf("User different than expected. Expect: %s. Got: %s", user.ID, result.User.ID)
	}

	if result.Stats != statsResponse.ProfileStats {
		t.Errorf("Stats different than expected")
	}

	lengthExpected := 2
	actualLength := len(result.Posts.Posts)
	if actualLength != lengthExpected {
		t.Errorf("User lenght different than the expected. Expect: %d. Got: %d", lengthExpected, actualLength)
	}
}
