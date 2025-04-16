package repository

import (
	"context"

	"github.com/VaneZ444/auth-service/internal/entity"
)

// Интерфейсы для работы с БД.
type (
	UserRepository interface {
		SaveUser(ctx context.Context, user *entity.User) (userID int64, err error)
		GetUserByEmail(ctx context.Context, email string) (*entity.User, error)
		IsAdmin(ctx context.Context, userID int64) (bool, error)
		GetUserByID(ctx context.Context, id int64) (*entity.User, error)
	}

	AppRepository interface {
		GetAppByID(ctx context.Context, appID int32) (name string, secret string, err error)
	}
)
