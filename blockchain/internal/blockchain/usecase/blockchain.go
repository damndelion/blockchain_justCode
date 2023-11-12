package usecase

import (
	"context"
	"fmt"
	"github.com/evrone/go-clean-template/config/blockchain"
	"github.com/evrone/go-clean-template/internal/blockchain/transport"
	"github.com/evrone/go-clean-template/internal/blockchain/usecase/repo"
	"github.com/golang-jwt/jwt"
	"strconv"
)

type Blockchain struct {
	repo              repo.BlockchainRepo
	cfg               *blockchain.Config
	userGrpcTransport *transport.UserGrpcTransport
}

func NewBlockchain(repo *repo.BlockchainRepo, cfg *blockchain.Config, userGrpcTransport *transport.UserGrpcTransport) *Blockchain {
	return &Blockchain{*repo, cfg, userGrpcTransport}
}

func (b *Blockchain) Wallets(ctx context.Context) ([]string, error) {
	return b.repo.GetWallets(ctx)
}

func (b *Blockchain) Wallet(ctx context.Context, userId string) (string, error) {
	return b.repo.GetWallet(ctx, userId)
}

func (b *Blockchain) GetBalance(ctx context.Context, userId string) (float64, error) {
	return b.repo.GetBalance(ctx, userId)
}

func (b *Blockchain) GetBalanceByAddress(ctx context.Context, address string) (float64, error) {
	return b.repo.GetBalanceByAddress(ctx, address)
}
func (b *Blockchain) GetBalanceUSD(ctx context.Context, userId string) (float64, error) {
	return b.repo.GetBalanceUSD(ctx, userId)
}

func (b *Blockchain) CreateWallet(ctx context.Context, userID string) (string, error) {
	wallet, err := b.repo.CreateWallet(ctx, userID)
	if err != nil {
		return "", err
	}
	return wallet, nil
}

func (b *Blockchain) Send(ctx context.Context, from string, to string, amount float64) error {
	err := b.repo.Send(ctx, from, to, amount)
	if err != nil {
		return err
	}
	return nil
}

func (b *Blockchain) TopUp(ctx context.Context, from string, to string, amount float64) error {
	err := b.repo.TopUp(ctx, b.cfg.GenesisAddress, to, amount)
	if err != nil {
		return err
	}
	return nil
}

func (b *Blockchain) CheckForIdInAccessToken(urlUserID string, accessToken string) bool {
	token, err := jwt.Parse(accessToken, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(b.cfg.SecretKey), nil
	})
	claims := token.Claims.(jwt.MapClaims)

	if err != nil || !token.Valid {
		return false
	}
	urlUserIDFormat, _ := strconv.ParseFloat(urlUserID, 8)
	if claims["user_id"] != urlUserIDFormat {
		return false
	}
	return true
}

func (b *Blockchain) GetIdFromToken(accessToken string) (string, error) {
	token, err := jwt.Parse(accessToken, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(b.cfg.SecretKey), nil
	})
	claims := token.Claims.(jwt.MapClaims)

	if err != nil || !token.Valid {
		return "", err
	}
	id := fmt.Sprintf("%v", claims["user_id"])
	return id, nil
}
