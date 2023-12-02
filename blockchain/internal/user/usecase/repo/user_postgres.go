package repo

import (
	"context"
	"github.com/damndelion/blockchain_justCode/internal/user/controller/http/v1/dto"
	userEntity "github.com/damndelion/blockchain_justCode/internal/user/entity"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type UserRepo struct {
	DB *gorm.DB
}

func NewUserRepo(db *gorm.DB) *UserRepo {
	return &UserRepo{db}
}

func (ur *UserRepo) GetUsers(_ context.Context) (users []*userEntity.User, err error) {
	res := ur.DB.Find(&users)
	if res.Error != nil {
		return nil, res.Error
	}

	return users, nil
}

func (ur *UserRepo) CreateUser(ctx context.Context, userRequest dto.UserUpdateRequest) (int, error) {
	generatedHash, err := bcrypt.GenerateFromPassword([]byte(userRequest.Password), bcrypt.DefaultCost)
	if err != nil {
		return 0, err
	}
	user := userEntity.User{
		Name:     userRequest.Name,
		Email:    userRequest.Email,
		Password: string(generatedHash),
		Wallet:   userRequest.Wallet,
		Role:     "user",
	}
	res := ur.DB.WithContext(ctx).Create(&user)
	if res.Error != nil {
		return 0, res.Error
	}

	return user.ID, nil
}

func (ur *UserRepo) UpdateUser(_ context.Context, userData dto.UserUpdateRequest, email string) error {
	generatedHash, err := bcrypt.GenerateFromPassword([]byte(userData.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	user := userEntity.User{
		Name:     userData.Name,
		Email:    userData.Email,
		Password: string(generatedHash),
		Wallet:   userData.Wallet,
		Valid:    userData.Valid,
		Role:     userData.Role,
	}

	err = ur.DB.Model(&user).Where("email = ?", email).Updates(&user).Error
	if err != nil {
		return err
	}

	return nil
}

func (ur *UserRepo) GetUserRole(_ context.Context, id int) (string, error) {
	var user userEntity.User
	result := ur.DB.First(&user, id)
	if result.Error != nil {
		return "", result.Error
	}

	return user.Role, nil
}

func (ur *UserRepo) GetUserByEmail(ctx context.Context, email string) (user *userEntity.User, err error) {
	res := ur.DB.Where("email = ?", email).WithContext(ctx).Find(&user)
	if res.Error != nil {
		return nil, res.Error
	}

	return user, nil
}

func (ur *UserRepo) GetUserByID(ctx context.Context, id string) (user *userEntity.User, err error) {
	res := ur.DB.Where("id = ?", id).WithContext(ctx).Find(&user)
	if res.Error != nil {
		return nil, res.Error
	}

	return user, nil
}

func (ur *UserRepo) DeleteUser(ctx context.Context, id string) error {
	err := ur.DB.Where("user_id = ?", id).Delete(&userEntity.UserCredentials{}).WithContext(ctx).Error
	err = ur.DB.Where("user_id = ?", id).Delete(&userEntity.UserInfo{}).WithContext(ctx).Error
	err = ur.DB.Where("id = ?", id).Delete(&userEntity.User{}).WithContext(ctx).Error
	if err != nil {
		return err
	}

	return nil
}

func (ur *UserRepo) DeleteUserInfo(ctx context.Context, id string) error {
	err := ur.DB.Where("user_id = ?", id).Delete(&userEntity.UserInfo{}).WithContext(ctx).Error
	if err != nil {
		return err
	}

	return nil
}

func (ur *UserRepo) DeleteUserCred(ctx context.Context, id string) error {
	err := ur.DB.Where("user_id = ?", id).Delete(&userEntity.UserCredentials{}).WithContext(ctx).Error
	if err != nil {
		return err
	}

	return nil
}

func (ur *UserRepo) SetUserWallet(_ context.Context, userID, address string) error {
	err := ur.DB.Model(&userEntity.User{}).Where("id = ?", userID).Update("wallet", address).Error
	if err != nil {
		return err
	}

	return nil
}
