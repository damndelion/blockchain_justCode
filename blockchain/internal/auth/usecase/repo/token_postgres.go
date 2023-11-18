package repo

import (
	"context"
	"fmt"
	authEntity "github.com/evrone/go-clean-template/internal/auth/entity"
	"github.com/evrone/go-clean-template/internal/auth/transport"
	userEntity "github.com/evrone/go-clean-template/internal/user/entity"
	"github.com/opentracing/opentracing-go"
	"gorm.io/gorm"
)

type AuthRepo struct {
	DB                *gorm.DB
	userGrpcTransport *transport.UserGrpcTransport
}

func NewAuthRepo(db *gorm.DB, userGrpcTransport *transport.UserGrpcTransport) *AuthRepo {
	return &AuthRepo{db, userGrpcTransport}
}

func (t *AuthRepo) CreateUserToken(_ context.Context, userToken authEntity.Token) error {
	if err := t.DB.Create(&userToken).Error; err != nil {
		return err
	}

	return nil
}

func (t *AuthRepo) CreateUser(ctx context.Context, user *userEntity.User) (int, error) {
	grpcUser, err := t.userGrpcTransport.CreateUser(ctx, user)
	if err != nil {
		return 0, err
	}

	return int(grpcUser.Id), nil
}

func (t *AuthRepo) GetUserByEmail(ctx context.Context, email string) (*userEntity.User, error) {
	span, _ := opentracing.StartSpanFromContext(ctx, "get user by email repo")
	defer span.Finish()
	grpcUser, err := t.userGrpcTransport.GetUserByEmail(ctx, email)
	if err != nil {
		return nil, err
	}
	user := &userEntity.User{
		ID:       int(grpcUser.Id),
		Name:     grpcUser.Name,
		Email:    grpcUser.Email,
		Password: grpcUser.Password,
		Wallet:   grpcUser.Wallet,
		Valid:    grpcUser.Valid,
		Role:     grpcUser.Role,
	}

	return user, nil
}

func (t *AuthRepo) ConfirmCode(_ context.Context, email string) (int, error) {
	var code int
	res := t.DB.Model(&authEntity.UserCode{}).Where("email = ?", email).Pluck("code", &code)
	if res.Error != nil {
		return 0, res.Error
	}

	return code, nil
}

func (t *AuthRepo) CheckForEmail(_ context.Context, email string) error {
	var exists bool
	t.DB.Model(userEntity.User{}).Where("email = ?", email).Find(&exists)
	if exists {
		return fmt.Errorf("user with this email alraedy exists")
	}
	return nil
}
