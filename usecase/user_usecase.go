package usecase

import (
	"context"
	"fmt"

	"github.com/adetxt/user/domain"
	"github.com/adetxt/user/utils/password"
)

type userUsecase struct {
	userRepo domain.UserRepository
}

func NewUserUsecase(userRepo domain.UserRepository) domain.UserUsecase {
	return &userUsecase{
		userRepo: userRepo,
	}
}

func (uc *userUsecase) GetUsers(ctx context.Context, params *domain.GetUsersParams) ([]*domain.User, *domain.PaginationInfo, error) {
	return uc.userRepo.GetUsers(ctx, params)
}

func (uc *userUsecase) GetUserByIdentifier(ctx context.Context, identifier string, value interface{}) (*domain.User, error) {
	return uc.userRepo.GetUserByIdentifier(ctx, identifier, value)
}

func (uc *userUsecase) CreateUser(ctx context.Context, data *domain.User) (int64, error) {
	hashed, err := password.HashPassword(data.Password)
	if err != nil {
		return 0, err
	}

	data.Password = hashed

	return uc.userRepo.CreateUser(ctx, data)
}

func (uc *userUsecase) UpdateUser(ctx context.Context, data *domain.User) error {
	if data.Password != "" {
		hashed, err := password.HashPassword(data.Password)
		if err != nil {
			return err
		}

		data.Password = hashed
	}

	return uc.userRepo.UpdateUser(ctx, data)
}

func (uc *userUsecase) DeleteUser(ctx context.Context, id int64) error {
	return uc.userRepo.DeleteUser(ctx, id)
}

func (uc *userUsecase) GetRoles(ctx context.Context) ([]*domain.Role, error) {
	return uc.userRepo.GetRoles(ctx)
}

func (uc *userUsecase) Granted(ctx context.Context, userID int64, permissions []string) error {
	user, err := uc.GetUserByIdentifier(ctx, "id", userID)
	if err != nil {
		return err
	}

	ps, err := uc.userRepo.GetPermissionsByRole(ctx, user.Roles)
	if err != nil {
		return err
	}

	for _, p := range permissions {
		for _, v := range ps {
			if p == v {
				return nil
			}
		}
	}

	return fmt.Errorf("not granted")
}
