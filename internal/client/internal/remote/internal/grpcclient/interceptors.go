package grpcclient

import (
	"context"
	"fmt"

	"google.golang.org/grpc"

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
		encryptedData, err := s.crypto.EncryptData(string(r.GetData()))
		if err != nil {
			return fmt.Errorf("error on encrypting data: %w", err)
		}
		newReq := &pbapi.SaveRequest{
			UserInfo: r.GetUserInfo(),
			Name:     r.GetName(),
			Data:     encryptedData,
		}

		req = newReq
	}
	err := invoker(ctx, method, req, reply, cc, opts...)
	return err
}
