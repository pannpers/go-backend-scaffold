package server

import (
	"log"
	"net/http"

	"github.com/pannpers/protobuf-scaffold/gen/go/proto/api/v1/v1connect"
)

// ConnectServer represents the Connect server.
type ConnectServer struct {
	server *http.Server
	port   string
}

// NewConnectServer creates a new Connect server instance.
func NewConnectServer(userHandler v1connect.UserServiceHandler, postHandler v1connect.PostServiceHandler) *ConnectServer {
	mux := http.NewServeMux()

	// Register Connect handlers.
	path, handler := v1connect.NewUserServiceHandler(userHandler)
	mux.Handle(path, handler)

	path, handler = v1connect.NewPostServiceHandler(postHandler)
	mux.Handle(path, handler)

	server := &http.Server{
		Addr:    ":9090",
		Handler: mux,
	}

	return &ConnectServer{
		server: server,
		port:   ":9090",
	}
}

// Start starts the Connect server.
func (s *ConnectServer) Start() error {
	log.Printf("Connect Server starting on port %s", s.port)
	return s.server.ListenAndServe()
}

// Stop gracefully stops the Connect server.
func (s *ConnectServer) Stop() error {
	if s.server != nil {
		return s.server.Close()
	}
	return nil
}
