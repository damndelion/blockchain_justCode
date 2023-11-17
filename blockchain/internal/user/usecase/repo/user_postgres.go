package repo

import (
	"context"
	"errors"
	"fmt"

	"github.com/evrone/go-clean-template/internal/user/controller/http/v1/dto"
	userEntity "github.com/evrone/go-clean-template/internal/user/entity"
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
		Role:     "",
	}
	res := ur.DB.WithContext(ctx).Create(&user)
	if res.Error != nil {
		return 0, res.Error
	}

	return user.ID, nil
}

func (ur *UserRepo) UpdateUser(_ context.Context, userData dto.UserUpdateRequest, id string) error {
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
	}

	err = ur.DB.Model(&user).Where("id = ?", id).Updates(&user).Error
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

func (ur *UserRepo) GetUsersDetailsInfo(_ context.Context) (usersInfo []*userEntity.UserInfo, err error) {
	info := ur.DB.Find(&usersInfo)
	if info.Error != nil {
		return nil, info.Error
	}

	return usersInfo, nil
}

func (ur *UserRepo) GetUsersCredentials(_ context.Context) (usersCred []*userEntity.UserCredentials, err error) {
	cred := ur.DB.Find(&usersCred)
	if cred.Error != nil {
		return nil, cred.Error
	}

	return usersCred, nil
}

func (ur *UserRepo) CreateUserDetailInfo(_ context.Context, userData dto.UserDetailRequest, id string) error {
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

func (ur *UserRepo) SetUserDetailInfo(_ context.Context, userData dto.UserDetailRequest, id string) error {
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

func (ur *UserRepo) UpdateUserInfo(_ context.Context, userData dto.UserInfoRequest, id string) error {
	userInfo := userEntity.UserInfo{
		UserID:  id,
		Age:     userData.Age,
		Phone:   userData.Phone,
		Address: userData.Address,
		Country: userData.Country,
		City:    userData.City,
	}

	err := ur.DB.Model(&userInfo).Where("user_id = ?", id).Updates(&userInfo).Error
	if err != nil {
		return err
	}

	return nil
}

func (ur *UserRepo) CreateUserInfo(_ context.Context, userData dto.UserInfoRequest) error {
	userInfo := userEntity.UserInfo{
		UserID:  userData.UserID,
		Age:     userData.Age,
		Phone:   userData.Phone,
		Address: userData.Address,
		Country: userData.Country,
		City:    userData.City,
	}

	err := ur.DB.Model(&userInfo).Create(&userInfo).Error
	if err != nil {
		return err
	}

	return nil
}

func (ur *UserRepo) CreateUserCred(_ context.Context, userData dto.UserCredRequest) error {
	generatedCardNum, err := bcrypt.GenerateFromPassword([]byte(userData.CardNum), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	generatedCVV, err := bcrypt.GenerateFromPassword([]byte(userData.CVV), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	userCred := userEntity.UserCredentials{
		UserID:  userData.UserID,
		CardNum: string(generatedCardNum),
		Type:    userData.CardType,
		CVV:     string(generatedCVV),
	}

	err = ur.DB.Model(&userCred).Create(&userCred).Error
	if err != nil {
		return err
	}

	return nil
}

func (ur *UserRepo) UpdateUserCredentials(_ context.Context, userData dto.UserCredRequest, id string) error {
	userCredentials := userEntity.UserCredentials{
		UserID:  id,
		CardNum: userData.CardNum,
		Type:    userData.CardType,
		CVV:     userData.CVV,
	}

	err := ur.DB.Model(&userCredentials).Where("user_id = ?", id).Updates(&userCredentials).Error
	if err != nil {
		return err
	}

	return nil
}

func (ur *UserRepo) GetUserWallet(_ context.Context, id string) (string, error) {
	var user userEntity.User
	if err := ur.DB.Model(&userEntity.User{}).Select("wallet").Where("id = ?", id).First(&user).Error; err != nil {
		if errors.Is(gorm.ErrRecordNotFound, err) {
			return "", errors.New(fmt.Sprintf("user does not have a wallet"))
		}

		return "", err
	}

	return user.Wallet, nil
}

func (ur *UserRepo) GetUsersWithFilter(_ context.Context, param, value string) (users []*userEntity.User, err error) {
	res := ur.DB.Where(fmt.Sprintf("%s = ?", param), value).Find(&users)
	if res.Error != nil {
		return nil, res.Error
	}

	return users, nil
}

func (ur *UserRepo) GetUsersInfoWithFilter(_ context.Context, param, value string) (users []*userEntity.UserInfo, err error) {
	res := ur.DB.Where(fmt.Sprintf("%s = ?", param), value).Find(&users)
	if res.Error != nil {
		return nil, res.Error
	}

	return users, nil
}

func (ur *UserRepo) GetUsersCredWithFilter(_ context.Context, param, value string) (users []*userEntity.UserCredentials, err error) {
	res := ur.DB.Where(fmt.Sprintf("%s = ?", param), value).Find(&users)
	if res.Error != nil {
		return nil, res.Error
	}

	return users, nil
}

func (ur *UserRepo) GetUsersWithSort(_ context.Context, sort, method string) (users []*userEntity.User, err error) {
	res := ur.DB.Order(fmt.Sprintf("%s %s", sort, method)).Find(&users)
	if res.Error != nil {
		return nil, res.Error
	}

	return users, nil
}

func (ur *UserRepo) GetUsersWithSearch(_ context.Context, param, value string) (users []*userEntity.User, err error) {
	res := ur.DB.Where(fmt.Sprintf("%s ILIKE ?", param), "%"+value+"%").
		Or(fmt.Sprintf("%s ILIKE ?", param), value+"%").
		Or(fmt.Sprintf("%s ILIKE ?", param), "%"+value).
		Find(&users)
	if res.Error != nil {
		return nil, res.Error
	}

	return users, nil
}

func (ur *UserRepo) GetUsersInfoWithSort(_ context.Context, sort, method string) (users []*userEntity.UserInfo, err error) {
	res := ur.DB.Order(fmt.Sprintf("%s %s", sort, method)).Find(&users)
	if res.Error != nil {
		return nil, res.Error
	}

	return users, nil
}

func (ur *UserRepo) GetUsersInfoWithSearch(_ context.Context, param, value string) (users []*userEntity.UserInfo, err error) {
	res := ur.DB.Where(fmt.Sprintf("%s ILIKE ?", param), "%"+value+"%").
		Or(fmt.Sprintf("%s ILIKE ?", param), value+"%").
		Or(fmt.Sprintf("%s ILIKE ?", param), "%"+value).
		Find(&users)
	if res.Error != nil {
		return nil, res.Error
	}

	return users, nil
}

func (ur *UserRepo) GetUsersCredWithSort(_ context.Context, sort, method string) (users []*userEntity.UserCredentials, err error) {
	res := ur.DB.Order(fmt.Sprintf("%s %s", sort, method)).Find(&users)
	if res.Error != nil {
		return nil, res.Error
	}

	return users, nil
}

func (ur *UserRepo) GetUsersCredWithSearch(_ context.Context, param, value string) (users []*userEntity.UserCredentials, err error) {
	res := ur.DB.Where(fmt.Sprintf("%s ILIKE ?", param), "%"+value+"%").
		Or(fmt.Sprintf("%s ILIKE ?", param), value+"%").
		Or(fmt.Sprintf("%s ILIKE ?", param), "%"+value).
		Find(&users)
	if res.Error != nil {
		return nil, res.Error
	}

	return users, nil
}

func (ur *UserRepo) GetUserInfoByID(ctx context.Context, id string) (userInfo *userEntity.UserInfo, err error) {
	res := ur.DB.Where("user_id = ?", id).WithContext(ctx).Find(&userInfo)
	if res.Error != nil {
		return nil, res.Error
	}

	return userInfo, nil
}

func (ur *UserRepo) GetUserCredByID(ctx context.Context, id string) (userCred *userEntity.UserCredentials, err error) {
	res := ur.DB.Where("user_id = ?", id).WithContext(ctx).Find(&userCred)
	if res.Error != nil {
		return nil, res.Error
	}

	return userCred, nil
}
