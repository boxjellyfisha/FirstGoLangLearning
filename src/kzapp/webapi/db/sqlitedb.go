package db

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

func InitSqliteDB(filepath string) *sql.DB {
	db, err := sql.Open("sqlite3", filepath)
	if err != nil {
		log.Fatal(err)
	}
	create := `
    CREATE TABLE IF NOT EXISTS users (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        name TEXT NOT NULL UNIQUE,
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

func NewFirstSqliteDB(filepath string) *FirstDB {
	db := InitSqliteDB(filepath)
	userDao := &UserSqliteDaoImpl{db: db}
	return &FirstDB{db: db, UserDao: userDao}
}
