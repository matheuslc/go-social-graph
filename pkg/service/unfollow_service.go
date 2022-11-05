package service

//go:generate mockgen -source=./unfollow_service.go -destination=../mock/service/unfollow_service.go

import (
	"errors"
	"gosocialgraph/pkg/repository"

	"github.com/google/uuid"
)

// UnfollowIntent defines what you need to execute the unfollow intent
type UnfollowIntent struct {
	To   uuid.UUID `json:"to"`
	From uuid.UUID `json:"from"`
}

// UnfollowService defines the struct holder for the use case
type UnfollowService struct {
	Repository repository.Unfollower
}

// UnfolowRunner
type UnfolowRunner interface {
	Run(intent UnfollowIntent) (bool, error)
}

// Run executes the use case
func (sv UnfollowService) Run(intent UnfollowIntent) (bool, error) {
	return sv.Repository.Unfollow(intent.To.String(), intent.From.String())
}

// NewUnfollowIntent validades and create a new unfollow intent
func NewUnfollowIntent(to, from uuid.UUID) (UnfollowIntent, error) {
	if to == from {
		return UnfollowIntent{}, errors.New("Can't unfollow yourself")
	}

	return UnfollowIntent{To: to, From: from}, nil
}
