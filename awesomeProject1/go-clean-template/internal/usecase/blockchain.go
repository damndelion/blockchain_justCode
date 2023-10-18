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
