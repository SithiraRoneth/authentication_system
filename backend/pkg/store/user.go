package store

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
)

var DB *sql.DB

func InitDB(dataSourceName string) error {
	var err error
	DB, err = sql.Open("mysql", dataSourceName)
	if err != nil {
		return err
	}
	return DB.Ping()
}

type User struct {
	ID           int
	Username     string
	PasswordHash string
}

func SaveUser(user User) error {
	query := "INSERT INTO users (username, password_hash) VALUES (?, ?)"
	_, err := DB.Exec(query, user.Username, user.PasswordHash)
	if err != nil {
		return fmt.Errorf("could not insert user: %v", err)
	}
	return nil
}

func GetUserByUsername(username string) (*User, error) {
	query := "SELECT id, username, password_hash FROM users WHERE username = ?"
	row := DB.QueryRow(query, username)

	var user User
	err := row.Scan(&user.ID, &user.Username, &user.PasswordHash)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}
