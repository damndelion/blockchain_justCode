package v1

import (
	"fmt"
	"github.com/evrone/go-clean-template/config/blockchain"
	"github.com/evrone/go-clean-template/internal/auth/controller/http/middleware"
	"github.com/evrone/go-clean-template/internal/blockchain/usecase"
	"github.com/evrone/go-clean-template/pkg/blockchain_logic"
	"github.com/evrone/go-clean-template/pkg/logger"
	"github.com/gin-gonic/gin"
	"net/http"
)

type chainRoutes struct {
	c usecase.ChainUseCase
	l logger.Interface
}

func newBlockchainRoutes(handler *gin.RouterGroup, c usecase.ChainUseCase, l logger.Interface, bc blockchain_logic.Blockchain, cfg *blockchain.Config) {
	r := &chainRoutes{c, l}

	blockchainHandler := handler.Group("/blockchain")
	{
		blockchainHandler.Use(middleware.JwtVerify(cfg.SecretKey))
		blockchainHandler.GET("wallet/all", r.GetWallets)
		blockchainHandler.GET("wallet/:userId", r.GetWallet)
		blockchainHandler.GET("wallet/balance", r.GetBalance)
		blockchainHandler.GET("wallet/usd/balance", r.GetBalanceUSD)
		blockchainHandler.POST("wallet/create", r.CreateWallet)
		blockchainHandler.POST("wallet/send", r.Send)

	}

}

// GetWallets godoc
// @Summary Get a list of wallets
// @Description Retrieve a list of wallets from the blockchain
// @Tags Blockchain
// @Accept json
// @Produce json
// @Param Authorization header string true "JWT Token"
// @Success 200 {array} []string "List of wallets"
// @Failure 401 {string} string "Unauthorized"
// @Failure 500 {string} string "Internal Server Error"
// @Router /v1/blockchain/wallet/all [get]
func (bc *chainRoutes) GetWallets(ctx *gin.Context) {
	wallets, err := bc.c.Wallets(ctx)
	if err != nil {
		bc.l.Error(fmt.Errorf("http - v1 - blockchain - getWallets: %w", err))
		errorResponse(ctx, http.StatusInternalServerError, "http - v1 - blockchain - getWallets error")
		return
	}

	ctx.JSON(http.StatusOK, wallets)
}

// GetWallet godoc
// @Summary Get a wallet by user ID
// @Description Retrieve a wallet from the blockchain for a specific user
// @Tags Blockchain
// @Accept json
// @Produce json
// @Param Authorization header string true "JWT Token"
// @Param userId path string true "User ID"
// @Success 200 {string} string "Wallet details"
// @Failure 400 {string} string "Invalid input"
// @Failure 401 {string} string "Unauthorized"
// @Failure 500 {string} string "Internal Server Error"
// @Router /v1/blockchain/wallets/{userId} [get]
func (bc *chainRoutes) GetWallet(ctx *gin.Context) {
	userId := ctx.Param("userId")
	wallet, err := bc.c.Wallet(ctx, userId)
	if err != nil {
		bc.l.Error(fmt.Errorf("http - v1 - blockchain - getWallet: %w", err))
		errorResponse(ctx, http.StatusInternalServerError, "http - v1 - blockchain - getWallet error")

		return
	}

	ctx.JSON(http.StatusOK, wallet)
}

// GetBalance godoc
// @Summary Get the balance of an address
// @Description Retrieve the balance of a specific address on the blockchain
// @Tags Blockchain
// @Accept json
// @Produce json
// @Param address query string true "Wallet Address" // Specify the wallet address as a query parameter
// @Param userId path string true "User ID"
// @Success 200 {number} float64 "Balance"
// @Failure 400 {string} string "Invalid input"
// @Failure 401 {string} string "Unauthorized"
// @Failure 500 {string} string "Internal Server Error"
// @Router /v1/blockchain/balance [get]
func (bc *chainRoutes) GetBalance(ctx *gin.Context) {
	address := ctx.Query("address")
	balance, err := bc.c.GetBalance(ctx, address)
	if err != nil {
		bc.l.Error(fmt.Errorf("http - v1 - blockchain - getBalance: %w", err))
		errorResponse(ctx, http.StatusInternalServerError, "http - v1 - blockchain - getBalance error")
		return
	}

	ctx.JSON(http.StatusOK, balance)
}

// GetBalanceUSD godoc
// @Summary Get the balance in USD of an address
// @Description Retrieve the balance in USD of a specific address on the blockchain
// @Tags Blockchain
// @Accept json
// @Produce json
// @Param address query string true "Wallet Address" // Specify the wallet address as a query parameter
// @Param userId path string true "User ID"
// @Success 200 {number} float64 "Balance in USD"
// @Failure 400 {string} string "Invalid input"
// @Failure 401 {string} string "Unauthorized"
// @Failure 500 {string} string "Internal Server Error"
// @Router /v1/blockchain/balanceUSD [get]
func (bc *chainRoutes) GetBalanceUSD(ctx *gin.Context) {
	address := ctx.Query("address")
	balance, err := bc.c.GetBalanceUSD(ctx, address)
	if err != nil {
		bc.l.Error(fmt.Errorf("http - v1 - blockchain - getBalanceUSD: %w", err))
		errorResponse(ctx, http.StatusInternalServerError, "http - v1 - blockchain - getBalance error")
		return
	}

	ctx.JSON(http.StatusOK, balance)
}

// CreateWallet godoc
// @Summary Create a new wallet
// @Description Create a new wallet on the blockchain
// @Tags Blockchain
// @Accept json
// @Produce json
// @Param userId path string true "User ID"
// @Success 200 {string} string "New Wallet"
// @Failure 400 {string} string "Invalid input"
// @Failure 401 {string} string "Unauthorized"
// @Failure 500 {string} string "Internal Server Error"
// @Router /v1/blockchain/createWallet [post]
func (bc *chainRoutes) CreateWallet(ctx *gin.Context) {
	wallet, err := bc.c.CreateWallet(ctx)
	if err != nil {
		bc.l.Error(fmt.Errorf("http - v1 - blockchain - createWallet: %w", err))
		errorResponse(ctx, http.StatusInternalServerError, "http - v1 - blockchain - createWallet error")
		return
	}
	ctx.JSON(http.StatusOK, wallet)
}

type SendData struct {
	From   string  `json:"from"`
	To     string  `json:"to"`
	Amount float64 `json:"amount"`
}

// Send godoc
// @Summary Send cryptocurrency to another address
// @Description Send cryptocurrency from one address to another on the blockchain
// @Tags Blockchain
// @Accept json
// @Produce json
// @Param sendRequest body SendData true "Send Request"
// @Param userId path string true "User ID"
// @Success 200 {string} string "Success"
// @Failure 400 {string} string "Invalid input"
// @Failure 401 {string} string "Unauthorized"
// @Failure 500 {string} string "Internal Server Error"
// @Router /v1/blockchain/send [post]
func (bc *chainRoutes) Send(ctx *gin.Context) {
	//TODO add dto

	var sendData SendData
	err := ctx.ShouldBindJSON(&sendData)

	if err != nil {
		bc.l.Error(fmt.Errorf("http - v1 - blockchain - send: %w", err))
		errorResponse(ctx, http.StatusBadRequest, "http - v1 - blockchain - send request error")
		return
	}
	err = bc.c.Send(ctx, sendData.From, sendData.To, sendData.Amount)
	if err != nil {
		bc.l.Error(fmt.Errorf("http - v1 - blockchain - send: %w", err))
		errorResponse(ctx, http.StatusInternalServerError, "http - v1 - blockchain - send error")

		return
	}
	ctx.JSON(http.StatusOK, "Success ")
}
