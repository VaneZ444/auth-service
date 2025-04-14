package usecase

import (
	"time"

	"github.com/VaneZ444/golang-forum/auth-service/internal/entity"
	"github.com/VaneZ444/golang-forum/auth-service/internal/repository"
)

type AuthUseCase struct {
	repo      repository.UserRepository
	secretKey string
	tokenTTL  time.Duration
}

func NewAuthUseCase(repo repository.UserRepository, secretKey string, tokenTTL time.Duration) *AuthUseCase {
	return &AuthUseCase{repo: repo, secretKey: secretKey, tokenTTL: tokenTTL}
}

func (uc *AuthUseCase) Register(email, password string) (int64, error) {
	user, err := entity.NewUser(email, password)
	if err != nil {
		return 0, err // Возвращаем ошибки валидации
	}
	return uc.repo.SaveUser(user)
}
