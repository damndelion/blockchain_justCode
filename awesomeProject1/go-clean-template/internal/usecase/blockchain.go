package usecase

import (
	"context"
	"github.com/evrone/go-clean-template/internal/usecase/repo"
)

type Blockchain struct {
	repo repo.BlockchainRepo
}

func NewBlockchain(repo *repo.BlockchainRepo) *Blockchain {
	return &Blockchain{repo: *repo}
}

func (b *Blockchain) Wallets(ctx context.Context) ([]string, error) {
	return b.repo.GetWallets(ctx)
}

func (b *Blockchain) GetBalance(ctx context.Context, address string) (float64, error) {
	return b.repo.GetBalance(ctx, address)
}

func (b *Blockchain) CreateWallet(ctx context.Context) (string, error) {
	wallet, err := b.repo.CreateWallet(ctx)
	if err != nil {
		return "", err
	}
	return wallet, nil
}

func (b *Blockchain) Send(ctx context.Context, from string, to string, amount float64) error {
	b.repo.Send(ctx, from, to, amount)
	return nil
}
