package grpc

import (
	"context"
	"encoding/json"
	"errors"
	"strconv"
	"time"

	"github.com/adetxt/user/config"
	"github.com/adetxt/user/domain"
	pbAccount "github.com/adetxt/user/gen/proto/go/account/v1"
	"github.com/adetxt/user/utils/auth"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"gorm.io/gorm"
)

type authHandler struct {
	cfg    config.Config
	authUc domain.AuthUsecase
}

func NewAuthHandler(cfg config.Config, authUc domain.AuthUsecase) pbAccount.AuthServiceServer {
	return &authHandler{
		cfg:    cfg,
		authUc: authUc,
	}
}

func (h *authHandler) Login(ctx context.Context, req *pbAccount.LoginRequest) (*pbAccount.LoginResponse, error) {
	if req.Email == "" {
		return nil, status.Error(codes.InvalidArgument, "email is required")
	}

	if req.Password == "" {
		return nil, status.Error(codes.InvalidArgument, "password is required")
	}

	loginInfo, err := h.authUc.Login(ctx, req.Email, req.Password)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, status.Error(codes.NotFound, "record not found")
		}

		return nil, err
	}

	return &pbAccount.LoginResponse{
		Token:                 loginInfo.Token,
		TokenExpiredAt:        loginInfo.TokenExpiredAt.Format(time.RFC3339),
		RefreshToken:          loginInfo.RefreshToken,
		RefreshTokenExpiredAt: loginInfo.TokenExpiredAt.Format(time.RFC3339),
	}, nil
}

func (h *authHandler) RefreshToken(ctx context.Context, req *pbAccount.RefreshTokenRequest) (*pbAccount.RefreshTokenResponse, error) {
	refreshClaims, err := auth.ParseToken(req.RefreshToken, h.cfg.JWTKey)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	id, err := strconv.ParseInt(refreshClaims["jti"].(string), 10, 64)
	if err != nil {
		return nil, err
	}

	t, err := h.authUc.RefreshToken(ctx, id)
	if err != nil {
		return nil, err
	}

	return &pbAccount.RefreshTokenResponse{
		Token:          t.Token,
		TokenExpiredAt: t.TokenExpiredAt.Format(time.RFC3339),
	}, nil
}

func getTokenInfo(ctx context.Context) (res domain.JWTClaims) {
	c := ctx.Value("claims").(string)
	_ = json.Unmarshal([]byte(c), &res)
	return
}
