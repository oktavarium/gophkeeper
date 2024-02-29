package grpcserver

import (
	"context"
	"fmt"
	"net"

	"google.golang.org/grpc"

	pbapi "github.com/oktavarium/gophkeeper/api"
)

type GrpcServer struct {
	ctx context.Context
	pbapi.UnimplementedGophKeeperServer
	addr string
	*grpc.Server
	storage Storage
}

func NewGrpcServer(ctx context.Context, s Storage, addr string) *GrpcServer {
	return &GrpcServer{
		ctx:     ctx,
		addr:    addr,
		storage: s,
	}
}

func (s *GrpcServer) ListenAndServe() error {
	s.Server = grpc.NewServer(grpc.UnaryInterceptor(s.cryptoUnaryInterceptor))
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
