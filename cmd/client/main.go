package main

import "github.com/oktavarium/gophkeeper/internal/client"
import "fmt"

func main() {
	if err := client.Run(); err != nil {
		panic(fmt.Errorf("error running client: %w", err))
	}
}
