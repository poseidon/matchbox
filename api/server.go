package api

import (
	"net/http"

	"github.com/coreos/pkg/capnslog"
)

var log = capnslog.NewPackageLogger("github.com/coreos/coreos-baremetal", "api")

// Config configures the api Server.
type Config struct {
	// Store for configs (boot, cloud)
	Store Store
	// Path to static image assets
	ImagePath string
}

// Server serves boot and cloud configs for PXE-based clients.
type Server struct {
	store     Store
	imagePath string
}

// NewServer returns a new Server.
func NewServer(config *Config) *Server {
	return &Server{
		store:     config.Store,
		imagePath: config.ImagePath,
	}
}

// HTTPHandler returns a HTTP handler for the server.
func (s *Server) HTTPHandler() http.Handler {
	mux := http.NewServeMux()
	// machines
	newMachineResource(mux, "/machine/", s.store)
	// named specs
	newSpecResource(mux, "/spec/", s.store)

	// Baremetal
	// iPXE
	mux.Handle("/boot.ipxe", logRequests(ipxeInspect()))
	mux.Handle("/boot.ipxe.0", logRequests(ipxeInspect()))
	mux.Handle("/ipxe", logRequests(ipxeHandler(s.store)))
	// Pixiecore
	mux.Handle("/pixiecore/v1/boot/", logRequests(pixiecoreHandler(s.store)))

	// cloud configs
	mux.Handle("/cloud", logRequests(cloudHandler(s.store)))
	// Kernel and Initrd Images
	mux.Handle("/images/", http.StripPrefix("/images/", http.FileServer(http.Dir(s.imagePath))))
	return mux
}
