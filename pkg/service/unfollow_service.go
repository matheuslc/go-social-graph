package service

//go:generate mockgen -source=./unfollow_service.go -destination=../mock/service/unfollow_service.go

import (
	"context"
	"gosocialgraph/pkg/repository"

	"github.com/google/uuid"
)

// UnfollowService defines the struct holder for the use case
type UnfollowService struct {
	Repository repository.Unfollower
}

// UnfolowRunner
type UnfolowRunner interface {
	Run(ctx context.Context, to, from uuid.UUID) error
}

// Run executes the use case
func (sv UnfollowService) Run(ctx context.Context, to, from uuid.UUID) error {
	return sv.Repository.Unfollow(ctx, to.String(), from.String())
}
