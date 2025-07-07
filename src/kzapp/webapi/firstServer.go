package webapi

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"

	"kzapp/webapi/db"
)

func RunServer() {

	userDao := ProvideUserDao()
	go userDao.CreateUser(db.User{
		Name:     "John Doe",
		Email:    "john.doe@example.com",
		Password: "password",
	})
	var users []db.User
	go func() { users, _ = userDao.GetUsers() }()
	defer fmt.Println(users)

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
