package handler

import (
	"encoding/json"
	"gosocialgraph/pkg/post"
	"gosocialgraph/pkg/timeline"
	"gosocialgraph/pkg/usecase"
	"gosocialgraph/pkg/user"
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
func (c AppContext) RepostHandler(w http.ResponseWriter, r *http.Request) {
	var intent post.CreateRepostIntent

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
func (context *AppContext) PostHandler(w http.ResponseWriter, r *http.Request) {
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

// FollowingHandler godoc
// @Summary      Starts to follow an user
// @Tags         follow
// @Accept       json
// @Produce      json
// @Param        user_id body string true "user_id"
// @Success 	 200 {object} timeline.FollowingResponse
// @Router       /follow [post]
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
	var invalidIntent usecase.FollowIntent

	err := json.NewDecoder(r.Body).Decode(&invalidIntent)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Request params are not the expected")
		return
	}

	intent, err := usecase.NewFollowIntent(invalidIntent.To, invalidIntent.From)

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
// @Success     200 {object} user.User
// @Router       /user [post]
func (context AppContext) CreateUserHandler(w http.ResponseWriter, r *http.Request) {
	username := r.FormValue("username")
	w.WriteHeader(http.StatusOK)

	intent := usecase.CreateUserIntent{
		Username: username,
	}

	persistedUser, err := context.CreateUserService.Run(intent)
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