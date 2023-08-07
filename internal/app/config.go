package app

import (
	"flag"
	"os"
)

type Configure struct {
	Host    string
	Address string
}

func (c *Configure) SetConfig() {
	serverAddress, isexist := os.LookupEnv("SERVER_ADDRESS")
	if isexist {
		c.Host = serverAddress
	} else {
		flag.StringVar(&c.Host, "a", "localhost:8080", "for listenandserve")
	}
	baseUrl, isexist := os.LookupEnv("BASE_URL")
	if isexist {
		c.Address = baseUrl
	} else {
		flag.StringVar(&c.Address, "b", "http://localhost:8080", "for response")
	}
	flag.Parse()
	c.Address += "/"
}
