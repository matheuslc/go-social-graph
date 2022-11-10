package service

//go:generate mockgen -source=./timeline_service.go -destination=../mock/service/timeline_service.go

import (
	"gosocialgraph/pkg/entity"
	"gosocialgraph/pkg/repository"

	"github.com/google/uuid"
)

// TimelineServive defines the struct holder for the service and its dependencies
type TimelineServive struct {
	Repository repository.TimelineReader
}

// Run executes the use case
func (sv TimelineServive) Run(userID uuid.UUID) (response []entity.UserPost, err error) {
	response, err = sv.Repository.TimelineFor(userID)
	if err != nil {
		return response, err
	}

	return response, err
}
