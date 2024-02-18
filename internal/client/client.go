package client

import (
	"fmt"

	"github.com/oktavarium/gophkeeper/internal/client/internal/cli"
)

func Run() error {
	if err := cli.Run(); err != nil {
		return fmt.Errorf("error on running cli: %w", err)
	}
	return nil
}
