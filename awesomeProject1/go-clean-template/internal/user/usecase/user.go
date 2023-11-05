package usecase

import (
	"context"
	"github.com/evrone/go-clean-template/internal/user/controller/http/v1/dto"
	"github.com/evrone/go-clean-template/internal/user/entity"
	"time"
)

type User struct {
	repo UserRepo
}

func NewUser(repo UserRepo) *User {
	return &User{repo}
}

func (u *User) Users(ctx context.Context) ([]*userEntity.User, error) {
	return u.repo.GetUsers(ctx)
}

func (u *User) CreateUser(ctx context.Context, user *userEntity.User) (int, error) {
	return u.repo.CreateUser(ctx, user)
}

func (u *User) GetUserByEmail(ctx context.Context, email string) (*userEntity.User, error) {
	time.Sleep(2 * time.Second)

	return u.repo.GetUserByEmail(ctx, email)
}

func (u *User) GetUserById(ctx context.Context, id int) (*userEntity.User, error) {

	return u.repo.GetUserByID(ctx, id)
}

func (u *User) CreateUserDetailInfo(ctx context.Context, userData dto.UserDetailRequest, id int) error {
	err := u.repo.CreateUserDetailInfo(ctx, userData, id)
	if err != nil {
		return err
	}
	return nil
}

func (u *User) SetUserDetailInfo(ctx context.Context, userData dto.UserDetailRequest, id int) error {
	err := u.repo.SetUserDetailInfo(ctx, userData, id)
	if err != nil {
		return err
	}
	return nil
}
