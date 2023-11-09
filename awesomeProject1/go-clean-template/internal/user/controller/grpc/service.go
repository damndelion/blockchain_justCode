package grpc

import (
	"context"
	"fmt"
	userEntity "github.com/evrone/go-clean-template/internal/user/entity"
	"github.com/evrone/go-clean-template/internal/user/usecase/repo"
	"github.com/evrone/go-clean-template/pkg/logger"
	pb "github.com/evrone/go-clean-template/pkg/protobuf/userService/gw"
)

type Service struct {
	pb.UnimplementedUserServiceServer
	logger *logger.Logger
	repo   *repo.UserRepo
}

func NewService(logger *logger.Logger, repo *repo.UserRepo) *Service {
	return &Service{
		logger: logger,
		repo:   repo,
	}
}

func (s *Service) GetUserByID(ctx context.Context, request *pb.GetUserByIDRequest) (*pb.User, error) {
	user, err := s.repo.GetUserByID(ctx, request.Id)
	if err != nil {
		s.logger.Error("failed to GetUserByID err: %v", err)
		return nil, fmt.Errorf("GetUserById err: %w", err)
	}

	return &pb.User{
		Id:       int32(user.Id),
		Name:     user.Name,
		Email:    user.Email,
		Password: user.Password,
		Wallet:   user.Wallet,
		Valid:    user.Valid,
	}, nil
}

func (s *Service) GetUserByEmail(ctx context.Context, request *pb.GetUserByEmailRequest) (*pb.User, error) {
	user, err := s.repo.GetUserByEmail(ctx, request.Email)
	if err != nil {
		s.logger.Error("failed to GetUserByEmail err: %v", err)
		return nil, fmt.Errorf("GetUserByEmail err: %w", err)
	}

	return &pb.User{
		Id:       int32(user.Id),
		Name:     user.Name,
		Email:    user.Email,
		Password: user.Password,
		Wallet:   user.Wallet,
		Valid:    user.Valid,
	}, nil
}

func (s *Service) GetUserWallet(ctx context.Context, request *pb.GetUserWalletRequest) (*pb.UserWalletResponse, error) {
	wallet, err := s.repo.GetUserWallet(ctx, request.Id)
	if err != nil {
		s.logger.Error("failed to GetUserByLogin err: %v", err)
		return nil, fmt.Errorf("GetUserByLogin err: %w", err)
	}

	return &pb.UserWalletResponse{
		Wallet: wallet,
	}, nil
}

func (s *Service) CreateUser(ctx context.Context, request *pb.CreateUserRequest) (*pb.CreateUserResponse, error) {
	user := request.GetUser()

	newUser := &userEntity.User{
		Name:     user.GetName(),
		Email:    user.GetEmail(),
		Password: user.GetPassword(),
		Wallet:   user.GetWallet(),
		Valid:    user.GetValid(),
	}

	id, err := s.repo.CreateUser(ctx, newUser)
	if err != nil {
		s.logger.Error("failed to CreateUser err: %v", err)
		return nil, fmt.Errorf("CreateUser err: %w", err)
	}

	return &pb.CreateUserResponse{Id: int32(id)}, nil
}

func (s *Service) SetUserWallet(ctx context.Context, request *pb.SetUserWalletRequest) (*pb.SetUserWalletResponse, error) {
	err := s.repo.SetUserWallet(ctx, request.UserId, request.Address)
	if err != nil {
		s.logger.Error("failed to SetUserWallet err: %v", err)
		grpcErr := &pb.SetUserWalletResponse{
			ErrorMessage: fmt.Sprintf("SetUserWallet err: %v", err),
		}
		return grpcErr, err
	}

	return &pb.SetUserWalletResponse{
		ErrorMessage: "",
	}, err
}
