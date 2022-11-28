package domain

import "context"

type UserUsecase interface {
	GetUsers(ctx context.Context, params *GetUsersParams) ([]*User, *PaginationInfo, error)
	GetUserByIdentifier(ctx context.Context, identifier string, value interface{}) (*User, error)
	CreateUser(ctx context.Context, data *User) (int64, error)
	UpdateUser(ctx context.Context, data *User) error
	DeleteUser(ctx context.Context, id int64) error
	GetRoles(ctx context.Context) ([]*Role, error)
	Granted(ctx context.Context, userID int64, permissions []string) error
}

type UserRepository interface {
	GetUsers(ctx context.Context, params *GetUsersParams) ([]*User, *PaginationInfo, error)
	GetUserByIdentifier(ctx context.Context, identifier string, value interface{}) (*User, error)
	CreateUser(ctx context.Context, data *User) (int64, error)
	UpdateUser(ctx context.Context, data *User) error
	DeleteUser(ctx context.Context, id int64) error
	GetRoles(ctx context.Context) ([]*Role, error)
	GetPermissionsByRole(ctx context.Context, roleNames []string) ([]string, error)
	Seeding(ctx context.Context) error
}

type User struct {
	ID       int64
	Name     string
	Email    string
	Password string
	Roles    []string
}

type Role struct {
	Name        string
	Permissions []string
}

type GetUsersParams struct {
	Page     int32
	PageSize int32
	Keyword  string
}
