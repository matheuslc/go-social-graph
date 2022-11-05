package service

//go:generate mockgen -source=./unfollow_service.go -destination=../mock/service/unfollow_service.go

import (
	"gosocialgraph/pkg/repository"

	"github.com/google/uuid"
)

// UnfollowService defines the struct holder for the use case
type UnfollowService struct {
	Repository repository.Unfollower
}

// UnfolowRunner
type UnfolowRunner interface {
	Run(to, from uuid.UUID) error
}

// Run executes the use case
func (sv UnfollowService) Run(to, from uuid.UUID) error {
	return sv.Repository.Unfollow(to.String(), from.String())
}
