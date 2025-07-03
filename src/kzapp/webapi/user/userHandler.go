package user

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	"net/http"
	"kzapp/webapi/pkg"
)

var (
	// key must be 16, 24 or 32 bytes long (AES-128, AES-192 or AES-256)
	key   = []byte("super-secret-key")
	store = sessions.NewCookieStore(key)
)

// example: curl -s --cookie "cookie-name=MTc0OTIwNzc2OHxEWDhFQVFMX2dBQUJFQUVRQUFBbF80QUFBUVp6ZEhKcGJtY01Ed0FOWVhWMGFHVnVkR2xqWVhSbFpBUmliMjlzQWdJQUFRPT184yNJN4tqSV2k9vtr72fgHJiib5ZUTwe7aeatyygo2ro=" http://localhost:80/info
func userSecretData(w http.ResponseWriter, r *http.Request) {
	session, err := store.Get(r, "cookie-name")

	if err != nil {
		w.WriteHeader(http.StatusForbidden)
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}

	// Check if user is authenticated
	if auth, isReturnValueAsBool := session.Values["authenticated"].(bool); !isReturnValueAsBool || !auth {
		w.WriteHeader(http.StatusForbidden)
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}

	// Print secret message
	w.WriteHeader(http.StatusOK)
	fmt.Fprintln(w, "The cake is a lie!")
}

// example: curl -s -I -XPOST http://localhost:80/login -d '{"name":"kzzz"}'
func login(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "cookie-name")

	// Authentication goes here
	// ...

	// Set user as authenticated
	session.Values["authenticated"] = true

	session.Save(r, w)

	//decode the body and create response
	var req UserResponse
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil || req.Name == "" {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	// create response
	greeting := "Hello, Web API with Go!"
	response := UserResponse{
		Name:    req.Name,
		Message: &greeting,
	}

	pkg.JsonResponse(w, response)
}

// example: curl -s --cookie "cookie-name=MTc0OTIwNzc2OHxEWDhFQVFMX2dBQUJFQUVRQUFBbF80QUFBUVp6ZEhKcGJtY01Ed0FOWVhWMGFHVnVkR2xqWVhSbFpBUmliMjlzQWdJQUFRPT184yNJN4tqSV2k9vtr72fgHJiib5ZUTwe7aeatyygo2ro=" -I -XPOST http://localhost:80/logout
func logout(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "cookie-name")

	// Revoke users authentication
	session.Values["authenticated"] = false
	session.Save(r, w)
}

func InitUserService(router *mux.Router) {
	router.HandleFunc("/login", pkg.Chain(login, pkg.Method("POST"), pkg.Logging()))
	router.HandleFunc("/logout", pkg.Chain(logout, pkg.Method("POST"), pkg.Logging()))
	router.HandleFunc("/info", pkg.Chain(userSecretData, pkg.Method("GET"), pkg.Logging()))
}
