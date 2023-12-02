package usecase

import (
	"context"
	"errors"
	"sync"
	"testing"

	"github.com/damndelion/blockchain_justCode/config/blockchain"
	"github.com/damndelion/blockchain_justCode/internal/blockchain/mocks"
	"github.com/damndelion/blockchain_justCode/internal/blockchain/transport"
)

//nolint:dupl  // not duplicate
func TestBlockchain_GetBalance(t *testing.T) {
	type fields struct {
		repoMock          *mocks.ChainRepo
		cfg               *blockchain.Config
		userGrpcTransport *transport.UserGrpcTransport
	}
	type args struct {
		ctx    context.Context
		userID string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    float64
		wantErr bool
	}{
		{
			name: "Success",
			fields: fields{
				repoMock:          &mocks.ChainRepo{},
				cfg:               &blockchain.Config{},
				userGrpcTransport: &transport.UserGrpcTransport{},
			},
			args: args{
				ctx:    context.Background(),
				userID: "14",
			},
			want:    10.00,
			wantErr: false,
		},
		{
			name: "Error",
			fields: fields{
				repoMock:          &mocks.ChainRepo{},
				cfg:               &blockchain.Config{},
				userGrpcTransport: &transport.UserGrpcTransport{},
			},
			args: args{
				ctx:    context.Background(),
				userID: "err",
			},
			want:    0.00,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if !tt.wantErr {
				tt.fields.repoMock.On("GetBalance", tt.args.ctx, tt.args.userID).Return(tt.want, nil)
			} else {
				tt.fields.repoMock.On("GetBalance", tt.args.ctx, tt.args.userID).Return(tt.want, errors.New("get wallet by address error"))
			}
			b := &Blockchain{
				repo:              tt.fields.repoMock,
				cfg:               tt.fields.cfg,
				userGrpcTransport: tt.fields.userGrpcTransport,
			}
			got, err := b.GetBalance(tt.args.ctx, tt.args.userID)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetBalance() error = %v, wantErr %v", err, tt.wantErr)

				return
			}
			if got != tt.want {
				t.Errorf("GetBalance() got = %v, want %v", got, tt.want)
			}
		})
	}
}

//nolint:dupl  // not duplicate
func TestBlockchain_GetBalanceByAddress(t *testing.T) {
	type fields struct {
		repoMock          *mocks.ChainRepo
		cfg               *blockchain.Config
		userGrpcTransport *transport.UserGrpcTransport
	}
	type args struct {
		ctx     context.Context
		address string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    float64
		wantErr bool
	}{
		{
			name: "Success",
			fields: fields{
				repoMock:          &mocks.ChainRepo{},
				cfg:               &blockchain.Config{},
				userGrpcTransport: &transport.UserGrpcTransport{},
			},
			args: args{
				ctx:     context.Background(),
				address: "1EzU4cx9yfdBC3X38MV3xYiNf3XVBmDBpW",
			},
			want:    10.00,
			wantErr: false,
		},
		{
			name: "Error",
			fields: fields{
				repoMock:          &mocks.ChainRepo{},
				cfg:               &blockchain.Config{},
				userGrpcTransport: &transport.UserGrpcTransport{},
			},
			args: args{
				ctx:     context.Background(),
				address: "err",
			},
			want:    0.00,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if !tt.wantErr {
				tt.fields.repoMock.On("GetBalanceByAddress", tt.args.ctx, tt.args.address).Return(tt.want, nil)
			} else {
				tt.fields.repoMock.On("GetBalanceByAddress", tt.args.ctx, tt.args.address).Return(tt.want, errors.New("get wallet by address error"))
			}
			b := &Blockchain{
				repo:              tt.fields.repoMock,
				cfg:               tt.fields.cfg,
				userGrpcTransport: tt.fields.userGrpcTransport,
			}
			got, err := b.GetBalanceByAddress(tt.args.ctx, tt.args.address)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetBalanceByAddress() error = %v, wantErr %v", err, tt.wantErr)

				return
			}
			if got != tt.want {
				t.Errorf("GetBalanceByAddress() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBlockchain_Send(t *testing.T) {
	type fields struct {
		repoMock          *mocks.ChainRepo
		cfg               *blockchain.Config
		userGrpcTransport *transport.UserGrpcTransport
	}
	type args struct {
		ctx    context.Context
		from   string
		to     string
		amount float64
		wg     *sync.WaitGroup
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "Success",
			fields: fields{
				repoMock:          &mocks.ChainRepo{},
				cfg:               &blockchain.Config{},
				userGrpcTransport: &transport.UserGrpcTransport{},
			},
			args: args{
				ctx:    context.Background(),
				from:   "1EzU4cx9yfdBC3X38MV3xYiNf3XVBmDBpW",
				to:     "1HM5Mom2VKzchdJToC1R4ji6K7XKt1Xf5B",
				amount: 1.00,
				wg:     &sync.WaitGroup{},
			},
			wantErr: false,
		},

		{
			name: "Not enough funds error",
			fields: fields{
				repoMock:          &mocks.ChainRepo{},
				cfg:               &blockchain.Config{},
				userGrpcTransport: &transport.UserGrpcTransport{},
			},
			args: args{
				ctx:    context.Background(),
				from:   "1EzU4cx9yfdBC3X38MV3xYiNf3XVBmDBpW",
				to:     "1HM5Mom2VKzchdJToC1R4ji6K7XKt1Xf5B",
				amount: 100.00,
				wg:     &sync.WaitGroup{},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if !tt.wantErr {
				tt.fields.repoMock.On("Send", tt.args.ctx, tt.args.from, tt.args.to, tt.args.amount, tt.args.wg).Return(nil)
			} else {
				tt.fields.repoMock.On("Send", tt.args.ctx, tt.args.from, tt.args.to, tt.args.amount, tt.args.wg).Return(errors.New("send error"))
			}
			b := &Blockchain{
				repo:              tt.fields.repoMock,
				cfg:               tt.fields.cfg,
				userGrpcTransport: tt.fields.userGrpcTransport,
			}
			if err := b.repo.Send(tt.args.ctx, tt.args.from, tt.args.to, tt.args.amount, tt.args.wg); (err != nil) != tt.wantErr {
				t.Errorf("Send() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestBlockchain_TopUp(t *testing.T) {
	type fields struct {
		repoMock          *mocks.ChainRepo
		cfg               *blockchain.Config
		userGrpcTransport *transport.UserGrpcTransport
	}
	type args struct {
		ctx    context.Context
		from   string
		to     string
		amount float64
		wg     *sync.WaitGroup
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "Success",
			fields: fields{
				repoMock:          &mocks.ChainRepo{},
				cfg:               &blockchain.Config{},
				userGrpcTransport: &transport.UserGrpcTransport{},
			},
			args: args{
				ctx:    context.Background(),
				from:   "1Pq4qTbgTH4KhmFiPQ91YXVyyK5oo6aX1G",
				to:     "1EzU4cx9yfdBC3X38MV3xYiNf3XVBmDBpW",
				amount: 1.00,
				wg:     &sync.WaitGroup{},
			},
			wantErr: false,
		},

		{
			name: "Error",
			fields: fields{
				repoMock:          &mocks.ChainRepo{},
				cfg:               &blockchain.Config{},
				userGrpcTransport: &transport.UserGrpcTransport{},
			},
			args: args{
				ctx:    context.Background(),
				from:   "1Pq4qTbgTH4KhmFiPQ91YXVyyK5oo6aX1G",
				to:     "1EzU4cx9yfdBC3X38MV3xYiNf3XVBmDBpW",
				amount: -1.00,
				wg:     &sync.WaitGroup{},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var err error
			if !tt.wantErr {
				tt.fields.repoMock.On("TopUp", tt.args.ctx, tt.args.from, tt.args.to, tt.args.amount, tt.args.wg).Return(err)
			} else {
				tt.fields.repoMock.On("TopUp", tt.args.ctx, tt.args.from, tt.args.to, tt.args.amount, tt.args.wg).Return(errors.New("error"))
			}
			b := &Blockchain{
				repo:              tt.fields.repoMock,
				cfg:               tt.fields.cfg,
				userGrpcTransport: tt.fields.userGrpcTransport,
			}
			if err = b.repo.TopUp(tt.args.ctx, tt.args.from, tt.args.to, tt.args.amount, tt.args.wg); (err != nil) != tt.wantErr {
				t.Errorf("TopUp() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestBlockchain_Wallet(t *testing.T) {
	type fields struct {
		repoMock          *mocks.ChainRepo
		cfg               *blockchain.Config
		userGrpcTransport *transport.UserGrpcTransport
	}
	type args struct {
		ctx    context.Context
		userID string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "Success",
			fields: fields{
				repoMock:          &mocks.ChainRepo{},
				cfg:               &blockchain.Config{},
				userGrpcTransport: &transport.UserGrpcTransport{},
			},
			args: args{
				ctx:    context.Background(),
				userID: "1",
			},
			want:    "1HM5Mom2VKzchdJToC1R4ji6K7XKt1Xf5B",
			wantErr: false,
		},
		{
			name: "GetWalletError",
			fields: fields{
				repoMock:          &mocks.ChainRepo{},
				cfg:               &blockchain.Config{},
				userGrpcTransport: &transport.UserGrpcTransport{},
			},
			args: args{
				ctx:    context.Background(),
				userID: "invalidUser",
			},
			want:    "",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.wantErr {
				tt.fields.repoMock.On("GetWallet", tt.args.ctx, tt.args.userID).Return("", errors.New("get wallet error"))
			} else {
				tt.fields.repoMock.On("GetWallet", tt.args.ctx, tt.args.userID).Return(tt.want, nil)
			}

			b := &Blockchain{
				repo:              tt.fields.repoMock,
				cfg:               tt.fields.cfg,
				userGrpcTransport: tt.fields.userGrpcTransport,
			}

			got, err := b.Wallet(tt.args.ctx, tt.args.userID)
			if (err != nil) != tt.wantErr {
				t.Errorf("Wallet() error = %v, wantErr %v", err, tt.wantErr)

				return
			}
			if got != tt.want {
				t.Errorf("Wallet() got = %v, want %v", got, tt.want)
			}
		})
	}
}
