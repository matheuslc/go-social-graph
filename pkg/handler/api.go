package handler

import (
	"encoding/json"
	"gosocialgraph/openapi"
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

	response, err := c.PostService.Run(intent)

	if err != nil {
		return err
	} else {
		return echoContext.JSON(http.StatusCreated, response)
	}
}

func (c AppContext) TimelineHandler(echoContext echo.Context, id uuid.UUID) error {
	response, err := c.TimelineServive.Run(id)
	if err != nil {
		return err
	}

	return echoContext.JSON(http.StatusOK, response)
}

// AllPostsHandler godoc
// @Summary      List all posts
// @Tags         follow
// @Accept       json
// @Produce      json
// @Param        user_id body string true "user_id"
// @Success 	 200 {object} service.AllPostResponse
// @Router       /all [get]
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

	return echoContext.JSON(http.StatusOK, response)
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
