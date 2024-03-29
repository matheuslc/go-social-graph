package service

//go:generate mockgen -source=./user_post_service.go -destination=../mock/service/user_post_service.go

import (
	"context"
	"gosocialgraph/pkg/entity"
	"gosocialgraph/pkg/repository"

	"github.com/google/uuid"
)

type UserPostService struct {
	Repository repository.TimelineReader
}

// UserPostRunner
type UserPostRunner interface {
	Run(ctx context.Context, userID uuid.UUID) (posts []entity.UserPost, err error)
}

func (sv UserPostService) Run(ctx context.Context, userID uuid.UUID) (posts []entity.UserPost, err error) {
	posts, err = sv.Repository.UserPosts(ctx, userID)

	if err != nil {
		return posts, err
	}

	return posts, nil
}
