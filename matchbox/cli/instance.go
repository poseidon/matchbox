package cli

import (
	"github.com/spf13/cobra"
)

// instanceCmd represents the instance command
var instanceCmd = &cobra.Command{
	Use:   "instance",
	Short: "View observed machine instances",
	Long:  `View observed machine instances`,
}

func init() {
	RootCmd.AddCommand(instanceCmd)
}
