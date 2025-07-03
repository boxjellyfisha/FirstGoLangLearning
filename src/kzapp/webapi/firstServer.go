package webapi

import (
	"net/http"

	"github.com/gorilla/mux"
)

func RunServer() {
	// Create the router
	router := mux.NewRouter()

	handlers := CreateHandler()

	// Handle incoming request
	for _, handler := range handlers {
		handler.InitService(router)
	}

	// Listen for HTTP Connections, Registering a Request Handler
	http.ListenAndServe(":80", router)
}
