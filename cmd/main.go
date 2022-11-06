package main

import (
	_ "gosocialgraph/docs"
	"gosocialgraph/pkg/handler"
	"gosocialgraph/server"
)

func main() {
	context := handler.NewAppContext()

	server.RegisterHandlers(context.Router, context)

	// context.Router.HandleFunc("/repost", context.RepostHandler).Methods("POST")
	// context.Router.HandleFunc("/all", context.AllPostsHandler).Methods("GET")
	// context.Router.HandleFunc("/following/{user_id}", context.FollowingHandler).Methods("GET")

	context.Router.Logger.Fatal(context.Router.Start(":3010"))
}
