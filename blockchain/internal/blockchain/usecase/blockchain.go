package usecase

import (
	"context"
	"sync"

	"github.com/evrone/go-clean-template/config/blockchain"
	"github.com/evrone/go-clean-template/internal/blockchain/transport"
	"github.com/evrone/go-clean-template/internal/blockchain/usecase/repo"
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
	wg.Add(1)
	go func() {
		err = b.repo.Send(ctx, from, to, amount, &wg)
	}()
	wg.Wait()
	if err != nil {
		return err
	}

	return nil
}

func (b *Blockchain) TopUp(ctx context.Context, to string, amount float64) error {
	var wg sync.WaitGroup
	var err error
	wg.Add(1)
	go func() {
		err = b.repo.TopUp(ctx, b.cfg.GenesisAddress, to, amount, &wg)
	}()
	wg.Wait()
	if err != nil {
		return err
	}

	return nil
}
