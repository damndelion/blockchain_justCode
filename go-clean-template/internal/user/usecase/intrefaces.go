// Package usecase implements application business logic. Each logic group in own file.
package usecase

import (
	"context"
	"github.com/evrone/go-clean-template/internal/user/controller/http/v1/dto"
	"github.com/evrone/go-clean-template/internal/user/entity"
)

//go:generate mockgen -source=interfaces.go -destination=./mocks_test.go -package=usecase_test

type (

	// User
	UserUseCase interface {
		Users(ctx context.Context) ([]*userEntity.User, error)
		CreateUser(ctx context.Context, user *userEntity.User) (int, error)
		GetUserByEmail(ctx context.Context, email string) (*userEntity.User, error)
		GetUserById(ctx context.Context, id string) (*userEntity.User, error)
		CreateUserDetailInfo(ctx context.Context, userData dto.UserDetailRequest, id string) error
		SetUserDetailInfo(ctx context.Context, userData dto.UserDetailRequest, id string) error
		GetIdFromToken(accessToken string) (string, error)
	}

	//UserRepo
	UserRepo interface {
		GetUsers(ctx context.Context) ([]*userEntity.User, error)
		CreateUser(ctx context.Context, user *userEntity.User) (int, error)
		GetUserByEmail(ctx context.Context, email string) (*userEntity.User, error)
		GetUserByID(ctx context.Context, id string) (*userEntity.User, error)
		SetUserWallet(ctx context.Context, userID string, address string) error
		CreateUserDetailInfo(ctx context.Context, userData dto.UserDetailRequest, id string) error
		SetUserDetailInfo(ctx context.Context, userData dto.UserDetailRequest, id string) error
	}
)
