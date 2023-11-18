package usecase

import (
	"context"
	"fmt"
	"strconv"
	"sync"

	"github.com/evrone/go-clean-template/config/blockchain"
	"github.com/evrone/go-clean-template/internal/blockchain/transport"
	"github.com/evrone/go-clean-template/internal/blockchain/usecase/repo"
	"github.com/golang-jwt/jwt"
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

func (b *Blockchain) Wallet(ctx context.Context, userID string) (string, error) {
	return b.repo.GetWallet(ctx, userID)
}

func (b *Blockchain) GetBalance(ctx context.Context, userID string) (float64, error) {
	return b.repo.GetBalance(ctx, userID)
}

func (b *Blockchain) GetBalanceByAddress(ctx context.Context, address string) (float64, error) {
	return b.repo.GetBalanceByAddress(ctx, address)
}

func (b *Blockchain) GetBalanceUSD(ctx context.Context, userID string) (float64, error) {
	return b.repo.GetBalanceUSD(ctx, userID)
}

func (b *Blockchain) CreateWallet(ctx context.Context, userID string) (string, error) {
	wallet, err := b.repo.CreateWallet(ctx, userID)
	if err != nil {
		return "", err
	}

	return wallet, nil
}

func (b *Blockchain) Send(ctx context.Context, from, to string, amount float64) error {
	var wg sync.WaitGroup
	var err error
	go func() {
		wg.Add(1)
		defer wg.Done()
		err = b.repo.Send(ctx, from, to, amount)
	}()
	wg.Wait()
	if err != nil {
		return err
	}

	return nil
}

func (b *Blockchain) TopUp(ctx context.Context, to string, amount float64) error {
	err := b.repo.TopUp(ctx, b.cfg.GenesisAddress, to, amount)
	if err != nil {
		return err
	}

	return nil
}

func (b *Blockchain) CheckForIDInAccessToken(urlUserID, accessToken string) bool {
	token, err := jwt.Parse(accessToken, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(b.cfg.SecretKey), nil
	})
	claims := token.Claims.(jwt.MapClaims)

	if err != nil || !token.Valid {
		return false
	}
	urlUserIDFormat, err := strconv.ParseFloat(urlUserID, 32)
	if err != nil {
		return false
	}
	if claims["user_id"] != urlUserIDFormat {
		return false
	}

	return true
}

func (b *Blockchain) GetIDFromToken(accessToken string) (string, error) {
	token, err := jwt.Parse(accessToken, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
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
