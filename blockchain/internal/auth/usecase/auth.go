package usecase

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"math/rand"
	"time"

	"github.com/evrone/go-clean-template/config/auth"
	dtoConsumer "github.com/evrone/go-clean-template/internal/auth/consumer/dto"
	"github.com/evrone/go-clean-template/internal/auth/controller/http/v1/dto"
	authEntity "github.com/evrone/go-clean-template/internal/auth/entity"
	"github.com/evrone/go-clean-template/internal/nats"
	userEntity "github.com/evrone/go-clean-template/internal/user/entity"
	"github.com/golang-jwt/jwt"
	"github.com/opentracing/opentracing-go"
	"golang.org/x/crypto/bcrypt"
)

var ConfirmationChannels = make(map[string]chan bool)

type Auth struct {
	repo                     AuthRepo
	cfg                      *auth.Config
	userVerificationProducer *nats.Producer
}

func NewAuth(repo AuthRepo, cfg *auth.Config, userVerificationProducer *nats.Producer) *Auth {
	return &Auth{repo, cfg, userVerificationProducer}
}

func (u *Auth) Register(ctx context.Context, name, email, password string) error {
	err := u.repo.CheckForEmail(ctx, email)
	if err != nil {
		return err
	}
	randomFloat := rand.Float64()
	randomNumber := int(randomFloat * 10000)
	if randomNumber < 1000 {
		randomNumber += 1000
	}
	msg := dtoConsumer.UserCode{Email: email, Code: fmt.Sprintf("%d", randomNumber)}
	b, err := json.Marshal(&msg)
	if err != nil {
		return err
	}

	err = u.userVerificationProducer.ProduceMessage(b)
	if err != nil {
		return err
	}

	userIdentifier := email

	confirmationChan := make(chan bool)
	ConfirmationChannels[userIdentifier] = confirmationChan

	success := <-confirmationChan

	delete(ConfirmationChannels, userIdentifier)

	if !success {
		return errors.New(fmt.Sprintf("user registration confirmation failed: %v", err))
	}

	if err != nil {
		return err
	}
	_, err = u.repo.CreateUser(ctx, &userEntity.User{
		Name:     name,
		Email:    email,
		Password: password,
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
	if err != nil {
		return nil, err
	}
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return nil, errors.New(fmt.Sprintf("passwords do not match %v", err))
	}
	accessToken, refreshToken, err := u.generateTokens(ctx, user)
	if err != nil {
		return nil, err
	}

	return &dto.LoginResponse{
		Name:         user.Name,
		Email:        user.Email,
		Role:         user.Role,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

func (u *Auth) Refresh(ctx context.Context, refreshToken string) (string, string, error) {
	claims := &jwt.StandardClaims{}
	token, err := jwt.ParseWithClaims(refreshToken, claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New(fmt.Sprintf("unexpected signing method: %v", token.Header["alg"]))
		}

		return []byte(u.cfg.SecretKey), nil
	})

	if err != nil || !token.Valid {
		return "", "", errors.New(fmt.Sprintf("invalid refresh token: %v", err))
	}

	userEmail := claims.Subject

	user, err := u.repo.GetUserByEmail(ctx, userEmail)
	if err != nil {
		return "", "", err
	}

	accessToken, refreshToken, err := u.generateTokens(ctx, user)
	if err != nil {
		return "", "", err
	}

	return accessToken, refreshToken, nil
}

func (u *Auth) generateTokens(ctx context.Context, user *userEntity.User) (string, string, error) {
	accessTokenClaims := jwt.MapClaims{
		"user_id": user.ID,
		"email":   user.Email,
		"name":    user.Name,
		"role":    user.Role,
		"exp":     time.Now().Add(time.Duration(u.cfg.AccessTokenTTL) * time.Second).Unix(),
	}

	access := jwt.NewWithClaims(jwt.SigningMethodHS256, accessTokenClaims)
	accessTokenString, err := access.SignedString([]byte(u.cfg.SecretKey))
	if err != nil {
		return "", "", err
	}

	refreshTokenClaims := jwt.MapClaims{
		"user_email": user.Email,
		"ExpiresAt":  time.Now().Add(time.Duration(u.cfg.RefreshTokenTTL) * time.Second).Unix(),
	}

	refreshToken := jwt.NewWithClaims(jwt.GetSigningMethod("HS256"), refreshTokenClaims)

	refreshTokenString, err := refreshToken.SignedString([]byte(u.cfg.SecretKey))
	if err != nil {
		return "", "", err
	}
	userToken := authEntity.Token{
		UserID:       user.ID,
		UserEmail:    user.Email,
		AccessToken:  accessTokenString,
		RefreshToken: refreshTokenString,
	}

	err = u.repo.CreateUserToken(ctx, userToken)

	if err != nil {
		return "", "", err
	}

	return accessTokenString, refreshTokenString, nil
}

func (u *Auth) ConfirmUserCode(ctx context.Context, email string, userCode int) error {
	code, err := u.repo.ConfirmCode(ctx, email)
	if err != nil {
		return err
	}
	if userCode == code {
		userIdentifier := email
		confirmationChan, exists := ConfirmationChannels[userIdentifier]

		if exists {
			confirmationChan <- true
			delete(ConfirmationChannels, userIdentifier)

			return nil
		}

		return errors.New(fmt.Sprintf("confirmation channel not found"))
	} else {
		return errors.New(fmt.Sprintf("invalid user code"))
	}
}
