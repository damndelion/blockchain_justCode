package v1

import (
	"fmt"
	"github.com/evrone/go-clean-template/config/user"
	"github.com/evrone/go-clean-template/internal/user/controller/http/middleware"
	"github.com/evrone/go-clean-template/internal/user/controller/http/v1/dto"
	_ "github.com/evrone/go-clean-template/internal/user/entity"
	userEntity "github.com/evrone/go-clean-template/internal/user/entity"
	"github.com/evrone/go-clean-template/internal/user/usecase"
	"github.com/evrone/go-clean-template/pkg/cache"
	"github.com/evrone/go-clean-template/pkg/logger"
	"github.com/gin-gonic/gin"
	"net/http"
)

type adminRoutes struct {
	u         usecase.UserUseCase
	l         logger.Interface
	userCache cache.User
}

func newAdminRoutes(handler *gin.RouterGroup, u usecase.UserUseCase, l logger.Interface, uc cache.User, cfg *user.Config) {
	r := &userRoutes{u, l, uc}

	adminHandler := handler.Group("admin")
	{
		adminHandler.Use(middleware.AdminVerify(cfg.SecretKey, u))

		adminHandler.GET("/all", r.GetUsers)
		adminHandler.GET("/all/sort", r.GetUsersWithSort)
		adminHandler.GET("/all/search", r.GetUsersWithSearch)
		adminHandler.GET("/email", r.GetUserByEmail)
		adminHandler.GET("/user/:id", r.GetUserById)
		adminHandler.PUT("/user/:id", r.UpdateUser)
		adminHandler.POST("/user", r.CreateUser)
		adminHandler.DELETE("/user/:id", r.DeleteUser)

		adminHandler.GET("/info", r.GetUsersDetailInfo)
		adminHandler.GET("/info/sort", r.GetUsersInfoWithSort)
		adminHandler.GET("/info/search", r.GetUsersInfoWithSearch)
		adminHandler.GET("/info/:id", r.GetUserDetailInfoById)
		adminHandler.PUT("/info/:id", r.UpdateUserInfo)
		adminHandler.POST("/info", r.CreateUserInfo)
		adminHandler.DELETE("/info/:id", r.DeleteUserInfo)

		adminHandler.GET("/cred", r.GetUsersDetailCred)
		adminHandler.GET("/cred/sort", r.GetUsersCredWithSort)
		adminHandler.GET("/cred/search", r.GetUsersCredWithSearch)
		adminHandler.GET("/cred/:id", r.GetUserDetailCredById)
		adminHandler.PUT("/cred/:id", r.UpdateUserCredentials)
		adminHandler.POST("/cred", r.CreateUserCred)
		adminHandler.DELETE("/cred/:id", r.DeleteUserCred)

	}

}

// GetUsers godoc
// @Summary Get a list of all users
// @Description Retrieve a list of users from the system
// @Tags Users
// @Accept json
// @Produce json
// @Success 200 users body entity.User true  "List of users"
// @Failure 500 {string} string "Internal Server Error"
// @Router /v1/users/all [get]
func (ur *userRoutes) GetUsers(ctx *gin.Context) {
	var users []*userEntity.User
	var err error
	if len(ctx.Request.URL.Query()) == 0 {
		users, err = ur.u.Users(ctx)
		if err != nil {
			ur.l.Error(fmt.Errorf("http - v1 - user - getUsers: %w", err))
			errorResponse(ctx, http.StatusInternalServerError, "http - v1 - user - getUsers error")

			return
		}
	} else {
		users, err = ur.u.UsersWithFilter(ctx, ctx.Request.URL.Query())
	}

	if err != nil {
		ur.l.Error(fmt.Errorf("http - v1 - user - getUsers: %w", err))
		errorResponse(ctx, http.StatusInternalServerError, "http - v1 - user - getUsers error")

		return
	}

	ctx.JSON(http.StatusOK, users)
}

func (ur *userRoutes) GetUsersWithSort(ctx *gin.Context) {
	sortParam := ctx.Query("sort")
	methodParam := ctx.Query("method")
	users, err := ur.u.UsersWithSort(ctx, sortParam, methodParam)

	if err != nil {
		ur.l.Error(fmt.Errorf("http - v1 - user - getUsers: %w", err))
		errorResponse(ctx, http.StatusInternalServerError, "http - v1 - user - getUsers error")

		return
	}

	ctx.JSON(http.StatusOK, users)
}

func (ur *userRoutes) GetUsersInfoWithSort(ctx *gin.Context) {
	sortParam := ctx.Query("sort")
	methodParam := ctx.Query("method")
	users, err := ur.u.UsersInfoWithSort(ctx, sortParam, methodParam)

	if err != nil {
		ur.l.Error(fmt.Errorf("http - v1 - user - getUsers: %w", err))
		errorResponse(ctx, http.StatusInternalServerError, "http - v1 - user - getUsers error")

		return
	}

	ctx.JSON(http.StatusOK, users)
}

func (ur *userRoutes) GetUsersCredWithSort(ctx *gin.Context) {
	sortParam := ctx.Query("sort")
	methodParam := ctx.Query("method")
	users, err := ur.u.UsersCredWithSort(ctx, sortParam, methodParam)

	if err != nil {
		ur.l.Error(fmt.Errorf("http - v1 - user - getUsers: %w", err))
		errorResponse(ctx, http.StatusInternalServerError, "http - v1 - user - getUsers error")

		return
	}

	ctx.JSON(http.StatusOK, users)
}

func (ur *userRoutes) GetUsersWithSearch(ctx *gin.Context) {
	users, err := ur.u.UsersWithSearch(ctx, ctx.Request.URL.Query())
	if err != nil {
		ur.l.Error(fmt.Errorf("http - v1 - user - getUsers: %w", err))
		errorResponse(ctx, http.StatusInternalServerError, "http - v1 - user - getUsers error")

		return
	}

	ctx.JSON(http.StatusOK, users)
}

func (ur *userRoutes) GetUsersInfoWithSearch(ctx *gin.Context) {
	users, err := ur.u.UsersInfoWithSearch(ctx, ctx.Request.URL.Query())
	if err != nil {
		ur.l.Error(fmt.Errorf("http - v1 - user - getUsers: %w", err))
		errorResponse(ctx, http.StatusInternalServerError, "http - v1 - user - getUsers error")

		return
	}

	ctx.JSON(http.StatusOK, users)
}

func (ur *userRoutes) GetUsersCredWithSearch(ctx *gin.Context) {
	users, err := ur.u.UsersCredWithSearch(ctx, ctx.Request.URL.Query())
	if err != nil {
		ur.l.Error(fmt.Errorf("http - v1 - user - getUsers: %w", err))
		errorResponse(ctx, http.StatusInternalServerError, "http - v1 - user - getUsers error")

		return
	}

	ctx.JSON(http.StatusOK, users)
}

func (ur *userRoutes) CreateUser(ctx *gin.Context) {
	var user *userEntity.User

	err := ctx.ShouldBindJSON(&user)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, err)
		return
	}

	insertedID, err := ur.u.CreateUser(ctx, user)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, err)
		return
	}

	ctx.JSON(http.StatusOK, insertedID)
}

func (ur *userRoutes) UpdateUser(ctx *gin.Context) {
	id := ctx.Param("id")
	var userData dto.UserUpdateRequest
	err := ctx.ShouldBindJSON(&userData)
	if err != nil {
		ur.l.Error(fmt.Errorf("http - v1 - user - set user detail: %w", err))
		errorResponse(ctx, http.StatusBadRequest, "http - v1 - user - set user detail info error")
		return
	}

	err = ur.u.UpdateUser(ctx, userData, id)
	if err != nil {
		ur.l.Error(fmt.Errorf("http - v1 - blockchain - set user detail: %w", err))
		errorResponse(ctx, http.StatusBadRequest, "http - v1 - user - set user detail info error")
		return
	}

	ctx.JSON(http.StatusOK, "Success")
}

func (ur *userRoutes) DeleteUser(ctx *gin.Context) {
	id := ctx.Param("id")

	err := ur.u.DeleteUser(ctx, id)
	if err != nil {
		ur.l.Error(fmt.Errorf("http - v1 - blockchain - set user detail: %w", err))
		errorResponse(ctx, http.StatusBadRequest, "http - v1 - user - set user detail info error")
		return
	}

	ctx.JSON(http.StatusOK, "Success")
}

func (ur *userRoutes) GetUsersDetailInfo(ctx *gin.Context) {
	var usersInfo []*userEntity.UserInfo
	var err error
	if len(ctx.Request.URL.Query()) == 0 {
		usersInfo, err = ur.u.UsersInfo(ctx)
		if err != nil {
			ur.l.Error(fmt.Errorf("http - v1 - user - getUsers: %w", err))
			errorResponse(ctx, http.StatusInternalServerError, "http - v1 - user - getUsers error")

			return
		}
	} else {
		usersInfo, err = ur.u.UsersInfoWithFilter(ctx, ctx.Request.URL.Query())
	}

	if err != nil {
		ur.l.Error(fmt.Errorf("http - v1 - user - getUsers: %w", err))
		errorResponse(ctx, http.StatusInternalServerError, "http - v1 - user - getUsers error")

		return
	}

	ctx.JSON(http.StatusOK, usersInfo)

}

func (ur *userRoutes) DeleteUserInfo(ctx *gin.Context) {
	id := ctx.Param("id")

	err := ur.u.DeleteUserInfo(ctx, id)
	if err != nil {
		ur.l.Error(fmt.Errorf("http - v1 - blockchain - set user detail: %w", err))
		errorResponse(ctx, http.StatusBadRequest, "http - v1 - user - set user detail info error")
		return
	}

	ctx.JSON(http.StatusOK, "Success")
}

func (ur *userRoutes) DeleteUserCred(ctx *gin.Context) {
	id := ctx.Param("id")

	err := ur.u.DeleteUserCred(ctx, id)
	if err != nil {
		ur.l.Error(fmt.Errorf("http - v1 - blockchain - set user detail: %w", err))
		errorResponse(ctx, http.StatusBadRequest, "http - v1 - user - set user detail info error")
		return
	}

	ctx.JSON(http.StatusOK, "Success")
}
func (ur *userRoutes) GetUsersDetailCred(ctx *gin.Context) {
	var usersCred []*userEntity.UserCredentials
	var err error
	if len(ctx.Request.URL.Query()) == 0 {
		usersCred, err = ur.u.UsersCred(ctx)
		if err != nil {
			ur.l.Error(fmt.Errorf("http - v1 - user - getUsers: %w", err))
			errorResponse(ctx, http.StatusInternalServerError, "http - v1 - user - getUsers error")

			return
		}
	} else {
		usersCred, err = ur.u.UsersCredWithFilter(ctx, ctx.Request.URL.Query())
	}

	if err != nil {
		ur.l.Error(fmt.Errorf("http - v1 - user - getUsers: %w", err))
		errorResponse(ctx, http.StatusInternalServerError, "http - v1 - user - getUsers error")

		return
	}

	ctx.JSON(http.StatusOK, usersCred)
}

func (ur *userRoutes) UpdateUserInfo(ctx *gin.Context) {
	id := ctx.Param("id")
	var userData dto.UserInfoRequest
	err := ctx.ShouldBindJSON(&userData)
	if err != nil {
		ur.l.Error(fmt.Errorf("http - v1 - user - set user detail: %w", err))
		errorResponse(ctx, http.StatusBadRequest, "http - v1 - user - set user detail info error")
		return
	}

	err = ur.u.UpdateUserInfo(ctx, userData, id)
	if err != nil {
		ur.l.Error(fmt.Errorf("http - v1 - blockchain - set user detail: %w", err))
		errorResponse(ctx, http.StatusBadRequest, "http - v1 - user - set user detail info error")
		return
	}

	ctx.JSON(http.StatusOK, "Success")
}

func (ur *userRoutes) CreateUserInfo(ctx *gin.Context) {
	var userData dto.UserInfoRequest
	err := ctx.ShouldBindJSON(&userData)
	if err != nil {
		ur.l.Error(fmt.Errorf("http - v1 - user - set user detail: %w", err))
		errorResponse(ctx, http.StatusBadRequest, "http - v1 - user - set user detail info error")
		return
	}

	err = ur.u.CreateUserInfo(ctx, userData)
	if err != nil {
		ur.l.Error(fmt.Errorf("http - v1 - blockchain - set user detail: %w", err))
		errorResponse(ctx, http.StatusBadRequest, "http - v1 - user - set user detail info error")
		return
	}

	ctx.JSON(http.StatusOK, "Success")
}

func (ur *userRoutes) UpdateUserCredentials(ctx *gin.Context) {
	id := ctx.Param("id")
	var userData dto.UserCredRequest
	err := ctx.ShouldBindJSON(&userData)
	if err != nil {
		ur.l.Error(fmt.Errorf("http - v1 - user - set user detail: %w", err))
		errorResponse(ctx, http.StatusBadRequest, "http - v1 - user - set user detail info error")
		return
	}

	err = ur.u.UpdateUserCredentials(ctx, userData, id)
	if err != nil {
		ur.l.Error(fmt.Errorf("http - v1 - blockchain - set user detail: %w", err))
		errorResponse(ctx, http.StatusBadRequest, "http - v1 - user - set user detail info error")
		return
	}

	ctx.JSON(http.StatusOK, "Success")
}

func (ur *userRoutes) CreateUserCred(ctx *gin.Context) {
	var userData dto.UserCredRequest
	err := ctx.ShouldBindJSON(&userData)
	if err != nil {
		ur.l.Error(fmt.Errorf("http - v1 - user - set user detail: %w", err))
		errorResponse(ctx, http.StatusBadRequest, "http - v1 - user - set user detail info error")
		return
	}

	err = ur.u.CreateUserCred(ctx, userData)
	if err != nil {
		ur.l.Error(fmt.Errorf("http - v1 - blockchain - set user detail: %w", err))
		errorResponse(ctx, http.StatusBadRequest, "http - v1 - user - set user detail info error")
		return
	}

	ctx.JSON(http.StatusOK, "Success")
}

func (ur *userRoutes) GetUserDetailInfoById(ctx *gin.Context) {
	id := ctx.Param("id")
	userById, err := ur.u.GetUserInfoById(ctx, id)
	if err != nil {
		ur.l.Error(fmt.Errorf("http - v1 - userById - getUsersById: %w", err))
		errorResponse(ctx, http.StatusInternalServerError, "http - v1 - userById - getUsersById error")
		return
	}

	ctx.JSON(http.StatusOK, userById)
}

func (ur *userRoutes) GetUserDetailCredById(ctx *gin.Context) {
	id := ctx.Param("id")
	userById, err := ur.u.GetUserCredById(ctx, id)
	if err != nil {
		ur.l.Error(fmt.Errorf("http - v1 - userById - getUsersById: %w", err))
		errorResponse(ctx, http.StatusInternalServerError, "http - v1 - userById - getUsersById error")
		return
	}

	ctx.JSON(http.StatusOK, userById)
}
