package handler

import (
	"fmt"
	p "gosocialgraph/pkg/persistence"
	"gosocialgraph/pkg/post"
	"gosocialgraph/pkg/timeline"
	"gosocialgraph/pkg/usecase"
	"gosocialgraph/pkg/user"
	"os"

	"github.com/gorilla/mux"
	"github.com/neo4j/neo4j-go-driver/neo4j"
)

type AppContext struct {
	Db     *neo4j.Driver
	Router *mux.Router
	Graph  *p.Graph

	FindUserService   *user.FindUserService
	CreateUserService *usecase.CreateUserService
	UserRepository    *user.Repository
	StatsService      *user.StatsService

	AllService       *timeline.AllService
	FollowingService *timeline.FollowingService
	ProfileService   *timeline.ProfileService

	FollowService   *usecase.FollowService
	UnfollowService *user.UnfollowService

	PostService   *post.PostService
	RepostService *post.RepostService
}

func NewAppContext() AppContext {
	fmt.Println(os.Getenv("NEO4J_HOST"))
	fmt.Println(os.Getenv("NEO4J_USERNAME"))
	fmt.Println(os.Getenv("NEO4J_PASSWORD"))

	db, err := p.New(os.Getenv("NEO4J_HOST"), os.Getenv("NEO4J_USERNAME"), os.Getenv("NEO4J_PASSWORD"))
	if err != nil {
		fmt.Printf("Can't connect to neo4j. Reason: %s", err)
	}

	r := mux.NewRouter()
	userRepository := user.Repository{
		Client: db,
	}

	postRepository := post.Repository{
		Client: db,
	}

	timelineRepository := timeline.Repository{
		Client: db,
	}

	return AppContext{
		Db:             &db,
		Router:         r,
		Graph:          &p.Graph{Client: db},
		UserRepository: &userRepository,
		StatsService: &usecase.StatsService{
			Repository: &userRepository,
		},
		ProfileService: &timeline.ProfileService{
			FindUserService: usecase.FindUserService{
				UserRepository: &userRepository,
			},
			StatsService: usecase.StatsService{
				Repository: &userRepository,
			},
			UserPostService: timeline.UserPostService{
				Repository: &timelineRepository,
			},
		},
		CreateUserService: &usecase.CreateUserService{
			UserRepository: &userRepository,
		},
		AllService: &timeline.AllService{
			Repository: &timelineRepository,
		},
		FollowingService: &timeline.FollowingService{
			Repository: &timelineRepository,
		},

		FollowService: &usecase.FollowService{
			UserRepository: &userRepository,
		},
		UnfollowService: &usecase.UnfollowService{
			Repository: &userRepository,
		},

		PostService: &post.PostService{
			Repository: &postRepository,
		},
		RepostService: &post.RepostService{
			Repository: &postRepository,
		},
	}
}
