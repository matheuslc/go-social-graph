package service

//go:generate mockgen -source=./follow_service.go -destination=../mock/service/follow_service.go

import (
	"gosocialgraph/pkg/repository"

	"github.com/google/uuid"
)

// FollowService defines the follow use case.
type FollowService struct {
	UserRepository repository.Follower
}

// FollowRunner defines the contract to run the follow usecase
type FollowRunner interface {
	Run(to, from uuid.UUID) (bool, error)
}

// Run execute the use case
func (sv *FollowService) Run(to, from uuid.UUID) (bool, error) {
	_, err := sv.UserRepository.Follow(to.String(), from.String())

	if err != nil {
		return false, err
	}

	return true, nil
}
