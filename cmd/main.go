package main

import (
	"fmt"
	"log"
	"net/http"

	_ "gosocialgraph/docs"
	"gosocialgraph/pkg/handler"

	httpSwagger "github.com/swaggo/http-swagger"
)

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
	context := handler.NewAppContext()

	context.Router.HandleFunc("/user", context.CreateUserHandler).Methods("POST")
	context.Router.HandleFunc("/profile/{user_id}", context.ProfileHandler).Methods("GET")
	context.Router.HandleFunc("/post", context.PostHandler).Methods("POST")
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
