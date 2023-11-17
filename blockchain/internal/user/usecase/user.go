package usecase

import (
	"context"
	"fmt"
	"net/url"

	"github.com/evrone/go-clean-template/config/user"
	"github.com/evrone/go-clean-template/internal/user/controller/http/v1/dto"
	userEntity "github.com/evrone/go-clean-template/internal/user/entity"
	"github.com/golang-jwt/jwt"
)

type User struct {
	repo UserRepo
	cfg  *user.Config
}

func NewUser(repo UserRepo, cfg *user.Config) *User {
	return &User{repo, cfg}
}

func (u *User) Users(ctx context.Context) ([]*userEntity.User, error) {
	return u.repo.GetUsers(ctx)
}

func (u *User) CreateUser(ctx context.Context, user dto.UserUpdateRequest) (int, error) {
	return u.repo.CreateUser(ctx, user)
}

func (u *User) UpdateUser(ctx context.Context, userData dto.UserUpdateRequest, id string) error {
	err := u.repo.UpdateUser(ctx, userData, id)
	if err != nil {
		return err
	}

	return nil
}

func (u *User) DeleteUser(ctx context.Context, id string) error {
	err := u.repo.DeleteUser(ctx, id)
	if err != nil {
		return err
	}

	return nil
}

func (u *User) GetUserByEmail(ctx context.Context, email string) (*userEntity.User, error) {
	return u.repo.GetUserByEmail(ctx, email)
}

func (u *User) GetUserRole(ctx context.Context, id int) (string, error) {
	return u.repo.GetUserRole(ctx, id)
}

func (u *User) GetUserByID(ctx context.Context, id string) (*userEntity.User, error) {
	return u.repo.GetUserByID(ctx, id)
}

func (u *User) UsersInfo(ctx context.Context) ([]*userEntity.UserInfo, error) {
	return u.repo.GetUsersDetailsInfo(ctx)
}

func (u *User) UsersInfoWithFilter(ctx context.Context, params url.Values) ([]*userEntity.UserInfo, error) {
	usersInfo, err := u.repo.GetUsersInfoWithFilter(ctx, params.Get("filter"), params.Get("value"))
	if err != nil {
		return nil, err
	}

	return usersInfo, err
}

func (u *User) UsersCredWithFilter(ctx context.Context, params url.Values) ([]*userEntity.UserCredentials, error) {
	usersCred, err := u.repo.GetUsersCredWithFilter(ctx, params.Get("filter"), params.Get("value"))
	if err != nil {
		return nil, err
	}

	return usersCred, err
}

func (u *User) UsersCred(ctx context.Context) ([]*userEntity.UserCredentials, error) {
	return u.repo.GetUsersCredentials(ctx)
}

func (u *User) CreateUserDetailInfo(ctx context.Context, userData dto.UserDetailRequest, id string) error {
	err := u.repo.CreateUserDetailInfo(ctx, userData, id)
	if err != nil {
		return err
	}

	return nil
}

func (u *User) SetUserDetailInfo(ctx context.Context, userData dto.UserDetailRequest, id string) error {
	err := u.repo.SetUserDetailInfo(ctx, userData, id)
	if err != nil {
		return err
	}

	return nil
}

func (u *User) DeleteUserInfo(ctx context.Context, id string) error {
	err := u.repo.DeleteUserInfo(ctx, id)
	if err != nil {
		return err
	}

	return nil
}

func (u *User) GetIDFromToken(accessToken string) (string, error) {
	token, err := jwt.Parse(accessToken, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(u.cfg.SecretKey), nil
	})
	claims := token.Claims.(jwt.MapClaims)
	if err != nil || !token.Valid {
		return "", err
	}
	id := fmt.Sprintf("%v", claims["user_id"])

	return id, nil
}

func (u *User) UpdateUserInfo(ctx context.Context, userData dto.UserInfoRequest, id string) error {
	err := u.repo.UpdateUserInfo(ctx, userData, id)
	if err != nil {
		return err
	}

	return nil
}

func (u *User) CreateUserInfo(ctx context.Context, userData dto.UserInfoRequest) error {
	err := u.repo.CreateUserInfo(ctx, userData)
	if err != nil {
		return err
	}

	return nil
}

func (u *User) CreateUserCred(ctx context.Context, userData dto.UserCredRequest) error {
	err := u.repo.CreateUserCred(ctx, userData)
	if err != nil {
		return err
	}

	return nil
}

func (u *User) UpdateUserCredentials(ctx context.Context, userData dto.UserCredRequest, id string) error {
	err := u.repo.UpdateUserCredentials(ctx, userData, id)
	if err != nil {
		return err
	}

	return nil
}

func (u *User) DeleteUserCred(ctx context.Context, id string) error {
	err := u.repo.DeleteUserCred(ctx, id)
	if err != nil {
		return err
	}

	return nil
}

func (u *User) UsersWithFilter(ctx context.Context, params url.Values) ([]*userEntity.User, error) {
	users, err := u.repo.GetUsersWithFilter(ctx, params.Get("filter"), params.Get("value"))
	if err != nil {
		return nil, err
	}

	return users, nil
}

func (u *User) UsersWithSearch(ctx context.Context, params url.Values) ([]*userEntity.User, error) {
	users, err := u.repo.GetUsersWithSearch(ctx, params.Get("search"), params.Get("value"))
	if err != nil {
		return nil, err
	}

	return users, nil
}

func (u *User) UsersWithSort(ctx context.Context, sort, method string) ([]*userEntity.User, error) {
	users, err := u.repo.GetUsersWithSort(ctx, sort, method)
	if err != nil {
		return nil, err
	}

	return users, nil
}

func (u *User) UsersInfoWithSearch(ctx context.Context, params url.Values) ([]*userEntity.UserInfo, error) {
	usersInfo, err := u.repo.GetUsersInfoWithSearch(ctx, params.Get("search"), params.Get("value"))
	if err != nil {
		return nil, err
	}

	return usersInfo, err
}

func (u *User) UsersInfoWithSort(ctx context.Context, sort, method string) ([]*userEntity.UserInfo, error) {
	usersInfo, err := u.repo.GetUsersInfoWithSort(ctx, sort, method)
	if err != nil {
		return nil, err
	}

	return usersInfo, nil
}

func (u *User) UsersCredWithSearch(ctx context.Context, params url.Values) ([]*userEntity.UserCredentials, error) {
	usersCred, err := u.repo.GetUsersCredWithSearch(ctx, params.Get("search"), params.Get("value"))
	if err != nil {
		return nil, err
	}

	return usersCred, err
}

func (u *User) UsersCredWithSort(ctx context.Context, sort, method string) ([]*userEntity.UserCredentials, error) {
	usersCred, err := u.repo.GetUsersCredWithSort(ctx, sort, method)
	if err != nil {
		return nil, err
	}

	return usersCred, err
}

func (u *User) GetUserInfoByID(ctx context.Context, id string) (*userEntity.UserInfo, error) {
	return u.repo.GetUserInfoByID(ctx, id)
}

func (u *User) GetUserCredByID(ctx context.Context, id string) (*userEntity.UserCredentials, error) {
	return u.repo.GetUserCredByID(ctx, id)
}
