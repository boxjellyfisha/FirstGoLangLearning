package webapi

import (
	"net/http"

	"github.com/gorilla/mux"

	"kzapp/webapi/count"
	"kzapp/webapi/hello"
	"kzapp/webapi/user"
)

func RunServer() {
	// Create the router
	router := mux.NewRouter()

	// Handle incoming request
	hello.InitGreeting(router)
	count.InitCaculator(router)
	user.InitUserService(router)

	// Listen for HTTP Connections, Registering a Request Handler
	http.ListenAndServe(":80", router)
}
