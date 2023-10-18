package repo

import (
	"context"
	"database/sql"
	"github.com/evrone/go-clean-template/internal/blockchain/blockchain_logic"
)

type BlockchainRepo struct {
	*sql.DB
}

func NewBlockchainRepo(db *sql.DB) *BlockchainRepo {
	return &BlockchainRepo{db}
}

func (ur *BlockchainRepo) GetWallets(ctx context.Context) (wallets []string, err error) {
	res := blockchain_logic.ListAddresses()

	return res, nil
}

func (ur *BlockchainRepo) GetBalance(ctx context.Context, address string) (balance float64, err error) {
	chain := blockchain_logic.CreateBlockchain(ur.DB, address)
	res := chain.GetBalance(address)
	return res, nil
}
