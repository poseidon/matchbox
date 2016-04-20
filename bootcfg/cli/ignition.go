package cli

import (
	"github.com/spf13/cobra"
)

// ignitionCmd represents the ignition command
var ignitionCmd = &cobra.Command{
	Use:   "ignition",
	Short: "Manage Ignition templates",
	Long:  `Manage Ignition templates`,
}

func init() {
	RootCmd.AddCommand(ignitionCmd)
}
