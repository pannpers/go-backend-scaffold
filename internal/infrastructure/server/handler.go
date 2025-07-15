package server

import (
	"fmt"
	"net/http"
)

// HTTPHandler implements the Handler interface.
type HTTPHandler struct {
	// Add your dependencies here.
	// For example:
	// userService *service.UserService.
	// authService *service.AuthService.
}

// NewHTTPHandler creates a new HTTP handler.
func NewHTTPHandler() *HTTPHandler {
	return &HTTPHandler{
		// Initialize your dependencies here.
	}
}

// RegisterRoutes registers all HTTP routes.
func (h *HTTPHandler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/", h.handleHome)
	// Add more routes here.
}

// handleHome handles the home endpoint.
func (h *HTTPHandler) handleHome(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)

		return
	}

	fmt.Fprintf(w, "Welcome to Go Backend Scaffold!")
}
