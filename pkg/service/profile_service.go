package service

//go:generate mockgen -source=./profile_service.go -destination=../mock/service/profile_service.go

import (
	"gosocialgraph/pkg/entity"

	"github.com/google/uuid"
)

// ProfileResponse defines the response from the use case, which in this case
// are going to be all information related to the user profile
type ProfileResponse struct {
	User  entity.User       `json:"user"`
	Stats ProfileStats      `json:"stats"`
	Posts []entity.UserPost `json:"posts"`
}

// ProfileService holds all dependencies a profile service must have
type ProfileService struct {
	FindUserService FindUserRunner
	StatsService    StatsRunner
	UserPostService UserPostRunner
}

// Run executes the use case. It will first find the user, than collec stats to show in his profile
func (sv ProfileService) Run(userID uuid.UUID) (ProfileResponse, error) {
	userFound, err := sv.FindUserService.Run(userID)
	stats, err := sv.StatsService.Run(userID)
	posts, err := sv.UserPostService.Run(userID)

	if err != nil {
		return ProfileResponse{}, err
	}

	return ProfileResponse{
		User:  userFound.User,
		Stats: stats.ProfileStats,
		Posts: posts,
	}, nil
}
