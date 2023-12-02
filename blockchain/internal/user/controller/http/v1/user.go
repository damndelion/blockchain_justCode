package v1

import (
	"fmt"
	"net/http"

	"github.com/damndelion/blockchain_justCode/config/user"
	authMiddlware "github.com/damndelion/blockchain_justCode/internal/user/controller/http/middleware"
	"github.com/damndelion/blockchain_justCode/internal/user/controller/http/v1/dto"
	_ "github.com/damndelion/blockchain_justCode/internal/user/entity"
	"github.com/damndelion/blockchain_justCode/internal/user/usecase"
	"github.com/damndelion/blockchain_justCode/pkg/cache"
	"github.com/damndelion/blockchain_justCode/pkg/logger"
	"github.com/gin-gonic/gin"
	"github.com/opentracing/opentracing-go"
)

type userRoutes struct {
	u         usecase.UserUseCase
	l         logger.Interface
	userCache cache.User
}

func newUserRoutes(handler *gin.RouterGroup, u usecase.UserUseCase, l logger.Interface, uc cache.User, cfg *user.Config) {
	r := &userRoutes{u, l, uc}

	userHandler := handler.Group("user")
	{
		userHandler.Use(authMiddlware.JwtVerify(cfg.SecretKey))

		userHandler.GET("/", r.GetUser)
		userHandler.GET("/info", r.GetUserDetailInfo)
		userHandler.GET("/cred", r.GetUserDetailCred)
		userHandler.POST("/info", r.CreateUserDetailInfo)
		userHandler.PUT("/info", r.SetUserDetailInfo)
	}
}

// GetUser godoc
// @Summary Get a user by jwt token
// @Description Retrieve a user by their authorization token
// @Tags Users
// @Accept json
// @Produce json
// @Param authorization header string true "JWT token"
// @Success 200 user body entity.User true "User information"
// @Failure 400 {string} Bad Request
// @Failure 401 {string} string "Unauthorized"
// @Failure 500 {string} string "Internal Server Error"
// @Router /v1/user [get].
func (ur *userRoutes) GetUser(ctx *gin.Context) {
	span := opentracing.StartSpan("get user handler")
	defer span.Finish()
	spanCtx := opentracing.ContextWithSpan(ctx.Request.Context(), span)
	userID, _ := ctx.Get("user_id")

	resUser, err := ur.u.GetUserByID(spanCtx, userID.(string))
	if err != nil {
		ur.l.Error(fmt.Errorf("http - v1 - user - getUsersById: %w", err))
		errorResponse(ctx, http.StatusInternalServerError, "getUsersById error")

		return
	}

	ctx.JSON(http.StatusOK, resUser)
}

// GetUserDetailInfo godoc
// @Summary Get a user information by jwt token
// @Description Retrieve a user information by their authorization token
// @Tags Users Information
// @Accept json
// @Produce json
// @Param authorization header string true "JWT token"
// @Success 200 user body entity.UserInfo true "User information"
// @Failure 400 {string} Bad Request
// @Failure 401 {string} string "Unauthorized"
// @Failure 500 {string} string "Internal Server Error"
// @Router /v1/user/info [get].
func (ur *userRoutes) GetUserDetailInfo(ctx *gin.Context) {
	span := opentracing.StartSpan("get user detail information handler")
	defer span.Finish()
	spanCtx := opentracing.ContextWithSpan(ctx.Request.Context(), span)
	userID, _ := ctx.Get("user_id")

	userByID, err := ur.u.GetUserInfoByID(spanCtx, userID.(string))
	if err != nil {
		ur.l.Error(fmt.Errorf("http - v1 - userById - getUsersById: %w", err))
		errorResponse(ctx, http.StatusInternalServerError, "getUsersById error")

		return
	}

	ctx.JSON(http.StatusOK, userByID)
}

// GetUserDetailCred godoc
// @Summary Get a user credentials by jwt token
// @Description Retrieve a user credentials by their authorization token
// @Tags Users Credentials
// @Accept json
// @Produce json
// @Param authorization header string true "JWT token"
// @Success 200 user body entity.UserCred true "User credentials"
// @Failure 400 {string} Bad Request
// @Failure 401 {string} string "Unauthorized"
// @Failure 500 {string} string "Internal Server Error"
// @Router /v1/user/cred [get].
func (ur *userRoutes) GetUserDetailCred(ctx *gin.Context) {
	span := opentracing.StartSpan("get user credentials handler")
	defer span.Finish()
	spanCtx := opentracing.ContextWithSpan(ctx.Request.Context(), span)
	userID, _ := ctx.Get("user_id")

	userByID, err := ur.u.GetUserCredByID(spanCtx, userID.(string))
	if err != nil {
		ur.l.Error(fmt.Errorf("http - v1 - userById - getUsersById: %w", err))
		errorResponse(ctx, http.StatusInternalServerError, "getUsersById error")

		return
	}

	ctx.JSON(http.StatusOK, userByID)
}

// CreateUserDetailInfo godoc
// @Summary Create both user info and cred
// @Description Creates user info and cred, sets users valid to true
// @Tags Users Information
// @Accept json
// @Produce json
// @Param authorization header string true "JWT token"
// @Param data body dto.UserDetailRequest true "JSON data"
// @Success 200 {string} string "Success""
// @Failure 400 {string} Bad Request
// @Failure 401 {string} string "Unauthorized"
// @Failure 500 {string} string "Internal Server Error"
// @Router /v1/user/info [post].
func (ur *userRoutes) CreateUserDetailInfo(ctx *gin.Context) {
	span := opentracing.StartSpan("create user detail information handler")
	defer span.Finish()
	spanCtx := opentracing.ContextWithSpan(ctx.Request.Context(), span)
	userID, _ := ctx.Get("user_id")

	var userData dto.UserDetailRequest
	err := ctx.ShouldBindJSON(&userData)

	if err != nil {
		ur.l.Error(fmt.Errorf("http - v1 - blockchain - create user detail: %w", err))
		errorResponse(ctx, http.StatusBadRequest, "create user detail info error")

		return
	}
	err = ur.u.CreateUserDetailInfo(spanCtx, userData, userID.(string))
	if err != nil {
		ur.l.Error(fmt.Errorf("http - v1 - blockchain - create user detail: %w", err))
		errorResponse(ctx, http.StatusBadRequest, "create user detail info error")

		return
	}
	ctx.JSON(http.StatusOK, "Success")
}

// SetUserDetailInfo godoc
// @Summary Update both user info and cred
// @Description Update user info and cred
// @Tags Users Information
// @Accept json
// @Produce json
// @Param authorization header string true "JWT token"
// @Param data body dto.UserDetailRequest true "JSON data"
// @Success 200 {string} string "Success""
// @Failure 400 {string} Bad Request
// @Failure 401 {string} string "Unauthorized"
// @Failure 500 {string} string "Internal Server Error"
// @Router /v1/user/info [put].
func (ur *userRoutes) SetUserDetailInfo(ctx *gin.Context) {
	span := opentracing.StartSpan("set user detail information handler")
	defer span.Finish()
	spanCtx := opentracing.ContextWithSpan(ctx.Request.Context(), span)
	userID, _ := ctx.Get("user_id")

	var userData dto.UserDetailRequest
	err := ctx.ShouldBindJSON(&userData)

	if err != nil {
		ur.l.Error(fmt.Errorf("http - v1 - user - set user detail: %w", err))
		errorResponse(ctx, http.StatusBadRequest, "set user detail info error")

		return
	}
	err = ur.u.SetUserDetailInfo(spanCtx, userData, userID.(string))
	if err != nil {
		ur.l.Error(fmt.Errorf("http - v1 - blockchain - set user detail: %w", err))
		errorResponse(ctx, http.StatusBadRequest, "set user detail info error")

		return
	}
	ctx.JSON(http.StatusOK, "Success")
}
