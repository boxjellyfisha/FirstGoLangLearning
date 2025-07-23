package pkg

import (
	"github.com/gin-gonic/gin"
	"github.com/gorilla/mux"
)

type Handler interface {
	InitService(route *mux.Router)
	InitServiceGin(route *gin.Engine)
}
