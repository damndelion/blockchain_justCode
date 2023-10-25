package repo

import (
	"context"
	"github.com/evrone/go-clean-template/internal/user/entity"
	"github.com/evrone/go-clean-template/pkg/postgres"
)

type UserRepo struct {
	*postgres.Postgres
}

func NewUserRepo(pg *postgres.Postgres) *UserRepo {
	return &UserRepo{pg}
}

func (ur *UserRepo) GetUsers(ctx context.Context) (users []*entity.User, err error) {

	res := ur.Find(&users)
	if res.Error != nil {
		return nil, res.Error
	}
	return users, nil
}

func (ur *UserRepo) CreateUser(ctx context.Context, user *entity.User) (int, error) {

	res := ur.WithContext(ctx).Create(user)
	if res.Error != nil {
		return 0, res.Error
	}
	return user.Id, nil
}

func (ur *UserRepo) GetUserByEmail(ctx context.Context, email string) (user *entity.User, err error) {
	res := ur.Where("email = ?", email).WithContext(ctx).Find(&user)
	if res.Error != nil {
		return nil, res.Error
	}
	return user, nil
}

func (ur *UserRepo) GetUserById(ctx context.Context, id int) (user *entity.User, err error) {
	res := ur.Where("id = ?", id).WithContext(ctx).Find(&user)
	if res.Error != nil {
		return nil, res.Error
	}
	return user, nil
}

func (ur *UserRepo) GetUserByID(ctx context.Context, id string) (user *entity.User, err error) {
	res := ur.WithContext(ctx).Where("id = ?", id).Find(&user)
	if res.Error != nil {
		return nil, res.Error
	}
	return user, nil
}
