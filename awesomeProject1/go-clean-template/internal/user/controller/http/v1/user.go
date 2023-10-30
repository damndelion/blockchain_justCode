package v1

import (
	"fmt"
	"github.com/evrone/go-clean-template/internal/user/usecase"
	"github.com/evrone/go-clean-template/pkg/cache"
	"github.com/evrone/go-clean-template/pkg/logger"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

type userRoutes struct {
	u         usecase.UserUseCase
	l         logger.Interface
	userCache cache.User
}

func newUserRoutes(handler *gin.RouterGroup, u usecase.UserUseCase, l logger.Interface, uc cache.User) {
	r := &userRoutes{u, l, uc}

	adminHandler := handler.Group("user")
	{
		adminHandler.GET("/all", r.GetUsers)
		adminHandler.GET("/", r.GetUserByEmail)
		adminHandler.GET("/:id", r.GetUserById)
	}

}

func (ur *userRoutes) GetUsers(ctx *gin.Context) {
	users, err := ur.u.Users(ctx)
	if err != nil {
		ur.l.Error(fmt.Errorf("http - v1 - user - getUsers: %w", err))
		errorResponse(ctx, http.StatusInternalServerError, "http - v1 - user - getUsers error")

		return
	}

	ctx.JSON(http.StatusOK, users)
}

//func (ur *userRoutes) CreateUser(ctx *gin.Context) {
//	var user *entity.User
//
//	err := ctx.ShouldBindJSON(&user)
//	if err != nil {
//		ctx.JSON(http.StatusBadRequest, err)
//		return
//	}
//
//	insertedID, err := ur.u.CreateUser(ctx, user)
//	if err != nil {
//		ctx.JSON(http.StatusInternalServerError, err)
//		return
//	}
//
//	ctx.JSON(http.StatusOK, insertedID)
//}

func (ur *userRoutes) GetUserByEmail(ctx *gin.Context) {

	email := ctx.Query("email")
	fmt.Println(email)
	user, err := ur.userCache.Get(ctx, email)
	if err != nil {
		return
	}

	if user == nil {
		user, err = ur.u.GetUserByEmail(ctx, email)
		if err != nil {
			ur.l.Error(fmt.Errorf("http - v1 - user - getUsersByEmail: %w", err))
			errorResponse(ctx, http.StatusInternalServerError, "http - v1 - user - getUsersByEmail error")
			return
		}

		err = ur.userCache.Set(ctx, email, user)
		if err != nil {
			ur.l.Error(fmt.Errorf("http - v1 - user - getUsersByEmail: %w", err))
			errorResponse(ctx, http.StatusInternalServerError, "http - v1 - user - getUsersByEmail cache error")
		}
	}

	ctx.JSON(http.StatusOK, user)
}

func (ur *userRoutes) GetUserById(ctx *gin.Context) {

	id, _ := strconv.Atoi(ctx.Param("id"))

	user, err := ur.userCache.Get(ctx, strconv.Itoa(id))
	if err != nil {
		return
	}

	if user == nil {
		user, err = ur.u.GetUserById(ctx, id)
		if err != nil {
			ur.l.Error(fmt.Errorf("http - v1 - user - getUsersById: %w", err))
			errorResponse(ctx, http.StatusInternalServerError, "http - v1 - user - getUsersById error")
			return
		}

		err = ur.userCache.Set(ctx, strconv.Itoa(id), user)
		if err != nil {
			ur.l.Error(fmt.Errorf("http - v1 - user - getUsersById: %w", err))
			errorResponse(ctx, http.StatusInternalServerError, "http - v1 - user - getUsersById cache error")
		}
	}

	ctx.JSON(http.StatusOK, user)
}
