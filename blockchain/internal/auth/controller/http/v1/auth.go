package v1

import (
	"fmt"
	"net/http"

	"github.com/evrone/go-clean-template/internal/auth/controller/http/v1/dto"
	"github.com/evrone/go-clean-template/internal/auth/usecase"
	"github.com/evrone/go-clean-template/pkg/logger"
	"github.com/gin-gonic/gin"
	"github.com/opentracing/opentracing-go"
)

type authRoutes struct {
	u usecase.AuthUseCase
	l logger.Interface
}

func newAuthRoutes(handler *gin.RouterGroup, u usecase.AuthUseCase, l logger.Interface) {
	r := &authRoutes{u, l}

	userHandler := handler.Group("/auth")
	{
		userHandler.POST("/register", r.Register)
		userHandler.POST("/login", r.Login)
		userHandler.POST("/refresh", r.Refresh)
		userHandler.POST("/confirm", r.Confirm)
	}
}

// Register godoc
// @Summary Register a new user
// @Description Register a new user in the system
// @Tags Auth
// @Accept json
// @Produce json
// @Param registerRequest body dto.RegisterRequest true "User registration request"
// @Success 200 {string} string "User successfully registered"
// @Failure 400 {string} string "Invalid input"
// @Failure 500 {string} string "Internal Server Error"
// @Router /v1/auth/register [post].
func (ar *authRoutes) Register(ctx *gin.Context) {
	span := opentracing.StartSpan("register handler")
	defer span.Finish()
	var registerRequest dto.RegisterRequest
	err := ctx.ShouldBindJSON(&registerRequest)

	if err != nil {
		ar.l.Error(fmt.Errorf("http - v1 - auth - register: %w", err))
		errorResponse(ctx, http.StatusBadRequest, "Registration form is not correct")

		return
	}
	spanCtx := opentracing.ContextWithSpan(ctx.Request.Context(), span)
	err = ar.u.Register(spanCtx, registerRequest.Name, registerRequest.Email, registerRequest.Password)
	if err != nil {
		ar.l.Error(fmt.Errorf("http - v1 - auth - register: %w", err))
		errorResponse(ctx, http.StatusInternalServerError, "Registration error")

		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "user successfully registered"})
}

// Login godoc
// @Summary User login
// @Description Authenticate a user and obtain an access token
// @Tags Auth
// @Accept json
// @Produce json
// @Param loginRequest body dto.LoginRequest true "User login request"
// @Success 200 {string} string "Access token"
// @Failure 400 {string} object dto.LoginResponse
// @Failure 500 {string} string "Internal Server Error"
// @Router /v1/auth/login [post].
func (ar *authRoutes) Login(ctx *gin.Context) {
	span := opentracing.StartSpan("login handler")
	defer span.Finish()

	var loginRequest dto.LoginRequest
	err := ctx.ShouldBindJSON(&loginRequest)
	if err != nil {
		ar.l.Error(fmt.Errorf("http - v1 - auth - login: %w", err))
		errorResponse(ctx, http.StatusBadRequest, "Login form error")

		return
	}
	spanCtx := opentracing.ContextWithSpan(ctx.Request.Context(), span)

	token, err := ar.u.Login(spanCtx, loginRequest.Email, loginRequest.Password)
	if err != nil {
		ar.l.Error(fmt.Errorf("http - v1 - auth - login: %w", err))
		errorResponse(ctx, http.StatusInternalServerError, "Login error")

		return
	}

	ctx.JSON(http.StatusOK, token)
}

// Refresh godoc
// @Summary User refresh token
// @Description Create new access token by refresh token
// @Tags Auth
// @Accept json
// @Produce json
// @Param refreshRequest body dto.RefreshRequest true "User refresh request"
// @Success 200 {string} string "Access token"
// @Failure 400 {string} string "Invalid input"
// @Failure 500 {string} string "Internal Server Error"
// @Router /v1/auth/refresh [post].
func (ar *authRoutes) Refresh(ctx *gin.Context) {
	span := opentracing.StartSpan("refresh handler")
	defer span.Finish()
	var refreshRequest dto.RefreshRequest
	err := ctx.ShouldBindJSON(&refreshRequest)
	spanCtx := opentracing.ContextWithSpan(ctx.Request.Context(), span)
	accessToken, refreshToken, err := ar.u.Refresh(spanCtx, refreshRequest.RefreshToken)
	if err != nil {
		ar.l.Error(fmt.Errorf("http - v1 - auth - refresh: %w", err))
		errorResponse(ctx, http.StatusInternalServerError, "Refresh error")

		return
	}
	res := dto.RefreshResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}

	ctx.JSON(http.StatusOK, res)
}

// Confirm godoc
// @Summary Confirm user
// @Description Confirm user by code
// @Tags Auth
// @Accept json
// @Produce json
// @Param refreshRequest body dto.ConfirmRequest true "User Confirm request"
// @Success 200 {string} string "User Confirmed"
// @Failure 400 {string} string "Invalid input"
// @Failure 500 {string} string "Internal Server Error"
// @Router /v1/auth/confirm [post].
func (ar *authRoutes) Confirm(ctx *gin.Context) {
	span := opentracing.StartSpan("confirm handler")
	defer span.Finish()
	var confirmRequest dto.ConfirmRequest
	err := ctx.ShouldBindJSON(&confirmRequest)
	spanCtx := opentracing.ContextWithSpan(ctx.Request.Context(), span)
	err = ar.u.ConfirmUserCode(spanCtx, confirmRequest.Email, confirmRequest.Code)
	if err != nil {
		ar.l.Error(fmt.Errorf("http - v1 - auth - refresh: %w", err))
		errorResponse(ctx, http.StatusInternalServerError, "Confirm error")

		return
	}

	ctx.JSON(http.StatusOK, "User Confirmed")
}
