package ui

import (
	"context"

	"github.com/oktavarium/gophkeeper/internal/client/internal/ui/internal/cli"
	"github.com/oktavarium/gophkeeper/internal/client/internal/ui/internal/cli/common"
)

func Run(ctx context.Context, s cli.Storage, c common.RemoteClient) error {
	return cli.Run(ctx, s, c)
}
