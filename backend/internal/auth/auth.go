package auth

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

var jwtKey = []byte("f5ee392ecd4fdb155a2b3e7b5551a143227b7cc424fcf139bc48b9bcd21a361fff1b58a66b60719fac5d4eb90564f91d8a65461cd4f5b9b82144bb9246a45fb746bc8390e417374dda82d17762b765d2dcd21e4cd10885c0ec2627ed0eace533cb66954af116753673a8a715e1ec37158ab9c5738dd4bd17c9dc3c356b7cfb08")

func GenerateJWT(userID int, username string) (string, error) {
	claims := jwt.MapClaims{
		"user_id":  userID,
		"username": username,
		"exp":      time.Now().Add(time.Hour * 72).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtKey)
}
