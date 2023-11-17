package transport

import (
	"context"
	"errors"
	"fmt"

	"github.com/evrone/go-clean-template/config/blockchain"
	pb "github.com/evrone/go-clean-template/pkg/protobuf/userService/gw"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type UserGrpcTransport struct {
	config blockchain.UserGrpcTransport
	client pb.UserServiceClient
}

func NewUserGrpcTransport(config blockchain.UserGrpcTransport) *UserGrpcTransport {
	opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}

	conn, err := grpc.Dial(config.Host, opts...)
	if err != nil {
		return nil
	}
	client := pb.NewUserServiceClient(conn)

	return &UserGrpcTransport{
		client: client,
		config: config,
	}
}

func (t *UserGrpcTransport) GetUserByID(ctx context.Context, id string) (*pb.User, error) {
	resp, err := t.client.GetUserByID(ctx, &pb.GetUserByIDRequest{
		Id: id,
	})
	if err != nil {
		return nil, errors.New(fmt.Sprintf("cannot GetUserByID: %v", err))
	}

	if resp == nil {
		return nil, errors.New(fmt.Sprintf("not found"))
	}

	return resp, nil
}

func (t *UserGrpcTransport) GetUserWallet(ctx context.Context, id string) (*pb.UserWalletResponse, error) {
	resp, err := t.client.GetUserWallet(ctx, &pb.GetUserWalletRequest{
		Id: id,
	})
	if err != nil {
		return nil, errors.New(fmt.Sprintf("cannot GetUserWallet: %v", err))
	}

	if resp == nil {
		return nil, errors.New(fmt.Sprintf("not found"))
	}

	return resp, nil
}

func (t *UserGrpcTransport) GetUserByEmail(ctx context.Context, email string) (*pb.User, error) {
	resp, err := t.client.GetUserByEmail(ctx, &pb.GetUserByEmailRequest{
		Email: email,
	})
	if err != nil {
		return nil, errors.New(fmt.Sprintf("cannot GetUserByID: %v", err))
	}

	if resp == nil {
		return nil, errors.New(fmt.Sprintf("not found"))
	}

	return resp, nil
}

func (t *UserGrpcTransport) SetUserWallet(ctx context.Context, userID, address string) (*pb.SetUserWalletResponse, error) {
	resp, err := t.client.SetUserWallet(ctx, &pb.SetUserWalletRequest{
		UserId:  userID,
		Address: address,
	})
	if err != nil {
		return resp, errors.New(fmt.Sprintf("cannot SetUserWallet: %v", err))
	}

	return resp, nil
}
