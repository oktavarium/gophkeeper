package grpcserver

import (
	"context"
	"fmt"
	"net"

	"google.golang.org/grpc"

	pbapi "github.com/oktavarium/gophkeeper/api"
	"github.com/oktavarium/gophkeeper/internal/server/internal/storage"
)

type GrpcServer struct {
	ctx context.Context
	pbapi.UnimplementedGophKeeperServer
	addr string
	*grpc.Server
	storage storage.Storage
}

func NewGrpcServer(ctx context.Context, s storage.Storage, addr string) *GrpcServer {
	return &GrpcServer{
		ctx:     ctx,
		addr:    addr,
		Server:  grpc.NewServer(),
		storage: s,
	}
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
