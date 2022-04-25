package user

import (
	"github.com/google/uuid"
)

type FindUserIntent struct {
	UserId uuid.UUID `json:"user_id"`
}

type FindUserResponse struct {
	User `json:"user"`
}

type FindUserService struct {
	UserRepository Reader
}

func (sv FindUserService) Run(intent FindUserIntent) (FindUserResponse, error) {
	user, err := sv.UserRepository.Find(intent.UserId.String())

	if err != nil {
		return FindUserResponse{}, err
	}

	return FindUserResponse{
		User: user,
	}, nil
}
