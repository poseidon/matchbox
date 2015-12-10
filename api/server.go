package api

import (
	"net/http"
)

// Server serves iPXE/Pixiecore boot configs and hosts images.
type Server struct {
	bootConfigs BootConfigProvider
}

// NewServer returns a new Server which uses the given BootConfigProvider.
func NewServer(bootConfigs BootConfigProvider) *Server {
	return &Server{
		bootConfigs: bootConfigs,
	}
}

// HTTPHandler returns a HTTP handler for the server.
func (s *Server) HTTPHandler() http.Handler {
	mux := http.NewServeMux()
	// iPXE
	mux.Handle("/ipxe/", ipxeMux(s.bootConfigs))
	// Pixiecore API Server
	mux.Handle(pixiecorePath, pixiecoreHandler(s.bootConfigs))
	// Kernel and Initrd Images
	mux.Handle("/images/", http.StripPrefix("/images/", http.FileServer(http.Dir("static"))))
	return mux
}
