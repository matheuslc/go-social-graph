package service

import (
	"gosocialgraph/pkg/repository"

	"github.com/google/uuid"
)

type CreateRepostIntent struct {
	UserId   uuid.UUID `json:"user_id"`
	ParentId uuid.UUID `json:"parent_id"`
	Quote    string    `json:"quote"`
}

type CreateRepostResponse struct {
	Id uuid.UUID `json:"id"`
}

type RepostService struct {
	Repository repository.Reposter
}

func (sv RepostService) Run(intent CreateRepostIntent) (bool, error) {
	ok, err := sv.Repository.Repost(intent.UserId.String(), intent.ParentId.String(), intent.Quote)

	if err != nil {
		return ok, err
	}

	return ok, nil
}
