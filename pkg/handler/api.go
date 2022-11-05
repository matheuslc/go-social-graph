package handler

import (
	"encoding/json"
	"gosocialgraph/pkg/service"
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
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

// PostHandler godoc
// @Summary      Creates a new post
// @Description  Creates a new post using the user id
// @Tags         post
// @Accept       json
// @Produce      json
// @Param        user_id body string true "user_id"
// @Param        content body string true "content"
// @Router       /post [post]
func (c *AppContext) PostHandler(w http.ResponseWriter, r *http.Request) {
	var intentToValidate service.CreatePostIntent

	err := json.NewDecoder(r.Body).Decode(&intentToValidate)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Request params are not the expected")
		return
	}

	intent, err := service.NewCreatePostIntent(intentToValidate.UserId, intentToValidate.Content)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	response, err := c.PostService.Run(intent)

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Could not create a new post")
	} else {
		respondWithJSON(w, http.StatusOK, response)
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

// FollowHandler godoc
// @Summary      follow a user
// @Tags         follow
// @Accept       json
// @Produce      json
// @Param        to body string true "to"
// @Param        from body string true "from"
// @Router       /follow [post]
func (c *AppContext) FollowHandler(w http.ResponseWriter, r *http.Request) {
	var invalidIntent service.FollowIntent

	err := json.NewDecoder(r.Body).Decode(&invalidIntent)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Request params are not the expected")
		return
	}

	intent, err := service.NewFollowIntent(invalidIntent.To, invalidIntent.From)

	if err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	response, err := c.FollowService.Run(intent)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
	} else {
		respondWithJSON(w, http.StatusOK, response)
	}
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

func (c *AppContext) ProfileHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	intent := service.ProfileIntent{
		UserID: uuid.MustParse(vars["user_id"]),
	}

	response, err := c.ProfileService.Run(intent)
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
// @Success     200 {object} entity.User
// @Router       /user [post]
func (c *AppContext) CreateUserHandler(w http.ResponseWriter, r *http.Request) {
	username := r.FormValue("username")
	w.WriteHeader(http.StatusOK)

	intent := service.CreateUserIntent{
		Username: username,
	}

	persistedUser, err := c.CreateUserService.Run(intent)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
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
