package store

import (
	"database/sql"
	"log"

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
		log.Printf("Error opening the database: %v", err)
		return err
	}

	err = DB.Ping()
	if err != nil {
		log.Printf("Error pingin+g the database: %v", err)
		return err
	}

	log.Println("Successfully connected to the database")
	return nil
}

func SaveUser(user User) error {
	_, err := DB.Exec("INSERT INTO users (username, password_hash) VALUES (?, ?)", user.Username, user.PasswordHash)
	if err != nil {
		log.Printf("Error saving user: %v", err)
	}
	return err
}

func GetUserByUsername(username string) (*User, error) {
	row := DB.QueryRow("SELECT id, username, password_hash FROM users WHERE username = ?", username)
	var user User
	err := row.Scan(&user.ID, &user.Username, &user.PasswordHash)
	if err == sql.ErrNoRows {
		return nil, nil
	} else if err != nil {
		log.Printf("Error querying user by username: %v", err)
	}
	return &user, err
}

func GetUserByID(id int) (*User, error) {

	row := DB.QueryRow("SELECT id, username, password_hash FROM users WHERE id = ?", id)
	var user User
	err := row.Scan(&user.ID, &user.Username, &user.PasswordHash)
	if err == sql.ErrNoRows {
		return nil, nil
	} else if err != nil {
		log.Printf("Error querying user by ID: %v", err)
	}
	return &user, err
}
