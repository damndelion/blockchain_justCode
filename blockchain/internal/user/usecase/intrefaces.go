// Package usecase implements application business logic. Each logic group in own file.
package usecase

import (
	"context"
	"net/url"

	"github.com/damndelion/blockchain_justCode/internal/user/controller/http/v1/dto"
	userEntity "github.com/damndelion/blockchain_justCode/internal/user/entity"
)

type (

	// UserUseCase -.
	UserUseCase interface {
		Users(ctx context.Context) ([]*userEntity.User, error)
		UsersWithFilter(ctx context.Context, params url.Values) ([]*userEntity.User, error)
		UsersWithSort(ctx context.Context, sort, method string) ([]*userEntity.User, error)
		UsersWithSearch(ctx context.Context, params url.Values) ([]*userEntity.User, error)
		CreateUser(ctx context.Context, user dto.UserCreateRequest) (int, error)
		UpdateUser(ctx context.Context, userData dto.UserUpdateRequest, email string) error
		GetUserByEmail(ctx context.Context, email string) (*userEntity.User, error)
		GetUserByID(ctx context.Context, id int) (*userEntity.User, error)
		DeleteUser(ctx context.Context, id int) error

		UsersInfo(ctx context.Context) ([]*userEntity.UserInfo, error)
		UsersInfoWithFilter(ctx context.Context, params url.Values) ([]*userEntity.UserInfo, error)
		UsersInfoWithSort(ctx context.Context, sort, method string) ([]*userEntity.UserInfo, error)
		UsersInfoWithSearch(ctx context.Context, params url.Values) ([]*userEntity.UserInfo, error)
		GetUserInfoByID(ctx context.Context, id int) (*userEntity.UserInfo, error)
		CreateUserDetailInfo(ctx context.Context, userData dto.UserDetailRequest, id int) error
		SetUserDetailInfo(ctx context.Context, userData dto.UserDetailRequest, id int) error
		UpdateUserInfo(ctx context.Context, userData dto.UserUpdateInfoRequest, id int) error
		CreateUserInfo(ctx context.Context, userData dto.UserCreateInfoRequest) error
		DeleteUserInfo(ctx context.Context, id int) error

		UsersCred(ctx context.Context) ([]*userEntity.UserCredentials, error)
		UsersCredWithFilter(ctx context.Context, params url.Values) ([]*userEntity.UserCredentials, error)
		UsersCredWithSort(ctx context.Context, sort, method string) ([]*userEntity.UserCredentials, error)
		UsersCredWithSearch(ctx context.Context, params url.Values) ([]*userEntity.UserCredentials, error)
		GetUserCredByID(ctx context.Context, id int) (*userEntity.UserCredentials, error) //
		CreateUserCred(ctx context.Context, userData dto.UserCreateCredRequest) error
		UpdateUserCredentials(ctx context.Context, userData dto.UserUpdateCredRequest, id int) error
		DeleteUserCred(ctx context.Context, id int) error
	}

	// UserRepo -.
	UserRepo interface {
		GetUsers(ctx context.Context) ([]*userEntity.User, error)
		GetUsersWithFilter(ctx context.Context, param, value string) (users []*userEntity.User, err error)
		GetUsersWithSort(ctx context.Context, param, method string) (users []*userEntity.User, err error)
		GetUsersWithSearch(ctx context.Context, param string, value string) (users []*userEntity.User, err error)
		GetUserByEmail(ctx context.Context, email string) (*userEntity.User, error)
		GetUserByID(ctx context.Context, id int) (*userEntity.User, error)
		CreateUser(ctx context.Context, user dto.UserCreateRequest) (int, error)
		UpdateUser(ctx context.Context, userData dto.UserUpdateRequest, email string) error
		SetUserWallet(ctx context.Context, userID, address string) error
		DeleteUser(ctx context.Context, id int) error

		GetUsersDetailsInfo(ctx context.Context) (usersInfo []*userEntity.UserInfo, err error)
		GetUsersInfoWithFilter(ctx context.Context, param, value string) (users []*userEntity.UserInfo, err error)
		GetUsersInfoWithSort(ctx context.Context, param, method string) (users []*userEntity.UserInfo, err error)
		GetUsersInfoWithSearch(ctx context.Context, param string, value interface{}) (users []*userEntity.UserInfo, err error)
		GetUserInfoByID(ctx context.Context, id int) (userInfo *userEntity.UserInfo, err error)
		CreateUserDetailInfo(ctx context.Context, userData dto.UserDetailRequest, id int) error
		SetUserDetailInfo(ctx context.Context, userData dto.UserDetailRequest, id int) error
		CreateUserInfo(ctx context.Context, userData dto.UserCreateInfoRequest) error
		UpdateUserInfo(ctx context.Context, userData dto.UserUpdateInfoRequest, id int) error
		DeleteUserInfo(ctx context.Context, id int) error

		GetUsersCredentials(ctx context.Context) (usersInfo []*userEntity.UserCredentials, err error)
		GetUsersCredWithFilter(ctx context.Context, param, value string) (users []*userEntity.UserCredentials, err error)
		GetUsersCredWithSort(ctx context.Context, param, method string) (users []*userEntity.UserCredentials, err error)
		GetUsersCredWithSearch(ctx context.Context, param, value string) (users []*userEntity.UserCredentials, err error)
		GetUserCredByID(ctx context.Context, id int) (userCred *userEntity.UserCredentials, err error)
		UpdateUserCredentials(ctx context.Context, userData dto.UserUpdateCredRequest, id int) error
		CreateUserCred(ctx context.Context, userData dto.UserCreateCredRequest) error
		DeleteUserCred(ctx context.Context, id int) error
	}
)
