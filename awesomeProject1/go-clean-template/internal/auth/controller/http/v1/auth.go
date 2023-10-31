package v1

import (
	"fmt"
	"github.com/evrone/go-clean-template/internal/auth/controller/http/v1/dto"
	"github.com/evrone/go-clean-template/internal/auth/usecase"
	"github.com/evrone/go-clean-template/pkg/logger"
	"github.com/gin-gonic/gin"
	"net/http"
)

type authRoutes struct {
	u usecase.AuthUseCase
	l logger.Interface
}

// @title Gin Swagger Example API
// @version 2.0
// @description This is a sample server server.
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:8082
// @BasePath /v1
// @schemes http
func newAuthRoutes(handler *gin.RouterGroup, u usecase.AuthUseCase, l logger.Interface) {
	r := &authRoutes{u, l}

	userHandler := handler.Group("/user")
	{
		userHandler.POST("/register", r.Register)
		userHandler.POST("/login", r.Login)
	}

}

// Register
// @summary     Show history
// @description Show all translation history
// @ID          history
// @Tags  	    translation
// @Accept      json
// @Produce     json
// @Success     200 {object} historyResponse
// @Failure     500 {object} response
// @Router /v1/user/register [post]
func (ar *authRoutes) Register(ctx *gin.Context) {
	var registerRequest dto.RegisterRequest
	err := ctx.ShouldBindJSON(&registerRequest)
	if err != nil {
		ar.l.Error(fmt.Errorf("http - v1 - auth - register: %w", err))
		errorResponse(ctx, http.StatusBadRequest, "http - v1 - auth - registration dto error")
		return
	}

	err = ar.u.Register(ctx, registerRequest.Name, registerRequest.Email, registerRequest.Password)
	if err != nil {
		ar.l.Error(fmt.Errorf("http - v1 - auth - register: %w", err))
		errorResponse(ctx, http.StatusInternalServerError, "http - v1 - auth - registration error")
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
// @Param loginRequest body LoginRequest true "User login request"
// @Success 200 {string} string "Access token"
// @Failure 400 {string} string "Invalid input"
// @Failure 500 {string} string "Internal Server Error"
// @Router /auth/login [post]
func (ar *authRoutes) Login(ctx *gin.Context) {
	var loginRequest dto.LoginRequest
	err := ctx.ShouldBindJSON(&loginRequest)
	if err != nil {
		ar.l.Error(fmt.Errorf("http - v1 - auth - login: %w", err))
		errorResponse(ctx, http.StatusBadRequest, "http - v1 - auth - login dto error")
		return
	}

	token, err := ar.u.Login(ctx, loginRequest.Email, loginRequest.Password)
	if err != nil {
		ar.l.Error(fmt.Errorf("http - v1 - auth - login: %w", err))
		errorResponse(ctx, http.StatusInternalServerError, "http - v1 - auth - login error")
		return
	}

	ctx.JSON(http.StatusOK, token)
}
