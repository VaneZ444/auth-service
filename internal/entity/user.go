package entity

import (
	"errors"
	"regexp"

	"golang.org/x/crypto/bcrypt"
)

type Role string

const (
	RoleUser  Role = "user"
	RoleAdmin Role = "admin"
)

type Status string

const (
	StatusActive Status = "active"
	StatusBanned Status = "banned"
)

type User struct {
	ID           int64
	Email        string
	Password     string // Нехешированный пароль
	PasswordHash string // Хешированный
	Role         Role   // 'user' или 'admin'
	Status       Status // 'active' или 'banned'
}

func NewUser(email, password string, role Role) (*User, error) {
	if err := validateEmail(email); err != nil {
		return nil, err
	}
	if len(password) < 8 {
		return nil, errors.New("password too short")
	}
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	return &User{
		Email:        email,
		Password:     password,
		PasswordHash: string(hashedPassword),
		Role:         role,
		Status:       StatusActive,
	}, nil
}

func validateEmail(email string) error {
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	if !emailRegex.MatchString(email) {
		return errors.New("invalid email")
	}
	return nil
}
