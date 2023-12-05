package repo

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/damndelion/blockchain_justCode/internal/blockchain/controller/http/v1/dto"
	"github.com/opentracing/opentracing-go"

	"github.com/damndelion/blockchain_justCode/internal/blockchain/transport"
	blockchainlogic "github.com/damndelion/blockchain_justCode/pkg/blockchain_logic"
)

type BlockchainRepo struct {
	*sql.DB
	chain             *blockchainlogic.Blockchain
	userGrpcTransport *transport.UserGrpcTransport
}

func NewBlockchainRepo(db *sql.DB, address string, userGrpcTransport *transport.UserGrpcTransport) *BlockchainRepo {
	chain := blockchainlogic.CreateBlockchain(db, address)
	fetchBTCPriceAndStoreInChannel()

	return &BlockchainRepo{db, chain, userGrpcTransport}
}

func (br *BlockchainRepo) GetWallet(ctx context.Context, userID string) (wallet string, err error) {
	span, _ := opentracing.StartSpanFromContext(ctx, "get wallet repo")
	defer span.Finish()
	address, err := br.userGrpcTransport.GetUserWallet(ctx, userID)
	if err != nil {
		return "", err
	}
	if address.Wallet == "" {
		return "", fmt.Errorf("user does not have a wallet")
	}

	return address.Wallet, nil
}

func (br *BlockchainRepo) GetBalance(ctx context.Context, userID string) (balance float64, err error) {
	span, _ := opentracing.StartSpanFromContext(ctx, "get balance repo")
	defer span.Finish()
	address, err := br.GetWallet(ctx, userID)
	if err != nil {
		return 0, err
	}

	res := br.chain.GetBalance(address)

	return res, nil
}

func (br *BlockchainRepo) GetBalanceByAddress(ctx context.Context, address string) (balance float64, err error) {
	span, _ := opentracing.StartSpanFromContext(ctx, "get balance by address repo")
	defer span.Finish()
	res := br.chain.GetBalance(address)

	return res, nil
}

var btcPrice float64

func (br *BlockchainRepo) GetBalanceUSD(ctx context.Context, userID string) (balance float64, err error) {
	span, _ := opentracing.StartSpanFromContext(ctx, "get balance is usd repo")
	defer span.Finish()
	address, err := br.GetWallet(ctx, userID)

	bitcoinBalance := br.chain.GetBalance(address)
	totalBalanceUSD := bitcoinBalance * btcPrice

	return totalBalanceUSD, nil
}

func (br *BlockchainRepo) CreateWallet(ctx context.Context, userID string) (string, error) {
	span, _ := opentracing.StartSpanFromContext(ctx, "create wallet repo")
	defer span.Finish()
	wallet, err := br.GetWallet(ctx, userID)
	if wallet != "" {
		return "", errors.New(fmt.Sprintf("user wallet already exists"))
	}
	user, err := br.userGrpcTransport.GetUserByID(ctx, userID)
	if !user.Valid {
		return "", errors.New(fmt.Sprintf("user is not valid"))
	}

	address := blockchainlogic.CreateWallet()

	_, err = br.userGrpcTransport.SetUserWallet(ctx, userID, address)
	if err != nil {
		return "", err
	}

	return address, nil
}

func (br *BlockchainRepo) Send(ctx context.Context, from, to string, amount float64, wg *sync.WaitGroup) error {
	span, _ := opentracing.StartSpanFromContext(ctx, "send repo")
	defer span.Finish()
	defer wg.Done()
	user, err := br.userGrpcTransport.GetUserByID(ctx, from)
	if err != nil {
		return err
	}
	err = br.chain.Send(user.Wallet, to, amount)
	if err != nil {
		return err
	}

	return nil
}

func (br *BlockchainRepo) TopUp(ctx context.Context, from, to string, amount float64, wg *sync.WaitGroup) error {
	span, _ := opentracing.StartSpanFromContext(ctx, "top up repo")
	defer span.Finish()
	defer wg.Done()
	if amount < 0 {
		return fmt.Errorf("top up amount can not be negative")
	}
	user, err := br.userGrpcTransport.GetUserByID(ctx, to)
	if err != nil {
		return err
	}
	err = br.chain.Send(from, user.Wallet, amount)
	if err != nil {
		return err
	}

	return nil
}

func fetchBTCPriceAndStoreInChannel() {
	go func() {
		for {
			url := "https://api.coingecko.com/api/v3/simple/price?ids=bitcoin&vs_currencies=usd"

			response, err := http.Get(url)
			if err != nil {
				log.Println("Error fetching BTC price:", err)

				continue
			}

			var data dto.CoinGeckoResponse
			if err := json.NewDecoder(response.Body).Decode(&data); err != nil {
				log.Println("Error decoding BTC price data:", err)

				continue
			}

			btcPrice = data.Bitcoin.USD
			response.Body.Close()

			time.Sleep(10 * time.Second)
		}
	}()
}
