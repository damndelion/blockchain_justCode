package repo

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/evrone/go-clean-template/pkg/blockchain_logic"
	"strconv"
)

type BlockchainRepo struct {
	*sql.DB
	chain *blockchain_logic.Blockchain
}

func NewBlockchainRepo(db *sql.DB, address string) *BlockchainRepo {
	chain := blockchain_logic.CreateBlockchain(db, address)
	return &BlockchainRepo{db, chain}
}

func (br *BlockchainRepo) GetWallets(ctx context.Context) (wallets []string, err error) {
	res := blockchain_logic.ListAddresses()

	return res, nil
}

func (br *BlockchainRepo) GetWallet(ctx context.Context, userId string) (wallet string, err error) {
	query := "SELECT wallet FROM users WHERE id = $1"

	err = br.DB.QueryRow(query, userId).Scan(&userId)
	if err != nil {
		return "", err
	}

	return userId, nil
}

func (br *BlockchainRepo) GetBalance(ctx context.Context, userId string) (balance float64, err error) {
	address, err := br.GetWallet(ctx, userId)
	if err != nil {
		return 0, err
	}
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
	address := blockchain_logic.CreateWallet()
	err := br.SetUserWallet(ctx, userID, address)
	if err != nil {
		return "", err
	}
	return address, nil
}

func (br *BlockchainRepo) Send(ctx context.Context, from string, to string, amount float64) error {
	//TODO grpc GetUserByID to get valid status
	address, err := br.GetWallet(ctx, from)
	if err != nil {
		return err
	}
	err = br.chain.Send(address, to, amount)
	if err != nil {
		return err
	}
	return nil
}

func (br *BlockchainRepo) TopUp(ctx context.Context, from string, to string, amount float64) error {
	//TODO grpc GetUserByID to get valid status
	address, err := br.GetWallet(ctx, to)
	if err != nil {
		return err
	}
	err = br.chain.Send(from, address, amount)
	if err != nil {
		return err
	}
	return nil
}

func (br *BlockchainRepo) SetUserWallet(ctx context.Context, userID string, address string) (err error) {
	id, _ := strconv.Atoi(userID)
	query := "SELECT wallet FROM users WHERE id = $1"
	res, err := br.DB.Exec(query, id)
	if res == nil {
		return fmt.Errorf("user already have wallet existing")
	}
	query = "UPDATE users SET wallet = $1 WHERE id = $2"
	_, err = br.DB.Exec(query, address, id)
	if err != nil {
		return err
	}
	return nil
}
