package service

//go:generate mockgen -source=./follow_service.go -destination=../mock/service/follow_service.go

import (
	"context"
	"gosocialgraph/pkg/repository"

	"github.com/google/uuid"
)

// FollowService defines the follow use case.
type FollowService struct {
	UserRepository repository.Follower
}

// FollowRunner defines the contract to run the follow usecase
type FollowRunner interface {
	Run(to, from uuid.UUID) error
}

// Run execute the use case
func (sv *FollowService) Run(ctx context.Context, to, from uuid.UUID) error {
	if err := sv.UserRepository.Follow(ctx, to.String(), from.String()); err != nil {
		return err
	}

	return nil
}
