package user

import "github.com/google/uuid"

type UnfollowIntent struct {
	To   uuid.UUID `json:"to"`
	From uuid.UUID `json:"from"`
}

type UnfollowService struct {
	UserRepository Writer
}

func (sv UnfollowService) Run(intent UnfollowIntent) (bool, error) {
	_, err := sv.UserRepository.Unfollow(intent.To.String(), intent.From.String())

	if err != nil {
		return false, err
	}

	return true, nil
}
