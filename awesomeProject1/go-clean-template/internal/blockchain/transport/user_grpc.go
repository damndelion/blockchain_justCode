package transport

import (
	"context"
	"fmt"
	"github.com/evrone/go-clean-template/config/blockchain"
	pb "github.com/evrone/go-clean-template/pkg/protobuf/userService/gw"
	"google.golang.org/grpc"
)

type UserGrpcTransport struct {
	config blockchain.UserGrpcTransport
	client pb.UserServiceClient
}

func NewUserGrpcTransport(config blockchain.UserGrpcTransport) *UserGrpcTransport {
	opts := []grpc.DialOption{grpc.WithInsecure()}

	conn, _ := grpc.Dial(config.Host, opts...)

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
		return nil, fmt.Errorf("cannot GetUserByID: %w", err)
	}

	if resp == nil {
		return nil, fmt.Errorf("not found")
	}

	return resp, nil
}

func (t *UserGrpcTransport) GetUserWallet(ctx context.Context, id string) (*pb.UserWalletResponse, error) {
	resp, err := t.client.GetUserWallet(ctx, &pb.GetUserWalletRequest{
		Id: id,
	})

	if err != nil {
		return nil, fmt.Errorf("cannot GetUserWallet: %w", err)
	}

	if resp == nil {
		return nil, fmt.Errorf("not found")
	}

	return resp, nil
}

func (t *UserGrpcTransport) GetUserByEmail(ctx context.Context, email string) (*pb.User, error) {
	resp, err := t.client.GetUserByEmail(ctx, &pb.GetUserByEmailRequest{
		Email: email,
	})

	if err != nil {
		return nil, fmt.Errorf("cannot GetUserByID: %w", err)
	}

	if resp == nil {
		return nil, fmt.Errorf("not found")
	}

	return resp, nil
}

func (t *UserGrpcTransport) SetUserWallet(ctx context.Context, userID string, address string) (*pb.SetUserWalletResponse, error) {
	resp, err := t.client.SetUserWallet(ctx, &pb.SetUserWalletRequest{
		UserId:  userID,
		Address: address,
	})

	if err != nil {
		return resp, fmt.Errorf("cannot SetUserWallet: %w", err)
	}

	return resp, nil
}
