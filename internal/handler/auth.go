package handler

import (
	"context"

	"github.com/VaneZ444/golang-forum/auth-service/internal/entity"
	"github.com/VaneZ444/golang-forum/auth-service/internal/usecase"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	ssov1 "github.com/VaneZ444/golang-forum/shared/protos/gen/go/sso"
)

type AuthHandler struct {
	authUC *usecase.AuthUseCase
	protos.UnimplementedAuthServer
}

func NewAuthHandler(authUC *usecase.AuthUseCase) *AuthHandler {
	return &AuthHandler{authUC: authUC}
}

func (h *AuthHandler) Register(ctx context.Context, req *protos.RegisterRequest) (*protos.RegisterResponse, error) {
	userID, err := h.authUC.Register(req.Email, req.Password)
	if err != nil {
		switch err {
		case entity.ErrInvalidEmail:
			return nil, status.Error(codes.InvalidArgument, "invalid email")
		case entity.ErrWeakPassword:
			return nil, status.Error(codes.InvalidArgument, "weak password")
		default:
			return nil, status.Error(codes.Internal, "registration failed")
		}
	}
	return &protos.RegisterResponse{UserId: userID}, nil
}
