package jwt

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type Claims struct {
	UserID   int64  `json:"user_id"`
	AppID    int32  `json:"app_id"`
	Role     string `json:"role"`
	Nickname string `json:"nickname"` // <-- добавили ник
	jwt.RegisteredClaims
}

type Service interface {
	GenerateToken(userID int64, appID int32, role, nickname string) (string, error)
	ParseToken(tokenStr string) (*Claims, error)
}

type jwtService struct {
	secret   string
	tokenTTL time.Duration
}

func NewService(secret string, tokenTTL time.Duration) Service {
	return &jwtService{
		secret:   secret,
		tokenTTL: tokenTTL,
	}
}

func (s *jwtService) GenerateToken(userID int64, appID int32, role, nickname string) (string, error) {
	now := time.Now()

	claims := &Claims{
		UserID:   userID,
		AppID:    appID,
		Role:     role,
		Nickname: nickname, // <-- сохраняем ник
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(s.tokenTTL)),
			Issuer:    "auth-service",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(s.secret))
}

func (s *jwtService) ParseToken(tokenStr string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(s.secret), nil
	})
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrInvalidToken, err)
	}

	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return nil, ErrInvalidToken
	}

	return claims, nil
}
