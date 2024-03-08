package server

import (
	"flag"
)

type config struct {
	serverAddr string
	dbURI      string
	certPath   string
	keyPath    string
}

func loadFlags() config {
	c := config{}
	flag.StringVar(&c.serverAddr, "a", "localhost:8888", "server address")
	flag.StringVar(&c.dbURI, "d", "mongodb://root:example@localhost:27018/", "mongo connection string")
	flag.StringVar(&c.certPath, "c", "./rootCACert.pem", "server CA certificate")
	flag.StringVar(&c.keyPath, "k", "./rootCAKey.pem", "server CA key")
	flag.Parse()

	return c
}
