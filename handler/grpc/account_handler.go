package grpc

import (
	"context"
	"errors"
	"fmt"
	"strconv"

	"github.com/adetxt/user/domain"
	pbAccount "github.com/adetxt/user/gen/proto/go/account/v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"gorm.io/gorm"
)

type accountHandler struct {
	pbAccount.UnimplementedAccountServiceServer
	userUsecase domain.UserUsecase
}

func NewAccountHandler(userUsecase domain.UserUsecase) pbAccount.AccountServiceServer {
	return &accountHandler{
		userUsecase: userUsecase,
	}
}

func (h *accountHandler) GetUsers(ctx context.Context, req *pbAccount.GetUsersRequest) (*pbAccount.GetUsersResponse, error) {
	c := getTokenInfo(ctx)

	id, err := strconv.ParseInt(c.ID, 10, 64)
	if err != nil {
		return nil, err
	}

	if err := h.userUsecase.Granted(ctx, id, []string{"user:list"}); err != nil {
		return nil, status.Error(codes.PermissionDenied, err.Error())
	}

	page := 1
	if req.Page > 1 {
		page = int(req.Page)
	}

	pageSize := 10
	if req.PageSize > 0 {
		pageSize = int(req.PageSize)
	}

	users, pagination, err := h.userUsecase.GetUsers(ctx, &domain.GetUsersParams{
		Page:     int32(page),
		PageSize: int32(pageSize),
		Keyword:  req.Keyword,
	})
	if err != nil {
		return nil, err
	}

	resUsers := make([]*pbAccount.User, len(users))
	for i := 0; i < len(users); i++ {
		resUsers[i] = &pbAccount.User{
			Id:    int32(users[i].ID),
			Name:  users[i].Name,
			Email: users[i].Email,
			Roles: users[i].Roles,
		}
	}

	return &pbAccount.GetUsersResponse{
		Items:    resUsers,
		Page:     int32(page),
		PageSize: int32(pageSize),
		Total:    int32(pagination.TotalData),
	}, nil
}

func (h *accountHandler) GetUser(ctx context.Context, req *pbAccount.GetUserRequest) (*pbAccount.GetUserResponse, error) {
	c := getTokenInfo(ctx)

	id, err := strconv.ParseInt(c.ID, 10, 64)
	if err != nil {
		return nil, err
	}

	if err := h.userUsecase.Granted(ctx, id, []string{"user:detail"}); err != nil {
		return nil, status.Error(codes.PermissionDenied, err.Error())
	}

	identifier := ""
	var value interface{}

	if req.Id != 0 {
		identifier = "id"
		value = req.Id
	} else if req.Email != "" {
		identifier = "email"
		value = req.Email
	} else {
		return nil, status.Error(codes.InvalidArgument, "id or email is required")
	}

	user, err := h.userUsecase.GetUserByIdentifier(ctx, identifier, value)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, status.Error(codes.NotFound, fmt.Sprintf("record with %s = %v not found", identifier, value))
		}

		return nil, err
	}

	return &pbAccount.GetUserResponse{
		User: &pbAccount.User{
			Id:    int32(user.ID),
			Name:  user.Name,
			Email: user.Email,
			Roles: user.Roles,
		},
	}, nil
}

func (h *accountHandler) GetCurrentUser(ctx context.Context, req *emptypb.Empty) (*pbAccount.GetUserResponse, error) {
	c := getTokenInfo(ctx)

	id, err := strconv.ParseInt(c.ID, 10, 64)
	if err != nil {
		return nil, err
	}

	user, err := h.userUsecase.GetUserByIdentifier(ctx, "id", id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, status.Error(codes.NotFound, err.Error())
		}

		return nil, err
	}

	return &pbAccount.GetUserResponse{
		User: &pbAccount.User{
			Id:    int32(user.ID),
			Name:  user.Name,
			Email: user.Email,
			Roles: user.Roles,
		},
	}, nil
}

func (h *accountHandler) CreateUser(ctx context.Context, req *pbAccount.CreateUserRequest) (*pbAccount.CreateUserResponse, error) {
	c := getTokenInfo(ctx)

	id, err := strconv.ParseInt(c.ID, 10, 64)
	if err != nil {
		return nil, err
	}

	if err := h.userUsecase.Granted(ctx, id, []string{"user:create"}); err != nil {
		return nil, status.Error(codes.PermissionDenied, err.Error())
	}

	if req.Name == "" {
		return nil, status.Error(codes.InvalidArgument, "name is required")
	}

	if req.Email == "" {
		return nil, status.Error(codes.InvalidArgument, "email is required")
	}

	if req.Password == "" {
		return nil, status.Error(codes.InvalidArgument, "password is required")
	}

	if req.Password != req.PasswordValidation {
		return nil, status.Error(codes.InvalidArgument, "password is not the same")
	}

	id, err = h.userUsecase.CreateUser(ctx, &domain.User{
		Name:     req.Name,
		Email:    req.Email,
		Password: req.Password,
	})
	if err != nil {
		return nil, err
	}

	return &pbAccount.CreateUserResponse{
		Id: int32(id),
	}, nil
}

func (h *accountHandler) UpdateUser(ctx context.Context, req *pbAccount.UpdateUserRequest) (*emptypb.Empty, error) {
	c := getTokenInfo(ctx)

	id, err := strconv.ParseInt(c.ID, 10, 64)
	if err != nil {
		return nil, err
	}

	if err := h.userUsecase.Granted(ctx, id, []string{"user:update"}); err != nil {
		return nil, status.Error(codes.PermissionDenied, err.Error())
	}

	if req.Id < 1 {
		return nil, status.Error(codes.InvalidArgument, "id is required")
	}

	if req.Password != "" && req.Password != req.PasswordValidation {
		return nil, status.Error(codes.InvalidArgument, "password is not the same")
	}

	if err := h.userUsecase.UpdateUser(ctx, &domain.User{
		ID:       int64(req.Id),
		Name:     req.Name,
		Email:    req.Email,
		Password: req.Password,
	}); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, status.Error(codes.NotFound, "record not found")
		}

		return nil, err
	}

	return &emptypb.Empty{}, nil
}

func (h *accountHandler) DeleteUser(ctx context.Context, req *pbAccount.DeleteUserRequest) (*emptypb.Empty, error) {
	c := getTokenInfo(ctx)

	id, err := strconv.ParseInt(c.ID, 10, 64)
	if err != nil {
		return nil, err
	}

	if err := h.userUsecase.Granted(ctx, id, []string{"user:delete"}); err != nil {
		return nil, status.Error(codes.PermissionDenied, err.Error())
	}

	if req.Id < 1 {
		return nil, status.Error(codes.InvalidArgument, "id is required")
	}

	if err := h.userUsecase.DeleteUser(ctx, int64(req.Id)); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, status.Error(codes.NotFound, "record not found")
		}

		return nil, err
	}

	return &emptypb.Empty{}, nil
}

func (h *accountHandler) GetRoles(ctx context.Context, req *emptypb.Empty) (*pbAccount.GetRolesResponse, error) {
	roles, err := h.userUsecase.GetRoles(ctx)
	if err != nil {
		return nil, err
	}

	res := make([]*pbAccount.Role, len(roles))
	for i := 0; i < len(roles); i++ {
		res[i] = &pbAccount.Role{
			Name:        roles[i].Name,
			Permissions: roles[i].Permissions,
		}
	}

	return &pbAccount.GetRolesResponse{
		Items: res,
	}, nil
}
