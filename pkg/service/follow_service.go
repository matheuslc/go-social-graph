package service

//go:generate mockgen -source=./follow_service.go -destination=../mock/service/follow_service.go

import (
	"fmt"
	"gosocialgraph/pkg/repository"

	"github.com/google/uuid"
)

// FollowIntent defines the input to execute the follow request
type FollowIntent struct {
	To   uuid.UUID `json:"to"`
	From uuid.UUID `json:"from"`
}

// FollowService defines the follow use case.
type FollowService struct {
	UserRepository repository.Follower
}

// FollowRunner defines the contract to run the follow usecase
type FollowRunner interface {
	Run(intent FollowIntent) (bool, error)
}

// Run execute the use case
func (sv *FollowService) Run(intent FollowIntent) (bool, error) {
	_, err := sv.UserRepository.Follow(intent.To.String(), intent.From.String())

	if err != nil {
		return false, err
	}

	return true, nil
}

// NewFollowIntent creates a new intent for follow. It validades some rules before creating
func NewFollowIntent(to, from uuid.UUID) (FollowIntent, error) {
	if to == from {
		return FollowIntent{}, fmt.Errorf("can't follow yourself")
	}

	return FollowIntent{To: to, From: from}, nil
}
