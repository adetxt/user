package usecase

import (
	"context"
	"fmt"
	"time"

	"github.com/adetxt/user/config"
	"github.com/adetxt/user/domain"
	authUtils "github.com/adetxt/user/utils/auth"
	passwordUtils "github.com/adetxt/user/utils/password"
)

type authUsecase struct {
	cfg      config.Config
	userRepo domain.UserRepository
}

func NewAuthUsecase(cfg config.Config, userRepo domain.UserRepository) domain.AuthUsecase {
	return &authUsecase{
		cfg:      cfg,
		userRepo: userRepo,
	}
}

func (uc *authUsecase) Login(ctx context.Context, email, password string) (*domain.FullToken, error) {
	user, err := uc.userRepo.GetUserByIdentifier(ctx, "email", email)
	if err != nil {
		return nil, err
	}

	if err := passwordUtils.ComparePassword(password, user.Password); err != nil {
		return nil, domain.ErrPasswordIncorrect
	}

	expDuration := time.Duration(1) * time.Hour
	exp := time.Now().Add(expDuration)

	token, err := authUtils.GetToken(user.ID, exp, uc.cfg.JWTKey)
	if err != nil {
		return nil, fmt.Errorf("failen when generating token : %v", err.Error())
	}

	refreshExpDuration := time.Duration((24 * 7)) * time.Hour
	refreshExp := time.Now().Add(refreshExpDuration)

	refreshToken, err := authUtils.GetRefreshToken(user.ID, refreshExp, uc.cfg.JWTKey)
	if err != nil {
		return nil, fmt.Errorf("failen when generating refresh token")
	}

	return &domain.FullToken{
		Token:                 token,
		TokenExpiredAt:        exp,
		RefreshToken:          refreshToken,
		RefreshTokenExpiredAt: refreshExp,
	}, nil
}

func (uc *authUsecase) RefreshToken(ctx context.Context, id int64) (*domain.AccessToken, error) {
	expDuration := time.Duration(1) * time.Hour
	exp := time.Now().Add(expDuration)

	token, err := authUtils.GetToken(id, exp, uc.cfg.JWTKey)
	if err != nil {
		return nil, fmt.Errorf("failen when generating token")
	}

	return &domain.AccessToken{
		Token:          token,
		TokenExpiredAt: exp,
	}, nil
}
