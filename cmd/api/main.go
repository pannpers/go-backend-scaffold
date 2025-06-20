package main

import (
	"log"

	"github.com/pannpers/go-backend-scaffold/internal/di"
)

func main() {
	server, err := di.InitializeConnectServer()
	if err != nil {
		log.Fatalf("Failed to initialize API: %v", err)
	}

	if err := server.Start(); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
