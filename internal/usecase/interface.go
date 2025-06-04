package usecase

import "context"

type AuthUseCase interface {
	Register(ctx context.Context, email, password string) (int64, error)
	CreateAdmin(ctx context.Context, email, password string, createdBy int64) (int64, error)
	Login(ctx context.Context, email, password string, appID int32) (string, error)
	IsAdmin(ctx context.Context, userID int64) (bool, error)
}
