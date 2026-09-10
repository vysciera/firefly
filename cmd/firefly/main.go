package main

import (
	"log"
	"net/http"
	"os"

	"firefly/internal/web"
)

func main() {
	addr := ":3000"

	if value := os.Getenv("FIREFLY_ADDR"); value != "" {
		addr = value
	}

	server, err := web.NewServer()
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("firefly listening on %s", addr)

	if err := http.ListenAndServe(addr, server.Handler()); err != nil {
		log.Fatal(err)
	}
}
