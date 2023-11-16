package usecase

import (
	"context"
	"github.com/evrone/go-clean-template/internal/auth/controller/http/v1/dto"
	authEntity "github.com/evrone/go-clean-template/internal/auth/entity"
	userEntity "github.com/evrone/go-clean-template/internal/user/entity"
)

type (

	// AuthUseCase -.
	AuthUseCase interface {
		Register(ctx context.Context, name, email, password string) error
		Login(ctx context.Context, email, password string) (*dto.LoginResponse, error)
		Refresh(ctx context.Context, refreshToken string) (string, error)
		ConfirmUserCode(ctx context.Context, email string, userCode int) error
	}

	//AuthRepo -.
	AuthRepo interface {
		CreateUserToken(ctx context.Context, userToken authEntity.Token) error
		CreateUser(ctx context.Context, user *userEntity.User) (int, error)
		GetUserByEmail(ctx context.Context, email string) (*userEntity.User, error)
		ConfirmCode(ctx context.Context, email string) (int, error)
	}
)
