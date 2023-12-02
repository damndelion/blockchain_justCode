package repo

import (
	"context"
	"errors"
	"fmt"

	blockchainlogic "github.com/evrone/go-clean-template/pkg/blockchain_logic"
)

type WalletBuilder struct {
	br      *BlockchainRepo
	ctx     context.Context
	userID  string
	address string
}

func NewWalletBuilder(br *BlockchainRepo, ctx context.Context, userID string) *WalletBuilder {
	return &WalletBuilder{
		br:     br,
		ctx:    ctx,
		userID: userID,
	}
}

func (wb *WalletBuilder) Build() error {
	err := wb.GetWallet()
	if err != nil {
		return err
	}
	err = wb.GetUser()
	if err != nil {
		return err
	}
	wb.CreateWallet()
	err = wb.SetWallet()
	if err != nil {
		return err
	}

	return nil
}

func (wb *WalletBuilder) GetWallet() error {
	wallet, _ := wb.br.GetWallet(wb.ctx, wb.userID)
	if wallet != "" {
		return errors.New(fmt.Sprintf("user wallet already exists"))
	}

	return nil
}

func (wb *WalletBuilder) GetUser() error {
	user, err := wb.br.userGrpcTransport.GetUserByID(wb.ctx, wb.userID)
	if !user.Valid {
		return errors.New(fmt.Sprintf("user is not valid"))
	}
	if err != nil {
		return err
	}

	return nil
}

func (wb *WalletBuilder) CreateWallet() {
	address := blockchainlogic.CreateWallet()
	wb.address = address
}

func (wb *WalletBuilder) SetWallet() error {
	_, err := wb.br.userGrpcTransport.SetUserWallet(wb.ctx, wb.userID, wb.address)
	if err != nil {
		return err
	}

	return nil
}
