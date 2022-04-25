package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	p "gosocialgraph/pkg/persistence"
	"gosocialgraph/pkg/post"
	"gosocialgraph/pkg/timeline"
	"gosocialgraph/pkg/user"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/neo4j/neo4j-go-driver/neo4j"

	_ "gosocialgraph/cmd/docs"

	httpSwagger "github.com/swaggo/http-swagger"
)

type AppContext struct {
	Db     *neo4j.Driver
	Router *mux.Router
	Graph  *p.Graph

	FindUserService *user.FindUserService
	UserRepository  *user.UserRepository
	StatsService    *user.StatsService

	AllService       *timeline.AllService
	FollowingService *timeline.FollowingService
	ProfileService   *timeline.ProfileService

	FollowService   *user.FollowService
	UnfollowService *user.UnfollowService

	PostService   *post.PostService
	RepostService *post.RepostService
}

func NewAppContext() AppContext {
	db, err := p.New(os.Getenv("NEO4J_HOST"), os.Getenv("NEO4J_USERNAME"), os.Getenv("NEO4J_PASSWORD"))
	if err != nil {
		fmt.Printf("Can't connect to neo4j. Reason: %s", err)
	}

	r := mux.NewRouter()
	userRepository := user.UserRepository{
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
		StatsService: &user.StatsService{
			Repository: &userRepository,
		},
		ProfileService: &timeline.ProfileService{
			FindUserService: user.FindUserService{
				UserRepository: &userRepository,
			},
			StatsService: user.StatsService{
				Repository: &userRepository,
			},
			UserPostService: timeline.UserPostService{
				Repository: &timelineRepository,
			},
		},
		AllService: &timeline.AllService{
			Repository: &timelineRepository,
		},
		FollowingService: &timeline.FollowingService{
			Repository: &timelineRepository,
		},

		FollowService: &user.FollowService{
			UserRepository: &userRepository,
		},
		UnfollowService: &user.UnfollowService{
			UserRepository: &userRepository,
		},

		PostService: &post.PostService{
			Repository: &postRepository,
		},
		RepostService: &post.RepostService{
			Repository: &postRepository,
		},
	}
}

// @title Swagger Example API
// @version 1.0
// @description This is a sample server Petstore server.
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:3010
func main() {
	context := NewAppContext()

	context.Router.HandleFunc("/user", context.CreateUserHandler).Methods("POST")
	context.Router.HandleFunc("/profile/{user_id}", context.ProfileHandler).Methods("GET")
	context.Router.HandleFunc("/tweet", context.PostHandler).Methods("POST")
	context.Router.HandleFunc("/repost", context.RepostHandler).Methods("POST")
	context.Router.HandleFunc("/follow", context.FollowHandler).Methods("POST")
	context.Router.HandleFunc("/unfollow", context.UnfollowHandler).Methods("DELETE")
	context.Router.HandleFunc("/all", context.AllPostsHandler).Methods("GET")
	context.Router.HandleFunc("/following/{user_id}", context.FollowingHandler).Methods("GET")
	context.Router.PathPrefix("/swagger/").Handler(httpSwagger.WrapHandler)

	server := &http.Server{
		Handler: context.Router,
		Addr:    "0.0.0.0:3010",
	}

	log.Fatal(server.ListenAndServe())
	fmt.Println("We are online! Running on 0.0.0.0:3010")
}

func (context AppContext) RepostHandler(w http.ResponseWriter, r *http.Request) {
	var intent post.CreateRepostIntent

	err := json.NewDecoder(r.Body).Decode(&intent)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Request params are not the expected")
	}

	_, err = context.RepostService.Run(intent)

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Could not repost")
	} else {
		w.WriteHeader(http.StatusCreated)
	}
}

func (context AppContext) PostHandler(w http.ResponseWriter, r *http.Request) {
	var intentToValidate post.CreatePostIntent

	err := json.NewDecoder(r.Body).Decode(&intentToValidate)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Request params are not the expected")
		return
	}

	intent, err := post.NewCreatePostIntent(intentToValidate.UserId, intentToValidate.Content)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	response, err := context.PostService.Run(intent)

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Could not create a new post")
	} else {
		respondWithJSON(w, http.StatusOK, response)
	}
}

func (context AppContext) FollowingHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	intent := timeline.FollowingIntent{
		UserId: uuid.MustParse(vars["user_id"]),
	}

	response, err := context.FollowingService.Run(intent)

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Could not list all posts")
	} else {
		respondWithJSON(w, http.StatusOK, response)
	}
}

func (context AppContext) AllPostsHandler(w http.ResponseWriter, r *http.Request) {
	response, err := context.AllService.Run()

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Could not list all posts")
	} else {
		respondWithJSON(w, http.StatusOK, response)
	}
}

func (context AppContext) FollowHandler(w http.ResponseWriter, r *http.Request) {
	var invalidIntent user.FollowIntent

	err := json.NewDecoder(r.Body).Decode(&invalidIntent)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Request params are not the expected")
		return
	}

	intent, err := user.NewFollowIntent(invalidIntent.To, invalidIntent.From)

	if err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	response, err := context.FollowService.Run(intent)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
	} else {
		respondWithJSON(w, http.StatusOK, response)
	}
}

func (context AppContext) UnfollowHandler(w http.ResponseWriter, r *http.Request) {
	var unfollowIntent user.UnfollowIntent

	err := json.NewDecoder(r.Body).Decode(&unfollowIntent)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Request params are not the expected")
	}

	_, err = context.UnfollowService.Run(unfollowIntent)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
	} else {
		respondWithJSON(w, http.StatusOK, "ok")
	}
}

func (context AppContext) ProfileHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	intent := timeline.ProfileIntent{
		UserId: uuid.MustParse(vars["user_id"]),
	}

	response, err := context.ProfileService.Run(intent)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
	} else {
		respondWithJSON(w, http.StatusOK, response)
	}
}

// CreateUserHandler godoc
// @Summary      Create a user
// @Description  creates a new user which is required to use all other resources
// @Tags         users
// @Accept       json
// @Produce      json
// @Param        username query string true "username"
// @Router       /user [post]
func (context AppContext) CreateUserHandler(w http.ResponseWriter, r *http.Request) {
	username := r.FormValue("username")
	w.WriteHeader(http.StatusOK)

	persistedUser, err := context.UserRepository.Create(username)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Could not create a new user")
	} else {
		respondWithJSON(w, http.StatusOK, persistedUser)
	}
}

func respondWithError(w http.ResponseWriter, code int, message string) {
	respondWithJSON(w, code, map[string]string{"error": message})
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}
