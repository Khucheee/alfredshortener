package app

import (
	"flag"
)

type Configure struct {
	Host    string
	Address string
}

func (c *Configure) ParseFlags() {
	flag.StringVar(&c.Host, "a", "localhost:8080", "for listenandserve")
	flag.StringVar(&c.Address, "b", "http://localhost:8080/", "for response")
	flag.Parse()
}
