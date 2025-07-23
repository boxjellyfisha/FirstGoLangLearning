package user

import (
	"encoding/json"
	"fmt"
	"kzapp/webapi/db"
	"kzapp/webapi/pkg"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
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

// 第一行：var _ pkg.Handler = (*UserHandler)(nil)
// 這行是用來在編譯時期檢查 UserHandler 是否有實作 pkg.Handler 介面，
// 這裡用的是 UserHandler 的指標型別（*UserHandler），代表你通常會用指標方式傳遞 handler。

// 第二行：var _ pkg.Handler = UserHandler{}
// 這行則是用來檢查 UserHandler 的值型別（非指標）是否有實作 pkg.Handler 介面，
// 這代表你也可以直接用值型別來傳遞 handler。

// 兩者差異：
//   - (*UserHandler)(nil) 是指標型別，UserHandler{} 是值型別。
//   - 這兩行都只是型別檢查，不會產生任何執行時的程式碼。
//   - 如果你的方法接收器是指標（func (h *UserHandler) ...），只有第一行會通過； -> 如果實作本身會更改 UserHandler 內的成員變數，則使用指標型別
//     如果是值接收器（func (h UserHandler) ...），兩行都會通過。 -> 如果實作本身不會更改 UserHandler 內的成員變數，則使用值型別
var _ pkg.Handler = (*UserHandler)(nil)

var _ pkg.Handler = UserHandler{}

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

func (h UserHandler) InitServiceGin(router *gin.Engine) {
	router.POST("/signup", pkg.ChainGin(h.signup), gin.Logger())
	router.POST("/login", pkg.ChainGin(h.login), gin.Logger())
	router.POST("/logout", pkg.ChainGin(h.logout), gin.Logger())
	router.GET("/info", pkg.ChainGin(h.userSecretData), gin.Logger())
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

// Signup godoc
//
//	@Summary		signup a new user
//	@Description	create new user into mongodb
//	@Tags			account
//	@Accept			json
//	@Produce		json
//	@Param			user body db.User true "User" example({"name":"kzzz", "email":"kzzz@gmail.com", "password":"123456"})
//	@Success		200		{object}	SignupResponse	"success"		example({"message":"User created successfully", "id":1})
//	@Failure		400		{string}	string			"error"		example("Invalid JSON")
//	@Failure		404		{string}	string			"error"		example("Failed to create user")
//	@Failure		500		{string}	string			"error"		example("Failed to create user")
//	@Router			/signup [POST]
//
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
	response := SignupResponse{
		Message: "User created successfully",
		ID:      id,
	}
	pkg.JsonResponse(w, response)
}

type SignupResponse struct {
	Message string `json:"message"`
	ID      int    `json:"id"`
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

// Login godoc
//
//	@Summary		user login
//	@Description	check the user exist and set the session
//	@Tags			account
//	@Accept			json
//	@Produce		json
//	@Param			user	body		UserResponse	true	"User"
//	@Success		200		{object}	UserResponse
//	@Failure		400		{string}	string
//	@Failure		500		{string}	string
//	@Router			/login [POST]
//
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

// Login godoc
//
//	@Summary		user logout
//	@Description	clear the user session
//	@Tags			account
//	@Accept			json
//	@Produce		json
//	@Param			cookie	header		string	true	"cookie-name"
//	@Success		200		{any}		any
//	@Failure		400		{string}	string
//	@Failure		500		{string}	string
//	@Router			/logout [GET]
//
// example: curl -s --cookie "cookie-name=MTc0OTIwNzc2OHxEWDhFQVFMX2dBQUJFQUVRQUFBbF80QUFBUVp6ZEhKcGJtY01Ed0FOWVhWMGFHVnVkR2xqWVhSbFpBUmliMjlzQWdJQUFRPT184yNJN4tqSV2k9vtr72fgHJiib5ZUTwe7aeatyygo2ro=" -I -XPOST http://localhost:80/logout
func (h UserHandler) logout(w http.ResponseWriter, r *http.Request) {
	session, _ := h.store.Get(r, "cookie-name")

	// Revoke users authentication
	session.Values["authenticated"] = false
	session.Save(r, w)
}
