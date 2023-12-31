package repo

import (
	"context"
	"errors"
	"fmt"

	"github.com/damndelion/blockchain_justCode/internal/user/controller/http/v1/dto"
	userEntity "github.com/damndelion/blockchain_justCode/internal/user/entity"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

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

func (ur *UserRepo) SetUserDetailInfo(_ context.Context, userData dto.UserDetailRequest, id int) error {
	userInfo := userEntity.UserInfo{
		UserID:  id,
		Age:     userData.Age,
		Phone:   userData.Phone,
		Address: userData.Address,
		Country: userData.Country,
		City:    userData.City,
	}
	generatedCardNum, err := bcrypt.GenerateFromPassword([]byte(userData.CardNum), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	generatedCVV, err := bcrypt.GenerateFromPassword([]byte(userData.CVV), bcrypt.DefaultCost)
	userCredentials := userEntity.UserCredentials{
		UserID:  id,
		CardNum: string(generatedCardNum),
		Type:    userData.CardType,
		CVV:     string(generatedCVV),
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

func (ur *UserRepo) UpdateUserInfo(_ context.Context, userData dto.UserUpdateInfoRequest, id int) error {
	userInfo := userEntity.UserInfo{
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

func (ur *UserRepo) CreateUserInfo(_ context.Context, userData dto.UserCreateInfoRequest) error {
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

func (ur *UserRepo) CreateUserCred(_ context.Context, userData dto.UserCreateCredRequest) error {
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

func (ur *UserRepo) UpdateUserCredentials(_ context.Context, userData dto.UserUpdateCredRequest, id int) error {
	generatedCardNum, err := bcrypt.GenerateFromPassword([]byte(userData.CardNum), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	generatedCVV, err := bcrypt.GenerateFromPassword([]byte(userData.CVV), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	userCredentials := userEntity.UserCredentials{
		UserID:  id,
		CardNum: string(generatedCardNum),
		Type:    userData.CardType,
		CVV:     string(generatedCVV),
	}

	err = ur.DB.Model(&userCredentials).Where("user_id = ?", id).Updates(&userCredentials).Error
	if err != nil {
		return err
	}

	return nil
}

func (ur *UserRepo) GetUserWallet(_ context.Context, id int) (string, error) {
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

func (ur *UserRepo) GetUsersWithSearch(_ context.Context, param string, value string) (users []*userEntity.User, err error) {
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

func (ur *UserRepo) GetUsersInfoWithSearch(_ context.Context, param string, value interface{}) (users []*userEntity.UserInfo, err error) {
	var condition string
	switch v := value.(type) {
	case string:
		condition = fmt.Sprintf("%s ILIKE ?", param)
		value = "%" + v + "%"
	case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64:
		condition = fmt.Sprintf("%s = ?", param)
	default:
		return nil, fmt.Errorf("unsupported type for search value: %T", value)
	}
	res := ur.DB.Where(condition, value).Find(&users)
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

func (ur *UserRepo) GetUserInfoByID(ctx context.Context, id int) (userInfo *userEntity.UserInfo, err error) {
	res := ur.DB.Where("user_id = ?", id).WithContext(ctx).Find(&userInfo)
	if res.Error != nil {
		return nil, res.Error
	}

	return userInfo, nil
}

func (ur *UserRepo) GetUserCredByID(ctx context.Context, id int) (userCred *userEntity.UserCredentials, err error) {
	res := ur.DB.Where("user_id = ?", id).WithContext(ctx).Find(&userCred)
	if res.Error != nil {
		return nil, res.Error
	}

	return userCred, nil
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

func (ur *UserRepo) DeleteUser(ctx context.Context, id int) error {
	err := ur.DB.Where("user_id = ?", id).Delete(&userEntity.UserCredentials{}).WithContext(ctx).Error
	err = ur.DB.Where("user_id = ?", id).Delete(&userEntity.UserInfo{}).WithContext(ctx).Error
	err = ur.DB.Where("id = ?", id).Delete(&userEntity.User{}).WithContext(ctx).Error
	if err != nil {
		return err
	}

	return nil
}

func (ur *UserRepo) DeleteUserInfo(ctx context.Context, id int) error {
	err := ur.DB.Where("user_id = ?", id).Delete(&userEntity.UserInfo{}).WithContext(ctx).Error
	if err != nil {
		return err
	}

	return nil
}

func (ur *UserRepo) DeleteUserCred(ctx context.Context, id int) error {
	err := ur.DB.Where("user_id = ?", id).Delete(&userEntity.UserCredentials{}).WithContext(ctx).Error
	if err != nil {
		return err
	}

	return nil
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
