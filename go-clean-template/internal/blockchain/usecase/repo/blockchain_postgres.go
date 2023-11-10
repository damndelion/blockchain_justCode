package repo

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/evrone/go-clean-template/internal/blockchain/transport"
	"github.com/evrone/go-clean-template/pkg/blockchain_logic"
)

type BlockchainRepo struct {
	*sql.DB
	chain             *blockchain_logic.Blockchain
	userGrpcTransport *transport.UserGrpcTransport
}

func NewBlockchainRepo(db *sql.DB, address string, userGrpcTransport *transport.UserGrpcTransport) *BlockchainRepo {
	chain := blockchain_logic.CreateBlockchain(db, address)
	return &BlockchainRepo{db, chain, userGrpcTransport}
}

func (br *BlockchainRepo) GetWallets(ctx context.Context) (wallets []string, err error) {
	res := blockchain_logic.ListAddresses()

	return res, nil
}

func (br *BlockchainRepo) GetWallet(ctx context.Context, userId string) (wallet string, err error) {
	address, err := br.userGrpcTransport.GetUserWallet(ctx, userId)
	if err != nil {
		return "", err
	}

	return address.Wallet, nil
}

func (br *BlockchainRepo) GetBalance(ctx context.Context, userId string) (balance float64, err error) {
	address, err := br.GetWallet(ctx, userId)
	if err != nil {
		return 0, err
	}
	res := br.chain.GetBalance(address)
	return res, nil
}
func (br *BlockchainRepo) GetBalanceByAddress(ctx context.Context, address string) (balance float64, err error) {
	res := br.chain.GetBalance(address)
	return res, nil
}

func (br *BlockchainRepo) GetBalanceUSD(ctx context.Context, userId string) (balance float64, err error) {
	address, err := br.GetWallet(ctx, userId)
	if err != nil {
		return -1, err
	}
	res, err := br.chain.GetBalanceInUSD(address)
	if err != nil {
		return -1, err
	}
	return res, nil
}

func (br *BlockchainRepo) CreateWallet(ctx context.Context, userID string) (string, error) {
	wallet, err := br.GetWallet(ctx, userID)
	if wallet != "" {
		return "", fmt.Errorf("User wallet already exists")
	}
	user, err := br.userGrpcTransport.GetUserByID(ctx, userID)
	if user.Valid == false {
		return "", fmt.Errorf("user is not valid")
	}

	address := blockchain_logic.CreateWallet()

	_, err = br.userGrpcTransport.SetUserWallet(ctx, userID, address)
	if err != nil {
		return "", err
	}
	return address, nil
}

func (br *BlockchainRepo) Send(ctx context.Context, from string, to string, amount float64) error {
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

func (br *BlockchainRepo) TopUp(ctx context.Context, from string, to string, amount float64) error {
	user, err := br.userGrpcTransport.GetUserByID(ctx, to)
	if err != nil {
		return err
	}
	fmt.Println(user.Wallet)
	err = br.chain.Send(from, user.Wallet, amount)
	if err != nil {
		return err
	}
	return nil
}
