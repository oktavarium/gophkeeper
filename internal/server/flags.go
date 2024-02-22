package server

import (
	"flag"
)

type config struct {
	serverAddr string
}

func loadFlags() config {
	c := config{}
	flag.StringVar(&c.serverAddr, "a", "localhost:8888", "server address")
	flag.Parse()

	return c
}
