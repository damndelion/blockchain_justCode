package usecase

import (
	"context"
	"errors"
	"fmt"
	"github.com/evrone/go-clean-template/internal/usecase"
	"github.com/evrone/go-clean-template/internal/user/controller/http/v1/dto"
	"github.com/evrone/go-clean-template/internal/user/entity"
	"github.com/golang-jwt/jwt"
	"github.com/jackc/pgx/v4"
	"golang.org/x/crypto/bcrypt"
	"time"
)

type User struct {
	repo usecase.UserRepo
}

func NewUser(repo usecase.UserRepo) *User {
	return &User{repo: repo}
}

func (u *User) Users(ctx context.Context) ([]*entity.User, error) {
	return u.repo.GetUsers(ctx)
}

func (u *User) CreateUser(ctx context.Context, user *entity.User) (int, error) {
	return u.repo.CreateUser(ctx, user)
}

func (u *User) Register(ctx context.Context, name, email, password string) error {
	//email password
	//is email exists return with message "go to login"

	generatedHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	_, err = u.repo.CreateUser(ctx, &entity.User{
		Name:     name,
		Email:    email,
		Password: string(generatedHash),
	})
	if err != nil {
		return err
	}

	return nil
}

func (u *User) GetUserByEmail(ctx context.Context, email string) (*entity.User, error) {
	time.Sleep(2 * time.Second)

	return u.repo.GetUserByEmail(ctx, email)
}

func (u *User) GetUserById(ctx context.Context, id int) (*entity.User, error) {
	time.Sleep(2 * time.Second)

	return u.repo.GetUserById(ctx, id)
}

func (u *User) Login(ctx context.Context, email, password string) (*dto.LoginResponse, error) {
	//есть ли такой аккаунт с email =  email
	user, err := u.repo.GetUserByEmail(ctx, email)
	switch {
	case err == nil:
	case err == pgx.ErrNoRows:
		return nil, errors.New("user is not exist")
	default:
		return nil, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return nil, errors.New(fmt.Sprintf("passwords do not match %v", err))
	}

	claims := jwt.MapClaims{
		"email": user.Email,
		"name":  user.Name,
		"exp":   time.Now().Add(time.Hour * 1).Unix(),
	}

	token := jwt.NewWithClaims(jwt.GetSigningMethod("HS256"), claims)

	tokenString, err := token.SignedString([]byte("practice_7"))
	if err != nil {
		return nil, err
	}

	err = Cre(ctx, claims)
	if err != nil {
		return nil, fmt.Errorf("CreateUserToken err: %w", err)
	}

	return &dto.LoginResponse{
		Name:  user.Name,
		Email: user.Email,
		Token: tokenString,
	}, nil
}
