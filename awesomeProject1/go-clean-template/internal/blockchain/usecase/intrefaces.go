// Package usecase implements application business logic. Each logic group in own file.
package usecase

import (
	"context"
)

//go:generate mockgen -source=interfaces.go -destination=./mocks_test.go -package=usecase_test

type (
	ChainUseCase interface {
		Wallets(ctx context.Context) ([]string, error)
		Wallet(ctx context.Context, userId string) (string, error)
		GetBalance(ctx context.Context, address string) (float64, error)
		GetBalanceUSD(ctx context.Context, address string) (float64, error)
		CreateWallet(ctx context.Context) (string, error)
		Send(ctx context.Context, from string, to string, amount float64) error
	}

	ChainRepo interface {
		GetWallets(ctx context.Context) ([]string, error)
		GetWallet(ctx context.Context, userId string) (string, error)
		GetBalance(ctx context.Context, address string) (float64, error)
		GetBalanceUSD(ctx context.Context, address string) (float64, error)
		CreateWallet(ctx context.Context) (string, error)
		Send(ctx context.Context, from string, to string, amount float64) error
	}
)
