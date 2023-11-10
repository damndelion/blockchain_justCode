package repo

import (
	"context"
	"database/sql"
	"github.com/evrone/go-clean-template/pkg/blockchain_logic"
	"github.com/evrone/go-clean-template/pkg/postgres"
	"testing"
)

func TestBlockchainRepo_GetWallet(t *testing.T) {
	_, db, _ := postgres.New("postgres://postgres:postgres@localhost:5432/postgres")
	defer db.Close()
	address := blockchain_logic.CreateWallet()

	blockchain := blockchain_logic.CreateBlockchain(db, address)

	type fields struct {
		DB    *sql.DB
		chain *blockchain_logic.Blockchain
	}
	type args struct {
		ctx    context.Context
		userId string
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
				userId: "9",
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
				userId: "8",
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
				userId: "err",
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
			gotWallet, err := br.GetWallet(tt.args.ctx, tt.args.userId)
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
	_, db, _ := postgres.New("postgres://postgres:postgres@localhost:5432/postgres")
	defer db.Close()
	address := blockchain_logic.CreateWallet()

	blockchain := blockchain_logic.CreateBlockchain(db, address)

	type fields struct {
		DB    *sql.DB
		chain *blockchain_logic.Blockchain
	}
	type args struct {
		ctx    context.Context
		userId string
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
				userId: "9",
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
				userId: "8",
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
				userId: "err",
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
			gotBalance, err := br.GetBalance(tt.args.ctx, tt.args.userId)
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
