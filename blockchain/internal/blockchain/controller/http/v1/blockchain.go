package v1

import (
	"fmt"
	"github.com/evrone/go-clean-template/config/blockchain"
	"github.com/evrone/go-clean-template/internal/auth/controller/http/middleware"
	"github.com/evrone/go-clean-template/internal/blockchain/controller/http/v1/dto"
	"github.com/evrone/go-clean-template/internal/blockchain/usecase"
	"github.com/evrone/go-clean-template/pkg/blockchain_logic"
	"github.com/evrone/go-clean-template/pkg/logger"
	"github.com/gin-gonic/gin"
	"github.com/skip2/go-qrcode"
	"net/http"
)

type chainRoutes struct {
	c   usecase.ChainUseCase
	l   logger.Interface
	cfg *blockchain.Config
}

func newBlockchainRoutes(handler *gin.RouterGroup, c usecase.ChainUseCase, l logger.Interface, bc blockchain_logic.Blockchain, cfg *blockchain.Config) {
	r := &chainRoutes{c, l, cfg}

	blockchainHandler := handler.Group("/blockchain/wallet")
	{
		blockchainHandler.Use(middleware.JwtVerify(cfg.SecretKey))
		blockchainHandler.GET("/all", r.GetWallets)
		blockchainHandler.GET("/", r.GetWallet)
		blockchainHandler.GET("/balance", r.GetBalance)
		blockchainHandler.GET("/balance/address", r.GetBalanceByAddress)
		blockchainHandler.GET("/usd/balance", r.GetBalanceUSD)
		blockchainHandler.POST("/create", r.CreateWallet)
		blockchainHandler.POST("/send", r.Send)
		blockchainHandler.POST("/topup", r.TopUp)
		blockchainHandler.GET("/qr", r.GetWalletQRCode)

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
		errorResponse(ctx, http.StatusInternalServerError, fmt.Errorf("%w ", err))
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
// @Router /v1/blockchain/wallet/{userId} [get]
func (bc *chainRoutes) GetWallet(ctx *gin.Context) {
	authHeader := ctx.GetHeader("Authorization")
	userId, err := bc.c.GetIdFromToken(authHeader)
	wallet, err := bc.c.Wallet(ctx, userId)
	if err != nil {
		bc.l.Error(fmt.Errorf("http - v1 - blockchain - getWallet: %w", err))
		errorResponse(ctx, http.StatusInternalServerError, fmt.Errorf("%w ", err))
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
// @Router /v1/blockchain/balance/{userId} [get]
func (bc *chainRoutes) GetBalance(ctx *gin.Context) {
	authHeader := ctx.GetHeader("Authorization")
	userId, err := bc.c.GetIdFromToken(authHeader)

	balance, err := bc.c.GetBalance(ctx, userId)
	if err != nil {
		bc.l.Error(fmt.Errorf("http - v1 - blockchain - getBalance: %w", err))
		errorResponse(ctx, http.StatusInternalServerError, fmt.Errorf("%w ", err))
		return
	}

	ctx.JSON(http.StatusOK, balance)
}

func (bc *chainRoutes) GetBalanceByAddress(ctx *gin.Context) {
	var sendData dto.AddressRequest
	err := ctx.ShouldBindJSON(&sendData)
	if err != nil {
		bc.l.Error(fmt.Errorf("http - v1 - blockchain - send: %w", err))
		errorResponse(ctx, http.StatusBadRequest, fmt.Errorf("%w ", err))
		return
	}

	balance, err := bc.c.GetBalanceByAddress(ctx, sendData.Address)
	if err != nil {
		bc.l.Error(fmt.Errorf("http - v1 - blockchain - getBalance: %w", err))
		errorResponse(ctx, http.StatusInternalServerError, fmt.Errorf("%w ", err))
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
// @Router /v1/blockchain/usd/balance/{userId} [get]
func (bc *chainRoutes) GetBalanceUSD(ctx *gin.Context) {
	authHeader := ctx.GetHeader("Authorization")
	userId, err := bc.c.GetIdFromToken(authHeader)

	balance, err := bc.c.GetBalanceUSD(ctx, userId)
	if err != nil {
		bc.l.Error(fmt.Errorf("http - v1 - blockchain - getBalanceUSD: %w", err))
		errorResponse(ctx, http.StatusInternalServerError, fmt.Errorf("%w ", err))
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
// @Router /v1/blockchain/wallet/create [post]
func (bc *chainRoutes) CreateWallet(ctx *gin.Context) {
	accessToken := ctx.GetHeader("Authorization")

	id, err := bc.c.GetIdFromToken(accessToken)
	if err != nil {
		bc.l.Error(fmt.Errorf("http - v1 - blockchain - createWallet: %w", err))
		errorResponse(ctx, http.StatusInternalServerError, fmt.Errorf("%w ", err))
		return
	}
	wallet, err := bc.c.CreateWallet(ctx, id)
	if err != nil {
		bc.l.Error(fmt.Errorf("http - v1 - blockchain - createWallet: %w", err))
		errorResponse(ctx, http.StatusInternalServerError, fmt.Errorf("%w ", err))
		return
	}
	ctx.JSON(http.StatusOK, wallet)
}

// Send godoc
// @Summary Send cryptocurrency to another address
// @Description Send cryptocurrency from one address to another on the blockchain
// @Tags Blockchain
// @Accept json
// @Produce json
// @Param sendRequest body dto.SendRequest true "Send Request"
// @Param userId path string true "User ID"
// @Success 200 {string} string "Success"
// @Failure 400 {string} string "Invalid input"
// @Failure 401 {string} string "Unauthorized"
// @Failure 500 {string} string "Internal Server Error"
// @Router /v1/blockchain/wallet/send [post]
func (bc *chainRoutes) Send(ctx *gin.Context) {

	var sendData dto.SendRequest
	err := ctx.ShouldBindJSON(&sendData)

	if err != nil {
		bc.l.Error(fmt.Errorf("http - v1 - blockchain - send: %w", err))
		errorResponse(ctx, http.StatusBadRequest, fmt.Errorf("%w ", err))
		return
	}
	authHeader := ctx.GetHeader("Authorization")
	userId, err := bc.c.GetIdFromToken(authHeader)

	err = bc.c.Send(ctx, userId, sendData.To, sendData.Amount)
	if err != nil {
		bc.l.Error(fmt.Errorf("http - v1 - blockchain - send: %w", err))
		errorResponse(ctx, http.StatusInternalServerError, err)

		return
	}
	ctx.JSON(http.StatusOK, "Success ")
}

// TopUp godoc
// @Summary TopUp top up of an account
// @Description TopUp top up of an account
// @Tags Blockchain
// @Accept json
// @Produce json
// @Param topUpRequest body dto.TopUpRequest true "Top up Request"
// @Param userId path string true "User ID"
// @Success 200 {string} string "Success"
// @Failure 400 {string} string "Invalid input"
// @Failure 401 {string} string "Unauthorized"
// @Failure 500 {string} string "Internal Server Error"
// @Router /v1/blockchain/wallet/topup [post]
func (bc *chainRoutes) TopUp(ctx *gin.Context) {

	var topupData dto.TopUpRequest
	err := ctx.ShouldBindJSON(&topupData)
	authHeader := ctx.GetHeader("Authorization")

	id, err := bc.c.GetIdFromToken(authHeader)

	if err != nil {
		bc.l.Error(fmt.Errorf("http - v1 - blockchain - topup: %w", err))
		errorResponse(ctx, http.StatusBadRequest, fmt.Errorf("%w ", err))
		return
	}

	err = bc.c.TopUp(ctx, "", id, topupData.Amount)
	if err != nil {
		bc.l.Error(fmt.Errorf("http - v1 - blockchain - send: %w", err))
		errorResponse(ctx, http.StatusInternalServerError, fmt.Errorf("%w ", err))
		return
	}
	ctx.JSON(http.StatusOK, "Success ")
}

// GetWalletQRCode godoc
// @Summary Get a wallet QR code by user ID
// @Description Return a wallet QR code by user ID
// @Tags Blockchain
// @Accept json
// @Produce json
// @Param Authorization header string true "JWT Token"
// @Param userId path string true "User ID"
// @Success 200 {string} string "Wallet QR code"
// @Failure 400 {string} string "Invalid input"
// @Failure 401 {string} string "Unauthorized"
// @Failure 500 {string} string "Internal Server Error"
// @Router /v1/blockchain/wallet/qr/{userId} [get]
func (bc *chainRoutes) GetWalletQRCode(ctx *gin.Context) {
	authHeader := ctx.GetHeader("Authorization")
	userId, err := bc.c.GetIdFromToken(authHeader)

	wallet, err := bc.c.Wallet(ctx, userId)
	if err != nil {
		bc.l.Error(fmt.Errorf("http - v1 - blockchain - GetWalletQRCode: %w", err))
		errorResponse(ctx, http.StatusInternalServerError, fmt.Errorf("%w ", err))
		return
	}

	qrCode, err := qrcode.Encode(wallet, qrcode.Medium, 256)
	if err != nil {
		bc.l.Error(fmt.Errorf("Failed to generate QR code: %w", err))
		errorResponse(ctx, http.StatusInternalServerError, fmt.Errorf("Failed to generate QR code: %w ", err))
		return
	}

	ctx.Header("Content-Type", "image/png")

	ctx.Data(http.StatusOK, "image/png", qrCode)
}
