package domain

import (
	"context"
	"errors"
	"time"
)

type AuthUsecase interface {
	Login(ctx context.Context, email, password string) (*FullToken, error)
	RefreshToken(ctx context.Context, id int64) (*AccessToken, error)
}

var (
	ErrPasswordIncorrect = errors.New("password incorrect")
)

type FullToken struct {
	Token                 string
	TokenExpiredAt        time.Time
	RefreshToken          string
	RefreshTokenExpiredAt time.Time
}

type AccessToken struct {
	Token          string
	TokenExpiredAt time.Time
}

type JWTClaims struct {
	ID  string `json:"jti"`
	Exp int64  `json:"exp"`
}
