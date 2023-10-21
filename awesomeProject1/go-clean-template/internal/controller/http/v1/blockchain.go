package v1

import (
	"github.com/evrone/go-clean-template/internal/usecase"
	"github.com/evrone/go-clean-template/pkg/blockchain_logic"
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
		blockchainHandler.GET("wallet/all", r.GetWallets)
		blockchainHandler.GET("wallet/balance", r.GetBalance)
		blockchainHandler.POST("wallet/create", r.CreateWallet)
		blockchainHandler.POST("wallet/send", r.Send)

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

func (bc *chainRoutes) CreateWallet(ctx *gin.Context) {
	wallet, err := bc.c.CreateWallet(ctx)
	if err != nil {
		return
	}
	ctx.JSON(http.StatusOK, wallet)
}

func (bc *chainRoutes) Send(ctx *gin.Context) {
	type SendData struct {
		From   string  `json:"from"`
		To     string  `json:"to"`
		Amount float64 `json:"amount"`
	}
	var sendData SendData
	err := ctx.ShouldBindJSON(&sendData)

	if err != nil {
		ctx.JSON(http.StatusBadRequest, err)
		return
	}
	err = bc.c.Send(ctx, sendData.From, sendData.To, sendData.Amount)
	if err != nil {
		return
	}
	ctx.JSON(http.StatusOK, "Success ")
}
