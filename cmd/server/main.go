package main

import (
	"fmt"

	"github.com/oktavarium/gophkeeper/internal/server"
)

func main() {
	if err := server.Run(); err != nil {
		panic(fmt.Errorf("error running server: %w", err))
	}
}
