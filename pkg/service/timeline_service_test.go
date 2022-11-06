package service_test

import (
	"errors"
	"gosocialgraph/pkg/entity"
	mock_repository "gosocialgraph/pkg/mock/repository"
	"gosocialgraph/pkg/service"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
)

func TestFollowingRun(t *testing.T) {
	controller := gomock.NewController(t)

	userID := uuid.New()
	timelineUser := entity.User{ID: userID}
	repo := mock_repository.NewMockTimelineReader(controller)
	posts := []entity.UserPost{}
	firstPost := entity.Post{ID: uuid.New(), Content: "foo text", CreatedAt: time.Now()}
	secondPost := entity.Post{ID: uuid.New(), Content: "fake content", CreatedAt: time.Now()}
	userPost := entity.UserPost{User: timelineUser, Post: firstPost}
	userSecondPost := entity.UserPost{User: timelineUser, Post: secondPost}

	posts = append(posts, userPost, userSecondPost)

	repo.EXPECT().TimelineFor(userID).Return(posts, nil)

	sv := service.TimelineServive{Repository: repo}

	result, err := sv.Run(userID)
	if err != nil {
		t.Errorf("Expected a timeline for user. Got error")
	}

	expectLength := len(posts)
	postLength := len(*result.Posts)
	if postLength != expectLength {
		t.Errorf("Expect %d, got %d ", expectLength, postLength)
	}
}

func TestFollowingFailRun(t *testing.T) {
	controller := gomock.NewController(t)

	userID := uuid.New()
	timelineUser := entity.User{ID: userID}
	repo := mock_repository.NewMockTimelineReader(controller)
	posts := []entity.UserPost{}
	firstPost := entity.Post{ID: uuid.New(), Content: "foo text", CreatedAt: time.Now()}
	userPost := entity.UserPost{User: timelineUser, Post: firstPost}

	posts = append(posts, userPost)

	repo.EXPECT().TimelineFor(userID).Return(posts, errors.New("Fake error"))

	sv := service.TimelineServive{Repository: repo}

	_, err := sv.Run(userID)
	if err == nil {
		t.Errorf("Expected an error for user time")
	}
}
