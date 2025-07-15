package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/pannpers/go-backend-scaffold/internal/di"
)

func main() {
	// Create a context that will be canceled when OS signals are received
	ctx, stop := signal.NotifyContext(context.Background(),
		os.Interrupt,    // SIGINT (Ctrl+C)
		syscall.SIGTERM, // SIGTERM (k8s termination signal)
		syscall.SIGQUIT, // SIGQUIT
	)
	defer stop()

	log.Println("Starting server...")

	app, err := di.InitializeApp()
	if err != nil {
		log.Fatalf("Failed to initialize API: %v", err)
	}

	// Start server in a goroutine
	errChan := make(chan error, 1)

	go func() {
		if err := app.Server.Start(); err != nil {
			errChan <- err
		}
	}()

	// Wait for either context cancellation (signal) or server error
	select {
	case <-ctx.Done():
		log.Println("Received shutdown signal, stopping server gracefully...")
		app.Shutdown(context.Background())

	case err := <-errChan:
		log.Printf("Server failed to start: %v", err)
		app.Shutdown(context.Background())
		os.Exit(1)
	}

	log.Println("Server stopped")
}
