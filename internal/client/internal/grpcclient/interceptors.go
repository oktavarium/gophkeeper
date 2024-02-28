package grpcclient

import (
	"context"
	"fmt"

	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/timestamppb"

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
		r.Token.Id = tokenId
		r.Token.ValidUntil = timestamppb.New(tokenValidUntil)
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
