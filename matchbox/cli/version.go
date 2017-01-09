package cli

import (
	"fmt"
	"runtime"

	"github.com/spf13/cobra"
)

var (
	// version provided by compile time flag: -ldflags "-X main.version $GIT_SHA"
	version = "was not built properly"
	// versionCmd represents the version command
	versionCmd = &cobra.Command{
		Use:   "version",
		Short: "Print the version and exit",
		Long:  `Print the version of the bootcmd client`,
		Run:   displayVersion,
	}
)

func displayVersion(cmd *cobra.Command, args []string) {
	fmt.Printf("bootcmd Version: %s\n", version)
	fmt.Printf("Go Version: %s\n", runtime.Version())
	fmt.Printf("Go OS/Arch: %s/%s\n", runtime.GOOS, runtime.GOARCH)
}

func init() {
	RootCmd.AddCommand(versionCmd)
}
