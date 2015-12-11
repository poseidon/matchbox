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
	Kernel:  "/images/kernel/coreos_production_pxe.vmlinuz",
	Initrd:  []string{"/images/initrd/coreos_production_pxe_image.cpio.gz"},
	Cmdline: map[string]interface{}{},
}

func main() {
	// initial boot configs
	bootAdapter := api.NewMapBootAdapter()
	bootAdapter.SetDefault(CoreOSStable)
	// api server
	server := api.NewServer(bootAdapter)
	h := server.HTTPHandler()
	log.Printf("Starting boot config server")
	err := http.ListenAndServe(address, h)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
