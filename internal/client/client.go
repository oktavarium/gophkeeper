package client

import (
	"context"
	"fmt"

	"github.com/oktavarium/gophkeeper/internal/client/internal/cli"
	"github.com/oktavarium/gophkeeper/internal/client/internal/storage"
)

func Run() error {
	ctx := context.Background()
	s := storage.NewRemoteStorage()
	if err := cli.Run(ctx, s); err != nil {
		return fmt.Errorf("error on running cli: %w", err)
	}
	return nil
}
