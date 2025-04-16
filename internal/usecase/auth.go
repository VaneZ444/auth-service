package usecase

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/VaneZ444/auth-service/internal/entity"
	"github.com/VaneZ444/auth-service/internal/repository"
	"github.com/golang-jwt/jwt/v5"
)

type AuthUseCase struct {
	userRepo repository.UserRepository
	appRepo  repository.AppRepository
	secret   string
	tokenTTL time.Duration
	logger   *slog.Logger
}

// Custom JWT claims
type claims struct {
	UserID int64  `json:"user_id"`
	AppID  int32  `json:"app_id"`
	Role   string `json:"role"`
	jwt.RegisteredClaims
}

func NewAuthUseCase(
	userRepo repository.UserRepository,
	appRepo repository.AppRepository,
	secret string,
	tokenTTL time.Duration,
	logger *slog.Logger,
) *AuthUseCase {
	return &AuthUseCase{
		userRepo: userRepo,
		appRepo:  appRepo,
		secret:   secret,
		tokenTTL: tokenTTL,
		logger:   logger,
	}
}

// Регистрация пользователя
func (uc *AuthUseCase) Register(ctx context.Context, email, password string) (int64, error) {
	user, err := entity.NewUser(email, password, entity.RoleUser)
	if err != nil {
		return 0, fmt.Errorf("%w: %v", entity.ErrInvalidCredentials, err)
	}
	return uc.userRepo.SaveUser(ctx, user)
}

// Аутентификация
func (uc *AuthUseCase) Login(ctx context.Context, email, password string, appID int32) (string, error) {
	user, err := uc.userRepo.GetUserByEmail(ctx, email)
	if err != nil {
		return "", fmt.Errorf("%w: %v", entity.ErrUserNotFound, err)
	}

	if user.Status == entity.StatusBanned {
		return "", entity.ErrUserBanned
	}

	if !checkPasswordHash(password, user.PasswordHash) {
		return "", entity.ErrInvalidCredentials
	}

	return uc.generateJWT(user.ID, appID, user.Role)
}

// Создание администратора
func (uc *AuthUseCase) CreateAdmin(ctx context.Context, email, password string, requestingUserID int64) (int64, error) {
	isAdmin, err := uc.IsAdmin(ctx, requestingUserID)
	if err != nil || !isAdmin {
		return 0, entity.ErrAccessDenied
	}

	user, err := entity.NewUser(email, password, entity.RoleAdmin)
	if err != nil {
		return 0, err
	}
	return uc.userRepo.SaveUser(ctx, user)
}

// Проверка прав администратора
func (uc *AuthUseCase) IsAdmin(ctx context.Context, userID int64) (bool, error) {
	user, err := uc.userRepo.GetUserByID(ctx, userID)
	if err != nil {
		return false, fmt.Errorf("%w: %v", entity.ErrUserNotFound, err)
	}
	return user.Role == entity.RoleAdmin, nil
}

// Парсинг JWT токена
func (uc *AuthUseCase) ParseToken(tokenString string) (*claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(uc.secret), nil
	})

	if err != nil {
		return nil, fmt.Errorf("token parse error: %w", err)
	}

	if claims, ok := token.Claims.(*claims); ok && token.Valid {
		return claims, nil
	}

	return nil, entity.ErrInvalidToken
}

// Генерация JWT токена
func (uc *AuthUseCase) generateJWT(userID int64, appID int32, role entity.Role) (string, error) {
	now := time.Now()
	claims := &claims{
		UserID: userID,
		AppID:  appID,
		Role:   string(role),
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(uc.tokenTTL)),
			Issuer:    "auth-service",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(uc.secret))
}

// Проверка пароля (bcrypt)
func checkPasswordHash(password, hash string) bool {
	// Реализация с bcrypt
	// err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	// return err == nil
	return true // Заглушка для примера
}
