package webapi

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
)

type Middleware func(http.HandlerFunc) http.HandlerFunc

type UserResponse struct {
	Name    string  `json:"name"`
	Message *string `json:"message,omitempty"` // Optional string field
}

// Logging logs all requests with its path and the time it took to process
func Logging() Middleware {

	// Create a new Middleware
	return func(f http.HandlerFunc) http.HandlerFunc {

		// Define the http.HandlerFunc
		return func(w http.ResponseWriter, r *http.Request) {

			// Do middleware things
			start := time.Now()
			defer func() { log.Println(r.URL.Path, time.Since(start)) }()

			// Call the next middleware/handler in chain
			f(w, r)
		}
	}
}

// Method ensures that url can only be requested with a specific method, else returns a 400 Bad Request
func Method(m string) Middleware {

	// Create a new Middleware
	return func(f http.HandlerFunc) http.HandlerFunc {

		// Define the http.HandlerFunc
		return func(w http.ResponseWriter, r *http.Request) {

			// Do middleware things
			if r.Method != m {
				http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
				return
			}

			// Call the next middleware/handler in chain
			f(w, r)
		}
	}
}

// Chain applies middlewares to a http.HandlerFunc
func Chain(f http.HandlerFunc, middlewares ...Middleware) http.HandlerFunc {
	for _, m := range middlewares {
		f = m(f)
	}
	return f
}

var (
	// key must be 16, 24 or 32 bytes long (AES-128, AES-192 or AES-256)
	key   = []byte("super-secret-key")
	store = sessions.NewCookieStore(key)
)

// example: curl -s --cookie "cookie-name=MTc0OTIwNzc2OHxEWDhFQVFMX2dBQUJFQUVRQUFBbF80QUFBUVp6ZEhKcGJtY01Ed0FOWVhWMGFHVnVkR2xqWVhSbFpBUmliMjlzQWdJQUFRPT184yNJN4tqSV2k9vtr72fgHJiib5ZUTwe7aeatyygo2ro=" http://localhost:80/info
func userSecretData(w http.ResponseWriter, r *http.Request) {
	session, err := store.Get(r, "cookie-name")

	if err != nil {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}

	fmt.Fprintln(w, "authenticated ", session.Values["authenticated"])

	// Check if user is authenticated
	if auth, isReturnValueAsBool := session.Values["authenticated"].(bool); !isReturnValueAsBool || !auth {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}

	// Print secret message
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

	encodeResponse(w, response)
}

// example: curl -s --cookie "cookie-name=MTc0OTIwNzc2OHxEWDhFQVFMX2dBQUJFQUVRQUFBbF80QUFBUVp6ZEhKcGJtY01Ed0FOWVhWMGFHVnVkR2xqWVhSbFpBUmliMjlzQWdJQUFRPT184yNJN4tqSV2k9vtr72fgHJiib5ZUTwe7aeatyygo2ro=" -I -XPOST http://localhost:80/logout
func logout(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "cookie-name")

	// Revoke users authentication
	session.Values["authenticated"] = false
	session.Save(r, w)
}

func greet(w http.ResponseWriter, r *http.Request) {
	userName := r.URL.Query().Get("name")
	fmt.Printf("userName: %s\n", userName)
	greeting := "Welcome Back!"

	response := UserResponse{
		Name:    userName,
		Message: &greeting,
	}

	encodeResponse(w, response)
}

type Plusable struct {
	Num1 int `json:"a"`
	Num2 int `json:"b"`
}
type Sum struct {
	Total int `json:"sum"`
}

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
	encodeResponse(w, response)
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

	encodeResponse(w, response)
}

func encodeResponse(w http.ResponseWriter, response any) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func Server() {
	// Create the router
	router := mux.NewRouter()

	handleFun := func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello, you've requested: %s\n", r.URL.Path)
	}
	// Handle incoming request
	router.HandleFunc("/", Chain(handleFun, Method("GET"), Logging()))
	router.HandleFunc("/greet", Chain(greet, Method("GET"), Logging()))
	router.HandleFunc("/square/{num}", Chain(square, Method("GET"), Logging()))
	router.HandleFunc("/add", Chain(add, Method("POST"), Logging()))

	router.HandleFunc("/login", Chain(login, Method("POST"), Logging()))
	router.HandleFunc("/logout", Chain(logout, Method("POST"), Logging()))
	router.HandleFunc("/info", Chain(userSecretData, Method("GET"), Logging()))

	// Listen for HTTP Connections, Registering a Request Handler
	http.ListenAndServe(":80", router)
}
