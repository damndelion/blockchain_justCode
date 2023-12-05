package v1

import (
	"net/http"

	"github.com/damndelion/blockchain_justCode/internal/blockchain/controller/http/middleware"
	"github.com/prometheus/client_golang/prometheus"

	"github.com/damndelion/blockchain_justCode/pkg/cache"

	"github.com/damndelion/blockchain_justCode/config/blockchain"
	_ "github.com/damndelion/blockchain_justCode/docs/blockchain"
	"github.com/damndelion/blockchain_justCode/internal/blockchain/usecase"
	blockchainlogic "github.com/damndelion/blockchain_justCode/pkg/blockchain_logic"
	"github.com/damndelion/blockchain_justCode/pkg/logger"
	"github.com/gin-contrib/pprof"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func NewBlockchainRouter(handler *gin.Engine, l logger.Interface, c usecase.ChainUseCase, bc blockchainlogic.Blockchain, cfg *blockchain.Config, cache cache.Blockchain) {
	// Options
	handler.Use(gin.Logger())
	handler.Use(gin.Recovery())
	pprof.Register(handler)

	// Swagger
	// Swagger
	url := ginSwagger.URL("http://localhost:8081/swagger/doc.json") // The url pointing to API definition
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
		newBlockchainRoutes(h, c, l, bc, cfg, cache)
	}
}
