package repository

import "github.com/VaneZ444/golang-forum/auth-service/internal/entity"

// UserRepository — интерфейс, который реализует PostgreSQL.
type UserRepository interface {
	SaveUser(user *entity.User) (userID int64, err error)
	GetUserByEmail(email string) (*entity.User, error)
	IsAdmin(userID int64) (bool, error)
}
