package grpcclient

import (
	"context"
	"crypto/tls"
	"fmt"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"

	pbapi "github.com/oktavarium/gophkeeper/api"
)

func (c *GrpcClient) Init(ctx context.Context, addr string) error {
	if err := c.isInited(); err == nil {
		if err := c.conn.Close(); err != nil {
			return fmt.Errorf("error on closing current conn: %w", err)
		}
		c.conn = nil
	}

	creds := credentials.NewTLS(&tls.Config{InsecureSkipVerify: true})

	conn, err := grpc.DialContext(ctx,
		addr,
		grpc.WithTransportCredentials(creds),
		grpc.WithUnaryInterceptor(c.cryptoUnaryInterceptor))
	if err != nil {
		return fmt.Errorf("error on dialing: %s: %w", addr, err)
	}

	c.conn = conn
	c.client = pbapi.NewGophKeeperClient(conn)

	return nil
}

func (c *GrpcClient) isInited() error {
	if c.conn == nil {
		return fmt.Errorf("client not inited")
	}

	return nil
}

func (c *GrpcClient) SetConn(conn *grpc.ClientConn) {
	c.conn = conn
	c.client = pbapi.NewGophKeeperClient(conn)
}
