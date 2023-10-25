// Package usecase implements application business logic. Each logic group in own file.
package usecase

import (
	"context"
	entity2 "github.com/evrone/go-clean-template/internal/auth/entity"
	"github.com/evrone/go-clean-template/internal/user/controller/http/v1/dto"
	"github.com/evrone/go-clean-template/internal/user/entity"
)

//go:generate mockgen -source=interfaces.go -destination=./mocks_test.go -package=usecase_test

type (

	//Token Use case
	TokenUseCase interface {
		Users(ctx context.Context) ([]*entity.User, error)
		CreateUser(ctx context.Context, user *entity.User) (int, error)
		GetUserByEmail(ctx context.Context, email string) (*entity.User, error)
		GetUserById(ctx context.Context, id int) (*entity.User, error)
		Register(ctx context.Context, name, email, password string) error
		Login(ctx context.Context, email, password string) (*dto.LoginResponse, error)
	}

	//Token repo
	TokenRepo interface {
		CreateUserToken(ctx context.Context, userToken entity2.Token) error
		UpdateUserToken(ctx context.Context, userToken entity2.Token) error
	}
)
