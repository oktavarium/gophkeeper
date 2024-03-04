package client

import (
	"context"
	"fmt"

	"github.com/oktavarium/gophkeeper/internal/client/internal/grpcclient"
	"github.com/oktavarium/gophkeeper/internal/client/internal/storage"
	"github.com/oktavarium/gophkeeper/internal/client/internal/ui"
)

func Run() error {
	ctx := context.Background()
	localStore := storage.NewStorage()
	remoteClient, err := grpcclient.NewGrpcClient(localStore)
	if err != nil {
		return fmt.Errorf("error on creating new remote storage; %w", err)
	}
	if err := ui.Run(ctx, localStore, remoteClient); err != nil {
		return fmt.Errorf("error on running cli: %w", err)
	}
	return nil
}
