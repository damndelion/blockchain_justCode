package usecase

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"math/rand"
	"time"

	"github.com/damndelion/blockchain_justCode/config/auth"
	dtoConsumer "github.com/damndelion/blockchain_justCode/internal/auth/consumer/dto"
	"github.com/damndelion/blockchain_justCode/internal/auth/controller/http/v1/dto"
	authEntity "github.com/damndelion/blockchain_justCode/internal/auth/entity"
	"github.com/damndelion/blockchain_justCode/internal/nats"
	userEntity "github.com/damndelion/blockchain_justCode/internal/user/entity"
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
	span, spanCtx := opentracing.StartSpanFromContext(ctx, "register use case")
	defer span.Finish()
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

	timeout := time.After(60 * time.Second)

	select {
	case success := <-confirmationChan:
		delete(ConfirmationChannels, userIdentifier)
		if !success {
			return errors.New(fmt.Sprintf("user registration confirmation failed"))
		}
	case <-timeout:
		return errors.New("user registration timed out")
	}

	_, err = u.repo.CreateUser(spanCtx, &userEntity.User{
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
	accessToken, refreshToken, err := u.generateTokens(spanCtx, user)
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
	span, spanCtx := opentracing.StartSpanFromContext(ctx, "refresh use case")
	defer span.Finish()

	// Parse the refresh token
	token, err := jwt.Parse(refreshToken, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(u.cfg.SecretKey), nil
	})

	// Verify the token's claims
	var claims jwt.MapClaims
	var ok bool
	if claims, ok = token.Claims.(jwt.MapClaims); ok && token.Valid {
		if exp, ok := claims["exp"].(float64); ok {
			expirationTime := time.Unix(int64(exp), 0)
			if time.Now().After(expirationTime) {
				return "", "", fmt.Errorf("refresh token expired")
			}
		} else {
			return "", "", fmt.Errorf("invalid refresh token")
		}
	} else {
		return "", "", fmt.Errorf("invalid refresh token")
	}

	// Retrieve the user associated with the refresh token
	userEmail := fmt.Sprintf("%v", claims["email"])
	user, err := u.repo.GetUserByEmail(spanCtx, userEmail)
	if err != nil {
		return "", "", err
	}

	// Check if the refresh token in the database matches the provided one
	dbRefreshToken, err := u.repo.GetRefreshToken(spanCtx, user.ID)
	if err != nil {
		return "", "", err
	}

	if dbRefreshToken != refreshToken {
		return "", "", fmt.Errorf("invalid refresh token")
	}

	// Generate new access and refresh tokens
	accessToken, newRefreshToken, err := u.generateTokens(spanCtx, user)
	if err != nil {
		return "", "", err
	}

	// Update the refresh token in the database
	err = u.repo.UpdateRefreshToken(spanCtx, user.ID, newRefreshToken)
	if err != nil {
		return "", "", err
	}

	return accessToken, newRefreshToken, nil
}

func (u *Auth) generateTokens(ctx context.Context, user *userEntity.User) (string, string, error) {
	span, _ := opentracing.StartSpanFromContext(ctx, "generate token")
	defer span.Finish()
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
		"email": user.Email,
		"exp":   time.Now().Add(time.Duration(u.cfg.RefreshTokenTTL) * time.Second).Unix(),
	}

	refreshToken := jwt.NewWithClaims(jwt.GetSigningMethod("HS256"), refreshTokenClaims)

	refreshTokenString, err := refreshToken.SignedString([]byte(u.cfg.SecretKey))
	if err != nil {
		return "", "", err
	}
	userToken := authEntity.Token{
		UserID:       user.ID,
		RefreshToken: refreshTokenString,
	}

	err = u.repo.CreateUserToken(ctx, userToken)

	if err != nil {
		return "", "", err
	}

	return accessTokenString, refreshTokenString, nil
}

func (u *Auth) ConfirmUserCode(ctx context.Context, email string, userCode int) error {
	span, spanCtx := opentracing.StartSpanFromContext(ctx, "confirm user code usecase")
	defer span.Finish()
	code, err := u.repo.ConfirmCode(spanCtx, email)
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
