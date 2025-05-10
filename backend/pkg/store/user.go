package store

import (
	"database/sql"

	_ "github.com/go-sql-driver/mysql"
)

var DB *sql.DB

type User struct {
	ID           int
	Username     string
	PasswordHash string
}

func InitDB(dataSourceName string) error {
	var err error
	DB, err = sql.Open("mysql", dataSourceName)
	if err != nil {
		return err
	}
	return DB.Ping()
}

func SaveUser(user User) error {
	_, err := DB.Exec("INSERT INTO users (username, password_hash) VALUES (?, ?)", user.Username, user.PasswordHash)
	return err
}

func GetUserByUsername(username string) (*User, error) {
	row := DB.QueryRow("SELECT id, username, password_hash FROM users WHERE username = ?", username)
	var user User
	err := row.Scan(&user.ID, &user.Username, &user.PasswordHash)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return &user, err
}
