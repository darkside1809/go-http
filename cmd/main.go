package main

import (
	"log"
	"net"

	"github.com/darkside1809/http/pkg/server"
)


func main() {

	host := "0.0.0.0"
	port := "9999"
	srv := server.NewServer(net.JoinHostPort(host,port))

	err := srv.Start()
	if err != nil {
		log.Print(err)
		return
	}
}