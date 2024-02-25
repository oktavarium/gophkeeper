package client

import (
	"context"
	"fmt"

	"github.com/oktavarium/gophkeeper/internal/client/internal/cli"
	"github.com/oktavarium/gophkeeper/internal/client/internal/remote"
	"github.com/oktavarium/gophkeeper/internal/client/internal/storage"
)

func Run() error {
	ctx := context.Background()
	localStore := storage.NewStorage()
	remoteClient, err := remote.NewGrpcClient()
	if err != nil {
		return fmt.Errorf("error on creating new remote storage; %w", err)
	}
	if err := cli.Run(ctx, localStore, remoteClient); err != nil {
		return fmt.Errorf("error on running cli: %w", err)
	}
	return nil
}
