package webapi

import (
	"github.com/gorilla/mux"

	"kzapp/webapi/hello"
	"kzapp/webapi/count"
	"kzapp/webapi/user"
)

type Handler interface {
	InitService(route *mux.Router)
}

func CreateHandler() []Handler {
	return []Handler{
		hello.GreetingHandler{},
		count.Calculator{},
		user.CreateUserHandler([]byte("super-secret-key")),
	}
}


