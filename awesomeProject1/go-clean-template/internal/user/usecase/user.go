package usecase

import (
	"context"
	"fmt"
	"github.com/evrone/go-clean-template/config/user"
	"github.com/evrone/go-clean-template/internal/user/controller/http/v1/dto"
	"github.com/evrone/go-clean-template/internal/user/entity"
	"github.com/golang-jwt/jwt"
	"time"
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

func (u *User) CreateUser(ctx context.Context, user *userEntity.User) (int, error) {
	return u.repo.CreateUser(ctx, user)
}

func (u *User) GetUserByEmail(ctx context.Context, email string) (*userEntity.User, error) {
	time.Sleep(2 * time.Second)

	return u.repo.GetUserByEmail(ctx, email)
}

func (u *User) GetUserById(ctx context.Context, id string) (*userEntity.User, error) {

	return u.repo.GetUserByID(ctx, id)
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

func (u *User) GetIdFromToken(accessToken string) (string, error) {
	token, err := jwt.Parse(accessToken, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
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
