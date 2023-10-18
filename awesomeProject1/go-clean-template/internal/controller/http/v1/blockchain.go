package v1

import (
	"github.com/evrone/go-clean-template/internal/blockchain/blockchain_logic"
	"github.com/evrone/go-clean-template/internal/usecase"
	"github.com/evrone/go-clean-template/pkg/logger"
	"github.com/gin-gonic/gin"
	"net/http"
)

type chainRoutes struct {
	c usecase.ChainUseCase
	l logger.Interface
}

func newBlockchainRoutes(handler *gin.RouterGroup, c usecase.ChainUseCase, l logger.Interface, bc blockchain_logic.Blockchain) {
	r := &chainRoutes{c, l}

	blockchainHandler := handler.Group("/blockchain")
	{
		blockchainHandler.GET("wallets/all", r.GetWallets)
		blockchainHandler.GET("wallet/balance", r.GetBalance)

	}

}

func (bc *chainRoutes) GetWallets(ctx *gin.Context) {

	wallets, err := bc.c.Wallets(ctx)
	if err != nil {
		bc.l.Error(err, "http - v1 - user - all")
		errorResponse(ctx, http.StatusInternalServerError, "database problems")

		return
	}

	ctx.JSON(http.StatusOK, wallets)
}

func (bc *chainRoutes) GetBalance(ctx *gin.Context) {
	address := ctx.Query("address")
	balance, err := bc.c.GetBalance(ctx, address)
	if err != nil {
		bc.l.Error(err, "http - v1 - user - all")
		errorResponse(ctx, http.StatusInternalServerError, "database problems")
		return
	}

	ctx.JSON(http.StatusOK, balance)
}
