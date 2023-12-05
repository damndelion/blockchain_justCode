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
