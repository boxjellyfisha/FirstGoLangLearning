package db

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

func InitDB(filepath string) *sql.DB {
	db, err := sql.Open("sqlite3", filepath)
	if err != nil {
		log.Fatal(err)
	}
	create := `
    CREATE TABLE IF NOT EXISTS users (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        name TEXT NOT NULL,
        email TEXT NOT NULL,
        password TEXT NOT NULL,
        created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
        updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
    );`
	if _, err := db.Exec(create); err != nil {
		log.Fatalf("table create error: %v", err)
	}
	return db
}

func NewFirstDB(filepath string) *FirstDB {
	db := InitDB(filepath)
	userDao := &UserDaoImpl{db: db}
	return &FirstDB{db: db, UserDao: userDao}
}

type FirstDB struct {
	db      *sql.DB
	UserDao UserDao
}

// Close 關閉數據庫連接
func (f *FirstDB) Close() error {
	if f.db != nil {
		return f.db.Close()
	}
	return nil
}
