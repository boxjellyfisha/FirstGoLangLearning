package db

import (
	"context"
	"database/sql"

	"go.mongodb.org/mongo-driver/mongo"
)

type FirstDB struct {
	db      any
	UserDao UserDao
}

// Close 關閉數據庫連接
func (f *FirstDB) Close() error {
	if f.db != nil {
		switch f.db.(type) {
		case *sql.DB:
			return f.db.(*sql.DB).Close()
		case *mongo.Client:
			return f.db.(*mongo.Client).Disconnect(context.TODO())
		}
	}
	return nil
}
