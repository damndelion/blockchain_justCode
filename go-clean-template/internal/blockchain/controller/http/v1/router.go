package v1

import (
	"github.com/evrone/go-clean-template/config/blockchain"
	_ "github.com/evrone/go-clean-template/docs/blockchain"
	"github.com/evrone/go-clean-template/internal/blockchain/usecase"
	"github.com/evrone/go-clean-template/pkg/blockchain_logic"
	"github.com/evrone/go-clean-template/pkg/logger"
	"github.com/gin-contrib/pprof"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"net/http"
)

func NewBlockchainRouter(handler *gin.Engine, l logger.Interface, c usecase.ChainUseCase, bc blockchain_logic.Blockchain, cfg *blockchain.Config) {
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

	// Routers
	h := handler.Group("/v1")
	{
		newBlockchainRoutes(h, c, l, bc, cfg)
	}
}
