package handler

import (
	"context"
	"errors"
	"log/slog"

	"github.com/VaneZ444/auth-service/internal/entity"
	"github.com/VaneZ444/auth-service/internal/usecase"

	ssov1 "github.com/VaneZ444/golang-forum-protos/gen/go/sso"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

type AuthHandler struct {
	authUC                        *usecase.AuthUseCase
	logger                        *slog.Logger
	ssov1.UnimplementedAuthServer // Обязательно для gRPC
}

func NewAuthHandler(authUC *usecase.AuthUseCase, logger *slog.Logger) *AuthHandler {
	return &AuthHandler{
		authUC: authUC,
		logger: logger,
	}
}

// Register — регистрация обычного пользователя.
func (h *AuthHandler) Register(ctx context.Context, req *ssov1.RegisterRequest) (*ssov1.RegisterResponse, error) {
	const op = "handler.Register"
	h.logger.Info("register request", slog.String("op", op), slog.String("email", req.Email))

	// Делегируем бизнес-логику UseCase
	userID, err := h.authUC.Register(ctx, req.Email, req.Password)
	if err != nil {
		h.logger.Error("registration failed", slog.String("op", op), slog.String("err", err.Error()))

		if errors.Is(err, entity.ErrInvalidCredentials) {
			return nil, status.Error(codes.InvalidArgument, "invalid email or password")
		}
		return nil, status.Error(codes.Internal, "internal error")
	}

	return &ssov1.RegisterResponse{UserId: userID}, nil
}

// CreateAdmin — создание админа (требует прав суперадмина).
func (h *AuthHandler) CreateAdmin(ctx context.Context, req *ssov1.CreateAdminRequest) (*ssov1.CreateAdminResponse, error) {
	const op = "handler.CreateAdmin"
	h.logger.Info("create admin request", slog.String("op", op))

	// Проверяем, что вызывающий — админ (через gRPC-метаданные или токен)
	callerRole, err := h.getCallerRole(ctx)
	if err != nil || callerRole != entity.RoleAdmin {
		return nil, status.Error(codes.PermissionDenied, "admin rights required")
	}
	callerUserID, err := h.getCallerUserID(ctx)
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, "unauthenticated")
	}
	userID, err := h.authUC.CreateAdmin(ctx, req.Email, req.Password, callerUserID)

	if err != nil {
		h.logger.Error("create admin failed", slog.String("op", op), slog.String("err", err.Error()))
		return nil, status.Error(codes.Internal, "internal error")
	}

	return &ssov1.CreateAdminResponse{UserId: userID}, nil
}

// Login — аутентификация пользователя.
func (h *AuthHandler) Login(ctx context.Context, req *ssov1.LoginRequest) (*ssov1.LoginResponse, error) {
	const op = "handler.Login"
	h.logger.Info("login request", slog.String("op", op), slog.String("email", req.Email))

	token, err := h.authUC.Login(ctx, req.Email, req.Password, req.AppId)
	if err != nil {
		h.logger.Error("login failed", slog.String("op", op), slog.String("err", err.Error()))

		switch {
		case errors.Is(err, entity.ErrInvalidCredentials):
			return nil, status.Error(codes.Unauthenticated, "invalid credentials")
		case errors.Is(err, entity.ErrUserBanned):
			return nil, status.Error(codes.PermissionDenied, "user banned")
		default:
			return nil, status.Error(codes.Internal, "internal error")
		}
	}

	return &ssov1.LoginResponse{Token: token}, nil
}

// IsAdmin — проверка прав администратора.
func (h *AuthHandler) IsAdmin(ctx context.Context, req *ssov1.IsAdminRequest) (*ssov1.IsAdminResponse, error) {
	const op = "handler.IsAdmin"
	h.logger.Debug("is_admin check", slog.String("op", op), slog.Int64("user_id", req.UserId))

	isAdmin, err := h.authUC.IsAdmin(ctx, req.UserId)
	if err != nil {
		h.logger.Error("is_admin check failed", slog.String("op", op), slog.String("err", err.Error()))
		return nil, status.Error(codes.Internal, "internal error")
	}

	return &ssov1.IsAdminResponse{IsAdmin: isAdmin}, nil
}

// getCallerRole — вспомогательный метод для проверки прав вызывающего.
func (h *AuthHandler) getCallerRole(ctx context.Context) (entity.Role, error) {
	// Пример: извлечение роли из JWT в метаданных gRPC
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return "", errors.New("metadata not found")
	}

	tokens := md.Get("authorization")
	if len(tokens) == 0 {
		return "", errors.New("token not provided")
	}

	claims, err := h.authUC.ParseToken(tokens[0])
	if err != nil {
		return "", err
	}

	return entity.Role(claims.Role), nil
}

func (h *AuthHandler) getCallerUserID(ctx context.Context) (int64, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return 0, errors.New("metadata not found")
	}

	tokens := md.Get("authorization")
	if len(tokens) == 0 {
		return 0, errors.New("token not provided")
	}

	claims, err := h.authUC.ParseToken(tokens[0])
	if err != nil {
		return 0, err
	}

	return claims.UserID, nil
}
