package api

import (
	"net/http"
)

// var globalAppcontext handler.AppContext

// func init() {
// 	globalAppcontext = handler.NewAppContext()
// }
func HelloWorld(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello World"))
}

// func Handler(w http.ResponseWriter, r *http.Request) {
// 	// appContext := handler.NewAppContext()
// 	// server.RegisterHandlers(globalAppcontext.Router, globalAppcontext)
// 	// server.RegisterHandlers(globalAppcontext.Router, globalAppcontext)

// 	// globalAppcontext.Router.ServeHTTP(w, r)
// }
