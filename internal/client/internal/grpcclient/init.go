package grpcclient

import (
	"context"
	"crypto/tls"
	"fmt"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"

	pbapi "github.com/oktavarium/gophkeeper/api"
)

func (s *GrpcClient) Init(ctx context.Context, addr string) error {
	if err := s.isInited(); err == nil {
		if err := s.conn.Close(); err != nil {
			return fmt.Errorf("error on closing current conn: %w", err)
		}

		s.conn = nil
	}

	creds := credentials.NewTLS(&tls.Config{InsecureSkipVerify: true})

	conn, err := grpc.DialContext(ctx,
		addr,
		grpc.WithTransportCredentials(creds),
		grpc.WithUnaryInterceptor(s.cryptoUnaryInterceptor))
	if err != nil {
		return fmt.Errorf("error on dialing: %s: %w", addr, err)
	}

	s.conn = conn
	s.client = pbapi.NewGophKeeperClient(conn)

	return nil
}

func (s *GrpcClient) isInited() error {
	if s.conn == nil {
		return fmt.Errorf("client not inited")
	}

	return nil
}
