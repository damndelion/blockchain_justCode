package usecase

import (
	"context"
	"errors"
	"fmt"
	"github.com/evrone/go-clean-template/config/auth"
	authEntity "github.com/evrone/go-clean-template/internal/auth/entity"
	"github.com/evrone/go-clean-template/internal/user/controller/http/v1/dto"
	userEntity "github.com/evrone/go-clean-template/internal/user/entity"
	"github.com/golang-jwt/jwt"
	"github.com/jackc/pgx/v4"
	"golang.org/x/crypto/bcrypt"
	"time"
)

type Auth struct {
	repo AuthRepo
	cfg  *auth.Config
}

func NewAuth(repo AuthRepo, cfg *auth.Config) *Auth {
	return &Auth{repo, cfg}
}

func (t *Auth) Register(ctx context.Context, name, email, password string) error {

	generatedHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	_, err = t.repo.CreateUser(ctx, &userEntity.User{
		Name:     name,
		Email:    email,
		Password: string(generatedHash),
	})
	if err != nil {
		return err
	}

	return nil
}

func (u *Auth) Login(ctx context.Context, email, password string) (*dto.LoginResponse, error) {
	user, err := u.repo.GetUserByEmail(ctx, email)
	switch {
	case err == nil:
	case errors.Is(err, pgx.ErrNoRows):
		return nil, errors.New("user is not exist")
	default:
		return nil, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return nil, errors.New(fmt.Sprintf("passwords do not match %v", err))
	}

	accessTokenClaims := jwt.MapClaims{
		"user_id": user.Id,
		"email":   user.Email,
		"name":    user.Name,
		"exp":     time.Now().Add(time.Duration(u.cfg.AccessTokenTTL) * time.Second).Unix(),
	}
	fmt.Println(time.Duration(u.cfg.AccessTokenTTL) * time.Second)

	accessToken := jwt.NewWithClaims(jwt.GetSigningMethod("HS256"), accessTokenClaims)

	accessTokenString, err := accessToken.SignedString([]byte(u.cfg.SecretKey))
	if err != nil {
		return nil, err
	}

	userToken := authEntity.Token{
		UserID: user.Id,
		Token:  accessTokenString,
	}
	err = u.repo.CreateUserToken(ctx, userToken)

	return &dto.LoginResponse{
		Name:        user.Name,
		Email:       user.Email,
		AccessToken: accessTokenString,
	}, nil
}
