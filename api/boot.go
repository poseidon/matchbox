package api

// BootConfig defines a kernel image, kernel options, and initrds to boot.
type BootConfig struct {
	// the URL of the kernel image
	Kernel string `json:"kernel"`
	// the init RAM filesystem URLs
	Initrd []string `json:"initrd"`
	// command line kernel options
	Cmdline map[string]interface{} `json:"cmdline"`
}
