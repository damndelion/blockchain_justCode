package usecase

import (
	"context"
	"github.com/evrone/go-clean-template/internal/user/entity"
	"time"
)

type User struct {
	repo UserRepo
}

func NewUser(repo UserRepo) *User {
	return &User{repo}
}

func (u *User) Users(ctx context.Context) ([]*entity.User, error) {
	return u.repo.GetUsers(ctx)
}

func (u *User) CreateUser(ctx context.Context, user *entity.User) (int, error) {
	return u.repo.CreateUser(ctx, user)
}

//func (u *User) Register(ctx context.Context, name, email, password string) error {
//
//	generatedHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
//	if err != nil {
//		return err
//	}
//
//	_, err = u.repo.CreateUser(ctx, &entity.User{
//		Name:     name,
//		Email:    email,
//		Password: string(generatedHash),
//	})
//	if err != nil {
//		return err
//	}
//
//	return nil
//}

func (u *User) GetUserByEmail(ctx context.Context, email string) (*entity.User, error) {
	time.Sleep(2 * time.Second)

	return u.repo.GetUserByEmail(ctx, email)
}

func (u *User) GetUserById(ctx context.Context, id int) (*entity.User, error) {

	return u.repo.GetUserById(ctx, id)
}

//func (u *User) Login(ctx context.Context, email, password string) (*dto.LoginResponse, error) {
//	user, err := u.repo.GetUserByEmail(ctx, email)
//	switch {
//	case err == nil:
//	case errors.Is(err, pgx.ErrNoRows):
//		return nil, errors.New("user is not exist")
//	default:
//		return nil, err
//	}
//
//	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
//	if err != nil {
//		return nil, errors.New(fmt.Sprintf("passwords do not match %v", err))
//	}
//
//	accessTokenClaims := jwt.MapClaims{
//		"user_id":   user.Id,
//		"email":     user.Email,
//		"name":      user.Name,
//		"ExpiresAt": time.Now().Add(time.Hour * 1).Unix(),
//	}
//
//	accessToken := jwt.NewWithClaims(jwt.GetSigningMethod("HS256"), accessTokenClaims)
//
//	accessTokenString, err := accessToken.SignedString([]byte(u.cfg.SecretKey))
//	if err != nil {
//		return nil, err
//	}
//
//	refreshTokenClaims := jwt.MapClaims{
//		"user_id":   user.Id,
//		"ExpiresAt": time.Now().Add(time.Hour * 1),
//	}
//
//	refreshToken := jwt.NewWithClaims(jwt.GetSigningMethod("HS256"), refreshTokenClaims)
//
//	resfreshTokenString, err := refreshToken.SignedString([]byte(u.cfg.SecretKey))
//	if err != nil {
//		return nil, err
//	}
//
//	return &dto.LoginResponse{
//		Name:         user.Name,
//		Email:        user.Email,
//		AccessToken:  accessTokenString,
//		RefreshToken: resfreshTokenString,
//	}, nil
//}
