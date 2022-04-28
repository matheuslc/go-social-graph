package handler

import (
	"fmt"
	p "gosocialgraph/pkg/persistence"
	"gosocialgraph/pkg/repository"
	"gosocialgraph/pkg/service"
	"os"

	"github.com/gorilla/mux"
	"github.com/neo4j/neo4j-go-driver/neo4j"
)

type AppContext struct {
	Db     *neo4j.Driver
	Router *mux.Router
	Graph  *p.Graph

	FindUserService   *service.FindUserService
	CreateUserService *service.CreateUserService
	UserRepository    *repository.UserRepository
	StatsService      *service.StatsService

	AllService       *service.AllService
	FollowingService *service.FollowingService
	ProfileService   *service.ProfileService

	FollowService   *service.FollowService
	UnfollowService *service.UnfollowService

	PostService   *service.PostService
	RepostService *service.RepostService
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
	userRepository := repository.UserRepository{
		Client: db,
	}

	postRepository := repository.PostRepository{
		Client: db,
	}

	timelineRepository := repository.TimelineRepository{
		Client: db,
	}

	return AppContext{
		Db:             &db,
		Router:         r,
		Graph:          &p.Graph{Client: db},
		UserRepository: &userRepository,
		StatsService: &service.StatsService{
			Repository: &userRepository,
		},
		ProfileService: &service.ProfileService{
			FindUserService: service.FindUserService{
				UserRepository: &userRepository,
			},
			StatsService: service.StatsService{
				Repository: &userRepository,
			},
			UserPostService: service.UserPostService{
				Repository: &timelineRepository,
			},
		},
		CreateUserService: &service.CreateUserService{
			UserRepository: &userRepository,
		},
		AllService: &service.AllService{
			Repository: &timelineRepository,
		},
		FollowingService: &service.FollowingService{
			Repository: &timelineRepository,
		},

		FollowService: &service.FollowService{
			UserRepository: &userRepository,
		},
		UnfollowService: &service.UnfollowService{
			Repository: &userRepository,
		},

		PostService: &service.PostService{
			Repository: &postRepository,
		},
		RepostService: &service.RepostService{
			Repository: &postRepository,
		},
	}
}
