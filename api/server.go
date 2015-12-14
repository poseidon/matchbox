package api

import (
	"net/http"

	"github.com/coreos/pkg/capnslog"
)

var log = capnslog.NewPackageLogger("github.com/coreos/coreos-baremetal", "api")

// Config configures the api Server.
type Config struct {
	// Path to static image assets
	ImagePath string
	// Adapter which provides BootConfigs
	BootAdapter BootAdapter
}

// Server serves iPXE/Pixiecore boot configs and hosts images.
type Server struct {
	imagePath string
	bootConfigs BootAdapter
}

// NewServer returns a new Server which uses the given BootAdapter.
func NewServer(config *Config) *Server {
	return &Server{
		imagePath: config.ImagePath,
		bootConfigs: config.BootAdapter,
	}
}

// HTTPHandler returns a HTTP handler for the server.
func (s *Server) HTTPHandler() http.Handler {
	mux := http.NewServeMux()
	// iPXE
	mux.Handle("/ipxe/", ipxeMux(s.bootConfigs))
	// Pixiecore
	mux.Handle(pixiecorePath, pixiecoreHandler(s.bootConfigs))
	// Kernel and Initrd Images
	mux.Handle("/images/", http.StripPrefix("/images/", http.FileServer(http.Dir(s.imagePath))))
	return mux
}
