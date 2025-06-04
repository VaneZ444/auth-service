package utils

import (
	"errors"

	"golang.org/x/crypto/bcrypt"
)

var (
	ErrHashMismatch = errors.New("password does not match hash")
)

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

func CheckPasswordHash(password, hash string) error {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	if err != nil {
		return ErrHashMismatch
	}
	return nil
}
