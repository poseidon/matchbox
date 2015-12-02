package server

import (
	"net/http"
	"encoding/json"
	"log"
	"strings"
)

// Server manages boot and cloud configs for hosts by MAC address or UUID.
type Server struct {
	bootConfigProvider BootConfigProvider
}

// NewServer returns a new Server which uses the given BootConfigProvider.
func NewServer(bootConfigProvider BootConfigProvider) *Server {
	return &Server{
		bootConfigProvider: bootConfigProvider,
	}
}

// HTTPHandler returns a HTTP handler for the server.
func (s *Server) HTTPHandler() http.Handler {
	mux := http.NewServeMux()
	// Pixiecore API Server
	mux.Handle("/v1/boot/", pixiecoreHandler(s.bootConfigProvider))
	// Kernel and File Server
	mux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	return mux
}

// pixiecoreHandler implements the Pixiecore API Server Spec.
func pixiecoreHandler(bootConfigs BootConfigProvider) http.Handler {
	fn := func(w http.ResponseWriter, req *http.Request) {
		remainder := strings.TrimPrefix(req.URL.String(), "/v1/boot/")
		log.Printf("render boot config for %s", remainder)
		bootConfig, err := bootConfigs.Get(remainder)
		if err != nil {
			http.Error(w, err.Error(), 404)
			return
		}
		json.NewEncoder(w).Encode(bootConfig)
	}
	return http.HandlerFunc(fn)
}
