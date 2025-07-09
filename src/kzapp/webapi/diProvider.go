package webapi

import (
	"errors"
	"log"
	"path/filepath"
	"sync"

	"kzapp/webapi/count"
	"kzapp/webapi/db"
	"kzapp/webapi/hello"
	"kzapp/webapi/pkg"
	"kzapp/webapi/user"
)

func CreateHandler() []pkg.Handler {
	userDao, err := ProvideUserDao()
	if err != nil {
		log.Fatal(err)
		return []pkg.Handler{
			hello.GreetingHandler{},
			count.Calculator{},
		}
	}
	return []pkg.Handler{
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
		currentDir, err := pkg.GetCurrentDir()
		if err != nil {
			initError = errors.New("failed to get current directory")
			return
		}
		dbPath := filepath.Join(currentDir, "db", "test.db")
		log.Println("資料庫路徑:", dbPath)

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
