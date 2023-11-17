package grpc

import (
	"errors"
	"fmt"
	"net"

	pb "github.com/evrone/go-clean-template/pkg/protobuf/userService/gw"
	"google.golang.org/grpc"
)

type Server struct {
	port       string
	service    *Service
	grpcServer *grpc.Server
}

func NewServer(
	port string,
	service *Service,
) *Server {
	var opts []grpc.ServerOption

	grpcServer := grpc.NewServer(opts...)

	return &Server{
		port:       port,
		service:    service,
		grpcServer: grpcServer,
	}
}

func (s *Server) Start() error {
	listener, err := net.Listen("tcp", s.port)
	if err != nil {
		return errors.New(fmt.Sprintf("failed to listen grpc port: %s", s.port))
	}

	pb.RegisterUserServiceServer(s.grpcServer, s.service)

	go func() {
		err = s.grpcServer.Serve(listener)
		if err != nil {
			return
		}
	}()

	return nil
}

func (s *Server) Close() {
	s.grpcServer.GracefulStop()
}
