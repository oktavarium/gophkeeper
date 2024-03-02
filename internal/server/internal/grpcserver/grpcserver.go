package grpcserver

import (
	"context"
	"fmt"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"

	pbapi "github.com/oktavarium/gophkeeper/api"
)

type GrpcServer struct {
	ctx context.Context
	pbapi.UnimplementedGophKeeperServer
	addr string
	*grpc.Server
	storage Storage
}

func NewGrpcServer(ctx context.Context, s Storage, addr string, certPath, keyPath string) (*GrpcServer, error) {
	creds, err := credentials.NewServerTLSFromFile(certPath, keyPath)
	if err != nil {
		return nil, fmt.Errorf("failed to create credentials: %v", err)
	}

	newServer := &GrpcServer{
		ctx:     ctx,
		addr:    addr,
		storage: s,
	}

	newServer.Server = grpc.NewServer(
		grpc.Creds(creds),
		grpc.UnaryInterceptor(newServer.cryptoUnaryInterceptor),
	)

	return newServer, nil
}

func (s *GrpcServer) ListenAndServe() error {
	listen, err := net.Listen("tcp", s.addr)
	if err != nil {
		return fmt.Errorf("error on listening socket for grpc: %w", err)
	}

	pbapi.RegisterGophKeeperServer(s, s)
	if err := s.Serve(listen); err != nil {
		return fmt.Errorf("error on serving grpc: %w", err)
	}

	return nil
}
