package v1

import (
	"github.com/evrone/go-clean-template/internal/auth/usecase"
	"github.com/evrone/go-clean-template/internal/user/controller/http/v1/dto"
	"github.com/evrone/go-clean-template/pkg/logger"
	"github.com/gin-gonic/gin"
	"net/http"
)

type authRoutes struct {
	u usecase.AuthUseCase
	l logger.Interface
}

func newAuthRoutes(handler *gin.RouterGroup, u usecase.AuthUseCase, l logger.Interface) {
	r := &authRoutes{u, l}

	userHandler := handler.Group("/user")
	{
		userHandler.POST("/register", r.Register)
		userHandler.POST("/login", r.Login)
	}

}

func (ar *authRoutes) Register(ctx *gin.Context) {
	var registerRequest dto.RegisterRequest

	err := ctx.ShouldBindJSON(&registerRequest)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, err)
		return
	}

	err = ar.u.Register(ctx, registerRequest.Name, registerRequest.Email, registerRequest.Password)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, err)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "user successfully registered"})
}

func (ar *authRoutes) Register2(ctx *gin.Context) {

	ctx.JSON(http.StatusOK, gin.H{"message": "user successfully registered"})
}

func (ar *authRoutes) Login(ctx *gin.Context) {
	var loginRequest dto.LoginRequest
	err := ctx.ShouldBindJSON(&loginRequest)
	if err != nil {

		ctx.JSON(http.StatusBadRequest, err)
		return
	}

	token, err := ar.u.Login(ctx, loginRequest.Email, loginRequest.Password)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, err)
		return
	}

	ctx.JSON(http.StatusOK, token)
}
