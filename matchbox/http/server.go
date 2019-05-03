package http

import (
	"net/http"

	"github.com/sirupsen/logrus"

	"github.com/poseidon/matchbox/matchbox/server"
	"github.com/poseidon/matchbox/matchbox/sign"
)

// Config configures a Server.
type Config struct {
	Core   server.Server
	Logger *logrus.Logger
	// Path to static assets
	AssetsPath string
	// config signers (.sig and .asc)
	Signer        sign.Signer
	ArmoredSigner sign.Signer
}

// Server serves boot and provisioning configs to machines via HTTP.
type Server struct {
	core          server.Server
	logger        *logrus.Logger
	assetsPath    string
	signer        sign.Signer
	armoredSigner sign.Signer
}

// NewServer returns a new Server.
func NewServer(config *Config) *Server {
	return &Server{
		core:          config.Core,
		logger:        config.Logger,
		assetsPath:    config.AssetsPath,
		signer:        config.Signer,
		armoredSigner: config.ArmoredSigner,
	}
}

// HTTPHandler returns a HTTP handler for the server.
func (s *Server) HTTPHandler() http.Handler {
	mux := http.NewServeMux()

	chain := func(next http.Handler) http.Handler {
		return s.logRequest(next)
	}
	// matchbox version
	mux.Handle("/", s.logRequest(homeHandler()))
	// Boot via GRUB
	mux.Handle("/grub", chain(s.selectProfile(s.core, s.grubHandler())))
	// Boot via iPXE
	mux.Handle("/boot.ipxe", chain(ipxeInspect()))
	mux.Handle("/boot.ipxe.0", chain(ipxeInspect()))
	mux.Handle("/ipxe", chain(s.selectProfile(s.core, s.ipxeHandler())))
	// Ignition Config
	mux.Handle("/ignition", chain(s.selectGroup(s.core, s.ignitionHandler(s.core))))
	// Cloud-Config
	mux.Handle("/cloud", chain(s.selectGroup(s.core, s.cloudHandler(s.core))))
	// Generic template
	mux.Handle("/generic", chain(s.selectGroup(s.core, s.genericHandler(s.core))))
	// Metadata
	mux.Handle("/metadata", chain(s.selectGroup(s.core, s.metadataHandler())))

	// Signatures
	if s.signer != nil {
		signerChain := func(next http.Handler) http.Handler {
			return s.logRequest(sign.SignatureHandler(s.signer, next))
		}
		mux.Handle("/grub.sig", signerChain(s.selectProfile(s.core, s.grubHandler())))
		mux.Handle("/boot.ipxe.sig", signerChain(ipxeInspect()))
		mux.Handle("/boot.ipxe.0.sig", signerChain(ipxeInspect()))
		mux.Handle("/ipxe.sig", signerChain(s.selectProfile(s.core, s.ipxeHandler())))
		mux.Handle("/ignition.sig", signerChain(s.selectGroup(s.core, s.ignitionHandler(s.core))))
		mux.Handle("/cloud.sig", signerChain(s.selectGroup(s.core, s.cloudHandler(s.core))))
		mux.Handle("/generic.sig", signerChain(s.selectGroup(s.core, s.genericHandler(s.core))))
		mux.Handle("/metadata.sig", signerChain(s.selectGroup(s.core, s.metadataHandler())))
	}
	if s.armoredSigner != nil {
		signerChain := func(next http.Handler) http.Handler {
			return s.logRequest(sign.SignatureHandler(s.armoredSigner, next))
		}
		mux.Handle("/grub.asc", signerChain(s.selectProfile(s.core, s.grubHandler())))
		mux.Handle("/boot.ipxe.asc", signerChain(ipxeInspect()))
		mux.Handle("/boot.ipxe.0.asc", signerChain(ipxeInspect()))
		mux.Handle("/ipxe.asc", signerChain(s.selectProfile(s.core, s.ipxeHandler())))
		mux.Handle("/ignition.asc", signerChain(s.selectGroup(s.core, s.ignitionHandler(s.core))))
		mux.Handle("/cloud.asc", signerChain(s.selectGroup(s.core, s.cloudHandler(s.core))))
		mux.Handle("/generic.asc", signerChain(s.selectGroup(s.core, s.genericHandler(s.core))))
		mux.Handle("/metadata.asc", signerChain(s.selectGroup(s.core, s.metadataHandler())))
	}

	// kernel, initrd, and TLS assets
	if s.assetsPath != "" {
		mux.Handle("/assets/", s.logRequest(http.StripPrefix("/assets/", http.FileServer(http.Dir(s.assetsPath)))))
	}
	return mux
}
