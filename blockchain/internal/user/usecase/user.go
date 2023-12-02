package usecase

import (
	"context"
	"github.com/opentracing/opentracing-go"

	"github.com/evrone/go-clean-template/config/user"
	"github.com/evrone/go-clean-template/internal/user/controller/http/v1/dto"
	userEntity "github.com/evrone/go-clean-template/internal/user/entity"
)

type User struct {
	repo UserRepo
	cfg  *user.Config
}

func NewUser(repo UserRepo, cfg *user.Config) *User {
	return &User{repo, cfg}
}

func (u *User) Users(ctx context.Context) ([]*userEntity.User, error) {
	return u.repo.GetUsers(ctx)
}

func (u *User) CreateUser(ctx context.Context, user dto.UserUpdateRequest) (int, error) {
	return u.repo.CreateUser(ctx, user)
}

func (u *User) UpdateUser(ctx context.Context, userData dto.UserUpdateRequest, id string) error {
	err := u.repo.UpdateUser(ctx, userData, id)
	if err != nil {
		return err
	}

	return nil
}

func (u *User) DeleteUser(ctx context.Context, id string) error {
	err := u.repo.DeleteUser(ctx, id)
	if err != nil {
		return err
	}

	return nil
}

func (u *User) GetUserByEmail(ctx context.Context, email string) (*userEntity.User, error) {
	return u.repo.GetUserByEmail(ctx, email)
}

func (u *User) GetUserByID(ctx context.Context, id string) (*userEntity.User, error) {
	span, spanCtx := opentracing.StartSpanFromContext(ctx, "get user by id use case")
	defer span.Finish()

	return u.repo.GetUserByID(spanCtx, id)
}
