package app

import (
	"flag"
	"os"
)

type Configure struct {
	Host    string
	Address string
}

func NewConfig() *Configure {
	return &Configure{Host: "", Address: ""}
}

func (c *Configure) SetConfig() {

	flag.StringVar(&c.Host, "a", "localhost:8080", "for listenandserve")
	flag.StringVar(&c.Address, "b", "http://localhost:8080", "for response")
	flag.Parse()

	serverAddress, isexist := os.LookupEnv("SERVER_ADDRESS")
	if isexist {
		c.Host = serverAddress
	}
	baseURL, isexist := os.LookupEnv("BASE_URL")
	if isexist {
		c.Address = baseURL
	}
	c.Address += "/"

}
