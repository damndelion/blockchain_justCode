package repo

import (
	"context"
	"fmt"
	"github.com/evrone/go-clean-template/internal/user/controller/http/v1/dto"
	"github.com/evrone/go-clean-template/internal/user/entity"
	"gorm.io/gorm"
	"strconv"
)

type UserRepo struct {
	DB *gorm.DB
}

func NewUserRepo(db *gorm.DB) *UserRepo {
	return &UserRepo{db}
}

func (ur *UserRepo) GetUsers(ctx context.Context) (users []*userEntity.User, err error) {

	res := ur.DB.Find(&users)
	if res.Error != nil {
		return nil, res.Error
	}
	return users, nil
}

func (ur *UserRepo) CreateUser(ctx context.Context, user *userEntity.User) (int, error) {

	res := ur.DB.WithContext(ctx).Create(user)
	if res.Error != nil {
		return 0, res.Error
	}
	return user.Id, nil
}

func (ur *UserRepo) GetUserByEmail(ctx context.Context, email string) (user *userEntity.User, err error) {
	res := ur.DB.Where("email = ?", email).WithContext(ctx).Find(&user)
	if res.Error != nil {
		return nil, res.Error
	}
	return user, nil
}

func (ur *UserRepo) GetUserByID(ctx context.Context, id int) (user *userEntity.User, err error) {
	res := ur.DB.Where("id = ?", id).WithContext(ctx).Find(&user)
	if res.Error != nil {
		return nil, res.Error
	}
	return user, nil
}

func (ur *UserRepo) SetUserWallet(ctx context.Context, userID string, address string) error {
	id, _ := strconv.Atoi(userID)
	query := "SELECT wallet FROM users WHERE id = $2"
	res := ur.DB.Exec(query, address, id)
	if res == nil {
		return fmt.Errorf("user already have wallet existing")
	}
	query = "UPDATE users SET wallet = $1 WHERE id = $2"
	res = ur.DB.Exec(query, address, id)
	return nil
}

func (ur *UserRepo) CreateUserDetailInfo(ctx context.Context, userData dto.UserDetailRequest, id int) error {
	userInfo := userEntity.UserInfo{
		UserID:  id,
		Age:     userData.Age,
		Phone:   userData.Phone,
		Address: userData.Address,
		Country: userData.Country,
		City:    userData.City,
	}
	userCredentials := userEntity.UserCredentials{
		UserID:  id,
		CardNum: userData.CardNum,
		Type:    userData.CardType,
		CVV:     userData.CVV,
	}
	tx := ur.DB.Begin()
	if err := tx.Create(&userCredentials).Error; err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Create(&userInfo).Error; err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()

	if err := ur.DB.Model(&userEntity.User{}).Where("id = ?", id).Update("valid", true).Error; err != nil {
		return err
	}
	return nil
}

func (ur *UserRepo) SetUserDetailInfo(ctx context.Context, userData dto.UserDetailRequest, id int) error {
	userInfo := userEntity.UserInfo{
		UserID:  id,
		Age:     userData.Age,
		Phone:   userData.Phone,
		Address: userData.Address,
		Country: userData.Country,
		City:    userData.City,
	}
	userCredentials := userEntity.UserCredentials{
		UserID:  id,
		CardNum: userData.CardNum,
		Type:    userData.CardType,
		CVV:     userData.CVV,
	}
	tx := ur.DB.Begin()
	if err := tx.Model(&userCredentials).Where("user_id = ?", id).Updates(&userCredentials).Error; err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Model(&userInfo).Where("user_id = ?", id).Updates(&userInfo).Error; err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()

	if err := ur.DB.Model(&userEntity.User{}).Where("id = ?", id).Update("valid", true).Error; err != nil {
		return err
	}
	return nil
}
