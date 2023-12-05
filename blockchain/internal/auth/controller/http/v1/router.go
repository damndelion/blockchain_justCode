package v1

import (
	"net/http"

	"github.com/damndelion/blockchain_justCode/internal/auth/controller/http/middleware"
	"github.com/prometheus/client_golang/prometheus"

	_ "github.com/damndelion/blockchain_justCode/docs/auth"
	"github.com/damndelion/blockchain_justCode/internal/auth/usecase"
	"github.com/damndelion/blockchain_justCode/pkg/logger"
	"github.com/gin-contrib/pprof"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func NewAuthRouter(handler *gin.Engine, l logger.Interface, u usecase.AuthUseCase) {
	// Options
	handler.Use(gin.Logger())
	handler.Use(gin.Recovery())
	pprof.Register(handler)

	// Swagger

	url := ginSwagger.URL("http://localhost:8082/swagger/doc.json") // The url pointing to API definition
	handler.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler, url))

	// K8s probe

	handler.GET("/healthz", func(c *gin.Context) { c.Status(http.StatusOK) })

	// Prometheus metrics
	handler.GET("/metrics", gin.WrapH(promhttp.Handler()))
	handler.Use(middleware.MetricsHandler())
	prometheus.MustRegister()

	// Routers
	h := handler.Group("/v1")
	{
		newAuthRoutes(h, u, l)
	}
}
