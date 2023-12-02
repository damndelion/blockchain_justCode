package usecase

import (
	"context"
	"sync"

	"github.com/opentracing/opentracing-go"

	"github.com/damndelion/blockchain_justCode/config/blockchain"
	"github.com/damndelion/blockchain_justCode/internal/blockchain/transport"
)

type Blockchain struct {
	repo              ChainRepo
	cfg               *blockchain.Config
	userGrpcTransport *transport.UserGrpcTransport
}

func NewBlockchain(repo ChainRepo, cfg *blockchain.Config, userGrpcTransport *transport.UserGrpcTransport) *Blockchain {
	return &Blockchain{repo, cfg, userGrpcTransport}
}

func (b *Blockchain) Wallets(ctx context.Context) ([]string, error) {
	span, spanCtx := opentracing.StartSpanFromContext(ctx, "wallets use case")
	defer span.Finish()

	return b.repo.GetWallets(spanCtx)
}

func (b *Blockchain) Wallet(ctx context.Context, userID string) (string, error) {
	span, spanCtx := opentracing.StartSpanFromContext(ctx, "wallet use case")
	defer span.Finish()

	return b.repo.GetWallet(spanCtx, userID)
}

func (b *Blockchain) GetBalance(ctx context.Context, userID string) (float64, error) {
	span, spanCtx := opentracing.StartSpanFromContext(ctx, "get balance use case")
	defer span.Finish()

	return b.repo.GetBalance(spanCtx, userID)
}

func (b *Blockchain) GetBalanceByAddress(ctx context.Context, address string) (float64, error) {
	span, spanCtx := opentracing.StartSpanFromContext(ctx, "get balance by address use case")
	defer span.Finish()

	return b.repo.GetBalanceByAddress(spanCtx, address)
}

func (b *Blockchain) GetBalanceUSD(ctx context.Context, userID string) (float64, error) {
	span, spanCtx := opentracing.StartSpanFromContext(ctx, "get balance is usd use case")
	defer span.Finish()

	return b.repo.GetBalanceUSD(spanCtx, userID)
}

func (b *Blockchain) CreateWallet(ctx context.Context, userID string) (string, error) {
	span, spanCtx := opentracing.StartSpanFromContext(ctx, "create wallet use case")
	defer span.Finish()
	wallet, err := b.repo.CreateWallet(spanCtx, userID)
	if err != nil {
		return "", err
	}

	return wallet, nil
}

func (b *Blockchain) Send(ctx context.Context, from, to string, amount float64) error {
	span, spanCtx := opentracing.StartSpanFromContext(ctx, "send use case")
	defer span.Finish()
	var wg sync.WaitGroup
	var err error
	wg.Add(1)
	go func() {
		err = b.repo.Send(spanCtx, from, to, amount, &wg)
	}()
	wg.Wait()
	if err != nil {
		return err
	}

	return nil
}

func (b *Blockchain) TopUp(ctx context.Context, to string, amount float64) error {
	span, spanCtx := opentracing.StartSpanFromContext(ctx, "top up use case")
	defer span.Finish()
	var wg sync.WaitGroup
	var err error
	wg.Add(1)
	go func() error {
		err = b.repo.TopUp(spanCtx, b.cfg.GenesisAddress, to, amount, &wg)
		if err != nil {
			return err
		}

		return nil
	}()
	wg.Wait()
	if err != nil {
		return err
	}

	return nil
}
