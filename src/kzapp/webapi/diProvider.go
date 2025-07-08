package webapi

import (
	"errors"
	"fmt"
	"log"
	"path/filepath"
	"sync"

	"github.com/gorilla/mux"

	"kzapp/webapi/count"
	"kzapp/webapi/db"
	"kzapp/webapi/hello"
	"kzapp/webapi/user"
)

type Handler interface {
	InitService(route *mux.Router)
}

func CreateHandler() []Handler {
	userDao, err := ProvideUserDao()
	if err != nil {
		log.Fatal(err)
		return []Handler{
			hello.GreetingHandler{},
			count.Calculator{},
		}
	}
	return []Handler{
		hello.GreetingHandler{},
		count.Calculator{},
		user.CreateUserHandler([]byte("super-secret-key"), userDao),
	}
}

var (
	once             sync.Once
	userDaoSingleton db.UserDao
	firstDBSingleton *db.FirstDB
	initError        error
)

func ProvideUserDao() (db.UserDao, error) {
	once.Do(func() {
		dbPath := filepath.Join("db", "test.db")
		fmt.Println(dbPath)

		firstDB := db.NewFirstDB(dbPath)
		if firstDB == nil {
			initError = errors.New("failed to initialize database")
			return
		}
		firstDBSingleton = firstDB
		userDaoSingleton = firstDB.UserDao
	})
	return userDaoSingleton, initError
}

// GetInitError 返回初始化錯誤（如果有的話）
func GetInitError() error {
	return initError
}

// Shutdown 優雅關閉所有資源
func Shutdown() error {
	if firstDBSingleton != nil {
		log.Println("正在關閉數據庫連接...")
		return firstDBSingleton.Close()
	}
	return nil
}

func newFunction(userDao db.UserDao) {
	// go userDao.CreateUser(db.User{
	// 	Name:     "John Doe",
	// 	Email:    "john.doe@example.com",
	// 	Password: "password",
	// })
	users := make([]db.User, 2)
	go func() {
		user, _ := userDao.GetUsers()
		copy(users, user)
	}()
	defer log.Println(users)
}
