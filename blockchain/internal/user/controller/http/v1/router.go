package v1

import (
	"net/http"

	"github.com/evrone/go-clean-template/internal/user/controller/http/middleware"

	"github.com/prometheus/client_golang/prometheus"

	"github.com/evrone/go-clean-template/config/user"
	_ "github.com/evrone/go-clean-template/docs/user"
	"github.com/evrone/go-clean-template/internal/user/usecase"
	"github.com/evrone/go-clean-template/pkg/cache"
	"github.com/evrone/go-clean-template/pkg/logger"
	_ "github.com/evrone/go-clean-template/pkg/protobuf/userService/gw"
	"github.com/gin-contrib/pprof"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func NewUserRouter(handler *gin.Engine, l logger.Interface, u usecase.UserUseCase, uc cache.User, cfg *user.Config) {
	// Options
	handler.Use(gin.Logger())
	handler.Use(gin.Recovery())
	pprof.Register(handler)

	// Swagger
	url := ginSwagger.URL("http://localhost:8080/swagger/doc.json") // The url pointing to API definition
	handler.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler, url))
	handler.Static("/grpc/swagger", "/Users/daniar/GolandProjects/blockchain/blockchain/pkg/protobuf/userService/gw")
	// K8s probe
	handler.GET("/healthz", func(c *gin.Context) { c.Status(http.StatusOK) })

	// Prometheus metrics
	handler.GET("/metrics", gin.WrapH(promhttp.Handler()))
	handler.Use(middleware.MetricsHandler())
	prometheus.MustRegister()
	// Routers
	h := handler.Group("/v1")
	{
		newUserRoutes(h, u, l, uc, cfg)
		newAdminRoutes(h, u, l, uc, cfg)
	}
}
