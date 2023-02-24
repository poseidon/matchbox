package main

import (
	"flag"
	"fmt"
	"net"
	"net/http"
	"os"
	"strings"

	"github.com/coreos/pkg/flagutil"
	web "github.com/poseidon/matchbox/matchbox/http"
	"github.com/poseidon/matchbox/matchbox/rpc"
	"github.com/poseidon/matchbox/matchbox/server"
	"github.com/poseidon/matchbox/matchbox/sign"
	"github.com/poseidon/matchbox/matchbox/storage"
	"github.com/poseidon/matchbox/matchbox/tlsutil"
	"github.com/poseidon/matchbox/matchbox/version"
	"github.com/sirupsen/logrus"
)

var (
	// Defaults to info logging
	log = logrus.New()
)

type CliFlags struct {
		address      string
		rpcAddress   string
		dataPath     string
		assetsPath   string
		etcdEndpoints string
		logLevel     string
		grpcCAFile   string
		grpcCertFile string
		grpcKeyFile  string
		storageType  string
		tlsCertFile  string
		tlsKeyFile   string
		tlsEnabled   bool
		keyRingPath  string
		version      bool
		help         bool
}


func main() {
	var flags CliFlags

	flag.StringVar(&flags.address, "address", "127.0.0.1:8080", "HTTP listen address")
	flag.StringVar(&flags.rpcAddress, "rpc-address", "", "RPC listen address")

	flag.StringVar(&flags.dataPath, "data-path", "/var/lib/matchbox", "Path to data directory")
	flag.StringVar(&flags.storageType, "storage-type", "file", "Type of storage to use for data (file, etcd)")
	flag.StringVar(&flags.etcdEndpoints, "etcd-endpoints", "127.0.0.1:2380", "Comme-separated list of host:port etcd endpoints")

	flag.StringVar(&flags.assetsPath, "assets-path", "/var/lib/matchbox/assets", "Path to static assets")

	// Log levels https://github.com/sirupsen/logrus/blob/master/logrus.go#L36
	flag.StringVar(&flags.logLevel, "log-level", "info", "Set the logging level")

	// gRPC Server TLS
	flag.StringVar(&flags.grpcCertFile, "cert-file", "/etc/matchbox/server.crt", "Path to the server TLS certificate file")
	flag.StringVar(&flags.grpcKeyFile, "key-file", "/etc/matchbox/server.key", "Path to the server TLS key file")

	// gRPC TLS Client Authentication
	flag.StringVar(&flags.grpcCAFile, "ca-file", "/etc/matchbox/ca.crt", "Path to the CA verify and authenticate client certificates")

	// Signing
	flag.StringVar(&flags.keyRingPath, "key-ring-path", "", "Path to a private keyring file")

	// SSL flags
	flag.StringVar(&flags.tlsCertFile, "web-cert-file", "/etc/matchbox/ssl/server.crt", "Path to the server TLS certificate file")
	flag.StringVar(&flags.tlsKeyFile, "web-key-file", "/etc/matchbox/ssl/server.key", "Path to the server TLS key file")
	flag.BoolVar(&flags.tlsEnabled, "web-ssl", false, "True to enable HTTPS")

	// subcommands
	flag.BoolVar(&flags.version, "version", false, "print version and exit")
	flag.BoolVar(&flags.help, "help", false, "print usage and exit")

	// parse command-line and environment variable arguments
	flag.Parse()
	if err := flagutil.SetFlagsFromEnv(flag.CommandLine, "MATCHBOX"); err != nil {
		log.Fatal(err.Error())
	}
	// restrict OpenPGP passphrase to pass via environment variable only
	passphrase := os.Getenv("MATCHBOX_PASSPHRASE")

	if flags.version {
		fmt.Println(version.Version)
		return
	}

	if flags.help {
		flag.Usage()
		return
	}

	// validate arguments
	if flags.storageType == "file" {
		if finfo, err := os.Stat(flags.dataPath); err != nil || !finfo.IsDir() {
			log.Fatal("A valid -data-path is required")
		}
	}
	if flags.assetsPath != "" {
		if finfo, err := os.Stat(flags.assetsPath); err != nil || !finfo.IsDir() {
			log.Fatalf("Provide a valid -assets-path or '' to disable asset serving: %s", flags.assetsPath)
		}
	}
	if flags.rpcAddress != "" {
		if _, err := os.Stat(flags.grpcCertFile); err != nil {
			log.Fatalf("Provide a valid TLS server certificate with -cert-file: %v", err)
		}
		if _, err := os.Stat(flags.grpcKeyFile); err != nil {
			log.Fatalf("Provide a valid TLS server key with -key-file: %v", err)
		}
		if _, err := os.Stat(flags.grpcCAFile); err != nil {
			log.Fatalf("Provide a valid TLS certificate authority for authorizing client certificates: %v", err)
		}
	}
	if flags.tlsEnabled {
		if _, err := os.Stat(flags.tlsCertFile); err != nil {
			log.Fatalf("Provide a valid SSL server certificate with -web-cert-file: %v", err)
		}
		if _, err := os.Stat(flags.tlsKeyFile); err != nil {
			log.Fatalf("Provide a valid SSL server key with -web-key-file: %v", err)
		}
	}

	// logging setup
	lvl, err := logrus.ParseLevel(flags.logLevel)
	if err != nil {
		log.Fatalf("invalid log-level: %v", err)
	}
	log.Level = lvl

	// (optional) signing
	var signer, armoredSigner sign.Signer
	if flags.keyRingPath != "" {
		entity, err := sign.LoadGPGEntity(flags.keyRingPath, passphrase)
		if err != nil {
			log.Fatal(err)
		}
		signer = sign.NewGPGSigner(entity)
		armoredSigner = sign.NewArmoredGPGSigner(entity)
	}

	// storage
	store := createStore(&flags, log)

	// core logic
	server := server.NewServer(&server.Config{
		Store: store,
	})

	// gRPC Server (feature disabled by default)
	if flags.rpcAddress != "" {
		log.Infof("Starting matchbox gRPC server on %s", flags.rpcAddress)
		log.Infof("Using TLS server certificate: %s", flags.grpcCertFile)
		log.Infof("Using TLS server key: %s", flags.grpcKeyFile)
		log.Infof("Using CA certificate: %s to authenticate client certificates", flags.grpcCAFile)
		lis, err := net.Listen("tcp", flags.rpcAddress)
		if err != nil {
			log.Fatalf("failed to start listening: %v", err)
		}
		tlsinfo := tlsutil.TLSInfo{
			CertFile: flags.grpcCertFile,
			KeyFile:  flags.grpcKeyFile,
			CAFile:   flags.grpcCAFile,
		}
		tlscfg, err := tlsinfo.ServerConfig()
		if err != nil {
			log.Fatalf("Invalid TLS credentials: %v", err)
		}
		grpcServer := rpc.NewServer(server, tlscfg)
		go grpcServer.Serve(lis)
		defer grpcServer.Stop()
	}

	config := &web.Config{
		Core:          server,
		Logger:        log,
		AssetsPath:    flags.assetsPath,
		Signer:        signer,
		ArmoredSigner: armoredSigner,
	}
	httpServer := web.NewServer(config)

	if flags.tlsEnabled {
		// HTTPS Server
		log.Infof("Starting matchbox HTTPS server on %s", flags.address)
		log.Infof("Using SSL server certificate: %s", flags.tlsCertFile)
		log.Infof("Using SSL server key: %s", flags.tlsKeyFile)
		err = http.ListenAndServeTLS(flags.address, flags.tlsCertFile, flags.tlsKeyFile, httpServer.HTTPHandler())
		if err != nil {
			log.Fatalf("failed to start listening: %v", err)
		}
	} else {
		// HTTP Server
		log.Infof("Starting matchbox HTTP server on %s", flags.address)
		err = http.ListenAndServe(flags.address, httpServer.HTTPHandler())
		if err != nil {
			log.Fatalf("failed to start listening: %v", err)
		}
	}

}

func createStore(flags *CliFlags, logger *logrus.Logger) storage.Store {
	var store storage.Store

	switch flags.storageType {
	case "file":
		store = storage.NewFileStore(&storage.Config{
			Root:   flags.dataPath,
			Logger: log,
		})
	case "etcd":
		endPoints := strings.Split(flags.etcdEndpoints, ",")
		store = storage.NewEtcdStore(endPoints,logger)
	}

	return store
}