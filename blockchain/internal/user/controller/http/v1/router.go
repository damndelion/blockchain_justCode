package v1

import (
	"net/http"

	"github.com/damndelion/blockchain_justCode/internal/user/controller/http/middleware"

	"github.com/prometheus/client_golang/prometheus"

	"github.com/damndelion/blockchain_justCode/config/user"
	_ "github.com/damndelion/blockchain_justCode/docs/user"
	"github.com/damndelion/blockchain_justCode/internal/user/usecase"
	"github.com/damndelion/blockchain_justCode/pkg/logger"
	_ "github.com/damndelion/blockchain_justCode/pkg/protobuf/userService/gw"
	"github.com/gin-contrib/pprof"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func NewUserRouter(handler *gin.Engine, l logger.Interface, u usecase.UserUseCase, cfg *user.Config) {
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
		newUserRoutes(h, u, l, cfg)
		newAdminRoutes(h, u, l, cfg)
	}
}
