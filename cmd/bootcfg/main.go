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

// Example Boot Configs

var CoreOSStable = &api.BootConfig{
	Kernel:  "http://stable.release.core-os.net/amd64-usr/current/coreos_production_pxe.vmlinuz",
	Initrd:  []string{"http://stable.release.core-os.net/amd64-usr/current/coreos_production_pxe_image.cpio.gz"},
	Cmdline: map[string]interface{}{},
}

var CoreOSBeta = &api.BootConfig{
	Kernel:  "http://beta.release.core-os.net/amd64-usr/current/coreos_production_pxe.vmlinuz",
	Initrd:  []string{"http://beta.release.core-os.net/amd64-usr/current/coreos_production_pxe_image.cpio.gz"},
	Cmdline: map[string]interface{}{},
}

var CoreOSAlpha = &api.BootConfig{
	Kernel:  "http://alpha.release.core-os.net/amd64-usr/current/coreos_production_pxe.vmlinuz",
	Initrd:  []string{"http://alpha.release.core-os.net/amd64-usr/current/coreos_production_pxe_image.cpio.gz"},
	Cmdline: map[string]interface{}{},
}

var CoreOSLocal = &api.BootConfig{
	Kernel: "/images/coreos_production_pxe.vmlinuz",
	Initrd: []string{"/images/coreos_production_pxe_image.cpio.gz"},
	Cmdline: map[string]interface{}{},
}

var CoreOSLocalAutoLogin = &api.BootConfig{
	Kernel: "/images/coreos_production_pxe.vmlinuz",
	Initrd: []string{"/images/coreos_production_pxe_image.cpio.gz"},
	Cmdline: map[string]interface{}{
		"coreos.autologin": "",
	},
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
	bootAdapter.SetUUID("8a549bf5-075c-4372-8b0d-ce7844faa48c", CoreOSLocalAutoLogin )
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
