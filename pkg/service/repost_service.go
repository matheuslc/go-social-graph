package service

import (
	"gosocialgraph/pkg/repository"

	"github.com/google/uuid"
)

// RepostIntent defines the input the use case expect to run
type RepostIntent struct {
	UserID       uuid.UUID `json:"user_id"`
	ParentPostID uuid.UUID `json:"parent_id"`
	Quote        string    `json:"quote"`
}

// RepostResponse defines the output for the usecase
type RepostResponse struct {
	Status bool `json:"status"`
}

// RepostService holds the structure and its dependencies
type RepostService struct {
	Repository repository.Reposter
}

// RepostRunner defines the signature to run this service
type RepostRunner interface {
	Run(intent RepostIntent) (bool, error)
}

// Run executes the usecase
func (sv RepostService) Run(intent RepostIntent) (RepostResponse, error) {
	ok, err := sv.Repository.Repost(intent.UserID.String(), intent.ParentPostID.String(), intent.Quote)
	response := RepostResponse{Status: ok}

	return response, err
}
