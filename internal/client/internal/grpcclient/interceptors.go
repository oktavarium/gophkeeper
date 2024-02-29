package grpcclient

import (
	"context"
	"fmt"
	"time"

	"google.golang.org/grpc"

	pbapi "github.com/oktavarium/gophkeeper/api"
)

// cryptoUnaryInterceptor перехватчик для шифрования данных.
func (c *GrpcClient) cryptoUnaryInterceptor(ctx context.Context,
	method string,
	req,
	reply any,
	cc *grpc.ClientConn,
	invoker grpc.UnaryInvoker,
	opts ...grpc.CallOption,
) error {
	tokenId, tokenValidUntil, err := c.storage.GetToken()

	switch r := req.(type) {
	case *pbapi.SyncRequest:
		if tokenValidUntil.Before(time.Now()) {
			return ErrTokenExpired
		}
		r.TokenID = tokenId
		req = r
	}

	methodErr := invoker(ctx, method, req, reply, cc, opts...)
	if methodErr != nil {
		return methodErr
	}

	var receivedToken *pbapi.Token
	tokenId, _, err = c.storage.GetToken()
	if err != nil {
		return fmt.Errorf("error on getting token: %w", err)
	}

	switch r := reply.(type) {
	case *pbapi.LoginResponse:
		receivedToken = r.GetToken()
	case *pbapi.RegisterResponse:
		receivedToken = r.GetToken()
	case *pbapi.SyncResponse:
		receivedToken = r.GetToken()
	}

	if tokenId != receivedToken.GetId() {
		if err := c.storage.UpdateToken(receivedToken.GetId(), receivedToken.GetValidUntil().AsTime()); err != nil {
			return fmt.Errorf("error on updating token: %w", err)
		}
	}

	return nil
}
