package repo

import (
	"context"
	"github.com/evrone/go-clean-template/internal/auth/entity"
	userEntity "github.com/evrone/go-clean-template/internal/user/entity"
	"gorm.io/gorm"
)

type AuthRepo struct {
	DB *gorm.DB
}

func NewAuthRepo(db *gorm.DB) *AuthRepo {
	return &AuthRepo{db}
}

func (t *AuthRepo) CreateUserToken(ctx context.Context, userToken entity.Token) error {
	if err := t.DB.Create(&userToken).Error; err != nil {
		return err
	}

	return nil

}

func (t *AuthRepo) UpdateUserToken(ctx context.Context, userToken entity.Token) error {
	return nil

}

func (t *AuthRepo) CreateUser(ctx context.Context, user *userEntity.User) (int, error) {
	result := t.DB.Create(&user)
	if result.Error != nil {
		return 0, result.Error
	}

	return user.Id, nil

}

func (t *AuthRepo) GetUserByEmail(ctx context.Context, email string) (user *userEntity.User, err error) {
	res := t.DB.Where("email = ?", email).WithContext(ctx).Find(&user)
	if res.Error != nil {
		return nil, res.Error
	}
	return user, nil
}
