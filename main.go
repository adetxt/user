package main

import (
	"context"
	"encoding/json"
	"errors"
	"strings"

	"github.com/adetxt/edison"
	"github.com/adetxt/user/config"
	pbAccount "github.com/adetxt/user/gen/proto/go/account/v1"
	grpcHdl "github.com/adetxt/user/handler/grpc"
	usermysql "github.com/adetxt/user/repository/user_mysql"
	"github.com/adetxt/user/usecase"
	"github.com/adetxt/user/utils/auth"
	"github.com/adetxt/user/utils/mysql"
	"github.com/golang-jwt/jwt/v4"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"gorm.io/gorm"
)

func main() {
	cfg := config.New()
	db := initMySQL(cfg)

	// DEVELOPMENT OPNLY
	db.AutoMigrate(usermysql.User{}, usermysql.Role{}, usermysql.Permission{}, usermysql.RolePermission{}, usermysql.UserRole{})

	// repository
	userRepo := usermysql.New(db)

	// usecase
	userUc := usecase.NewUserUsecase(userRepo)
	authUc := usecase.NewAuthUsecase(cfg, userRepo)

	// handler
	accountHdl := grpcHdl.NewAccountHandler(userUc)
	authHdl := grpcHdl.NewAuthHandler(cfg, authUc)

	// init edison
	ed := edison.New()

	// register interceptor
	ed.UnaryServerInterceptor(
		localAuthInterceptor(cfg.JWTKey),
	)

	ed.Prepare(
		edison.RestPort(cfg.RestPort),
		edison.GrpcPort(cfg.GrpcPort),
		edison.GracefullShutdown(),
	)

	ed.RestRouter("GET", "/seed", func(ctx context.Context, clientCtx edison.RestContext) error {
		if err := userRepo.Seeding(ctx); err != nil {
			return err
		}

		clientCtx.EchoContext.JSON(200, "ok")
		return nil
	})

	pbAccount.RegisterAccountService(ed, accountHdl)
	pbAccount.RegisterAuthService(ed, authHdl)

	ed.Start()
}

func initMySQL(cfg config.Config) *gorm.DB {
	return mysql.Init(&mysql.Config{
		Host:     cfg.DBHost,
		Port:     cfg.DBPort,
		DBName:   cfg.DBName,
		Username: cfg.DBUsername,
		Password: cfg.DBPassword,
	})
}

func localAuthInterceptor(JWTKey string) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		data := map[string]interface{}{}

		// Skip authorize when GetJWT is requested
		if info.FullMethod != "/account.v1.AuthService/Login" {
			claims, err := authorize(ctx, JWTKey, info.FullMethod)
			if err != nil {
				return nil, err
			}

			data = claims
		}

		b, err := json.Marshal(data)
		if err != nil {
			return nil, err
		}

		ctx = context.WithValue(ctx, "claims", string(b))

		// Calls the handler
		h, err := handler(ctx, req)

		return h, err
	}
}

func authorize(ctx context.Context, JWTKey, fullMth string) (map[string]interface{}, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, status.Errorf(codes.InvalidArgument, "Retrieving metadata is failed")
	}

	authHeader, ok := md["authorization"]
	if !ok {
		return nil, status.Errorf(codes.Unauthenticated, "Authorization token is not supplied")
	}

	token := authHeader[0]

	if strings.Contains(token, " ") {
		tokenParts := strings.Split(token, " ")
		token = tokenParts[1]
	}

	// validateToken function validates the token
	claims, err := auth.ParseToken(token, JWTKey)
	if err != nil && !errors.Is(err, jwt.ErrTokenExpired) {
		return nil, status.Errorf(codes.Unauthenticated, err.Error())
	}

	if fullMth != "/account.v1.AuthService/RefreshToken" && errors.Is(err, jwt.ErrTokenExpired) {
		return nil, status.Errorf(codes.Unavailable, "token is expired")
	}

	return claims, nil
}
