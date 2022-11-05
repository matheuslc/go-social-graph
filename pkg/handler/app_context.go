package handler

import (
	"fmt"
	p "gosocialgraph/pkg/persistence"
	"gosocialgraph/pkg/repository"
	"gosocialgraph/pkg/service"
	"os"

	"github.com/labstack/echo/v4"
	middleware "github.com/labstack/echo/v4/middleware"
	"github.com/neo4j/neo4j-go-driver/neo4j"
)

// AppContext is the application container, where dependencies are defined
type AppContext struct {
	Db     *neo4j.Driver
	Router *echo.Echo
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

// NewAppContext creates a new AppContext within all dependencies builded
func NewAppContext() AppContext {
	db, err := p.New(os.Getenv("NEO4J_HOST"), os.Getenv("NEO4J_USERNAME"), os.Getenv("NEO4J_PASSWORD"))
	if err != nil {
		fmt.Printf("Can't connect to neo4j. Reason: %s", err)
	}

	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

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
		Router:         e,
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
