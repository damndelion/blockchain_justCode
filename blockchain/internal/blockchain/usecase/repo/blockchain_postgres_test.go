package repo

import (
	"context"
	"database/sql"
	"testing"

	blockchainlogic "github.com/evrone/go-clean-template/pkg/blockchain_logic"
	"github.com/evrone/go-clean-template/pkg/postgres"
)

func TestBlockchainRepo_GetWallet(t *testing.T) {
	_, db, err := postgres.New("postgres://postgres:postgres@localhost:5432/postgres")
	if err != nil {
		return
	}
	defer func(db *sql.DB) {
		err = db.Close()
		if err != nil {
			return
		}
	}(db)
	address := blockchainlogic.CreateWallet()

	blockchain := blockchainlogic.CreateBlockchain(db, address)

	type fields struct {
		DB    *sql.DB
		chain *blockchainlogic.Blockchain
	}
	type args struct {
		ctx    context.Context
		userID string
	}
	tests := []struct {
		name       string
		fields     fields
		args       args
		wantWallet string
		wantErr    bool
	}{
		{
			name: "Test case 1",
			fields: fields{
				DB:    db,
				chain: blockchain,
			},
			args: args{
				ctx:    context.TODO(),
				userID: "9",
			},
			wantWallet: "1FwwR4SLCKBxRT7nMcAPsYfb5LgYLLPpTy",
			wantErr:    false,
		},
		{
			name: "Test case 2",
			fields: fields{
				DB:    db,
				chain: blockchain,
			},
			args: args{
				ctx:    context.TODO(),
				userID: "8",
			},
			wantWallet: "13MAGqpqLKDCv96g9CHLFeL2rFLuCkTYHR",
			wantErr:    false,
		},
		{
			name: "Test case 3",
			fields: fields{
				DB:    db,
				chain: blockchain,
			},
			args: args{
				ctx:    context.TODO(),
				userID: "err",
			},
			wantWallet: "",
			wantErr:    true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			br := &BlockchainRepo{
				DB:    tt.fields.DB,
				chain: tt.fields.chain,
			}
			gotWallet, err := br.GetWallet(tt.args.ctx, tt.args.userID)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetWallet() error = %v, wantErr %v", err, tt.wantErr)

				return
			}
			if gotWallet != tt.wantWallet {
				t.Errorf("GetWallet() gotWallet = %v, want %v", gotWallet, tt.wantWallet)
			}
		})
	}
}

func TestBlockchainRepo_GetBalance(t *testing.T) {
	_, db, err := postgres.New("postgres://postgres:postgres@localhost:5432/postgres")
	defer func(db *sql.DB) {
		err = db.Close()
		if err != nil {
			return
		}
	}(db)
	address := blockchainlogic.CreateWallet()

	blockchain := blockchainlogic.CreateBlockchain(db, address)

	type fields struct {
		DB    *sql.DB
		chain *blockchainlogic.Blockchain
	}
	type args struct {
		ctx    context.Context
		userID string
	}
	tests := []struct {
		name        string
		fields      fields
		args        args
		wantBalance float64
		wantErr     bool
	}{
		{
			name: "Test case 1",
			fields: fields{
				DB:    db,
				chain: blockchain,
			},
			args: args{
				ctx:    context.TODO(),
				userID: "9",
			},
			wantBalance: 0.20000000000000284,
			wantErr:     false,
		},
		{
			name: "Test case 2",
			fields: fields{
				DB:    db,
				chain: blockchain,
			},
			args: args{
				ctx:    context.TODO(),
				userID: "8",
			},
			wantBalance: 999659.5999999996,
			wantErr:     false,
		},
		{
			name: "Test case 3",
			fields: fields{
				DB:    db,
				chain: blockchain,
			},
			args: args{
				ctx:    context.TODO(),
				userID: "err",
			},
			wantBalance: 0,
			wantErr:     true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			br := &BlockchainRepo{
				DB:    tt.fields.DB,
				chain: tt.fields.chain,
			}
			gotBalance, err := br.GetBalance(tt.args.ctx, tt.args.userID)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetBalance() error = %v, wantErr %v", err, tt.wantErr)

				return
			}
			if gotBalance != tt.wantBalance {
				t.Errorf("GetBalance() gotBalance = %v, want %v", gotBalance, tt.wantBalance)
			}
		})
	}
}
