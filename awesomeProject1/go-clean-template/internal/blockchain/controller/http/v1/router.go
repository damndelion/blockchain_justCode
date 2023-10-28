package v1

import (
	"github.com/evrone/go-clean-template/config/blockchain"
	"github.com/evrone/go-clean-template/internal/blockchain/usecase"
	"github.com/evrone/go-clean-template/pkg/blockchain_logic"
	"github.com/evrone/go-clean-template/pkg/logger"
	"github.com/gin-gonic/gin"
)

func NewBlockchainRouter(handler *gin.Engine, l logger.Interface, c usecase.ChainUseCase, bc blockchain_logic.Blockchain, cfg *blockchain.Config) {
	// Options
	handler.Use(gin.Logger())
	handler.Use(gin.Recovery())

	// Routers
	h := handler.Group("/v1")
	{
		newBlockchainRoutes(h, c, l, bc, cfg)
	}
}
