package grpcclient

import (
	"google.golang.org/grpc"

	pbapi "github.com/oktavarium/gophkeeper/api"
)

type GrpcClient struct {
	conn    *grpc.ClientConn
	client  pbapi.GophKeeperClient
	storage storage
}

func NewGrpcClient(s storage) (*GrpcClient, error) {
	return &GrpcClient{
		storage: s,
	}, nil
}
