package ui

import (
	"context"

	"github.com/oktavarium/gophkeeper/internal/client/internal/ui/internal/cli"
)

func Run(ctx context.Context, s cli.Storage, c cli.RemoteClient) error {
	return cli.Run(ctx, s, c)
}
