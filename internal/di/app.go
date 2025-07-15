package di

import (
	"context"
	"io"
	"log"

	"github.com/pannpers/go-backend-scaffold/internal/infrastructure/server"
)

type App struct {
	Server  *server.ConnectServer
	Closers []io.Closer
}

func (a *App) Shutdown(ctx context.Context) error {
	log.Println("Starting application shutdown...")

	// First, stop the server gracefully
	if err := a.Server.Stop(); err != nil {
		log.Printf("Error during server graceful shutdown: %v", err)
	}

	// Then close all other resources
	for _, closer := range a.Closers {
		if err := closer.Close(); err != nil {
			log.Printf("Error closing resource: %v", err)
		}
	}

	log.Println("Application shutdown complete")
	return nil
}
