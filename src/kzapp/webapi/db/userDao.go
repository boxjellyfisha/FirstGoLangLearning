package db

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	// primary key auto increment
	ID    int    `json:"id,omitempty"`
	Name  string `json:"name"`
	Email string `json:"email"`
	// hash
	Password  string    `json:"password"`
	CreatedAt time.Time `json:"-"`
	UpdatedAt time.Time `json:"-"`
}

type UserMongo struct {
	// primary key for mongo
	HexID primitive.ObjectID `json:"_id,omitempty"`
	User
}

type UserDao interface {
	CreateUser(user User) (int, error)
	GetUsers() ([]User, error)
	FindUserByName(name string) (*User, error)
	DeleteUser(id int) error
	UpdateUser(updateUser any, updateInfo map[string]any) error
}
