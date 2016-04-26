package http

import (
	"net/http"

	"github.com/coreos/pkg/capnslog"

	"github.com/coreos/coreos-baremetal/bootcfg/server"
	"github.com/coreos/coreos-baremetal/bootcfg/sign"
	"github.com/coreos/coreos-baremetal/bootcfg/storage"
)

var log = capnslog.NewPackageLogger("github.com/coreos/coreos-baremetal/bootcfg", "http")

// Config configures the api Server.
type Config struct {
	// Store for configs
	Store storage.Store
	// Path to static assets
	AssetsPath string
	// config signers (.sig and .asc)
	Signer        sign.Signer
	ArmoredSigner sign.Signer
}

// Server serves boot and provisioning configs to machines via HTTP.
type Server struct {
	store         storage.Store
	assetsPath    string
	signer        sign.Signer
	armoredSigner sign.Signer
}

// NewServer returns a new Server.
func NewServer(config *Config) *Server {
	return &Server{
		store:         config.Store,
		assetsPath:    config.AssetsPath,
		signer:        config.Signer,
		armoredSigner: config.ArmoredSigner,
	}
}

// HTTPHandler returns a HTTP handler for the server.
func (s *Server) HTTPHandler() http.Handler {
	mux := http.NewServeMux()
	srv := server.NewServer(&server.Config{s.store})

	// Boot via GRUB
	mux.Handle("/grub", logRequests(NewHandler(selectProfile(srv, grubHandler()))))
	// Boot via iPXE
	mux.Handle("/boot.ipxe", logRequests(ipxeInspect()))
	mux.Handle("/boot.ipxe.0", logRequests(ipxeInspect()))
	mux.Handle("/ipxe", logRequests(NewHandler(selectProfile(srv, ipxeHandler()))))
	// Boot via Pixiecore
	mux.Handle("/pixiecore/v1/boot/", logRequests(NewHandler(pixiecoreHandler(srv))))
	// Ignition Config
	mux.Handle("/ignition", logRequests(NewHandler(selectGroup(srv, ignitionHandler(srv)))))
	// Cloud-Config
	mux.Handle("/cloud", logRequests(NewHandler(selectGroup(srv, cloudHandler(srv)))))
	// metadata
	mux.Handle("/metadata", logRequests(NewHandler(selectGroup(srv, metadataHandler()))))

	// Signatures
	if s.signer != nil {
		signerChain := func(next http.Handler) http.Handler {
			return logRequests(sign.SignatureHandler(s.signer, next))
		}
		mux.Handle("/grub.sig", signerChain(NewHandler(selectProfile(srv, grubHandler()))))
		mux.Handle("/boot.ipxe.sig", signerChain(ipxeInspect()))
		mux.Handle("/boot.ipxe.0.sig", signerChain(ipxeInspect()))
		mux.Handle("/ipxe.sig", signerChain(NewHandler(selectProfile(srv, ipxeHandler()))))
		mux.Handle("/pixiecore/v1/boot.sig/", signerChain(NewHandler(pixiecoreHandler(srv))))
		mux.Handle("/ignition.sig", signerChain(NewHandler(selectGroup(srv, ignitionHandler(srv)))))
		mux.Handle("/cloud.sig", signerChain(NewHandler(selectGroup(srv, cloudHandler(srv)))))
		mux.Handle("/metadata.sig", signerChain(NewHandler(selectGroup(srv, metadataHandler()))))
	}
	if s.armoredSigner != nil {
		signerChain := func(next http.Handler) http.Handler {
			return logRequests(sign.SignatureHandler(s.armoredSigner, next))
		}
		mux.Handle("/grub.asc", signerChain(NewHandler(selectProfile(srv, grubHandler()))))
		mux.Handle("/boot.ipxe.asc", signerChain(ipxeInspect()))
		mux.Handle("/boot.ipxe.0.asc", signerChain(ipxeInspect()))
		mux.Handle("/ipxe.asc", signerChain(NewHandler(selectProfile(srv, ipxeHandler()))))
		mux.Handle("/pixiecore/v1/boot.asc/", signerChain(NewHandler(pixiecoreHandler(srv))))
		mux.Handle("/ignition.asc", signerChain(NewHandler(selectGroup(srv, ignitionHandler(srv)))))
		mux.Handle("/cloud.asc", signerChain(NewHandler(selectGroup(srv, cloudHandler(srv)))))
		mux.Handle("/metadata.asc", signerChain(NewHandler(selectGroup(srv, metadataHandler()))))
	}

	// kernel, initrd, and TLS assets
	if s.assetsPath != "" {
		mux.Handle("/assets/", http.StripPrefix("/assets/", http.FileServer(http.Dir(s.assetsPath))))
	}
	return mux
}
