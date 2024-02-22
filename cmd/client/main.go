package main

import (
	"fmt"

	"github.com/oktavarium/gophkeeper/internal/client"
)

func main() {
	if err := client.Run(); err != nil {
		panic(fmt.Errorf("error running client: %w", err))
	}
}
