package usermysql

import (
	"context"
	"fmt"

	"github.com/adetxt/user/domain"
	"github.com/adetxt/user/utils/password"
	"golang.org/x/sync/errgroup"
	"gorm.io/gorm"
)

type repository struct {
	db *gorm.DB
}

func New(db *gorm.DB) domain.UserRepository {
	return &repository{
		db: db,
	}
}

func (r *repository) GetUsers(ctx context.Context, params *domain.GetUsersParams) ([]*domain.User, *domain.PaginationInfo, error) {
	pagination := domain.MakePaginationInfo(params.Page, params.PageSize)

	users := []User{}
	db := r.db.Model(User{}).Preload("Roles")

	if params.Keyword != "" {
		db.Where("name LIKE ?", "%"+params.Keyword+"%")
	}

	var totalData int64

	if err := db.Count(&totalData).Error; err != nil {
		return nil, nil, err
	}

	pagination.SetTotalData(totalData)

	db.Offset(int(pagination.GetOffset())).Limit(int(pagination.PageSize))

	if err := db.Order("id desc").Find(&users).Error; err != nil {
		return nil, nil, err
	}

	res := make([]*domain.User, len(users))
	for i := 0; i < len(users); i++ {
		res[i] = users[i].ToEntity()
	}

	return res, pagination, nil
}

func (r *repository) GetUserByIdentifier(ctx context.Context, identifier string, value interface{}) (*domain.User, error) {
	user := User{}
	db := r.db.Model(User{}).Preload("Roles")

	switch identifier {
	case "id":
		db.Where("id = ?", value)
	case "email":
		db.Where("email = ?", value)
	default:
		return nil, fmt.Errorf("invalid identifier")
	}

	if err := db.First(&user).Error; err != nil {
		return nil, err
	}

	return user.ToEntity(), nil
}

func (r *repository) CreateUser(ctx context.Context, data *domain.User) (int64, error) {
	user := MakeUser(data)

	if err := r.db.Create(user).Error; err != nil {
		return 0, err
	}

	return user.ID, nil
}

func (r *repository) UpdateUser(ctx context.Context, data *domain.User) error {
	user, err := r.GetUserByIdentifier(ctx, "id", data.ID)
	if err != nil {
		return err
	}

	updateData := make(map[string]interface{})

	if data.Name != "" {
		updateData["name"] = data.Name
	}

	if data.Email != "" {
		updateData["email"] = data.Email
	}

	if data.Password != "" {
		updateData["password"] = data.Password
	}

	if err := r.db.Model(user).Updates(updateData).Error; err != nil {
		return err
	}

	return nil
}

func (r *repository) DeleteUser(ctx context.Context, id int64) error {
	if err := r.db.Where("id = ?", id).Delete(&User{}).Error; err != nil {
		return err
	}

	return nil
}

func (r *repository) GetRoles(ctx context.Context) ([]*domain.Role, error) {
	roles := []Role{}

	if err := r.db.Model(&Role{}).Preload("Permissions").
		Find(&roles).Error; err != nil {
		return nil, err
	}

	res := make([]*domain.Role, len(roles))
	for i := 0; i < len(roles); i++ {
		res[i] = &domain.Role{
			Name:        roles[i].Name,
			Permissions: roles[i].PermissionNames(),
		}
	}

	return res, nil
}

func (r *repository) GetPermissionsByRole(ctx context.Context, roleNames []string) ([]string, error) {
	roles := []Role{}

	if err := r.db.Model(&Role{}).Preload("Permissions").
		Where("name IN ?", roleNames).
		Find(&roles).Error; err != nil {
		return nil, err
	}

	res := []string{}
	for i := 0; i < len(roles); i++ {
		res = append(res, roles[i].PermissionNames()...)
	}

	return res, nil
}

func (r *repository) Seeding(ctx context.Context) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		eg, _ := errgroup.WithContext(ctx)

		eg.Go(func() error {
			return tx.Model(Role{}).Create([]Role{
				{
					ID:   1,
					Name: "admin",
				},
				{
					ID:   2,
					Name: "user",
				},
			}).Error
		})

		eg.Go(func() error {
			return tx.Model(Permission{}).Create([]Permission{
				{
					ID:   1,
					Name: "user:list",
				},
				{
					ID:   2,
					Name: "user:detail",
				},
				{
					ID:   3,
					Name: "user:create",
				},
				{
					ID:   4,
					Name: "user:update",
				},
				{
					ID:   5,
					Name: "user:delete",
				},
			}).Error
		})

		eg.Go(func() error {
			hashed, _ := password.HashPassword("password")

			return tx.Model(User{}).Create([]User{
				{
					ID:       1,
					Name:     "admin",
					Email:    "admin@mail.com",
					Password: hashed,
				},
				{
					ID:       2,
					Name:     "user",
					Email:    "user@mail.com",
					Password: hashed,
				},
			}).Error
		})

		eg.Go(func() error {
			return tx.Model(UserRole{}).Create([]UserRole{
				{
					UserID: 1,
					RoleID: 1,
				},
				{
					UserID: 2,
					RoleID: 2,
				},
			}).Error
		})

		eg.Go(func() error {
			return tx.Model(RolePermission{}).Create([]RolePermission{
				{
					RoleID:       1,
					PermissionID: 1,
				},
				{
					RoleID:       1,
					PermissionID: 2,
				},
				{
					RoleID:       1,
					PermissionID: 3,
				},
				{
					RoleID:       1,
					PermissionID: 4,
				},
				{
					RoleID:       1,
					PermissionID: 5,
				},
				{
					RoleID:       2,
					PermissionID: 1,
				},
				{
					RoleID:       2,
					PermissionID: 2,
				},
			}).Error
		})

		return eg.Wait()
	})
}
