package usecase

import "context"

type AuthUseCase interface {
	Register(ctx context.Context, email, nickname, password string) (int64, string, error)
	CreateAdmin(ctx context.Context, email, password, nickname string, createdBy int64) (int64, error)
	Login(ctx context.Context, email, password string, appID int32) (string, string, error)
	IsAdmin(ctx context.Context, userID int64) (bool, error)
	RegisterAdmin(ctx context.Context, email, nickname, password string) (int64, string, error)
}
