package cli

import (
	"github.com/spf13/cobra"
)

// groupCmd represents the group command
var groupCmd = &cobra.Command{
	Use:   "group",
	Short: "Manage machine groups",
	Long:  `List and describe machine groups`,
}

func init() {
	RootCmd.AddCommand(groupCmd)
}
