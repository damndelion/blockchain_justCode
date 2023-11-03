package usecase

import (
	"context"
	"errors"
	"fmt"
	"github.com/evrone/go-clean-template/config/auth"
	"github.com/evrone/go-clean-template/internal/auth/controller/http/v1/dto"
	authEntity "github.com/evrone/go-clean-template/internal/auth/entity"
	userEntity "github.com/evrone/go-clean-template/internal/user/entity"
	"github.com/golang-jwt/jwt"
	"github.com/jackc/pgx/v4"
	"github.com/opentracing/opentracing-go"
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
	span, spanCtx := opentracing.StartSpanFromContext(ctx, "login use case")
	defer span.Finish()

	user, err := u.repo.GetUserByEmail(spanCtx, email)
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

	accessToken := jwt.NewWithClaims(jwt.GetSigningMethod("HS256"), accessTokenClaims)

	accessTokenString, err := accessToken.SignedString([]byte(u.cfg.SecretKey))
	if err != nil {
		return nil, err
	}

	refreshTokenClaims := jwt.MapClaims{
		"user_email": user.Email,
		"ExpiresAt":  time.Now().Add(time.Duration(u.cfg.RefreshTokenTTL) * time.Second).Unix(),
	}

	refreshToken := jwt.NewWithClaims(jwt.GetSigningMethod("HS256"), refreshTokenClaims)

	resfreshTokenString, err := refreshToken.SignedString([]byte(u.cfg.SecretKey))
	if err != nil {
		return nil, err
	}

	userToken := authEntity.Token{
		UserID:       user.Id,
		UserEmail:    user.Email,
		AccessToken:  accessTokenString,
		RefreshToken: resfreshTokenString,
	}

	err = u.repo.CreateUserToken(ctx, userToken)

	return &dto.LoginResponse{
		Name:         user.Name,
		Email:        user.Email,
		AccessToken:  accessTokenString,
		RefreshToken: resfreshTokenString,
	}, nil
}

func (u *Auth) Refresh(ctx context.Context, refreshToken string) (string, error) {
	// Step 1: Decode the refresh token
	claims := &jwt.StandardClaims{}
	token, err := jwt.ParseWithClaims(refreshToken, claims, func(token *jwt.Token) (interface{}, error) {
		// Check the signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(u.cfg.SecretKey), nil
	})

	// Check for errors during token parsing
	if err != nil || !token.Valid {
		return "", fmt.Errorf("Invalid refresh token: %v", err)
	}

	// Step 2: Get the user email from the claims
	userEmail := claims.Subject

	// Step 3: Fetch the user from the database using the userEmail
	user, err := u.repo.GetUserByEmail(ctx, userEmail)
	if err != nil {
		return "", err
	}

	// Step 4: Generate a new access token using the user data
	accessToken, err := u.generateAccessToken(user)
	if err != nil {
		return "", err
	}

	return accessToken, nil
}

func (u *Auth) generateAccessToken(user *userEntity.User) (string, error) {
	// Create a new set of claims for the access token
	accessTokenClaims := jwt.MapClaims{
		"user_id": user.Id,
		"email":   user.Email,
		"name":    user.Name,
		"exp":     time.Now().Add(time.Duration(u.cfg.AccessTokenTTL) * time.Second).Unix(),
	}

	// Sign the access token with the secret key
	access := jwt.NewWithClaims(jwt.SigningMethodHS256, accessTokenClaims)
	accessToken, err := access.SignedString([]byte(u.cfg.SecretKey))
	if err != nil {
		return "", err
	}

	return accessToken, nil
}
