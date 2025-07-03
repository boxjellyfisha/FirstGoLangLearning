package hello

import (
	"fmt"
	"github.com/gorilla/mux"
	"kzapp/webapi/pkg"
	"kzapp/webapi/user"
	"net/http"
)

func sayhello() func(w http.ResponseWriter, r *http.Request) {
	handleFun := func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello, you've requested: %s\n", r.URL.Path)
	}
	return handleFun
}

func greet(w http.ResponseWriter, r *http.Request) {
	userName := r.URL.Query().Get("name")
	fmt.Printf("userName: %s\n", userName)
	greeting := "Welcome Back!"

	response := user.UserResponse{
		Name:    userName,
		Message: &greeting,
	}

	pkg.JsonResponse(w, response)
}

func InitGreeting(router *mux.Router) {
	router.HandleFunc("/", pkg.Chain(sayhello(), pkg.Method("GET"), pkg.Logging()))
	router.HandleFunc("/greet", pkg.Chain(greet, pkg.Method("GET"), pkg.Logging()))
}
