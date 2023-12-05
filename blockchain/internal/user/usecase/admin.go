package usecase

import (
	"context"
	"github.com/damndelion/blockchain_justCode/internal/user/controller/http/v1/dto"
	userEntity "github.com/damndelion/blockchain_justCode/internal/user/entity"
	"github.com/opentracing/opentracing-go"
	"net/url"
	"strconv"
)

func (u *User) UsersInfo(ctx context.Context) ([]*userEntity.UserInfo, error) {
	return u.repo.GetUsersDetailsInfo(ctx)
}

func (u *User) UsersInfoWithFilter(ctx context.Context, params url.Values) ([]*userEntity.UserInfo, error) {
	return u.repo.GetUsersInfoWithFilter(ctx, params.Get("filter"), params.Get("value"))
}

func (u *User) UsersCredWithFilter(ctx context.Context, params url.Values) ([]*userEntity.UserCredentials, error) {
	return u.repo.GetUsersCredWithFilter(ctx, params.Get("filter"), params.Get("value"))
}

func (u *User) UsersCred(ctx context.Context) ([]*userEntity.UserCredentials, error) {
	return u.repo.GetUsersCredentials(ctx)
}

func (u *User) CreateUserDetailInfo(ctx context.Context, userData dto.UserDetailRequest, id int) error {
	span, spanCtx := opentracing.StartSpanFromContext(ctx, "create user detail information use case")
	defer span.Finish()

	return u.repo.CreateUserDetailInfo(spanCtx, userData, id)
}

func (u *User) SetUserDetailInfo(ctx context.Context, userData dto.UserDetailRequest, id int) error {
	span, spanCtx := opentracing.StartSpanFromContext(ctx, "set user detail information use case")
	defer span.Finish()

	return u.repo.SetUserDetailInfo(spanCtx, userData, id)
}

func (u *User) DeleteUserInfo(ctx context.Context, id int) error {
	return u.repo.DeleteUserInfo(ctx, id)
}

func (u *User) UpdateUserInfo(ctx context.Context, userData dto.UserUpdateInfoRequest, id int) error {
	return u.repo.UpdateUserInfo(ctx, userData, id)
}

func (u *User) CreateUserInfo(ctx context.Context, userData dto.UserCreateInfoRequest) error {
	return u.repo.CreateUserInfo(ctx, userData)
}

func (u *User) CreateUserCred(ctx context.Context, userData dto.UserCreateCredRequest) error {
	return u.repo.CreateUserCred(ctx, userData)
}

func (u *User) UpdateUserCredentials(ctx context.Context, userData dto.UserUpdateCredRequest, id int) error {
	return u.repo.UpdateUserCredentials(ctx, userData, id)
}

func (u *User) DeleteUserCred(ctx context.Context, id int) error {
	return u.repo.DeleteUserCred(ctx, id)
}

func (u *User) UsersWithFilter(ctx context.Context, params url.Values) ([]*userEntity.User, error) {
	return u.repo.GetUsersWithFilter(ctx, params.Get("filter"), params.Get("value"))
}

func (u *User) UsersWithSearch(ctx context.Context, params url.Values) ([]*userEntity.User, error) {
	searchParam := params.Get("search")
	searchValueStr := params.Get("value")

	return u.repo.GetUsersWithSearch(ctx, searchParam, searchValueStr)
}

func (u *User) UsersWithSort(ctx context.Context, sort, method string) ([]*userEntity.User, error) {
	return u.repo.GetUsersWithSort(ctx, sort, method)
}

func (u *User) UsersInfoWithSearch(ctx context.Context, params url.Values) ([]*userEntity.UserInfo, error) {
	searchParam := params.Get("search")
	searchValueStr := params.Get("value")

	var searchValue interface{}
	var err error

	searchValue, err = strconv.Atoi(searchValueStr)
	if err != nil {
		searchValue = searchValueStr
	}

	return u.repo.GetUsersInfoWithSearch(ctx, searchParam, searchValue)
}

func (u *User) UsersInfoWithSort(ctx context.Context, sort, method string) ([]*userEntity.UserInfo, error) {
	return u.repo.GetUsersInfoWithSort(ctx, sort, method)
}

func (u *User) UsersCredWithSearch(ctx context.Context, params url.Values) ([]*userEntity.UserCredentials, error) {
	return u.repo.GetUsersCredWithSearch(ctx, params.Get("search"), params.Get("value"))
}

func (u *User) UsersCredWithSort(ctx context.Context, sort, method string) ([]*userEntity.UserCredentials, error) {
	return u.repo.GetUsersCredWithSort(ctx, sort, method)
}

func (u *User) GetUserInfoByID(ctx context.Context, id int) (*userEntity.UserInfo, error) {
	span, spanCtx := opentracing.StartSpanFromContext(ctx, "get user detail information by id use case")
	defer span.Finish()

	return u.repo.GetUserInfoByID(spanCtx, id)
}

func (u *User) GetUserCredByID(ctx context.Context, id int) (*userEntity.UserCredentials, error) {
	span, spanCtx := opentracing.StartSpanFromContext(ctx, "get user credentials by id use case")
	defer span.Finish()

	return u.repo.GetUserCredByID(spanCtx, id)
}
