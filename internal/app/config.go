package app

import (
	"flag"
	"os"
)

type Configure struct {
	Host     string
	Address  string
	FilePath string
	Dblink   string
}

func NewConfig() *Configure {
	return &Configure{}
}

func (c *Configure) SetConfig() {
	flag.StringVar(&c.Host, "a", "localhost:8080", "for listenandserve")
	flag.StringVar(&c.Address, "b", "http://localhost:8080", "for response")
	flag.StringVar(&c.FilePath, "f", "../../tmp/short-url-db.json", "for saving data")
	flag.StringVar(&c.Dblink, "d", "localhost", "for database link")
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

	File, isexist := os.LookupEnv("FILE_STORAGE_PATH")
	if isexist {
		c.FilePath = File
	}

	db, isexist := os.LookupEnv("DATABASE_DSN")
	if isexist {
		c.Dblink = db
	}
}
