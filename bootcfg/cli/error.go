package cmd

import (
	"fmt"
	"os"
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
