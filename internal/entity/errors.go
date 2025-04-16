package entity

import "errors"

//доменные ошибки
var (
	ErrAppNotFound        = errors.New("app not found")
	ErrInvalidToken       = errors.New("invalid token")
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrUserBanned         = errors.New("user banned")
	ErrUserNotFound       = errors.New("user not found")
	ErrEmailExists        = errors.New("email already exists")
	ErrAccessDenied       = errors.New("access denied")
	ErrTokenExpired       = errors.New("token expired")
)
