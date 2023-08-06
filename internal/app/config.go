package app

import "flag"

func SetHost() *string {
	Host := flag.String("host", "localhost:8080", "for listenandserve")

	return Host
}
func setAddress() *string {
	Address := flag.String("address", "http://localhost:8080/", "for response")
	flag.Parse()
	return Address
}
