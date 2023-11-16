// Package usecase implements application business logic. Each logic group in own file.
package usecase

import (
	"context"
	"github.com/evrone/go-clean-template/internal/user/controller/http/v1/dto"
	"github.com/evrone/go-clean-template/internal/user/entity"
	"net/url"
)

//go:generate mockgen -source=interfaces.go -destination=./mocks_test.go -package=usecase_test

type (

	// UserUseCase -.
	UserUseCase interface {
		Users(ctx context.Context) ([]*userEntity.User, error)
		UsersWithFilter(ctx context.Context, params url.Values) ([]*userEntity.User, error)
		UsersWithSort(ctx context.Context, sort string, method string) ([]*userEntity.User, error)
		UsersWithSearch(ctx context.Context, params url.Values) ([]*userEntity.User, error)
		CreateUser(ctx context.Context, user dto.UserUpdateRequest) (int, error)
		UpdateUser(ctx context.Context, userData dto.UserUpdateRequest, id string) error
		GetUserByEmail(ctx context.Context, email string) (*userEntity.User, error)
		GetUserRole(ctx context.Context, id int) (string, error)
		GetUserById(ctx context.Context, id string) (*userEntity.User, error)
		DeleteUser(ctx context.Context, id string) error

		UsersInfo(ctx context.Context) ([]*userEntity.UserInfo, error)
		UsersInfoWithFilter(ctx context.Context, params url.Values) ([]*userEntity.UserInfo, error)
		UsersInfoWithSort(ctx context.Context, sort string, method string) ([]*userEntity.UserInfo, error)
		UsersInfoWithSearch(ctx context.Context, params url.Values) ([]*userEntity.UserInfo, error)
		GetUserInfoById(ctx context.Context, id string) (*userEntity.UserInfo, error)
		CreateUserDetailInfo(ctx context.Context, userData dto.UserDetailRequest, id string) error
		SetUserDetailInfo(ctx context.Context, userData dto.UserDetailRequest, id string) error
		UpdateUserInfo(ctx context.Context, userData dto.UserInfoRequest, id string) error
		CreateUserInfo(ctx context.Context, userData dto.UserInfoRequest) error
		DeleteUserInfo(ctx context.Context, id string) error

		UsersCred(ctx context.Context) ([]*userEntity.UserCredentials, error)
		UsersCredWithFilter(ctx context.Context, params url.Values) ([]*userEntity.UserCredentials, error)
		UsersCredWithSort(ctx context.Context, sort string, method string) ([]*userEntity.UserCredentials, error)
		UsersCredWithSearch(ctx context.Context, params url.Values) ([]*userEntity.UserCredentials, error)
		GetUserCredById(ctx context.Context, id string) (*userEntity.UserCredentials, error)
		CreateUserCred(ctx context.Context, userData dto.UserCredRequest) error
		UpdateUserCredentials(ctx context.Context, userData dto.UserCredRequest, id string) error
		DeleteUserCred(ctx context.Context, id string) error

		GetIdFromToken(accessToken string) (string, error)
	}

	//UserRepo -.
	UserRepo interface {
		GetUsers(ctx context.Context) ([]*userEntity.User, error)
		GetUsersWithFilter(ctx context.Context, param string, value string) (users []*userEntity.User, err error)
		GetUsersWithSort(ctx context.Context, param string, method string) (users []*userEntity.User, err error)
		GetUsersWithSearch(ctx context.Context, param string, value string) (users []*userEntity.User, err error)
		GetUserByEmail(ctx context.Context, email string) (*userEntity.User, error)
		GetUserByID(ctx context.Context, id string) (*userEntity.User, error)
		GetUserRole(ctx context.Context, id int) (string, error)
		CreateUser(ctx context.Context, user dto.UserUpdateRequest) (int, error)
		UpdateUser(ctx context.Context, userData dto.UserUpdateRequest, id string) error
		SetUserWallet(ctx context.Context, userID string, address string) error
		DeleteUser(ctx context.Context, id string) error

		GetUsersDetailsInfo(ctx context.Context) (usersInfo []*userEntity.UserInfo, err error)
		GetUsersInfoWithFilter(ctx context.Context, param string, value string) (users []*userEntity.UserInfo, err error)
		GetUsersInfoWithSort(ctx context.Context, param string, method string) (users []*userEntity.UserInfo, err error)
		GetUsersInfoWithSearch(ctx context.Context, param string, value string) (users []*userEntity.UserInfo, err error)
		GetUserInfoByID(ctx context.Context, id string) (userInfo *userEntity.UserInfo, err error)
		CreateUserDetailInfo(ctx context.Context, userData dto.UserDetailRequest, id string) error
		SetUserDetailInfo(ctx context.Context, userData dto.UserDetailRequest, id string) error
		CreateUserInfo(ctx context.Context, userData dto.UserInfoRequest) error
		UpdateUserInfo(ctx context.Context, userData dto.UserInfoRequest, id string) error
		DeleteUserInfo(ctx context.Context, id string) error

		GetUsersCredentials(ctx context.Context) (usersInfo []*userEntity.UserCredentials, err error)
		GetUsersCredWithFilter(ctx context.Context, param string, value string) (users []*userEntity.UserCredentials, err error)
		GetUsersCredWithSort(ctx context.Context, param string, method string) (users []*userEntity.UserCredentials, err error)
		GetUsersCredWithSearch(ctx context.Context, param string, value string) (users []*userEntity.UserCredentials, err error)
		GetUserCredByID(ctx context.Context, id string) (userCred *userEntity.UserCredentials, err error)
		UpdateUserCredentials(ctx context.Context, userData dto.UserCredRequest, id string) error
		CreateUserCred(ctx context.Context, userData dto.UserCredRequest) error
		DeleteUserCred(ctx context.Context, id string) error
	}
)
