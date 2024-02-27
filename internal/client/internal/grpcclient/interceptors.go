package grpcclient

import (
	"context"

	"google.golang.org/grpc"
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
	// switch r := req.(type) {
	// case *pbapi.SaveRequest:
	// 	encryptedData, err := s.crypto.EncryptData(string(r.GetData()))
	// 	if err != nil {
	// 		return fmt.Errorf("error on encrypting data: %w", err)
	// 	}
	// 	newReq := &pbapi.SaveRequest{
	// 		Name: r.GetName(),
	// 		Data: encryptedData,
	// 	}

	// 	req = newReq
	// }
	err := invoker(ctx, method, req, reply, cc, opts...)
	return err
}
