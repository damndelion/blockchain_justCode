// Package usecase implements application business logic. Each logic group in own file.
package usecase

import (
	"context"
	"sync"
)

//go:generate mockgen -source=interfaces.go -destination=./mocks_test.go -package=usecase_test

type (
	ChainUseCase interface {
		Wallet(ctx context.Context, userID string) (string, error)
		GetBalance(ctx context.Context, address string) (float64, error)
		GetBalanceUSD(ctx context.Context, address string) (float64, error)
		CreateWallet(ctx context.Context, userID string) (string, error)
		Send(ctx context.Context, from, to string, amount float64) error
		TopUp(ctx context.Context, to string, amount float64) error
		GetBalanceByAddress(ctx context.Context, address string) (float64, error)
	}

	ChainRepo interface {
		GetWallet(ctx context.Context, userID string) (string, error)
		GetBalance(ctx context.Context, userID string) (float64, error)
		GetBalanceUSD(ctx context.Context, userID string) (float64, error)
		CreateWallet(ctx context.Context, userID string) (string, error)
		Send(ctx context.Context, from, to string, amount float64, wg *sync.WaitGroup) error
		TopUp(ctx context.Context, from, to string, amount float64, wg *sync.WaitGroup) error
		GetBalanceByAddress(_ context.Context, address string) (float64, error)
	}
)
