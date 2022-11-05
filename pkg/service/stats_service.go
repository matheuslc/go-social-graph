package service

//go:generate mockgen -source=./stats_service.go -destination=../mock/service/stats_service.go

import (
	"gosocialgraph/pkg/repository"
	"sync"

	"github.com/google/uuid"
)

// StatsResponse defines the usecase response
type StatsResponse struct {
	ProfileStats `json:"user_stats"`
}

// ProfileStats groups structs related to user stats
type ProfileStats struct {
	Followers  int `json:"followers"`
	Following  int `json:"following"`
	PostsCount int `json:"posts_count"`
}

// StatsService defines the service and its dependencies
type StatsService struct {
	Repository repository.Stats
}

// StatsRunner
type StatsRunner interface {
	Run(userID uuid.UUID) (StatsResponse, error)
}

// Run executes the usecase
func (sv StatsService) Run(userID uuid.UUID) (StatsResponse, error) {
	var wg sync.WaitGroup
	var userStats ProfileStats
	var runningErros []error

	wg.Add(3)

	go func() {
		defer wg.Done()
		followers, err := sv.Repository.CountFollowers(userID)

		if err != nil {
			runningErros = append(runningErros, err)
		}

		userStats.Followers = int(followers)
	}()

	go func() {
		defer wg.Done()
		following, err := sv.Repository.CountFollowing(userID)

		if err != nil {
			runningErros = append(runningErros, err)
		}

		userStats.Following = int(following)
	}()

	go func() {
		defer wg.Done()
		posts, err := sv.Repository.CountPosts(userID)

		if err != nil {
			runningErros = append(runningErros, err)
		}

		userStats.PostsCount = int(posts)
	}()

	wg.Wait()

	if len(runningErros) > 0 {
		return StatsResponse{}, runningErros[0]
	}

	return StatsResponse{ProfileStats: userStats}, nil
}
