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
	server_address, isexist := os.LookupEnv("SERVER_ADDRESS")
	if isexist == false {
		flag.StringVar(&c.Host, "a", "localhost:8080", "for listenandserve")
	} else {
		c.Host = server_address
	}
	base_url, isexist := os.LookupEnv("BASE_URL")
	if isexist == false {
		flag.StringVar(&c.Address, "b", "http://localhost:8080", "for response")
	} else {
		c.Address = base_url
	}
	flag.Parse()
	c.Address += "/"
}
