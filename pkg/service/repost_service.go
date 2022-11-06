package service

//go:generate mockgen -source=./repost_service.go -destination=../mock/service/repost_service.go

import (
	"gosocialgraph/pkg/repository"

	"github.com/google/uuid"
)

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
	Run(userID, parentID uuid.UUID, quote string) (bool, error)
}

// Run executes the usecase
func (sv RepostService) Run(userID, parentID uuid.UUID, quote string) error {
	if err := sv.Repository.Repost(userID.String(), parentID.String(), quote); err != nil {
		return err
	}

	return nil
}
