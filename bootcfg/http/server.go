package http

import (
	"net/http"

	"github.com/Sirupsen/logrus"

	"github.com/coreos/coreos-baremetal/bootcfg/server"
	"github.com/coreos/coreos-baremetal/bootcfg/sign"
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

	// bootcfg version
	mux.Handle("/", s.logRequest(versionHandler()))
	// Boot via GRUB
	mux.Handle("/grub", s.logRequest(NewHandler(s.selectProfile(s.core, s.grubHandler()))))
	// Boot via iPXE
	mux.Handle("/boot.ipxe", s.logRequest(ipxeInspect()))
	mux.Handle("/boot.ipxe.0", s.logRequest(ipxeInspect()))
	mux.Handle("/ipxe", s.logRequest(NewHandler(s.selectProfile(s.core, s.ipxeHandler()))))
	// Boot via Pixiecore
	mux.Handle("/pixiecore/v1/boot/", s.logRequest(NewHandler(s.pixiecoreHandler(s.core))))
	// Ignition Config
	mux.Handle("/ignition", s.logRequest(NewHandler(s.selectGroup(s.core, s.ignitionHandler(s.core)))))
	// Cloud-Config
	mux.Handle("/cloud", s.logRequest(NewHandler(s.selectGroup(s.core, s.cloudHandler(s.core)))))
	// Generic template
	mux.Handle("/generic", s.logRequest(NewHandler(s.selectGroup(s.core, s.genericHandler(s.core)))))
	// Metadata
	mux.Handle("/metadata", s.logRequest(NewHandler(s.selectGroup(s.core, s.metadataHandler()))))

	// Signatures
	if s.signer != nil {
		signerChain := func(next http.Handler) http.Handler {
			return s.logRequest(sign.SignatureHandler(s.signer, next))
		}
		mux.Handle("/grub.sig", signerChain(NewHandler(s.selectProfile(s.core, s.grubHandler()))))
		mux.Handle("/boot.ipxe.sig", signerChain(ipxeInspect()))
		mux.Handle("/boot.ipxe.0.sig", signerChain(ipxeInspect()))
		mux.Handle("/ipxe.sig", signerChain(NewHandler(s.selectProfile(s.core, s.ipxeHandler()))))
		mux.Handle("/pixiecore/v1/boot.sig/", signerChain(NewHandler(s.pixiecoreHandler(s.core))))
		mux.Handle("/ignition.sig", signerChain(NewHandler(s.selectGroup(s.core, s.ignitionHandler(s.core)))))
		mux.Handle("/cloud.sig", signerChain(NewHandler(s.selectGroup(s.core, s.cloudHandler(s.core)))))
		mux.Handle("/generic.sig", signerChain(NewHandler(s.selectGroup(s.core, s.genericHandler(s.core)))))
		mux.Handle("/metadata.sig", signerChain(NewHandler(s.selectGroup(s.core, s.metadataHandler()))))
	}
	if s.armoredSigner != nil {
		signerChain := func(next http.Handler) http.Handler {
			return s.logRequest(sign.SignatureHandler(s.armoredSigner, next))
		}
		mux.Handle("/grub.asc", signerChain(NewHandler(s.selectProfile(s.core, s.grubHandler()))))
		mux.Handle("/boot.ipxe.asc", signerChain(ipxeInspect()))
		mux.Handle("/boot.ipxe.0.asc", signerChain(ipxeInspect()))
		mux.Handle("/ipxe.asc", signerChain(NewHandler(s.selectProfile(s.core, s.ipxeHandler()))))
		mux.Handle("/pixiecore/v1/boot.asc/", signerChain(NewHandler(s.pixiecoreHandler(s.core))))
		mux.Handle("/ignition.asc", signerChain(NewHandler(s.selectGroup(s.core, s.ignitionHandler(s.core)))))
		mux.Handle("/cloud.asc", signerChain(NewHandler(s.selectGroup(s.core, s.cloudHandler(s.core)))))
		mux.Handle("/generic.asc", signerChain(NewHandler(s.selectGroup(s.core, s.genericHandler(s.core)))))
		mux.Handle("/metadata.asc", signerChain(NewHandler(s.selectGroup(s.core, s.metadataHandler()))))
	}

	// kernel, initrd, and TLS assets
	if s.assetsPath != "" {
		mux.Handle("/assets/", s.logRequest(http.StripPrefix("/assets/", http.FileServer(http.Dir(s.assetsPath)))))
	}
	return mux
}
