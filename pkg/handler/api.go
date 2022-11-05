package handler

import (
	"encoding/json"
	"gosocialgraph/openapi"
	"gosocialgraph/pkg/service"
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/labstack/echo/v4"
)

// RepostHandler godoc
// @Summary      Repost a post from someone
// @Description  Creates a respost from a post
// @Tags         repost
// @Accept       json
// @Produce      json
// @Param        user_id body string true "user_id"
// @Param        parent_id body string true "parent_id"
// @Param        quote body string false "string"
// @Router       /repost [post]
func (c *AppContext) RepostHandler(w http.ResponseWriter, r *http.Request) {
	var intent service.RepostIntent

	err := json.NewDecoder(r.Body).Decode(&intent)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Request params are not the expected")
	}

	_, err = c.RepostService.Run(intent)

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Could not repost")
	} else {
		w.WriteHeader(http.StatusCreated)
	}
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

// FollowingHandler godoc
// @Summary      Starts to follow an user
// @Tags         follow
// @Accept       json
// @Produce      json
// @Param        user_id body string true "user_id"
// @Success 	 200 {object} service.FollowingResponse
// @Router       /follow [post]
func (c *AppContext) FollowingHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	intent := service.FollowingIntent{
		UserID: uuid.MustParse(vars["user_id"]),
	}

	response, err := c.FollowingService.Run(intent)

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Could not list all posts")
	} else {
		respondWithJSON(w, http.StatusOK, response)
	}
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
	response, err := c.FollowService.Run(id, from)
	if err != nil {
		return err
	}

	return echoContext.JSON(http.StatusOK, response)
}

func (c *AppContext) UnfollowHandler(w http.ResponseWriter, r *http.Request) {
	var unfollowIntent service.UnfollowIntent

	err := json.NewDecoder(r.Body).Decode(&unfollowIntent)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Request params are not the expected")
	}

	_, err = c.UnfollowService.Run(unfollowIntent)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
	} else {
		respondWithJSON(w, http.StatusOK, "ok")
	}
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
