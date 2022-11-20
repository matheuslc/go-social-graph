package pages

import (
	"fmt"
	"gosocialgraph/pkg/handler"
	"gosocialgraph/server"
	"net/http"
)

var globalAppcontext handler.AppContext

func init() {
	globalAppcontext = handler.NewAppContext()
}

func Handler(w http.ResponseWriter, r *http.Request) {
	// appContext := handler.NewAppContext()
	server.RegisterHandlers(globalAppcontext.Router, globalAppcontext)
	server.RegisterHandlers(globalAppcontext.Router, globalAppcontext)

	fmt.Println("daeee")

	globalAppcontext.Router.ServeHTTP(w, r)
}
