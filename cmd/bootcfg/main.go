package main

import (
	"flag"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/coreos/pkg/capnslog"
	"github.com/coreos/pkg/flagutil"

	"github.com/coreos/coreos-baremetal/bootcfg/api"
	"github.com/coreos/coreos-baremetal/bootcfg/config"
	"github.com/coreos/coreos-baremetal/bootcfg/sign"
	"github.com/coreos/coreos-baremetal/bootcfg/storage"
)

var (
	// version provided by compile time flag: -ldflags "-X main.version $GIT_SHA"
	version = "was not built properly"
	log     = capnslog.NewPackageLogger("github.com/coreos/coreos-baremetal/cmd/bootcfg", "main")
)

func main() {
	flags := struct {
		address     string
		configPath  string
		dataPath    string
		assetsPath  string
		keyRingPath string
		logLevel    string
		version     bool
		help        bool
	}{}
	flag.StringVar(&flags.address, "address", "127.0.0.1:8080", "HTTP listen address")
	flag.StringVar(&flags.configPath, "config", "./data/config.yaml", "Path to config file")
	flag.StringVar(&flags.dataPath, "data-path", "./data", "Path to data directory")
	flag.StringVar(&flags.assetsPath, "assets-path", "./assets", "Path to static assets")
	flag.StringVar(&flags.keyRingPath, "key-ring-path", "", "Path to a private keyring file")
	// available log levels https://godoc.org/github.com/coreos/pkg/capnslog#LogLevel
	flag.StringVar(&flags.logLevel, "log-level", "info", "Set the logging level")
	// subcommands
	flag.BoolVar(&flags.version, "version", false, "print version and exit")
	flag.BoolVar(&flags.help, "help", false, "print usage and exit")

	// parse command-line and environment variable arguments
	flag.Parse()
	if err := flagutil.SetFlagsFromEnv(flag.CommandLine, "BOOTCFG"); err != nil {
		log.Fatal(err.Error())
	}
	// restrict OpenPGP passphrase to pass via environment variable only
	passphrase := os.Getenv("BOOTCFG_PASSPHRASE")

	if flags.version {
		fmt.Println(version)
		return
	}

	if flags.help {
		flag.Usage()
		return
	}

	// validate arguments
	if url, err := url.Parse(flags.address); err != nil || url.String() == "" {
		log.Fatal("A valid HTTP listen address is required")
	}
	if finfo, err := os.Stat(flags.configPath); err != nil || finfo.IsDir() {
		log.Fatal("A path to a config file is required")
	}
	if finfo, err := os.Stat(flags.dataPath); err != nil || !finfo.IsDir() {
		log.Fatal("A path to a data directory is required")
	}
	if finfo, err := os.Stat(flags.assetsPath); err != nil || !finfo.IsDir() {
		log.Fatal("A path to an assets directory is required")
	}

	// logging setup
	lvl, err := capnslog.ParseLevel(strings.ToUpper(flags.logLevel))
	if err != nil {
		log.Fatalf("invalid log-level: %v", err)
	}
	capnslog.SetGlobalLogLevel(lvl)
	capnslog.SetFormatter(capnslog.NewPrettyFormatter(os.Stdout, false))

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

	// load bootstrap config
	cfg, err := config.LoadConfig(flags.configPath)
	if err != nil {
		log.Fatal(err)
	}

	// storage
	store := storage.NewFileStore(&storage.Config{
		Dir:    flags.dataPath,
		Groups: cfg.Groups,
	})

	// HTTP server
	config := &api.Config{
		Store:         store,
		AssetsPath:    flags.assetsPath,
		Signer:        signer,
		ArmoredSigner: armoredSigner,
	}
	server := api.NewServer(config)
	log.Infof("starting bootcfg HTTP server on %s", flags.address)
	err = http.ListenAndServe(flags.address, server.HTTPHandler())
	if err != nil {
		log.Fatalf("failed to start listening: %v", err)
	}
}
