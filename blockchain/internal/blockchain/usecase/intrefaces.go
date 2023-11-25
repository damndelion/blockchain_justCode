// Package usecase implements application business logic. Each logic group in own file.
package usecase

import (
	"context"
	"sync"
)

//go:generate mockgen -source=interfaces.go -destination=./mocks_test.go -package=usecase_test

type (
	ChainUseCase interface {
		Wallets(ctx context.Context) ([]string, error)
		Wallet(ctx context.Context, userID string) (string, error)
		GetBalance(ctx context.Context, address string) (float64, error)
		GetBalanceUSD(ctx context.Context, address string) (float64, error)
		CreateWallet(ctx context.Context, userID string) (string, error)
		Send(ctx context.Context, from, to string, amount float64) error
		TopUp(ctx context.Context, to string, amount float64) error
		GetBalanceByAddress(ctx context.Context, address string) (float64, error)
	}

	ChainRepo interface {
		GetWallets(_ context.Context) ([]string, error)
		GetWallet(ctx context.Context, userID string) (wallet string, err error)
		GetBalance(ctx context.Context, userID string) (balance float64, err error)
		GetBalanceUSD(ctx context.Context, userID string) (balance float64, err error)
		CreateWallet(ctx context.Context, userID string) (string, error)
		Send(ctx context.Context, from, to string, amount float64, wg sync.WaitGroup) error
		TopUp(ctx context.Context, from, to string, amount float64, wg sync.WaitGroup) error
		GetBalanceByAddress(_ context.Context, address string) (balance float64, err error)
	}
)
