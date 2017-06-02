package cli

import (
	"github.com/spf13/cobra"
)

// genericCmd represents the generic command
var genericCmd = &cobra.Command{
	Use:   "generic",
	Short: "Manage Generic templates",
	Long:  `Manage Generic templates`,
}

func init() {
	RootCmd.AddCommand(genericCmd)
}
