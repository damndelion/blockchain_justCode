// Package usecase implements application business logic. Each logic group in own file.
package usecase

import (
	"context"
	authEntity "github.com/evrone/go-clean-template/internal/auth/entity"
	"github.com/evrone/go-clean-template/internal/user/controller/http/v1/dto"
	userEntity "github.com/evrone/go-clean-template/internal/user/entity"
)

//go:generate mockgen -source=interfaces.go -destination=./mocks_test.go -package=usecase_test

type (

	//Auth Use case
	AuthUseCase interface {
		Register(ctx context.Context, name, email, password string) error
		Login(ctx context.Context, email, password string) (*dto.LoginResponse, error)
	}

	//Auth repo
	AuthRepo interface {
		CreateUserToken(ctx context.Context, userToken authEntity.Token) error
		UpdateUserToken(ctx context.Context, userToken authEntity.Token) error
		CreateUser(ctx context.Context, user *userEntity.User) (int, error)
		GetUserByEmail(ctx context.Context, email string) (*userEntity.User, error)
	}
)
