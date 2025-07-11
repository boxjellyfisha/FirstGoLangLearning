package db

import (
	"database/sql"
	"time"
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

type UserDao interface {
	CreateUser(user User) (int64, error)
	GetUsers() ([]User, error)
	FindUserByName(name string) (*User, error)
	DeleteUser(id int) error
}

type UserDaoImpl struct {
	db *sql.DB
}

func (u *UserDaoImpl) CreateUser(user User) (int64, error) {
	stmt, err := u.db.Prepare("INSERT INTO users (name, email, password) VALUES (?, ?, ?)")

	if err != nil {
		return -1, err
	}
	defer stmt.Close()

	result, err := stmt.Exec(user.Name, user.Email, user.Password)
	if err != nil {
		return -1, err
	}

	return result.LastInsertId()
}

func (u *UserDaoImpl) FindUserByName(name string) (*User, error) {
	rows, err := u.db.Query("SELECT * FROM users WHERE name = ?", name)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var user User
		err = rows.Scan(&user.ID, &user.Name, &user.Email, &user.Password, &user.CreatedAt, &user.UpdatedAt)
		if err != nil {
			return nil, err
		} else {
			return &user, nil
		}
	}
	return nil, nil
}

func (u *UserDaoImpl) GetUsers() ([]User, error) {
	rows, err := u.db.Query("SELECT * FROM users")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	users := []User{}
	for rows.Next() {
		var user User
		err = rows.Scan(&user.ID, &user.Name, &user.Email, &user.Password, &user.CreatedAt, &user.UpdatedAt)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}
	return users, nil
}

func (u *UserDaoImpl) DeleteUser(id int) error {
	stmt, err := u.db.Prepare("DELETE FROM users WHERE id = ?")
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(id)
	if err != nil {
		return err
	}

	return nil
}
