package hello

import (
	"fmt"
	"github.com/gorilla/mux"
	"kzapp/webapi/pkg"
	"net/http"
)

type GreetingHandler struct {}

func (h GreetingHandler) InitService(router *mux.Router) {
	router.HandleFunc("/", pkg.Chain(h.sayhello, pkg.Method("GET"), pkg.Logging()))
	router.HandleFunc("/greet", pkg.Chain(h.greet, pkg.Method("GET"), pkg.Logging()))
}

func (h GreetingHandler) sayhello(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello, you've requested: %s\n", r.URL.Path)
}

func (h GreetingHandler) greet(w http.ResponseWriter, r *http.Request) {
	userName := r.URL.Query().Get("name")
	fmt.Printf("userName: %s\n", userName)
	greeting := "Welcome Back!"

	response := GreetingResponse{
		Name:    userName,
		Message: &greeting,
	}

	pkg.JsonResponse(w, response)
}