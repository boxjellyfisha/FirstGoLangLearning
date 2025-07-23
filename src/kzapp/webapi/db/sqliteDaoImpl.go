package db

import (
	"database/sql"
)

type UserSqliteDaoImpl struct {
	db *sql.DB
}

var _ UserDao = &UserSqliteDaoImpl{}

func (u *UserSqliteDaoImpl) CreateUser(user User) (int, error) {
	stmt, err := u.db.Prepare("INSERT INTO users (name, email, password) VALUES (?, ?, ?)")

	if err != nil {
		return -1, err
	}
	defer stmt.Close()

	result, err := stmt.Exec(user.Name, user.Email, user.Password)
	if err != nil {
		return -1, err
	}

	lastInsertId, err := result.LastInsertId()
	return int(lastInsertId), err
}

func (u *UserSqliteDaoImpl) FindUserByName(name string) (*User, error) {
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

func (u *UserSqliteDaoImpl) GetUsers() ([]User, error) {
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

func (u *UserSqliteDaoImpl) DeleteUser(id int) error {
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
