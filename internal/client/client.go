package client

import (
	"context"
	"fmt"

	"github.com/oktavarium/gophkeeper/internal/client/grpcclient"
	"github.com/oktavarium/gophkeeper/internal/client/storage"
	"github.com/oktavarium/gophkeeper/internal/client/ui"
)

func Run() error {
	ctx := context.Background()
	localStore := storage.NewStorage()
	remoteClient, err := grpcclient.NewGrpcClient(localStore)
	if err != nil {
		return fmt.Errorf("error on creating new remote storage; %w", err)
	}
	err = ui.Run(ctx, localStore, remoteClient)
	if err != nil {
		return fmt.Errorf("error on running cli: %w", err)
	}
	return nil
}
