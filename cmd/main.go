package main

import (
	_ "gosocialgraph/docs"
	"gosocialgraph/pkg/handler"
	"gosocialgraph/server"
)

func main() {
	context := handler.NewAppContext()

	server.RegisterHandlers(context.Router, context)

	// context.Router.HandleFunc("/profile/{user_id}", context.ProfileHandler).Methods("GET")
	// context.Router.HandleFunc("/post", context.PostHandler).Methods("POST")
	// context.Router.HandleFunc("/repost", context.RepostHandler).Methods("POST")
	// context.Router.HandleFunc("/follow", context.FollowHandler).Methods("POST")
	// context.Router.HandleFunc("/unfollow", context.UnfollowHandler).Methods("DELETE")
	// context.Router.HandleFunc("/all", context.AllPostsHandler).Methods("GET")
	// context.Router.HandleFunc("/following/{user_id}", context.FollowingHandler).Methods("GET")

	context.Router.Logger.Fatal(context.Router.Start(":3010"))
}
