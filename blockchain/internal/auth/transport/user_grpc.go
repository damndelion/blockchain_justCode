package transport

import (
	"context"
	"fmt"

	"github.com/evrone/go-clean-template/config/auth"
	userEntity "github.com/evrone/go-clean-template/internal/user/entity"
	pb "github.com/evrone/go-clean-template/pkg/protobuf/userService/gw"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type UserGrpcTransport struct {
	config auth.UserGrpcTransport
	client pb.UserServiceClient
}

func NewUserGrpcTransport(config auth.UserGrpcTransport) *UserGrpcTransport {
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

func (t *UserGrpcTransport) CreateUser(ctx context.Context, user *userEntity.User) (*pb.CreateUserResponse, error) {
	grpcUser := &pb.CreateUserRequest{
		User: &pb.User{
			Id:       int32(user.ID),
			Name:     user.Name,
			Email:    user.Email,
			Password: user.Password,
			Wallet:   user.Wallet,
			Valid:    user.Valid,
			Role:     user.Role,
		},
	}
	resp, err := t.client.CreateUser(ctx, grpcUser)
	if err != nil {
		return nil, fmt.Errorf("cannot CreateUser: %w", err)
	}

	if resp == nil {
		return nil, fmt.Errorf("not found")
	}

	return resp, nil
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
