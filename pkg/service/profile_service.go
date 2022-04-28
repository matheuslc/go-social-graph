package service

import (
	"gosocialgraph/pkg/entity"

	"github.com/google/uuid"
)

// ProfileIntent defines what is needed to run the usecase
type ProfileIntent struct {
	UserID uuid.UUID `json:"user_id"`
}

// ProfileResponse defines the response from the use case, which in this case
// are going to be all information related to the user profile
type ProfileResponse struct {
	User  entity.User      `json:"user"`
	Stats ProfileStats     `json:"stats"`
	Posts UserPostResponse `json:"posts"`
}

// ProfileService holds all dependencies a profile service must have
type ProfileService struct {
	FindUserService FindUserRunner
	StatsService    StatsRunner
	UserPostService UserPostRunner
}

// Run executes the use case. It will first find the user, than collec stats to show in his profile
func (sv ProfileService) Run(intent ProfileIntent) (ProfileResponse, error) {
	userFound, err := sv.FindUserService.Run(FindUserIntent{UserID: intent.UserID})
	stats, err := sv.StatsService.Run(StatsIntent{UserID: intent.UserID})
	posts, err := sv.UserPostService.Run(UserPostsIntent{UserID: intent.UserID})

	if err != nil {
		return ProfileResponse{}, err
	}

	return ProfileResponse{
		User:  userFound.User,
		Stats: stats.ProfileStats,
		Posts: posts,
	}, nil
}
