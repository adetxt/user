package usermysql

import (
	"github.com/adetxt/user/domain"
)

type User struct {
	ID       int64  `gorm:"column:id;primaryKey"`
	Name     string `gorm:"column:name"`
	Email    string `gorm:"column:email;unique"`
	Password string `gorm:"column:password"`
	Roles    []Role `gorm:"many2many:user_roles"`
}

type Role struct {
	ID          int64        `gorm:"column:id;primaryKey"`
	Name        string       `gorm:"column:name;unique"`
	Permissions []Permission `gorm:"many2many:role_permissions"`
}

type Permission struct {
	ID   int64  `gorm:"column:id;primaryKey"`
	Name string `gorm:"column:name;unique"`
}

type UserRole struct {
	UserID int64 `gorm:"column:user_id;uniqueIndex:idx_id"`
	RoleID int64 `gorm:"column:role_id;uniqueIndex:idx_id"`
}

type RolePermission struct {
	RoleID       int64 `gorm:"column:role_id;uniqueIndex:idx_id"`
	PermissionID int64 `gorm:"column:permission_id;uniqueIndex:idx_id"`
}

func (User) TableName() string {
	return "users"
}

func (Role) TableName() string {
	return "roles"
}

func (Permission) TableName() string {
	return "permissions"
}

func (UserRole) TableName() string {
	return "user_roles"
}

func (RolePermission) TableName() string {
	return "role_permissions"
}

func (i *User) ToEntity() *domain.User {
	return &domain.User{
		ID:       i.ID,
		Name:     i.Name,
		Email:    i.Email,
		Password: i.Password,
		Roles:    i.RoleNames(),
	}
}

func (i *User) RoleNames() (res []string) {
	for _, v := range i.Roles {
		res = append(res, v.Name)
	}

	return
}

func (i *Role) PermissionNames() (res []string) {
	for _, v := range i.Permissions {
		res = append(res, v.Name)
	}

	return
}

func MakeUser(i *domain.User) *User {
	return &User{
		ID:       i.ID,
		Name:     i.Name,
		Email:    i.Email,
		Password: i.Password,
	}
}
