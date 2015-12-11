package main

import (
	"log"
	"net/http"

	"github.com/coreos/coreos-baremetal/api"
)

const address = ":8080"

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
	// load some boot configs
	bootAdapter := api.NewMapBootAdapter()
	bootAdapter.AddUUID("8a549bf5-075c-4372-8b0d-ce7844faa48c", CoreOSLocalAutoLogin )
	bootAdapter.SetDefault(CoreOSLocal)
	// api server
	server := api.NewServer(bootAdapter)
	log.Printf("Starting boot config server")
	err := http.ListenAndServe(address, server.HTTPHandler())
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
