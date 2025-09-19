package usecase

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/VaneZ444/auth-service/internal/entity"
	"github.com/VaneZ444/auth-service/internal/jwt"
	"github.com/VaneZ444/auth-service/internal/pkg/utils"
	"github.com/VaneZ444/auth-service/internal/pkg/validator"
	"github.com/VaneZ444/auth-service/internal/repository"
)

type authUseCase struct {
	userRepo   repository.UserRepository
	appRepo    repository.AppRepository
	jwtService jwt.Service
	logger     *slog.Logger
}

func NewAuthUseCase(
	userRepo repository.UserRepository,
	appRepo repository.AppRepository,
	jwtService jwt.Service,
	logger *slog.Logger,
) *authUseCase {
	return &authUseCase{
		userRepo:   userRepo,
		appRepo:    appRepo,
		jwtService: jwtService,
		logger:     logger,
	}
}

func (uc *authUseCase) Register(ctx context.Context, email, nickname, password string) (int64, string, error) {
	if err := validator.ValidateEmail(email); err != nil {
		return 0, "", fmt.Errorf("%w: %v", ErrInvalidCredentials, err)
	}
	if len(password) < 8 {
		return 0, "", fmt.Errorf("%w: password too short", ErrInvalidCredentials)
	}
	if len(nickname) < 3 {
		return 0, "", fmt.Errorf("%w: nickname too short", ErrInvalidCredentials)
	}

	hashedPassword, err := utils.HashPassword(password)
	if err != nil {
		return 0, "", err
	}

	user := entity.NewUser(email, nickname, hashedPassword, entity.RoleUser)
	id, err := uc.userRepo.SaveUser(ctx, user)
	if err != nil {
		return 0, "", err
	}
	return id, nickname, nil
}
func (uc *authUseCase) CreateAdmin(ctx context.Context, email, password, nickname string, requesterID int64) (int64, error) {
	// Проверяем, что requester имеет права
	isAdmin, err := uc.IsAdmin(ctx, requesterID)
	if err != nil || !isAdmin {
		return 0, fmt.Errorf("permission denied")
	}

	hashedPassword, err := utils.HashPassword(password)
	if err != nil {
		return 0, err
	}

	user := entity.NewUser(email, hashedPassword, nickname, entity.RoleAdmin)
	return uc.userRepo.SaveUser(ctx, user)
}

func (uc *authUseCase) Login(ctx context.Context, email, password string, appID int32) (string, string, string, error) {
	user, err := uc.userRepo.GetUserByEmail(ctx, email)
	uc.logger.Debug("User data",
		slog.Int64("id", user.ID),
		slog.String("email", user.Email),
		slog.String("nickname", user.Nickname),
		slog.String("role", string(user.Role)), // Проверяем роль
		slog.String("status", string(user.Status)),
	)

	if err != nil {
		return "", "", "", fmt.Errorf("%w: %v", ErrUserNotFound, err)
	}

	if user.Status == entity.StatusBanned {
		return "", "", "", ErrUserBanned
	}

	if err := utils.CheckPasswordHash(password, user.Hash); err != nil {
		return "", "", "", ErrInvalidCredentials
	}

	token, err := uc.jwtService.GenerateToken(user.ID, appID, string(user.Role), string(user.Nickname))
	if err != nil {
		return "", "", "", err
	}

	return token, user.Nickname, string(user.Role), nil
}

func (uc *authUseCase) IsAdmin(ctx context.Context, userID int64) (bool, error) {
	user, err := uc.userRepo.GetUserByID(ctx, userID)
	if err != nil {
		return false, fmt.Errorf("%w: %v", ErrUserNotFound, err)
	}
	return user.Role == entity.RoleAdmin, nil
}
func (uc *authUseCase) RegisterAdmin(ctx context.Context, email, nickname, password string) (int64, string, error) {
	if err := validator.ValidateEmail(email); err != nil {
		return 0, "", fmt.Errorf("%w: %v", ErrInvalidCredentials, err)
	}
	if len(password) < 8 {
		return 0, "", fmt.Errorf("%w: password too short", ErrInvalidCredentials)
	}
	if len(nickname) < 3 {
		return 0, "", fmt.Errorf("%w: nickname too short", ErrInvalidCredentials)
	}

	hashedPassword, err := utils.HashPassword(password)
	if err != nil {
		return 0, "", err
	}

	user := entity.NewUser(email, nickname, hashedPassword, entity.RoleAdmin)
	id, err := uc.userRepo.SaveUser(ctx, user)
	if err != nil {
		return 0, "", err
	}

	return id, nickname, nil
}
