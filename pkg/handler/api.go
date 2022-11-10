package handler

import (
	"encoding/json"
	"gosocialgraph/openapi"
	"gosocialgraph/pkg/handler/rest"
	"net/http"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

func (c AppContext) RepostHandler(echoContext echo.Context, id uuid.UUID) error {
	var intent openapi.RepostIntent
	if err := echoContext.Bind(&intent); err != nil {
		return err
	}

	if err := c.RepostService.Run(id, intent.Parent, *intent.Quote); err != nil {
		return err
	}

	return echoContext.NoContent(http.StatusOK)
}

func (c AppContext) PostHandler(echoContext echo.Context) error {
	var intent openapi.CreatePostRequest
	err := echoContext.Bind(&intent)
	if err != nil {
		return err
	}

	post, err := c.PostService.Run(intent.UserId, intent.Content)
	if err != nil {
		return err
	}

	restPost := openapi.CreatePostResponse{Id: post.ID, Content: post.Content}

	return echoContext.JSON(http.StatusCreated, restPost)
}

func (c AppContext) TimelineHandler(echoContext echo.Context, id uuid.UUID) error {
	response, err := c.TimelineServive.Run(id)
	if err != nil {
		return err
	}

	openapiResponse, err := rest.MapUserPostsToOpenAPI(response)
	if err != nil {
		return err
	}

	return echoContext.JSON(http.StatusOK, openapi.TimelineResponse{Posts: &openapiResponse})
}

func (c *AppContext) AllPostsHandler(w http.ResponseWriter, r *http.Request) {
	response, err := c.AllService.Run()

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Could not list all posts")
	} else {
		respondWithJSON(w, http.StatusOK, response)
	}
}

func (c AppContext) FollowHandler(echoContext echo.Context, id uuid.UUID, from uuid.UUID) error {
	if err := c.FollowService.Run(id, from); err != nil {
		return err
	}

	return echoContext.NoContent(http.StatusOK)
}

func (c AppContext) UnfollowHandler(echoContext echo.Context, id uuid.UUID, from uuid.UUID) error {
	if err := c.UnfollowService.Run(id, from); err != nil {
		return err
	}

	return echoContext.NoContent(http.StatusOK)
}

func (c AppContext) ProfileHandler(echoContext echo.Context, userID uuid.UUID) error {
	response, err := c.ProfileService.Run(userID)
	if err != nil {
		return err
	}

	openapiPosts, err := rest.MapUserPostsToOpenAPI(response.Posts)
	if err != nil {
		return err
	}

	openapiResponse := openapi.ProfileResponse{
		Posts: &openapiPosts,
		Stats: &openapi.UserStats{
			Followers:  &response.Stats.Followers,
			Following:  &response.Stats.Following,
			PostsCount: &response.Stats.PostsCount,
		},
		User: &openapi.User{
			CreatedAt: response.User.CreatedAt,
			Id:        response.User.ID,
			Username:  response.User.Username,
		},
	}

	return echoContext.JSON(http.StatusOK, openapiResponse)
}

func (c AppContext) CreateUser(echoContext echo.Context) error {
	username := echoContext.FormValue("username")

	persistedUser, err := c.CreateUserService.Run(username)
	if err != nil {
		return err
	}

	restResponse := openapi.CreateUserResponse{
		Id:        persistedUser.ID,
		CreatedAt: persistedUser.CreatedAt,
		Username:  persistedUser.Username,
	}

	return echoContext.JSON(http.StatusCreated, restResponse)
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
