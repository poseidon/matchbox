package main

import (
	"log"
	"net/http"

	"github.com/coreos/coreos-baremetal/server"
)

const address = ":8080"

func main() {
	bootConfigProvider := server.NewBootConfigProvider()
	bootConfigProvider.Add(server.DefaultAddr, server.CoreOSBootConfig)
	srv := server.NewServer(bootConfigProvider)
	h := srv.HTTPHandler()
	log.Printf("Starting coreos-baremetal metadata server")
	err := http.ListenAndServe(address, h)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}