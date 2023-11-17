package v1

import (
	"fmt"
	"net/http"

	"github.com/evrone/go-clean-template/config/user"
	"github.com/evrone/go-clean-template/internal/user/controller/http/middleware"
	"github.com/evrone/go-clean-template/internal/user/controller/http/v1/dto"
	_ "github.com/evrone/go-clean-template/internal/user/entity"
	userEntity "github.com/evrone/go-clean-template/internal/user/entity"
	"github.com/evrone/go-clean-template/internal/user/usecase"
	"github.com/evrone/go-clean-template/pkg/cache"
	"github.com/evrone/go-clean-template/pkg/logger"
	"github.com/gin-gonic/gin"
)

type adminRoutes struct {
	u         usecase.UserUseCase
	l         logger.Interface
	userCache cache.User
}

func newAdminRoutes(handler *gin.RouterGroup, u usecase.UserUseCase, l logger.Interface, uc cache.User, cfg *user.Config) {
	r := &adminRoutes{u, l, uc}

	adminHandler := handler.Group("admin")
	{
		adminHandler.Use(middleware.AdminVerify(cfg.SecretKey, u))

		adminHandler.GET("/all", r.GetUsers)
		adminHandler.GET("/user/:id", r.GetUserByID)
		adminHandler.GET("/email", r.GetUserByEmail)
		adminHandler.GET("/all/sort", r.GetUsersWithSort)
		adminHandler.GET("/all/search", r.GetUsersWithSearch)
		adminHandler.POST("/user/:id", r.UpdateUser)
		adminHandler.PUT("/user", r.CreateUser)
		adminHandler.DELETE("/user/:id", r.DeleteUser)

		adminHandler.GET("/info", r.GetUsersDetailInfo)
		adminHandler.GET("/info/sort", r.GetUsersInfoWithSort)
		adminHandler.GET("/info/search", r.GetUsersInfoWithSearch)
		adminHandler.GET("/info/:id", r.GetUserDetailInfoByID)
		adminHandler.POST("/info/:id", r.UpdateUserInfo)
		adminHandler.PUT("/info", r.CreateUserInfo)
		adminHandler.DELETE("/info/:id", r.DeleteUserInfo)

		adminHandler.GET("/cred", r.GetUsersDetailCred)
		adminHandler.GET("/cred/sort", r.GetUsersCredWithSort)
		adminHandler.GET("/cred/search", r.GetUsersCredWithSearch)
		adminHandler.GET("/cred/:id", r.GetUserDetailCredByID)
		adminHandler.POST("/cred/:id", r.UpdateUserCred)
		adminHandler.PUT("/cred", r.CreateUserCred)
		adminHandler.DELETE("/cred/:id", r.DeleteUserCred)
	}
}

// GetUsers godoc
// @Summary Get a list of all users
// @Description Retrieve a list of users from the system
// @Tags Users
// @Accept json
// @Produce json
// @Param authorization header string true "JWT token"
// @Param filter query string false "Filter column name"
// @Param value query string false "Filter value"
// @Success 200 users body entity.User true  "List of users"
// @Failure 400 {string} string "Bad Request"
// @Failure 401 {string} string "Unauthorized"
// @Failure 500 {string} string "Internal Server Error"
// @Router /v1/admin/all [get].
func (ur *adminRoutes) GetUsers(ctx *gin.Context) {
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

// GetUserByID godoc
// @Summary Get a user by ID
// @Description Retrieve a user by their unique ID
// @Tags Users
// @Accept json
// @Produce json
// @Param authorization header string true "JWT token"
// @Param id path int true "ID of the item"
// @Success 200 user body entity.User true "User"
// @Failure 400 {string} string "Bad Request"
// @Failure 401 {string} string "Unauthorized"
// @Failure 500 {string} string "Internal Server Error"
// @Router /v1/admin/user/{id} [get].
func (ur *adminRoutes) GetUserByID(ctx *gin.Context) {
	id := ctx.Param("id")

	resUser, err := ur.userCache.Get(ctx, id)
	if err != nil {
		return
	}

	if resUser == nil {
		resUser, err = ur.u.GetUserByID(ctx, id)
		if err != nil {
			ur.l.Error(fmt.Errorf("http - v1 - user - getUsersById: %w", err))
			errorResponse(ctx, http.StatusInternalServerError, "http - v1 - user - getUsersById error")

			return
		}

		err = ur.userCache.Set(ctx, id, resUser)
		if err != nil {
			ur.l.Error(fmt.Errorf("http - v1 - user - getUsersById: %w", err))
			errorResponse(ctx, http.StatusInternalServerError, "http - v1 - user - getUsersById cache error")
		}
	}

	ctx.JSON(http.StatusOK, resUser)
}

// GetUserByEmail godoc
// @Summary Get a user by email
// @Description Retrieve a user by their email address
// @Tags Users
// @Accept json
// @Produce json
// @Param authorization header string true "JWT token"
// @Param email query string true "Email address of the user"
// @Success 200 users body entity.User true "User"
// @Failure 400 {string} string "Bad Request"
// @Failure 401 {string} string "Unauthorized"
// @Failure 500 {string} string "Internal Server Error"
// @Router /v1/admin/email [get].
func (ur *adminRoutes) GetUserByEmail(ctx *gin.Context) {
	email := ctx.Query("email")
	resUser, err := ur.userCache.Get(ctx, email)
	if err != nil {
		return
	}

	if resUser == nil {
		resUser, err = ur.u.GetUserByEmail(ctx, email)
		if err != nil {
			ur.l.Error(fmt.Errorf("http - v1 - user - getUsersByEmail: %w", err))
			errorResponse(ctx, http.StatusInternalServerError, "http - v1 - user - getUsersByEmail error")

			return
		}

		err = ur.userCache.Set(ctx, email, resUser)
		if err != nil {
			ur.l.Error(fmt.Errorf("http - v1 - user - getUsersByEmail: %w", err))
			errorResponse(ctx, http.StatusInternalServerError, "http - v1 - user - getUsersByEmail cache error")
		}
	}

	ctx.JSON(http.StatusOK, resUser)
}

// GetUsersWithSort godoc
// @Summary Get a users with sort
// @Description Retrieve a user by sorting
// @Tags Users
// @Accept json
// @Produce json
// @Param authorization header string true "JWT token"
// @Param sort query string true "Sort value"
// @Param method query string true "Ascending or Descending"
// @Success 200 users body entity.User true "User"
// @Failure 400 {string} string "Bad Request"
// @Failure 401 {string} string "Unauthorized"
// @Failure 500 {string} string "Internal Server Error"
// @Router /v1/admin/all/sort [get].
func (ur *adminRoutes) GetUsersWithSort(ctx *gin.Context) {
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

// GetUsersWithSearch godoc
// @Summary Get a users with search
// @Description Retrieve a user by searching value
// @Tags Users
// @Accept json
// @Produce json
// @Param authorization header string true "JWT token"
// @Param search query string true "Search column"
// @Param value query string true "Search value"
// @Success 200 users body entity.User true "User"
// @Failure 400 {string} string "Bad Request"
// @Failure 401 {string} string "Unauthorized"
// @Failure 500 {string} string "Internal Server Error"
// @Router /v1/admin/all/search [get].
func (ur *adminRoutes) GetUsersWithSearch(ctx *gin.Context) {
	users, err := ur.u.UsersWithSearch(ctx, ctx.Request.URL.Query())
	if err != nil {
		ur.l.Error(fmt.Errorf("http - v1 - user - getUsers: %w", err))
		errorResponse(ctx, http.StatusInternalServerError, "http - v1 - user - getUsers error")

		return
	}

	ctx.JSON(http.StatusOK, users)
}

// UpdateUser godoc
// @Summary Update user
// @Description Update user
// @Tags Users
// @Accept json
// @Produce json
// @Param authorization header string true "JWT token"
// @Param id path int true "ID of the item"
// @Param data body dto.UserUpdateRequest true "JSON data"
// @Success 200 {string} string "Success"
// @Failure 400 {string} Bad Request
// @Failure 401 {string} string "Unauthorized"
// @Failure 500 {string} string "Internal Server Error"
// @Router /v1/admin/user/{id} [post].
func (ur *adminRoutes) UpdateUser(ctx *gin.Context) {
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

// CreateUser godoc
// @Summary Create user
// @Description Create user
// @Tags Users
// @Accept json
// @Produce json
// @Param authorization header string true "JWT token"
// @Param data body dto.UserUpdateRequest true "JSON data"
// @Success 200 {int} int id
// @Failure 400 {string} Bad Request
// @Failure 401 {string} string "Unauthorized"
// @Failure 500 {string} string "Internal Server Error"
// @Router /v1/admin/user [put].
func (ur *adminRoutes) CreateUser(ctx *gin.Context) {
	var userData dto.UserUpdateRequest

	err := ctx.ShouldBindJSON(&userData)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, err)

		return
	}

	insertedID, err := ur.u.CreateUser(ctx, userData)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, err)

		return
	}

	ctx.JSON(http.StatusOK, insertedID)
}

// DeleteUser godoc
// @Summary Delete user by id
// @Description Delete user by id and linked info and cred
// @Tags Users
// @Param authorization header string true "JWT token"
// @Param id path int true "ID of the item"
// @Success 200 {string} string "Success"
// @Failure 400 {string} Bad Request
// @Failure 401 {string} string "Unauthorized"
// @Failure 500 {string} string "Internal Server Error"
// @Router /v1/admin/user/{id} [delete].
func (ur *adminRoutes) DeleteUser(ctx *gin.Context) {
	id := ctx.Param("id")

	err := ur.u.DeleteUser(ctx, id)
	if err != nil {
		ur.l.Error(fmt.Errorf("http - v1 - blockchain - set user detail: %w", err))
		errorResponse(ctx, http.StatusBadRequest, "http - v1 - user - set user detail info error")

		return
	}

	ctx.JSON(http.StatusOK, "Success")
}

// GetUsersDetailInfo godoc
// @Summary Get a list of all user information
// @Description Get a list of all user information
// @Tags Users Information
// @Produce json
// @Param authorization header string true "JWT token"
// @Param filter query string false "Filter column name"
// @Param value query string false "Filter value"
// @Success 200 users body entity.UserInfo true  "List of users information"
// @Failure 400 {string} Bad Request
// @Failure 401 {string} string "Unauthorized"
// @Failure 500 {string} string "Internal Server Error"
// @Router /v1/admin/info [get].
func (ur *adminRoutes) GetUsersDetailInfo(ctx *gin.Context) {
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

// GetUsersInfoWithSort godoc
// @Summary Get a users information with sort
// @Description Retrieve a user information by sorting
// @Tags Users Information
// @Produce json
// @Param authorization header string true "JWT token"
// @Param sort query string true "Sort value"
// @Param method query string true "Ascending or Descending"
// @Success 200 users body entity.UserInfo true "User information"
// @Failure 400 {string} string "Bad Request"
// @Failure 401 {string} string "Unauthorized"
// @Failure 500 {string} string "Internal Server Error"
// @Router /v1/admin/info/sort [get].
func (ur *adminRoutes) GetUsersInfoWithSort(ctx *gin.Context) {
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

// GetUsersInfoWithSearch godoc
// @Summary Get a users information with search
// @Description Retrieve a user information by searching value
// @Tags Users Information
// @Produce json
// @Param authorization header string true "JWT token"
// @Param search query string true "Search column"
// @Param value query string true "Search value"
// @Success 200 users body entity.UserInfo true "User information"
// @Failure 400 {string} string "Bad Request"
// @Failure 401 {string} string "Unauthorized"
// @Failure 500 {string} string "Internal Server Error"
// @Router /v1/admin/info/search [get].
func (ur *adminRoutes) GetUsersInfoWithSearch(ctx *gin.Context) {
	users, err := ur.u.UsersInfoWithSearch(ctx, ctx.Request.URL.Query())
	if err != nil {
		ur.l.Error(fmt.Errorf("http - v1 - user - getUsers: %w", err))
		errorResponse(ctx, http.StatusInternalServerError, "http - v1 - user - getUsers error")

		return
	}

	ctx.JSON(http.StatusOK, users)
}

// GetUserDetailInfoByID godoc
// @Summary Get a user information by ID
// @Description Retrieve a user information by their unique ID
// @Tags Users Information
// @Produce json
// @Param authorization header string true "JWT token"
// @Param id path int true "ID of the item"
// @Success 200 user body entity.UserInfo true "User information"
// @Failure 400 {string} Bad Request
// @Failure 401 {string} string "Unauthorized"
// @Failure 500 {string} string "Internal Server Error"
// @Router /v1/admin/info/{id} [get].
func (ur *adminRoutes) GetUserDetailInfoByID(ctx *gin.Context) {
	id := ctx.Param("id")
	userByID, err := ur.u.GetUserInfoByID(ctx, id)
	if err != nil {
		ur.l.Error(fmt.Errorf("http - v1 - userById - getUsersById: %w", err))
		errorResponse(ctx, http.StatusInternalServerError, "http - v1 - userById - getUsersById error")

		return
	}

	ctx.JSON(http.StatusOK, userByID)
}

// UpdateUserInfo godoc
// @Summary Update user information
// @Description Update user information
// @Tags Users Information
// @Accept json
// @Param authorization header string true "JWT token"
// @Param id path int true "ID of the item"
// @Param data body dto.UserInfoRequest true "JSON data"
// @Success 200 {string} string "Success"
// @Failure 400 {string} Bad Request
// @Failure 401 {string} string "Unauthorized"
// @Failure 500 {string} string "Internal Server Error"
// @Router /v1/admin/info/{id} [post].
func (ur *adminRoutes) UpdateUserInfo(ctx *gin.Context) {
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

// CreateUserInfo godoc
// @Summary Create user information
// @Description Create user information
// @Tags Users Information
// @Accept json
// @Produce json
// @Param authorization header string true "JWT token"
// @Param data body dto.UserInfoRequest true "JSON data"
// @Success 200 {int} int id
// @Failure 400 {string} Bad Request
// @Failure 401 {string} string "Unauthorized"
// @Failure 500 {string} string "Internal Server Error"
// @Router /v1/admin/info [put].
func (ur *adminRoutes) CreateUserInfo(ctx *gin.Context) {
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

// DeleteUserInfo godoc
// @Summary Delete user information by user id
// @Description Delete user information by user id
// @Tags Users Information
// @Param authorization header string true "JWT token"
// @Param id path int true "ID of the item"
// @Success 200 {string} string "Success"
// @Failure 400 {string} Bad Request
// @Failure 401 {string} string "Unauthorized"
// @Failure 500 {string} string "Internal Server Error"
// @Router /v1/admin/info/{id} [delete].
func (ur *adminRoutes) DeleteUserInfo(ctx *gin.Context) {
	id := ctx.Param("id")

	err := ur.u.DeleteUserInfo(ctx, id)
	if err != nil {
		ur.l.Error(fmt.Errorf("http - v1 - blockchain - set user detail: %w", err))
		errorResponse(ctx, http.StatusBadRequest, "http - v1 - user - set user detail info error")

		return
	}

	ctx.JSON(http.StatusOK, "Success")
}

// GetUsersDetailCred godoc
// @Summary Get a list of all user credentials
// @Description Get a list of all user credentials
// @Tags Users Credentials
// @Produce json
// @Param authorization header string true "JWT token"
// @Param filter query string false "Filter column name"
// @Param value query string false "Filter value"
// @Success 200 users body entity.UserCred true  "List of users // @Tags Users credentials"
// @Failure 400 {string} Bad Request
// @Failure 401 {string} string "Unauthorized"
// @Failure 500 {string} string "Internal Server Error"
// @Router /v1/admin/cred [get].
func (ur *adminRoutes) GetUsersDetailCred(ctx *gin.Context) {
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

// GetUsersCredWithSort godoc
// @Summary Get a users credentials with sort
// @Description Retrieve a user credentials by sorting
// @Tags Users Credentials
// @Produce json
// @Param authorization header string true "JWT token"
// @Param sort query string true "Sort column"
// @Param method query string true "Ascending or Descending"
// @Success 200 users body entity.UserCred true "User credentials"
// @Failure 400 {string} string "Bad Request"
// @Failure 401 {string} string "Unauthorized"
// @Failure 500 {string} string "Internal Server Error"
// @Router /v1/admin/cred/sort [get].
func (ur *adminRoutes) GetUsersCredWithSort(ctx *gin.Context) {
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

// GetUsersCredWithSearch godoc
// @Summary Get a users credentials with search
// @Description Retrieve a user credentials by searching value
// @Tags Users Credentials
// @Produce json
// @Param authorization header string true "JWT token"
// @Param search query string true "Search column"
// @Param value query string true "Search value"
// @Success 200 users body entity.UserCred true "User credentials"
// @Failure 400 {string} string "Bad Request"
// @Failure 401 {string} string "Unauthorized"
// @Failure 500 {string} string "Internal Server Error"
// @Router /v1/admin/cred/search [get].
func (ur *adminRoutes) GetUsersCredWithSearch(ctx *gin.Context) {
	users, err := ur.u.UsersCredWithSearch(ctx, ctx.Request.URL.Query())
	if err != nil {
		ur.l.Error(fmt.Errorf("http - v1 - user - getUsers: %w", err))
		errorResponse(ctx, http.StatusInternalServerError, "http - v1 - user - getUsers error")

		return
	}

	ctx.JSON(http.StatusOK, users)
}

// GetUserDetailCredByID godoc
// @Summary Get a user credentials by ID
// @Description Retrieve a user credentials by their unique ID
// @Tags Users Credentials
// @Produce json
// @Param authorization header string true "JWT token"
// @Param id path int true "ID of the item"
// @Success 200 user body entity.UserCred true "User credentials"
// @Failure 400 {string} Bad Request
// @Failure 401 {string} string "Unauthorized"
// @Failure 500 {string} string "Internal Server Error"
// @Router /v1/admin/cred/{id} [get].
func (ur *adminRoutes) GetUserDetailCredByID(ctx *gin.Context) {
	id := ctx.Param("id")
	userByID, err := ur.u.GetUserCredByID(ctx, id)
	if err != nil {
		ur.l.Error(fmt.Errorf("http - v1 - userById - getUsersById: %w", err))
		errorResponse(ctx, http.StatusInternalServerError, "http - v1 - userById - getUsersById error")

		return
	}

	ctx.JSON(http.StatusOK, userByID)
}

// UpdateUserCred godoc
// @Summary Update user credentials
// @Description Update user credentials
// @Tags Users Credentials
// @Accept json
// @Param authorization header string true "JWT token"
// @Param id path int true "ID of the item"
// @Param data body dto.UserCredRequest true "JSON data"
// @Success 200 {string} string "Success"
// @Failure 400 {string} Bad Request
// @Failure 401 {string} string "Unauthorized"
// @Failure 500 {string} string "Internal Server Error"
// @Router /v1/admin/cred/{id} [post].
func (ur *adminRoutes) UpdateUserCred(ctx *gin.Context) {
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

// CreateUserCred godoc
// @Summary Create user credentials
// @Description Create user credentials
// @Tags Users Credentials
// @Accept json
// @Produce json
// @Param authorization header string true "JWT token"
// @Param data body dto.UserCredRequest true "JSON data"
// @Success 200 {int} int id
// @Failure 400 {string} Bad Request
// @Failure 401 {string} string "Unauthorized"
// @Failure 500 {string} string "Internal Server Error"
// @Router /v1/admin/cred [put].
func (ur *adminRoutes) CreateUserCred(ctx *gin.Context) {
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

// DeleteUserCred godoc
// @Summary Delete user credentials by user id
// @Description Delete user credentials by user id
// @Tags Users Credentials
// @Param authorization header string true "JWT token"
// @Param id path int true "ID of the item"
// @Success 200 {string} string "Success"
// @Failure 400 {string} Bad Request
// @Failure 401 {string} string "Unauthorized"
// @Failure 500 {string} string "Internal Server Error"
// @Router /v1/admin/cred/{id} [delete].
func (ur *adminRoutes) DeleteUserCred(ctx *gin.Context) {
	id := ctx.Param("id")

	err := ur.u.DeleteUserCred(ctx, id)
	if err != nil {
		ur.l.Error(fmt.Errorf("http - v1 - blockchain - set user detail: %w", err))
		errorResponse(ctx, http.StatusBadRequest, "http - v1 - user - set user detail info error")

		return
	}

	ctx.JSON(http.StatusOK, "Success")
}
