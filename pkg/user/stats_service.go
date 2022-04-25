package user

import (
	"sync"

	"github.com/google/uuid"
)

type StatsService struct {
	Repository Stats
}

type StatsIntent struct {
	UserId uuid.UUID `json:"user_id"`
}

type StatsResponse struct {
	UserStats `json:"user_stats"`
}

func (sv StatsService) Run(intent StatsIntent) (UserStats, error) {
	var wg sync.WaitGroup
	var userStats UserStats
	var runningErros []error

	wg.Add(3)

	go func() {
		defer wg.Done()
		followers, err := sv.Repository.CountFollowers(intent.UserId.String())

		if err != nil {
			runningErros = append(runningErros, err)
		}

		userStats.Followers = int(followers)
	}()

	go func() {
		defer wg.Done()
		following, err := sv.Repository.CountFollowing(intent.UserId.String())

		if err != nil {
			runningErros = append(runningErros, err)
		}

		userStats.Following = int(following)
	}()

	go func() {
		defer wg.Done()
		posts, err := sv.Repository.CountPosts(intent.UserId.String())

		if err != nil {
			runningErros = append(runningErros, err)
		}

		userStats.PostsCount = int(posts)
	}()

	wg.Wait()

	if len(runningErros) > 0 {
		return UserStats{}, runningErros[0]
	}

	return userStats, nil
}
