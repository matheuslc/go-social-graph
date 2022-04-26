package timeline

import (
	"gosocialgraph/pkg/user"

	"github.com/google/uuid"
)

type ProfileIntent struct {
	UserId uuid.UUID `json:"user_id"`
}

type ProfileResponse struct {
	User  user.User        `json:"user"`
	Stats user.Stats       `json:"stats"`
	Posts UserPostResponse `json:"posts"`
}

type ProfileService struct {
	FindUserService user.FindUserService
	StatsService    user.StatsService
	UserPostService UserPostService
}

func (sv ProfileService) Run(intent ProfileIntent) (ProfileResponse, error) {
	userFound, err := sv.FindUserService.Run(user.FindUserIntent{
		UserId: intent.UserId,
	})

	stats, err := sv.StatsService.Run(user.StatsIntent{
		UserId: intent.UserId,
	})

	posts, err := sv.UserPostService.Run(UserPostsIntent{
		UserId: intent.UserId,
	})

	if err != nil {
		return ProfileResponse{}, err
	}

	return ProfileResponse{
		User:  userFound.User,
		Stats: stats,
		Posts: posts,
	}, nil
}
