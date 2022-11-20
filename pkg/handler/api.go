package handler

import (
	"gosocialgraph/openapi"
	"gosocialgraph/pkg/handler/rest"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

func (c AppContext) LoginHandler(echoContext echo.Context) error {
	username := echoContext.FormValue("username")
	passowrd := echoContext.FormValue("password")

	token, refresh, err := c.AuthService.Run(username, passowrd)
	if err != nil {
		return err
	}

	tokenCookie := new(http.Cookie)
	tokenCookie.Name = "access_token"
	tokenCookie.Value = token

	refreshCookie := new(http.Cookie)
	refreshCookie.Name = "refresh_token"
	refreshCookie.Value = refresh

	tokenCookie.Expires = time.Now().Add(240 * time.Hour)
	refreshCookie.Expires = time.Now().Add(240 * time.Hour)

	echoContext.SetCookie(tokenCookie)
	echoContext.SetCookie(refreshCookie)

	return echoContext.JSON(http.StatusOK, openapi.LoginResponse{AccessToken: &token, RefreshToken: &refresh})
}

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

	userID := uuid.MustParse(echoContext.Get("userID").(string))

	post, err := c.PostService.Run(userID, intent.Content)
	if err != nil {
		return err
	}

	restPost := openapi.CreatePostResponse{Id: post.ID, Content: post.Content}

	return echoContext.JSON(http.StatusCreated, restPost)
}

func (c AppContext) TimelineHandler(echoContext echo.Context) error {
	userID := uuid.MustParse(echoContext.Get("userID").(string))

	response, err := c.TimelineServive.Run(userID)
	if err != nil {
		return err
	}

	openapiResponse, err := rest.MapUserPostsToOpenAPI(response)
	if err != nil {
		return err
	}

	return echoContext.JSON(http.StatusOK, openapi.TimelineResponse{Posts: &openapiResponse})
}

// func (c *AppContext) AllPostsHandler(w http.ResponseWriter, r *http.Request) {
// 	response, err := c.AllService.Run()

// 	if err != nil {
// 		respondWithError(w, http.StatusInternalServerError, "Could not list all posts")
// 	} else {
// 		respondWithJSON(w, http.StatusOK, response)
// 	}
// }

func (c AppContext) FollowHandler(echoContext echo.Context, from uuid.UUID) error {
	userID := uuid.MustParse(echoContext.Get("userID").(string))
	if err := c.FollowService.Run(userID, from); err != nil {
		return err
	}

	return echoContext.NoContent(http.StatusOK)
}

func (c AppContext) UnfollowHandler(echoContext echo.Context, from uuid.UUID) error {
	userID := uuid.MustParse(echoContext.Get("userID").(string))
	if err := c.UnfollowService.Run(userID, from); err != nil {
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
	email := echoContext.FormValue("email")
	password := echoContext.FormValue("password")

	persistedUser, err := c.CreateUserService.Run(username, email, password)
	if err != nil {
		return err
	}

	restResponse := openapi.CreateUserResponse{
		Id:        persistedUser.ID,
		CreatedAt: persistedUser.CreatedAt,
		Username:  persistedUser.Username,
		Email:     persistedUser.Email,
	}

	return echoContext.JSON(http.StatusCreated, restResponse)
}
