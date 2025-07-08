package user

import (
	"encoding/json"
	"fmt"
	"kzapp/webapi/db"
	"kzapp/webapi/pkg"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
)

type UserHandler struct {
	// key must be 16, 24 or 32 bytes long (AES-128, AES-192 or AES-256)
	// example: []byte("super-secret-key")
	key     []byte
	store   *sessions.CookieStore
	userDao db.UserDao
}

func CreateUserHandler(key []byte, userDao db.UserDao) UserHandler {
	return UserHandler{
		key:     key,
		store:   sessions.NewCookieStore(key),
		userDao: userDao,
	}
}

func (h UserHandler) InitService(router *mux.Router) {
	router.HandleFunc("/signup", pkg.Chain(h.signup, pkg.Method("POST"), pkg.Logging()))

	router.HandleFunc("/login", pkg.Chain(h.login, pkg.Method("POST"), pkg.Logging()))
	router.HandleFunc("/logout", pkg.Chain(h.logout, pkg.Method("POST"), pkg.Logging()))
	router.HandleFunc("/info", pkg.Chain(h.userSecretData, pkg.Method("GET"), pkg.Logging()))
}

// example: curl -s --cookie "cookie-name=MTc0OTIwNzc2OHxEWDhFQVFMX2dBQUJFQUVRQUFBbF80QUFBUVp6ZEhKcGJtY01Ed0FOWVhWMGFHVnVkR2xqWVhSbFpBUmliMjlzQWdJQUFRPT184yNJN4tqSV2k9vtr72fgHJiib5ZUTwe7aeatyygo2ro=" http://localhost:80/info
func (h UserHandler) userSecretData(w http.ResponseWriter, r *http.Request) {
	session, err := h.store.Get(r, "cookie-name")

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
	fmt.Fprintln(w, "The cake is a lie!")
}

// example: curl -XPOST http://localhost:80/signup -d '{"name":"kzzz", "email":"kzzz@gmail.com", "password":"123456"}'
func (h UserHandler) signup(w http.ResponseWriter, r *http.Request) {
	var req db.User
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil || req.Name == "" {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	id, err := h.userDao.CreateUser(req)
	if err != nil {
		http.Error(w, "Failed to create user", http.StatusInternalServerError)
		return
	}

	log.Println("created user id: ", id)

	// 返回成功響應
	response := map[string]interface{}{
		"message": "User created successfully",
		"id":      id,
	}
	pkg.JsonResponse(w, response)
}

// func createNewUser(userDao db.UserDao) {
// 	go userDao.CreateUser(db.User{
// 		Name:     "John Doe",
// 		Email:    "john.doe@example.com",
// 		Password: "password",
// 	})
// 	users := make([]db.User, 2)
// 	go func() {
// 		user, _ := userDao.GetUsers()
// 		copy(users, user)
// 	}()
// 	defer log.Println(users)
// }

// example: curl -s -I -XPOST http://localhost:80/login -d '{"name":"kzzz"}'
func (h UserHandler) login(w http.ResponseWriter, r *http.Request) {
	session, _ := h.store.Get(r, "cookie-name")

	// Authentication goes here
	//decode the body and create response
	var req UserResponse
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil || req.Name == "" {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	user, err := h.userDao.FindUserByName(req.Name)
	if err != nil || user == nil {
		http.Error(w, "Unknown user", http.StatusBadRequest)
		return
	}

	// Set user as authenticated
	session.Values["authenticated"] = true

	session.Save(r, w)

	// create response
	greeting := "Hello, Web API with Go!"
	response := UserResponse{
		Name:    req.Name,
		Message: &greeting,
	}

	pkg.JsonResponse(w, response)
}

// example: curl -s --cookie "cookie-name=MTc0OTIwNzc2OHxEWDhFQVFMX2dBQUJFQUVRQUFBbF80QUFBUVp6ZEhKcGJtY01Ed0FOWVhWMGFHVnVkR2xqWVhSbFpBUmliMjlzQWdJQUFRPT184yNJN4tqSV2k9vtr72fgHJiib5ZUTwe7aeatyygo2ro=" -I -XPOST http://localhost:80/logout
func (h UserHandler) logout(w http.ResponseWriter, r *http.Request) {
	session, _ := h.store.Get(r, "cookie-name")

	// Revoke users authentication
	session.Values["authenticated"] = false
	session.Save(r, w)
}
