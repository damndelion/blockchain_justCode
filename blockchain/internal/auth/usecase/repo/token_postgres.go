package repo

import (
	"context"
	"errors"
	"fmt"

	authEntity "github.com/damndelion/blockchain_justCode/internal/auth/entity"
	"github.com/damndelion/blockchain_justCode/internal/auth/transport"
	userEntity "github.com/damndelion/blockchain_justCode/internal/user/entity"
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

//func (t *AuthRepo) CreateUserToken(_ context.Context, userToken authEntity.Token) error {
//	if err := t.DB.Create(&userToken).Error; err != nil {
//		return err
//	}
//
//	return nil
//}

func (t *AuthRepo) CreateUserToken(_ context.Context, userToken authEntity.Token) error {
	existingToken := authEntity.Token{}
	err := t.DB.Model(&userToken).Where("user_id = ?", userToken.UserID).First(&existingToken).Error

	if err == nil {
		existingToken.RefreshToken = userToken.RefreshToken
		err = t.DB.Save(&existingToken).Error
	} else if errors.Is(err, gorm.ErrRecordNotFound) {
		err = t.DB.Create(&userToken).Error
	} else {
		return err
	}

	if err != nil {
		return err
	}

	return nil
}

func (t *AuthRepo) CreateUser(ctx context.Context, user *userEntity.User) (int, error) {
	span, spanCtx := opentracing.StartSpanFromContext(ctx, "create user repo")
	defer span.Finish()
	grpcUser, err := t.userGrpcTransport.CreateUser(spanCtx, user)
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

func (t *AuthRepo) ConfirmCode(ctx context.Context, email string) (int, error) {
	span, _ := opentracing.StartSpanFromContext(ctx, "confirm code repo")
	defer span.Finish()
	var code int
	res := t.DB.Model(&authEntity.UserCode{}).Where("email = ?", email).Pluck("code", &code)
	if res.Error != nil {
		return 0, res.Error
	}

	return code, nil
}

func (t *AuthRepo) CheckForEmail(ctx context.Context, email string) error {
	grpcUser, _ := t.userGrpcTransport.GetUserByEmail(ctx, email)
	if grpcUser.Id != 0 {
		return fmt.Errorf("user with this email alraedy exists")
	}

	return nil
}

// GetRefreshToken retrieves the refresh token associated with the specified user ID from the database.
func (r *AuthRepo) GetRefreshToken(ctx context.Context, userID int) (string, error) {
	span, _ := opentracing.StartSpanFromContext(ctx, "get refresh token")
	defer span.Finish()

	var refreshToken string
	err := r.DB.Model(&authEntity.Token{}).Where("user_id = ?", userID).Select("refresh_token").First(&refreshToken).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return "", nil
		}

		return "", err
	}

	return refreshToken, nil
}

// UpdateRefreshToken updates the refresh token associated with the specified user ID in the database.
func (r *AuthRepo) UpdateRefreshToken(ctx context.Context, userID int, refreshToken string) error {
	span, _ := opentracing.StartSpanFromContext(ctx, "update refresh token")
	defer span.Finish()

	err := r.DB.Model(&authEntity.Token{}).Where("user_id = ?", userID).Update("refresh_token", refreshToken).Error
	if err != nil {
		return err
	}

	return nil
}
