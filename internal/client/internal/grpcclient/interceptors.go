package grpcclient

import (
	"context"
	"fmt"

	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/timestamppb"

	pbapi "github.com/oktavarium/gophkeeper/api"
)

// cryptoUnaryInterceptor перехватчик для шифрования данных.
func (s *GrpcClient) cryptoUnaryInterceptor(ctx context.Context,
	method string,
	req,
	reply any,
	cc *grpc.ClientConn,
	invoker grpc.UnaryInvoker,
	opts ...grpc.CallOption,
) error {
	switch r := req.(type) {
	case *pbapi.SaveRequest:
		token, err := s.storage.GetToken()
		if err != nil {
			return fmt.Errorf("error on getting token: %w", err)
		}

		r.Token.Id = token.Id
		r.Token.ValidUntil = timestamppb.New(token.ValidUntil)

		req = r
	}

	err := invoker(ctx, method, req, reply, cc, opts...)
	return err
}
