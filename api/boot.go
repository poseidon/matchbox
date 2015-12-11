package api

// A BootAdapter maps MachineAttrs to a BootConfig which should be used.
type BootAdapter interface {
	// Get returns the BootConfig to boot the machine with the given attributes
	Get(attrs MachineAttrs) (*BootConfig, error)
}

// BootConfig defines the kernel image, kernel options, and initrds to boot
// on a client machine.
type BootConfig struct {
	// the URL of the kernel boot image
	Kernel string `json:"kernel"`
	// the initrd URLs which will be flattened into a single filesystem
	Initrd []string `json:"initrd"`
	// command line arguments to the kernel
	Cmdline map[string]interface{} `json:"cmdline"`
}
