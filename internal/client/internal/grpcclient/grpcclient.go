package grpcclient

import (
	"google.golang.org/grpc"

	pbapi "github.com/oktavarium/gophkeeper/api"
	"github.com/oktavarium/gophkeeper/internal/client/internal/crypto"
)

type GrpcClient struct {
	conn    *grpc.ClientConn
	client  pbapi.GophKeeperClient
	storage storage
	crypto  *crypto.Crypto
}

func NewGrpcClient(s storage) (*GrpcClient, error) {
	return &GrpcClient{
		storage: s,
	}, nil
}
