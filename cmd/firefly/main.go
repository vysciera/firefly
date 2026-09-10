package main

import (
	"log"
	"net/http"
	"os"

	"firefly/internal/flowerpress"
	"firefly/internal/web"
)

func main() {
	addr := ":3000"

	if value := os.Getenv("FIREFLY_ADDR"); value != "" {
		addr = value
	}

	flowerpressURL := "http://localhost:8080"

	if value := os.Getenv("FLOWERPRESS_URL"); value != "" {
		flowerpressURL = value
	}

	flowerpressClient := flowerpress.NewClient(flowerpressURL)

	server, err := web.NewServer(flowerpressClient)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("firefly listening on %s", addr)
	log.Printf("flowerpress: %s", flowerpressURL)

	if err := http.ListenAndServe(addr, server.Handler()); err != nil {
		log.Fatal(err)
	}
}
