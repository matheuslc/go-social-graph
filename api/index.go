package api

import (
	"gosocialgraph/pkg/handler"
	"gosocialgraph/server"
	"net/http"
)

var globalContext handler.AppContext

func init() {
	globalContext = handler.NewAppContext()
}

func Handler(w http.ResponseWriter, r *http.Request) {
	server.RegisterHandlers(globalContext.Router, globalContext)

	globalContext.Router.ServeHTTP(w, r)
}
