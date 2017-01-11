package cli

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

// CLI Tool errors
// http://tldp.org/LDP/abs/html/exitcodes.html
const (
	ExitSuccess = iota
	ExitError
	ExitBadConnection
	ExitBadArgs = 128
)

func exitWithError(code int, err error) {
	fmt.Fprintln(os.Stderr, "Error: ", err)
	os.Exit(code)
}

func usageError(cmd *cobra.Command, format string, args ...interface{}) error {
	msg := fmt.Sprintf(format, args...)
	return fmt.Errorf("%s\nSee '%s -h' for help", msg, cmd.CommandPath())
}
