package handler

import (
	"context"
	"errors"
	"log/slog"

	"github.com/VaneZ444/auth-service/internal/entity"
	"github.com/VaneZ444/auth-service/internal/jwt"
	"github.com/VaneZ444/auth-service/internal/usecase"
	ssov1 "github.com/VaneZ444/golang-forum-protos/gen/go/sso"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

type AuthHandler struct {
	authUC     usecase.AuthUseCase
	jwtService jwt.Service
	logger     *slog.Logger // TODO: заменить на интерфейс
	ssov1.UnimplementedAuthServer
}

func NewAuthHandler(authUC usecase.AuthUseCase, jwtService jwt.Service, logger *slog.Logger) *AuthHandler {
	return &AuthHandler{
		authUC:     authUC,
		jwtService: jwtService,
		logger:     logger,
	}
}

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

func (h *AuthHandler) getCallerRole(ctx context.Context) (entity.Role, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return "", errors.New("metadata not found")
	}

	tokens := md.Get("authorization")
	if len(tokens) == 0 {
		return "", errors.New("token not provided")
	}

	claims, err := h.jwtService.ParseToken(tokens[0])
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

	claims, err := h.jwtService.ParseToken(tokens[0])
	if err != nil {
		return 0, err
	}

	return claims.UserID, nil
}

func (h *AuthHandler) Register(ctx context.Context, req *ssov1.RegisterRequest) (*ssov1.RegisterResponse, error) {
	const op = "handler.Register"
	h.logger.Info("register request", slog.String("op", op), slog.String("email", req.Email))

	userID, nickname, err := h.authUC.Register(ctx, req.Email, req.Nickname, req.Password)
	if err != nil {
		h.logger.Error("registration failed", slog.String("op", op), slog.String("err", err.Error()))
		if errors.Is(err, usecase.ErrInvalidCredentials) {
			return nil, status.Error(codes.InvalidArgument, "invalid email, nickname or password")
		}
		return nil, status.Error(codes.Internal, "internal error")
	}

	return &ssov1.RegisterResponse{
		UserId:   userID,
		Nickname: nickname,
	}, nil
}

func (h *AuthHandler) Login(ctx context.Context, req *ssov1.LoginRequest) (*ssov1.LoginResponse, error) {
	const op = "handler.Login"
	h.logger.Info("login request", slog.String("op", op), slog.String("email", req.Email))

	token, nickname, err := h.authUC.Login(ctx, req.Email, req.Password, req.AppId)
	if err != nil {
		h.logger.Error("login failed", slog.String("op", op), slog.String("err", err.Error()))
		switch {
		case errors.Is(err, usecase.ErrInvalidCredentials):
			return nil, status.Error(codes.Unauthenticated, "invalid credentials")
		case errors.Is(err, usecase.ErrUserBanned):
			return nil, status.Error(codes.PermissionDenied, "user banned")
		default:
			return nil, status.Error(codes.Internal, "internal error")
		}
	}

	return &ssov1.LoginResponse{
		Token:    token,
		Nickname: nickname,
	}, nil
}

func (h *AuthHandler) ValidateToken(ctx context.Context, req *ssov1.ValidateTokenRequest) (*ssov1.ValidateTokenResponse, error) {
	claims, err := h.jwtService.ParseToken(req.Token)
	if err != nil {
		return &ssov1.ValidateTokenResponse{Valid: false}, nil
	}

	return &ssov1.ValidateTokenResponse{
		Valid:    true,
		UserId:   claims.UserID,
		Role:     claims.Role,
		Nickname: claims.Nickname,
	}, nil
}
