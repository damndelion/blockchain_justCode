package repo

import (
	"context"
	"database/sql"
	blockchain_logic "github.com/evrone/go-clean-template/pkg/blockchain_logic"
)

type BlockchainRepo struct {
	*sql.DB
	chain *blockchain_logic.Blockchain
}

func NewBlockchainRepo(db *sql.DB, address string) *BlockchainRepo {
	chain := blockchain_logic.CreateBlockchain(db, address)
	return &BlockchainRepo{db, chain}
}

func (ur *BlockchainRepo) GetWallets(ctx context.Context) (wallets []string, err error) {
	res := blockchain_logic.ListAddresses()

	return res, nil
}

func (ur *BlockchainRepo) GetBalance(ctx context.Context, address string) (balance float64, err error) {
	res := ur.chain.GetBalance(address)
	return res, nil
}

func (ur *BlockchainRepo) GetBalanceUSD(ctx context.Context, address string) (balance float64, err error) {
	res := ur.chain.GetBalanceInUSD(address)
	return res, nil
}

func (ur *BlockchainRepo) CreateWallet(ctx context.Context) (string, error) {
	address := blockchain_logic.CreateWallet()
	return address, nil
}

func (ur *BlockchainRepo) Send(ctx context.Context, from string, to string, amount float64) error {
	ur.chain.Send(from, to, amount)
	return nil
}
