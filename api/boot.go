package api

// BootConfig defines the kernel image, kernel options, and initrds to boot
// a client machine.
type BootConfig struct {
	// the URL of the kernel boot image
	Kernel string `json:"kernel"`
	// the initrd URLs which will be flattened into a single filesystem
	Initrd []string `json:"initrd"`
	// command line arguments to the kernel
	Cmdline map[string]interface{} `json:"cmdline"`
}
