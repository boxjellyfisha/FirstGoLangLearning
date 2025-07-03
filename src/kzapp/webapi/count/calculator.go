package count

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"kzapp/webapi/pkg"
	"net/http"
	"strconv"
	"strings"
)

// example: curl -XPOST http://localhost:80/add -d '{"a": 3, "b": 5}'
func add(w http.ResponseWriter, r *http.Request) {

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
func square(w http.ResponseWriter, r *http.Request) {
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

func InitCaculator(router *mux.Router) {
	router.HandleFunc("/square/{num}", pkg.Chain(square, pkg.Method("GET"), pkg.Logging()))
	router.HandleFunc("/add", pkg.Chain(add, pkg.Method("POST"), pkg.Logging()))
}
