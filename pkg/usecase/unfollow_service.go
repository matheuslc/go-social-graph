package usecase

import (
	"gosocialgraph/pkg/user"

	"github.com/google/uuid"
)

// UnfollowIntent defines what you need to execute the unfollow intent
type UnfollowIntent struct {
	To   uuid.UUID `json:"to"`
	From uuid.UUID `json:"from"`
}

// UnfollowService defines the struct holder for the use case
type UnfollowService struct {
	Repository user.Unfollower
}

// Run executes the use case
func (sv UnfollowService) Run(intent UnfollowIntent) (bool, error) {
	return sv.Repository.Unfollow(intent.To.String(), intent.From.String())
}
