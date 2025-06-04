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

func (uc *authUseCase) Register(ctx context.Context, email, password string) (int64, error) {
	if err := validator.ValidateEmail(email); err != nil {
		return 0, fmt.Errorf("%w: %v", ErrInvalidCredentials, err)
	}
	if len(password) < 8 {
		return 0, fmt.Errorf("%w: password too short", ErrInvalidCredentials)
	}
	hashedPassword, err := utils.HashPassword(password)
	if err != nil {
		return 0, err
	}

	user := entity.NewUser(email, hashedPassword, entity.RoleUser)
	return uc.userRepo.SaveUser(ctx, user)
}

func (uc *authUseCase) Login(ctx context.Context, email, password string, appID int32) (string, error) {
	user, err := uc.userRepo.GetUserByEmail(ctx, email)
	if err != nil {
		return "", fmt.Errorf("%w: %v", ErrUserNotFound, err)
	}

	if user.Status == entity.StatusBanned {
		return "", ErrUserBanned
	}

	if err := utils.CheckPasswordHash(password, user.Hash); err != nil {
		return "", ErrInvalidCredentials
	}

	token, err := uc.jwtService.GenerateToken(user.ID, appID, string(user.Role))
	if err != nil {
		return "", err
	}

	return token, nil
}

func (uc *authUseCase) CreateAdmin(ctx context.Context, email, password string, requestingUserID int64) (int64, error) {
	isAdmin, err := uc.IsAdmin(ctx, requestingUserID)
	if err != nil || !isAdmin {
		return 0, ErrAccessDenied
	}

	if err := validator.ValidateEmail(email); err != nil {
		return 0, fmt.Errorf("%w: %v", ErrInvalidCredentials, err)
	}
	if len(password) < 8 {
		return 0, fmt.Errorf("%w: password too short", ErrInvalidCredentials)
	}
	hashedPassword, err := utils.HashPassword(password)
	if err != nil {
		return 0, err
	}

	user := entity.NewUser(email, hashedPassword, entity.RoleAdmin)
	return uc.userRepo.SaveUser(ctx, user)
}

func (uc *authUseCase) IsAdmin(ctx context.Context, userID int64) (bool, error) {
	user, err := uc.userRepo.GetUserByID(ctx, userID)
	if err != nil {
		return false, fmt.Errorf("%w: %v", ErrUserNotFound, err)
	}
	return user.Role == entity.RoleAdmin, nil
}
