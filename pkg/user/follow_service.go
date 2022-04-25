package user

import (
	"fmt"

	"github.com/google/uuid"
)

type FollowIntent struct {
	To   uuid.UUID `json:"to"`
	From uuid.UUID `json:"from"`
}

type FollowService struct {
	UserRepository Writer
}

func (sv FollowService) Run(intent FollowIntent) (bool, error) {
	_, err := sv.UserRepository.Follow(intent.To.String(), intent.From.String())

	if err != nil {
		return false, err
	}

	return true, nil
}

func NewFollowIntent(to, from uuid.UUID) (FollowIntent, error) {
	if to == from {
		return FollowIntent{}, fmt.Errorf("can't follow yourself")
	}

	return FollowIntent{
		To:   to,
		From: from,
	}, nil
}
