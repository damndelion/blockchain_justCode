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

func (ur *UserRepo) CreateUser(ctx context.Context, userRequest dto.UserCreateRequest) (int, error) {
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

func (ur *UserRepo) SetUserWallet(_ context.Context, userID, address string) error {
	err := ur.DB.Model(&userEntity.User{}).Where("id = ?", userID).Update("wallet", address).Error
	if err != nil {
		return err
	}

	return nil
}

func (ur *UserRepo) CreateUserDetailInfo(_ context.Context, userData dto.UserDetailRequest, id int) error {
	generatedCardNum, err := bcrypt.GenerateFromPassword([]byte(userData.CardNum), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	generatedCVV, err := bcrypt.GenerateFromPassword([]byte(userData.CVV), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
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
		CardNum: string(generatedCardNum),
		Type:    userData.CardType,
		CVV:     string(generatedCVV),
	}
	tx := ur.DB.Begin()
	if err = tx.Create(&userCredentials).Error; err != nil {
		tx.Rollback()

		return err
	}

	if err = tx.Create(&userInfo).Error; err != nil {
		tx.Rollback()

		return err
	}

	tx.Commit()

	if err = ur.DB.Model(&userEntity.User{}).Where("id = ?", id).Update("valid", true).Error; err != nil {
		return err
	}

	return nil
}
