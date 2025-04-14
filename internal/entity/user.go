package entity

import (
	"errors"
	"regexp"
)

var (
	ErrInvalidEmail = errors.New("invalid email format")
	ErrWeakPassword = errors.New("password must be at least 8 characters")
)

type User struct {
	Email    string
	Password string
}

func NewUser(email, password string) (*User, error) {
	if !isValidEmail(email) {
		return nil, ErrInvalidEmail
	}
	if len(password) < 8 {
		return nil, ErrWeakPassword
	}
	return &User{Email: email, Password: password}, nil
}

func isValidEmail(email string) bool {
	const pattern = `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
	matched, _ := regexp.MatchString(pattern, email)
	return matched
}
