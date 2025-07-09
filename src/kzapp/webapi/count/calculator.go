package count

import (
	"encoding/json"
	"kzapp/webapi/pkg"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/mux"
)

type Calculator struct{}

var _ pkg.Handler = (*Calculator)(nil)

func (c Calculator) InitService(router *mux.Router) {
	router.HandleFunc("/square/{num}", pkg.Chain(c.square, pkg.Method("GET"), pkg.Logging()))
	router.HandleFunc("/add", pkg.Chain(c.add, pkg.Method("POST"), pkg.Logging()))
}

func (c Calculator) InitServiceGin(router *gin.Engine) {
	router.GET("/square/:num", pkg.ChainGin(c.square), gin.Logger())
	router.POST("/add", pkg.ChainGin(c.add), gin.Logger())
}

// example: curl -XPOST http://localhost:80/add -d '{"a": 3, "b": 5}'
func (c Calculator) add(w http.ResponseWriter, r *http.Request) {

	//decode the body and create response
	var req Plusable
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}
	// encode the sum to response
	response := Sum{req.Num1 + req.Num2}
	pkg.JsonResponse(w, response)
}

// example: curl http://localhost:80/square/9
func (c Calculator) square(w http.ResponseWriter, r *http.Request) {
	paths := strings.Split(r.URL.Path, "/")
	lastIndex := len(paths) - 1
	lastPath := paths[lastIndex]
	var request, err = strconv.Atoi(lastPath)
	if err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	response := Sum{
		Total: request * request,
	}

	pkg.JsonResponse(w, response)
}
