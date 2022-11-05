package main

import (
	_ "gosocialgraph/docs"
	"gosocialgraph/pkg/handler"
	"gosocialgraph/server"
)

func main() {
	context := handler.NewAppContext()

	server.RegisterHandlers(context.Router, context)
	context.Router.Logger.Fatal(context.Router.Start(":3010"))
}
