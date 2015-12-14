package main

import (
	"net/http"
	"flag"
	"os"
	"net/url"
	"strings"

	"github.com/coreos/coreos-baremetal/api"
	"github.com/coreos/pkg/flagutil"
	"github.com/coreos/pkg/capnslog"
)

var log = capnslog.NewPackageLogger("github.com/coreos/coreos-baremetal/cmd/bootcfg", "main")

var CoreOSLocal = &api.BootConfig{
	Kernel: "/images/stable/coreos_production_pxe.vmlinuz",
	Initrd: []string{"/images/stable/coreos_production_pxe_image.cpio.gz"},
	Cmdline: map[string]interface{}{},
}

func main() {
	flags := flag.NewFlagSet("bootcfg", flag.ExitOnError)
	address := flags.String("address", "127.0.0.1:8080", "HTTP listen address")
	imagesPath := flags.String("images-path", "./images", "Path to static image assets")
	// available log levels https://godoc.org/github.com/coreos/pkg/capnslog#LogLevel
	logLevel := flags.String("log-level", "info", "Set the logging level")

	// parse command-line and environment variable arguments
	if err := flags.Parse(os.Args[1:]); err != nil {
		log.Fatal(err.Error())
	}
	if err := flagutil.SetFlagsFromEnv(flags, "BOOTCFG"); err != nil {
		log.Fatal(err.Error())
	}

	// validate arguments
	if url, err := url.Parse(*address); err != nil || url.String() == "" {
		log.Fatal("A valid HTTP listen address is required")
	}
	if finfo, err := os.Stat(*imagesPath); err != nil || !finfo.IsDir() {
		log.Fatal("A path to an image assets directory is required")
	}

	// logging setup
	lvl, err := capnslog.ParseLevel(strings.ToUpper(*logLevel))
	if err != nil {
		log.Fatalf("Invalid log-level: %s", err.Error())
	}
	capnslog.SetGlobalLogLevel(lvl)
	capnslog.SetFormatter(capnslog.NewPrettyFormatter(os.Stdout, false))

	// load some boot configs
	bootAdapter := api.NewMapBootAdapter()
	bootAdapter.SetDefault(CoreOSLocal)

	config := &api.Config{
		ImagePath: *imagesPath,
		BootAdapter: bootAdapter,
	}

	// API server
	server := api.NewServer(config)
	log.Infof("Starting bootcfg API Server on %s", *address)
	err = http.ListenAndServe(*address, server.HTTPHandler())
	if err != nil {
		log.Fatalf("failed to start listening: %s", err)
	}
}
