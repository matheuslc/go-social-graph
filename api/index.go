package api

import (
	"gosocialgraph/pkg/handler"
	"gosocialgraph/server"
	"net/http"
)

func Handler(w http.ResponseWriter, r *http.Request) {
	appContext := handler.NewAppContext()
	server.RegisterHandlers(appContext.Router, appContext)

	appContext.Router.ServeHTTP(w, r)
}
