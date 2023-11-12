package usecase

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/evrone/go-clean-template/config/auth"
	dtoConsumer "github.com/evrone/go-clean-template/internal/auth/consumer/dto"
	"github.com/evrone/go-clean-template/internal/auth/controller/http/v1/dto"
	authEntity "github.com/evrone/go-clean-template/internal/auth/entity"
	"github.com/evrone/go-clean-template/internal/nats"
	userEntity "github.com/evrone/go-clean-template/internal/user/entity"
	"github.com/golang-jwt/jwt"
	"github.com/opentracing/opentracing-go"
	"golang.org/x/crypto/bcrypt"
	"math/rand"
	"time"
)

var confirmationChannels = make(map[string]chan bool)

type Auth struct {
	repo                     AuthRepo
	cfg                      *auth.Config
	userVerificationProducer *nats.Producer
}

func NewAuth(repo AuthRepo, cfg *auth.Config, userVerificationProducer *nats.Producer) *Auth {
	return &Auth{repo, cfg, userVerificationProducer}
}

func (t *Auth) Register(ctx context.Context, name, email, password string) error {
	randomFloat := rand.Float64()
	randomNumber := int(randomFloat * 10000)
	if randomNumber < 1000 {
		randomNumber += 1000
	}
	msg := dtoConsumer.UserCode{Email: email, Code: fmt.Sprintf("%d", randomNumber)}
	b, err := json.Marshal(&msg)
	if err != nil {
		return fmt.Errorf("failed to marshall UserCode err: %w", err)
	}

	t.userVerificationProducer.ProduceMessage(b)

	userIdentifier := email

	confirmationChan := make(chan bool)
	confirmationChannels[userIdentifier] = confirmationChan

	success := <-confirmationChan

	delete(confirmationChannels, userIdentifier)

	if !success {
		return errors.New("User registration confirmation failed")
	}

	if err != nil {
		return err
	}
	_, err = t.repo.CreateUser(ctx, &userEntity.User{
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

	return &dto.LoginResponse{
		Name:         user.Name,
		Email:        user.Email,
		Role:         user.Role,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

func (u *Auth) Refresh(ctx context.Context, refreshToken string) (string, error) {
	claims := &jwt.StandardClaims{}
	token, err := jwt.ParseWithClaims(refreshToken, claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(u.cfg.SecretKey), nil
	})

	if err != nil || !token.Valid {
		return "", fmt.Errorf("Invalid refresh token: %v", err)
	}

	userEmail := claims.Subject

	user, err := u.repo.GetUserByEmail(ctx, userEmail)
	if err != nil {
		return "", err
	}

	accessToken, refreshToken, err := u.generateTokens(ctx, user)
	if err != nil {
		return "", err
	}

	return accessToken, nil
}

func (u *Auth) generateTokens(ctx context.Context, user *userEntity.User) (string, string, error) {
	accessTokenClaims := jwt.MapClaims{
		"user_id": user.Id,
		"email":   user.Email,
		"name":    user.Name,
		"role":    user.Role,
		"exp":     time.Now().Add(time.Duration(u.cfg.AccessTokenTTL) * time.Second).Unix(),
	}

	access := jwt.NewWithClaims(jwt.SigningMethodHS256, accessTokenClaims)
	accessTokenString, err := access.SignedString([]byte(u.cfg.SecretKey))

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
		UserID:       user.Id,
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
		// Codes match, unlock the confirmation channel
		userIdentifier := email
		confirmationChan, exists := confirmationChannels[userIdentifier]

		if exists {
			confirmationChan <- true
			delete(confirmationChannels, userIdentifier)
			return nil
		} else {
			return errors.New("confirmation channel not found")
		}
	} else {
		return errors.New("Invalid user code")
	}

	if userCode == code {
		return nil
	} else {
		return errors.New("Code is mismatched")
	}
}
